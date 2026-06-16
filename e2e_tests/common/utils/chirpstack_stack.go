package utils

import (
	"context"
	"fmt"
	"os/exec"
	"path/filepath"
	"runtime"
)

// ChirpStackStack drives the lifecycle of the pinned ChirpStack docker stack the
// LoRaWAN journey provisions against. The journey brings it Up at the start and
// Down at the end, so each run gets a fresh LNS (which also wipes ChirpStack's
// remembered OTAA DevNonces, letting the journey reuse deterministic keys).
type ChirpStackStack struct {
	composeFile string
}

// NewChirpStackStack locates deployment/chirpstack/chirpstack.yml relative to
// this source file, so it resolves regardless of the test's working directory.
func NewChirpStackStack() *ChirpStackStack {
	_, thisFile, _, _ := runtime.Caller(0)
	// this file: <root>/e2e_tests/common/utils/chirpstack_stack.go
	root := filepath.Join(filepath.Dir(thisFile), "..", "..", "..")
	return &ChirpStackStack{
		composeFile: filepath.Join(root, "deployment", "chirpstack", "chirpstack.yml"),
	}
}

// Up starts the stack in the background and waits for the containers to start.
// Readiness of the API itself is owned by the caller (a gRPC login poll), since
// the chirpstack image is minimal and carries no usable container healthcheck.
func (s *ChirpStackStack) Up(ctx context.Context) error {
	return s.compose(ctx, "up", "-d", "--wait")
}

// Down tears the stack down and removes its volumes, leaving no state behind.
func (s *ChirpStackStack) Down(ctx context.Context) error {
	return s.compose(ctx, "down", "-v")
}

// ComposeFile is the absolute path to the stack's compose file.
func (s *ChirpStackStack) ComposeFile() string { return s.composeFile }

// compose runs `docker compose -f <file> <args...>`, surfacing the combined
// output on failure so a broken stack is diagnosable from the test log.
func (s *ChirpStackStack) compose(ctx context.Context, args ...string) error {
	full := append([]string{"compose", "-f", s.composeFile}, args...)
	cmd := exec.CommandContext(ctx, "docker", full...)
	out, err := cmd.CombinedOutput()
	if err != nil {
		return fmt.Errorf("docker %v: %w\n%s", args, err, out)
	}
	return nil
}
