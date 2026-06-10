package udp

import (
	"bytes"
	"encoding/json"
	"testing"
)

func TestMarshalPushDataHeaderAndPayload(t *testing.T) {
	eui := [8]byte{0x00, 0x16, 0xC0, 0x01, 0xF1, 0x50, 0x00, 0x01}
	pk := RxPk{Freq: 868.1, Modu: "LORA", DatR: "SF7BW125", Data: EncodePHY([]byte{0x40, 0x01, 0x02})}

	out, err := MarshalPushData(0xBEEF, eui, pk)
	if err != nil {
		t.Fatalf("marshal: %v", err)
	}
	if out[0] != protocolVersion || out[3] != idPushData {
		t.Fatalf("header wrong: ver=%#x id=%#x", out[0], out[3])
	}
	if out[1] != 0xBE || out[2] != 0xEF {
		t.Fatalf("token not big-endian: %#x %#x", out[1], out[2])
	}
	if !bytes.Equal(out[4:12], eui[:]) {
		t.Fatalf("gateway EUI mismatch")
	}
	var p pushPayload
	if err := json.Unmarshal(out[12:], &p); err != nil {
		t.Fatalf("json body: %v", err)
	}
	if len(p.RxPk) != 1 || p.RxPk[0].Freq != 868.1 {
		t.Fatalf("rxpk not encoded: %+v", p.RxPk)
	}
}

func TestMarshalStatIsPushDataWithStat(t *testing.T) {
	eui := [8]byte{0x00, 0x16, 0xC0, 0x01, 0xF1, 0x50, 0x00, 0x01}
	out, err := MarshalStat(0x0001, eui, Stat{Time: "2026-06-10 19:00:00 GMT", RxNb: 5, AckR: 100})
	if err != nil {
		t.Fatalf("marshal: %v", err)
	}
	if out[0] != protocolVersion || out[3] != idPushData {
		t.Fatalf("not a PUSH_DATA: ver=%#x id=%#x", out[0], out[3])
	}
	var p statPayload
	if err := json.Unmarshal(out[12:], &p); err != nil {
		t.Fatalf("json body: %v", err)
	}
	if p.Stat.RxNb != 5 || p.Stat.AckR != 100 {
		t.Fatalf("stat not encoded: %+v", p.Stat)
	}
}

func TestClassify(t *testing.T) {
	eui := [8]byte{1, 2, 3, 4, 5, 6, 7, 8}
	phy := []byte{0x60, 0xAA, 0xBB}
	resp := append(header(7, idPullResp), []byte(`{"txpk":{"data":"`+EncodePHY(phy)+`"}}`)...)

	tests := []struct {
		name     string
		in       []byte
		wantKind Kind
		wantPHY  []byte
	}{
		{"push ack", header(1, idPushAck), KindPushAck, nil},
		{"pull ack", header(2, idPullAck), KindPullAck, nil},
		{"pull resp carries phy", resp, KindPullResp, phy},
		{"keepalive echoed back is unknown", MarshalPullData(3, eui), KindUnknown, nil},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			kind, gotPHY, err := Classify(tt.in)
			if err != nil {
				t.Fatalf("classify: %v", err)
			}
			if kind != tt.wantKind {
				t.Fatalf("kind = %d, want %d", kind, tt.wantKind)
			}
			if !bytes.Equal(gotPHY, tt.wantPHY) {
				t.Fatalf("phy = %x, want %x", gotPHY, tt.wantPHY)
			}
		})
	}
}
