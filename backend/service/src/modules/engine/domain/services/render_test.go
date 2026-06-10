package services

import (
	"strconv"
	"strings"
	"testing"
)

func TestRender_DeterministicTokens(t *testing.T) {
	out := Render(`{"id":"{{deviceId}}","n":{{counter}}}`, "dev-1", 7)
	if out != `{"id":"dev-1","n":7}` {
		t.Fatalf("render = %q", out)
	}
}

func TestRender_RandIntInRange(t *testing.T) {
	for i := 0; i < 100; i++ {
		out := Render(`{{randInt(10,12)}}`, "d", 0)
		n, err := strconv.Atoi(out)
		if err != nil || n < 10 || n > 12 {
			t.Fatalf("randInt out = %q (n=%d err=%v)", out, n, err)
		}
	}
}

func TestRender_RandFloatInRange(t *testing.T) {
	for i := 0; i < 100; i++ {
		out := Render(`{{randFloat(1,2)}}`, "d", 0)
		f, err := strconv.ParseFloat(out, 64)
		if err != nil || f < 1 || f > 2 {
			t.Fatalf("randFloat out = %q", out)
		}
	}
}

func TestRender_UnknownTokenKept(t *testing.T) {
	if out := Render(`{{bogus}}`, "d", 0); out != "{{bogus}}" {
		t.Fatalf("unknown token = %q, want kept", out)
	}
}

func TestRender_UUIDAndNow(t *testing.T) {
	if out := Render(`{{uuid}}`, "d", 0); len(out) != 36 || strings.Count(out, "-") != 4 {
		t.Fatalf("uuid = %q", out)
	}
	if out := Render(`{{now}}`, "d", 0); !strings.Contains(out, "T") || !strings.HasSuffix(out, "Z") {
		t.Fatalf("now = %q", out)
	}
}
