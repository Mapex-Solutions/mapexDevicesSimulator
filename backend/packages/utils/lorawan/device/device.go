// Package device is the stateful LoRaWAN "brain" of a simulated end device: it
// holds the keys and frame counters, performs the OTAA join (or ABP activation),
// builds data uplinks, and decodes downlinks — all over the vendored crypto + our
// codec, so it is transport-agnostic and works against any LNS. It implements both
// LoRaWAN 1.0.x (legacy single NwkSKey + legacy MIC) and 1.1 OTAA (separate
// FNwkSIntKey/SNwkSIntKey/NwkSEncKey/AppSKey + the 1.1 uplink MIC). ABP uses the
// legacy path.
package device

import (
	"errors"
	"strings"

	"simulator/packages/utils/lorawan/band"
	"simulator/packages/utils/lorawan/codec"
	"simulator/packages/utils/lorawan/crypto"
	"simulator/packages/utils/lorawan/types"
)

// Config is the device identity + keys needed to activate and run.
type Config struct {
	MACVersion string // "1.0.2".."1.1.0"; "1.1*" selects the 1.1 key schedule
	Region     string // for the 1.1 uplink MIC (data-rate/channel) and radio meta

	JoinEUI types.EUI64
	DevEUI  types.EUI64
	AppKey  types.AES128Key // application root key (AppSKey derivation; 1.0.x join MIC)
	NwkKey  types.AES128Key // network root key (1.1 join MIC + network session keys)
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
	cfg    Config
	is11   bool        // device is provisioned for LoRaWAN 1.1
	use11  bool        // 1.1 session keys are active (OTAA 1.1 joined)
	region band.Region

	joined   bool
	devNonce uint16
	devAddr  types.DevAddr

	nwkSKey     types.AES128Key // 1.0.x network session key
	fNwkSIntKey types.AES128Key // 1.1 network session keys
	sNwkSIntKey types.AES128Key
	nwkSEncKey  types.AES128Key
	appSKey     types.AES128Key

	fCntUp   uint32
	fCntDown uint32
}

// New builds a session. An ABP device is immediately activated (legacy path); an
// OTAA device must JoinRequest() then ProcessJoinAccept() before sending.
func New(cfg Config) *Session {
	s := &Session{
		cfg:    cfg,
		region: band.Get(cfg.Region),
		is11:   strings.HasPrefix(cfg.MACVersion, "1.1"),
	}
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
// (monotonic per LoRaWAN 1.0.4 / 1.1). The MIC is signed with the NwkKey for 1.1 and
// the AppKey for 1.0.x.
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
	joinKey := s.cfg.AppKey
	if s.is11 {
		joinKey = s.cfg.NwkKey
	}
	return codec.MarshalJoinRequest(jr, joinKey)
}

// ProcessJoinAccept decrypts the join accept, derives the session keys for the
// device's MAC version, and activates the session.
func (s *Session) ProcessJoinAccept(phy []byte) error {
	dn := types.DevNonce{byte(s.devNonce), byte(s.devNonce >> 8)}
	if s.is11 {
		ja, err := codec.UnmarshalJoinAccept(phy, s.cfg.NwkKey)
		if err != nil {
			return err
		}
		s.fNwkSIntKey = crypto.DeriveFNwkSIntKey(s.cfg.NwkKey, ja.JoinNonce, s.cfg.JoinEUI, dn)
		s.sNwkSIntKey = crypto.DeriveSNwkSIntKey(s.cfg.NwkKey, ja.JoinNonce, s.cfg.JoinEUI, dn)
		s.nwkSEncKey = crypto.DeriveNwkSEncKey(s.cfg.NwkKey, ja.JoinNonce, s.cfg.JoinEUI, dn)
		s.appSKey = crypto.DeriveAppSKey(s.cfg.AppKey, ja.JoinNonce, s.cfg.JoinEUI, dn)
		s.devAddr = ja.DevAddr
		s.use11 = true
	} else {
		ja, err := codec.UnmarshalJoinAccept(phy, s.cfg.AppKey)
		if err != nil {
			return err
		}
		s.nwkSKey = crypto.DeriveLegacyNwkSKey(s.cfg.AppKey, ja.JoinNonce, ja.NetID, dn)
		s.appSKey = crypto.DeriveLegacyAppSKey(s.cfg.AppKey, ja.JoinNonce, ja.NetID, dn)
		s.devAddr = ja.DevAddr
	}
	s.fCntUp = 0
	s.fCntDown = 0
	s.joined = true
	return nil
}

// BuildUplink assembles the next data uplink, advancing the uplink frame counter.
// 1.1 sessions sign with the two network session keys plus the region's data-rate
// and channel; 1.0.x and ABP sign with the legacy NwkSKey.
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
	var phy []byte
	var err error
	if s.use11 {
		phy, err = codec.MarshalDataUplink11(up, s.sNwkSIntKey, s.fNwkSIntKey, s.appSKey,
			byte(s.region.UplinkDR), byte(s.region.UplinkChannel))
	} else {
		phy, err = codec.MarshalDataUplink(up, s.nwkSKey, s.appSKey)
	}
	if err != nil {
		return nil, err
	}
	s.fCntUp++
	return phy, nil
}

// DecodeDownlink parses and decrypts a data downlink, advancing the downlink frame
// counter to the received value. FRMPayload is decrypted with appSKey in both
// versions.
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
