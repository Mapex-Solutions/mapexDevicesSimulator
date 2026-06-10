package session

import "testing"

func TestDownlinkRouterRoutesByMHDRAndDevAddr(t *testing.T) {
	r := newDownlinkRouter()

	var gotJoin, gotData []byte
	r.setPending(func(phy []byte) { gotJoin = phy })

	// A join accept (MHDR 0x20) goes to the pending handler.
	ja := []byte{0x20, 1, 2, 3, 4, 5}
	r.route(ja)
	if gotJoin == nil {
		t.Fatal("join accept was not routed to the pending handler")
	}

	// After binding by DevAddr, a data downlink for that address is delivered.
	addr := [4]byte{0x26, 0x01, 0x1F, 0x88}
	r.bind(addr, func(phy []byte) { gotData = phy })
	// FHDR DevAddr is little-endian on the wire: reversed address bytes after MHDR.
	dl := []byte{0x60, addr[3], addr[2], addr[1], addr[0], 0x00, 0x00, 0x00}
	r.route(dl)
	if gotData == nil {
		t.Fatal("data downlink was not routed by DevAddr")
	}

	// A downlink for an unknown address is dropped (no panic, no delivery).
	gotData = nil
	other := []byte{0x60, 0xFF, 0xFF, 0xFF, 0xFF, 0x00, 0x00, 0x00}
	r.route(other)
	if gotData != nil {
		t.Fatal("downlink for an unbound address must be dropped")
	}
}
