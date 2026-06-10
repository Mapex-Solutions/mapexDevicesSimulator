// Package udp frames the Semtech UDP packet-forwarder protocol (the legacy
// gateway-to-LNS transport). It encodes uplinks (PUSH_DATA) and keepalives
// (PULL_DATA) and decodes downlinks (PULL_RESP), carrying the raw PHYPayload as
// base64. It is transport framing only: opening the socket and timing live in the
// connector. Any LNS speaking Semtech UDP (ChirpStack, TTS, mapexLNS) interoperates.
package udp

import (
	"encoding/base64"
	"encoding/json"
	"errors"
)

// protocolVersion is the Semtech UDP protocol version this gateway speaks.
const protocolVersion byte = 0x02

// Packet identifiers exchanged between the gateway and the LNS.
const (
	idPushData byte = 0x00
	idPushAck  byte = 0x01
	idPullData byte = 0x02
	idPullResp byte = 0x03
	idPullAck  byte = 0x04
	idTxAck    byte = 0x05
)

// errShort is returned when a datagram is too short to carry a header.
var errShort = errors.New("udp: datagram too short")

// RxPk is one received uplink reported to the LNS inside a PUSH_DATA.
type RxPk struct {
	Time string  `json:"time,omitempty"`
	Tmst uint32  `json:"tmst"`
	Freq float64 `json:"freq"`
	Chan uint8   `json:"chan"`
	RFCh uint8   `json:"rfch"`
	Stat int8    `json:"stat"`
	Modu string  `json:"modu"`
	DatR string  `json:"datr"`
	CodR string  `json:"codr"`
	RSSI int     `json:"rssi"`
	LSNR float64 `json:"lsnr"`
	Size int     `json:"size"`
	Data string  `json:"data"` // base64 PHYPayload
}

// TxPk is one downlink the LNS schedules inside a PULL_RESP.
type TxPk struct {
	Imme bool    `json:"imme"`
	Tmst uint32  `json:"tmst"`
	Freq float64 `json:"freq"`
	RFCh uint8   `json:"rfch"`
	Powe int     `json:"powe"`
	Modu string  `json:"modu"`
	DatR string  `json:"datr"`
	CodR string  `json:"codr"`
	Size int     `json:"size"`
	Data string  `json:"data"` // base64 PHYPayload
}

type pushPayload struct {
	RxPk []RxPk `json:"rxpk"`
}

// Stat is the periodic gateway-statistics report a real packet forwarder sends so
// the LNS marks the gateway online (its "last seen"). It rides in a PUSH_DATA.
type Stat struct {
	Time string  `json:"time"`
	Lati float64 `json:"lati"`
	Long float64 `json:"long"`
	Alti int     `json:"alti"`
	RxNb uint32  `json:"rxnb"`
	RxOK uint32  `json:"rxok"`
	RxFw uint32  `json:"rxfw"`
	AckR float64 `json:"ackr"`
	DwNb uint32  `json:"dwnb"`
	TxNb uint32  `json:"txnb"`
}

type statPayload struct {
	Stat Stat `json:"stat"`
}

type pullRespPayload struct {
	TxPk TxPk `json:"txpk"`
}

// header builds the 4-byte (PULL_DATA also appends the gateway EUI) datagram head.
func header(token uint16, id byte) []byte {
	return []byte{protocolVersion, byte(token >> 8), byte(token), id}
}

// MarshalPushData frames an uplink: header + gateway EUI + JSON rxpk array. The
// PHYPayload is base64-encoded into the rxpk Data field by the caller-supplied pk.
func MarshalPushData(token uint16, gatewayEUI [8]byte, pk RxPk) ([]byte, error) {
	body, err := json.Marshal(pushPayload{RxPk: []RxPk{pk}})
	if err != nil {
		return nil, err
	}
	out := header(token, idPushData)
	out = append(out, gatewayEUI[:]...)
	return append(out, body...), nil
}

// MarshalPullData frames the keepalive that tells the LNS where to send downlinks.
func MarshalPullData(token uint16, gatewayEUI [8]byte) []byte {
	out := header(token, idPullData)
	return append(out, gatewayEUI[:]...)
}

// MarshalStat frames a gateway-statistics PUSH_DATA; the LNS uses it to mark the
// gateway online and record its last-seen time.
func MarshalStat(token uint16, gatewayEUI [8]byte, stat Stat) ([]byte, error) {
	body, err := json.Marshal(statPayload{Stat: stat})
	if err != nil {
		return nil, err
	}
	out := header(token, idPushData)
	out = append(out, gatewayEUI[:]...)
	return append(out, body...), nil
}

// EncodePHY base64-encodes a raw PHYPayload for an rxpk/txpk Data field.
func EncodePHY(phy []byte) string { return base64.StdEncoding.EncodeToString(phy) }

// Kind classifies an inbound datagram so the connector can route it.
type Kind int

const (
	KindUnknown Kind = iota
	KindPushAck
	KindPullAck
	KindPullResp
)

// Classify reads the datagram id and, for a PULL_RESP, decodes the downlink
// PHYPayload. Other server packets (acks) carry no payload.
func Classify(datagram []byte) (Kind, []byte, error) {
	if len(datagram) < 4 {
		return KindUnknown, nil, errShort
	}
	switch datagram[3] {
	case idPushAck:
		return KindPushAck, nil, nil
	case idPullAck:
		return KindPullAck, nil, nil
	case idPullResp:
		var p pullRespPayload
		if err := json.Unmarshal(datagram[4:], &p); err != nil {
			return KindPullResp, nil, err
		}
		phy, err := base64.StdEncoding.DecodeString(p.TxPk.Data)
		if err != nil {
			return KindPullResp, nil, err
		}
		return KindPullResp, phy, nil
	default:
		return KindUnknown, nil, nil
	}
}

// MarshalTxAck frames the acknowledgement a gateway sends after a PULL_RESP.
func MarshalTxAck(token uint16, gatewayEUI [8]byte) []byte {
	out := header(token, idTxAck)
	return append(out, gatewayEUI[:]...)
}
