// Package utils holds the shared e2e helpers: the readiness check and the test
// fixtures (an HTTP echo) the journeys stand up.
package utils

import (
	"context"
	"fmt"
	"net/http"
	"sync"
	"time"

	"github.com/Mapex-Solutions/mapexGoKit/infrastructure/httpclient"

	"github.com/Mapex-Solutions/mapexDevicesSimulator/e2eTests/common/constants"
)

var (
	setupOnce sync.Once
	setupErr  error
)

// SetupE2EEnvironment is a thin liveness check: the simulator sidecar must be
// up. It is idempotent — every journey that calls it observes the same outcome
// without mutating state. Unlike a full platform, the simulator needs no seed
// data, so readiness is simply a green /api/health.
func SetupE2EEnvironment() error {
	setupOnce.Do(func() {
		setupErr = checkSidecarReady()
	})
	return setupErr
}

// checkSidecarReady probes the simulator's /api/health. A 200 confirms the
// sidecar is serving; it stays read-only and side-effect-free.
func checkSidecarReady() error {
	client := httpclient.New(httpclient.Config{
		BaseURL: constants.SimURL,
		Timeout: 5 * time.Second,
	})
	ctx, cancel := context.WithTimeout(context.Background(), 8*time.Second)
	defer cancel()

	resp, err := client.Raw(ctx, http.MethodGet, "/api/health", nil)
	if err != nil {
		return fmt.Errorf("e2e readiness: simulator health check failed (is it running on %s?): %w", constants.SimURL, err)
	}
	defer resp.Body.Close()
	if resp.StatusCode != http.StatusOK {
		return fmt.Errorf("e2e readiness: unexpected status %d from /api/health", resp.StatusCode)
	}
	return nil
}
