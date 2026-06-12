package session

import "testing"

func TestFormatDevAddr(t *testing.T) {
	if got := formatDevAddr([4]byte{0x26, 0x0b, 0xac, 0x12}); got != "26:0B:AC:12" {
		t.Fatalf("formatDevAddr = %q, want 26:0B:AC:12", got)
	}
}

func TestLoRaWANSession_JoinedDetail(t *testing.T) {
	s := &lorawanSession{addr: [4]byte{0x26, 0x0b, 0xac, 0x12}, class: "C"}
	if got := s.joinedDetail(); got != "DevAddr 26:0B:AC:12 · Class C" {
		t.Fatalf("joinedDetail = %q", got)
	}
}

func TestClassOrA(t *testing.T) {
	if got := classOrA(""); got != "A" {
		t.Fatalf("empty class should default to A, got %q", got)
	}
	if got := classOrA("C"); got != "C" {
		t.Fatalf("class C should pass through, got %q", got)
	}
}
