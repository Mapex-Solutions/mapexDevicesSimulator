package session

import (
	"net"
	"sync"

	"github.com/gorilla/websocket"

	mqttclient "github.com/Mapex-Solutions/mapexGoKit/infrastructure/mqttclient"

	"simulator/packages/utils/lorawan/band"
	"simulator/packages/utils/lorawan/device"
	"simulator/service/src/modules/engine/application/ports"
)

// connectorRegistry maps a protocol id to its connector. A protocol absent from
// the map is not session-capable and falls back to the one-shot dispatcher.
type connectorRegistry struct {
	byProtocol map[string]ports.Connector
}

// mqttConnector opens persistent MQTT sessions (one live client per device).
type mqttConnector struct{}

// mqttSession is a live MQTT connection for one device: it publishes uplinks and
// receives downlinks on the same open client.
type mqttSession struct {
	client *mqttclient.Client
}

// lorawanConnector opens LoRaWAN device sessions, sharing one gateway link across
// every device that transmits through the same gateway.
type lorawanConnector struct {
	mu    sync.Mutex
	links map[string]*sharedLink
	devs  map[string]*device.Session // keyed by DevEUI; persists DevNonce across reconnects
}

// lorawanSession is one device's live LoRaWAN session over a shared gateway link.
type lorawanSession struct {
	dev    *device.Session
	link   *sharedLink
	conn   *lorawanConnector
	region band.Region
	addr   [4]byte

	// Identity carried for the console status frames.
	devEui  string
	joinEui string
	class   string // LoRaWAN device class: A or C
}

// linkTransport is one live connection to the LNS (Basics Station WS or Semtech
// UDP). It forwards uplinks and pumps received downlink PHYPayloads to onDownlink.
type linkTransport interface {
	sendUp(phy []byte, dr int, freq uint64, datr string, rssi int, snr float64) error
	connected() bool
	close()
}

// sharedLink is a gateway-level link shared by every LoRaWAN device on that gateway.
// Devices register a downlink handler; received frames are routed by the connector's
// router. A reference count tears the link down when the last device disconnects.
type sharedLink struct {
	key       string
	transport linkTransport
	router    *downlinkRouter
	refs      int
}

// downlinkRouter dispatches a received downlink to the right device: a join accept
// goes to the device currently awaiting one; a data downlink goes by DevAddr.
type downlinkRouter struct {
	mu      sync.RWMutex
	byAddr  map[[4]byte]func([]byte)
	pending func([]byte)
}

// udpTransport speaks the Semtech UDP packet-forwarder protocol.
type udpTransport struct {
	conn  *net.UDPConn
	gwEUI [8]byte
	token uint16
	tmst  uint32
	rxnb  uint32
	mu    sync.Mutex
	up    bool
}

// wsTransport speaks the Basics Station LNS protocol over a WebSocket.
type wsTransport struct {
	conn  *websocket.Conn
	xtime int64
	mu    sync.Mutex
	up    bool
}
