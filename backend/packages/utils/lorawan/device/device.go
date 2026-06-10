// Package device is the stateful LoRaWAN "brain" of a simulated end device: it
// holds the keys and frame counters, performs the OTAA join (or ABP activation),
// builds data uplinks, and decodes downlinks — all over the vendored crypto + our
// codec, so it is transport-agnostic and works against any LNS. This v1 implements
// LoRaWAN 1.0.x (legacy) session keys + MIC, the dominant case; the JoinEUI is kept
// for the 1.1 derivation path.
package device

import (
	"errors"

	"simulator/packages/utils/lorawan/codec"
	"simulator/packages/utils/lorawan/crypto"
	"simulator/packages/utils/lorawan/types"
)

// Config is the device identity + keys needed to activate and run.
type Config struct {
	JoinEUI types.EUI64
	DevEUI  types.EUI64
	AppKey  types.AES128Key // root key (OTAA)
	NetID   types.NetID     // used for 1.0.x session-key derivation

	// ABP pre-provisioned session (used when Activation == "abp").
	Activation string // "otaa" | "abp"
	DevAddr    types.DevAddr
	NwkSKey    types.AES128Key
	AppSKey    types.AES128Key
}

// Session is a live device session. Not safe for concurrent use; the engine drives
// one device's sends/receives serially.
type Session struct {
	cfg Config

	joined   bool
	devNonce uint16
	devAddr  types.DevAddr
	nwkSKey  types.AES128Key
	appSKey  types.AES128Key

	fCntUp   uint32
	fCntDown uint32
}

// New builds a session. An ABP device is immediately activated; an OTAA device must
// JoinRequest() then ProcessJoinAccept() before sending.
func New(cfg Config) *Session {
	s := &Session{cfg: cfg}
	if cfg.Activation == "abp" {
		s.joined = true
		s.devAddr = cfg.DevAddr
		s.nwkSKey = cfg.NwkSKey
		s.appSKey = cfg.AppSKey
	}
	return s
}

// Joined reports whether the device has an active session.
func (s *Session) Joined() bool { return s.joined }

// DevAddr returns the current session device address (zero before join).
func (s *Session) DevAddr() types.DevAddr { return s.devAddr }

// JoinRequest builds the next OTAA join-request PHYPayload, advancing the DevNonce
// (monotonic per LoRaWAN 1.0.4 / 1.1).
func (s *Session) JoinRequest() ([]byte, error) {
	if s.cfg.Activation == "abp" {
		return nil, errors.New("device: ABP device does not join")
	}
	s.devNonce++
	jr := codec.JoinRequest{
		JoinEUI:  s.cfg.JoinEUI,
		DevEUI:   s.cfg.DevEUI,
		DevNonce: types.DevNonce{byte(s.devNonce), byte(s.devNonce >> 8)},
	}
	return codec.MarshalJoinRequest(jr, s.cfg.AppKey)
}

// ProcessJoinAccept decrypts the join accept, derives the 1.0.x session keys, and
// activates the session.
func (s *Session) ProcessJoinAccept(phy []byte) error {
	ja, err := codec.UnmarshalJoinAccept(phy, s.cfg.AppKey)
	if err != nil {
		return err
	}
	dn := types.DevNonce{byte(s.devNonce), byte(s.devNonce >> 8)}
	s.nwkSKey = crypto.DeriveLegacyNwkSKey(s.cfg.AppKey, ja.JoinNonce, ja.NetID, dn)
	s.appSKey = crypto.DeriveLegacyAppSKey(s.cfg.AppKey, ja.JoinNonce, ja.NetID, dn)
	s.devAddr = ja.DevAddr
	s.fCntUp = 0
	s.fCntDown = 0
	s.joined = true
	return nil
}

// BuildUplink assembles the next data uplink, advancing the uplink frame counter.
func (s *Session) BuildUplink(fPort byte, payload []byte, confirmed bool) ([]byte, error) {
	if !s.joined {
		return nil, errors.New("device: not joined")
	}
	up := codec.DataUplink{
		DevAddr:    s.devAddr,
		FCnt:       s.fCntUp,
		FPort:      fPort,
		FRMPayload: payload,
		Confirmed:  confirmed,
	}
	phy, err := codec.MarshalDataUplink(up, s.nwkSKey, s.appSKey)
	if err != nil {
		return nil, err
	}
	s.fCntUp++
	return phy, nil
}

// DecodeDownlink parses and decrypts a data downlink, advancing the downlink frame
// counter to the received value.
func (s *Session) DecodeDownlink(phy []byte) (codec.DataDownlink, error) {
	if !s.joined {
		return codec.DataDownlink{}, errors.New("device: not joined")
	}
	dl, err := codec.UnmarshalDataDownlink(phy, s.fCntDown, s.appSKey)
	if err != nil {
		return codec.DataDownlink{}, err
	}
	s.fCntDown = dl.FCnt + 1
	return dl, nil
}

// FCntUp / FCntDown expose the counters for console/metadata.
func (s *Session) FCntUp() uint32   { return s.fCntUp }
func (s *Session) FCntDown() uint32 { return s.fCntDown }
