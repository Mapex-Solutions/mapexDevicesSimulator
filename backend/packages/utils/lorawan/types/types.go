// Package types holds the fixed-size LoRaWAN identifier and key types used by the
// vendored crypto package. Names and byte widths mirror the LoRaWAN specification
// (and The Things Stack pkg/types) so the vendored crypto compiles unchanged; only
// the types the crypto package needs are defined here, as plain byte arrays.
package types

// AES128Key is a 128-bit AES key (root or session).
type AES128Key [16]byte

// DevAddr is the 32-bit device address assigned during activation.
type DevAddr [4]byte

// EUI64 is a 64-bit globally unique identifier (DevEUI / JoinEUI).
type EUI64 [8]byte

// DevNonce is the 16-bit device nonce in a join request.
type DevNonce [2]byte

// JoinNonce is the 24-bit join-server nonce in a join accept.
type JoinNonce [3]byte

// NetID is the 24-bit network identifier.
type NetID [3]byte
