// Package client is a thin wrapper over the ChirpStack v4 gRPC API used by the
// LoRaWAN journey to provision the LNS: a tenant, an application, an OTAA
// device profile, the device and its keys, and the gateway the simulator
// transmits through. Everything the journey needs from ChirpStack goes through
// here so the steps stay declarative.
//
// The API is consumed against a pinned server image (see
// deployment/chirpstack/chirpstack.yml) so these calls never drift.
package client

import (
	"context"
	"fmt"
	"time"

	csapi "github.com/chirpstack/chirpstack/api/go/v4/api"
	"github.com/chirpstack/chirpstack/api/go/v4/common"
	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials/insecure"
	"google.golang.org/grpc/metadata"
)

// Client holds the gRPC connection, the login JWT, and one typed service client
// per ChirpStack service the journey touches.
type Client struct {
	conn *grpc.ClientConn
	jwt  string

	tenant  csapi.TenantServiceClient
	app     csapi.ApplicationServiceClient
	profile csapi.DeviceProfileServiceClient
	device  csapi.DeviceServiceClient
	gateway csapi.GatewayServiceClient
}

// Connect dials the ChirpStack gRPC API and logs in, retrying the login until
// it succeeds or ctx is done. The retry loop doubles as the stack-readiness
// gate: a fresh stack only answers login once migrations and the API are up.
func Connect(ctx context.Context, addr, user, pass string) (*Client, error) {
	conn, err := grpc.NewClient(addr, grpc.WithTransportCredentials(insecure.NewCredentials()))
	if err != nil {
		return nil, fmt.Errorf("chirpstack dial: %w", err)
	}
	c := &Client{
		conn:    conn,
		tenant:  csapi.NewTenantServiceClient(conn),
		app:     csapi.NewApplicationServiceClient(conn),
		profile: csapi.NewDeviceProfileServiceClient(conn),
		device:  csapi.NewDeviceServiceClient(conn),
		gateway: csapi.NewGatewayServiceClient(conn),
	}
	internal := csapi.NewInternalServiceClient(conn)

	var lastErr error
	for {
		resp, err := internal.Login(ctx, &csapi.LoginRequest{Email: user, Password: pass})
		if err == nil {
			c.jwt = resp.GetJwt()
			return c, nil
		}
		lastErr = err
		select {
		case <-ctx.Done():
			conn.Close()
			return nil, fmt.Errorf("chirpstack login not ready before timeout: %w", lastErr)
		case <-time.After(time.Second):
		}
	}
}

// Close tears down the gRPC connection.
func (c *Client) Close() error { return c.conn.Close() }

// auth returns a context carrying the login JWT as the bearer the API expects.
func (c *Client) auth(ctx context.Context) context.Context {
	return metadata.AppendToOutgoingContext(ctx, "authorization", "Bearer "+c.jwt)
}

// EnsureTenant returns a tenant id to provision under, reusing the first tenant
// the server already has (a fresh ChirpStack seeds a default tenant) and
// creating one only if none exist.
func (c *Client) EnsureTenant(ctx context.Context, name string) (string, error) {
	list, err := c.tenant.List(c.auth(ctx), &csapi.ListTenantsRequest{Limit: 1})
	if err != nil {
		return "", fmt.Errorf("list tenants: %w", err)
	}
	if len(list.GetResult()) > 0 {
		return list.GetResult()[0].GetId(), nil
	}
	resp, err := c.tenant.Create(c.auth(ctx), &csapi.CreateTenantRequest{
		Tenant: &csapi.Tenant{Name: name, CanHaveGateways: true, MaxGatewayCount: 10, MaxDeviceCount: 100},
	})
	if err != nil {
		return "", fmt.Errorf("create tenant: %w", err)
	}
	return resp.GetId(), nil
}

// CreateApplication creates an application under the tenant and returns its id.
func (c *Client) CreateApplication(ctx context.Context, tenantID, name string) (string, error) {
	resp, err := c.app.Create(c.auth(ctx), &csapi.CreateApplicationRequest{
		Application: &csapi.Application{TenantId: tenantID, Name: name},
	})
	if err != nil {
		return "", fmt.Errorf("create application: %w", err)
	}
	return resp.GetId(), nil
}

// CreateDeviceProfile creates an EU868 / LoRaWAN 1.0.3 OTAA profile and returns
// its id. These match the simulator's LoRaWAN device defaults so the join keys
// line up on both sides.
func (c *Client) CreateDeviceProfile(ctx context.Context, tenantID, name string) (string, error) {
	resp, err := c.profile.Create(c.auth(ctx), &csapi.CreateDeviceProfileRequest{
		DeviceProfile: &csapi.DeviceProfile{
			TenantId:          tenantID,
			Name:              name,
			Region:            common.Region_EU868,
			RegionConfigId:    "eu868",
			MacVersion:        common.MacVersion_LORAWAN_1_0_3,
			RegParamsRevision: common.RegParamsRevision_A,
			AdrAlgorithmId:    "default",
			UplinkInterval:    3600,
			SupportsOtaa:      true,
		},
	})
	if err != nil {
		return "", fmt.Errorf("create device profile: %w", err)
	}
	return resp.GetId(), nil
}

// CreateGateway registers the gateway the simulator transmits through. ChirpStack
// drops frames from gateways it does not know, so this must precede any uplink.
func (c *Client) CreateGateway(ctx context.Context, tenantID, gatewayEUI, name string) error {
	_, err := c.gateway.Create(c.auth(ctx), &csapi.CreateGatewayRequest{
		Gateway: &csapi.Gateway{GatewayId: gatewayEUI, Name: name, TenantId: tenantID, StatsInterval: 30},
	})
	if err != nil {
		return fmt.Errorf("create gateway %s: %w", gatewayEUI, err)
	}
	return nil
}

// DeleteGateway removes the gateway. It is safe to call during rollback.
func (c *Client) DeleteGateway(ctx context.Context, gatewayEUI string) error {
	_, err := c.gateway.Delete(c.auth(ctx), &csapi.DeleteGatewayRequest{GatewayId: gatewayEUI})
	return err
}

// CreateDevice creates the OTAA device under the application and profile.
func (c *Client) CreateDevice(ctx context.Context, appID, profileID, devEUI, joinEUI, name string) error {
	_, err := c.device.Create(c.auth(ctx), &csapi.CreateDeviceRequest{
		Device: &csapi.Device{
			DevEui:          devEUI,
			JoinEui:         joinEUI,
			Name:            name,
			ApplicationId:   appID,
			DeviceProfileId: profileID,
		},
	})
	if err != nil {
		return fmt.Errorf("create device %s: %w", devEUI, err)
	}
	return nil
}

// SetDeviceKeys sets the OTAA root key. For LoRaWAN 1.0.x the AppKey is carried
// in the nwk_key field (the API documents this explicitly).
func (c *Client) SetDeviceKeys(ctx context.Context, devEUI, appKey string) error {
	_, err := c.device.CreateKeys(c.auth(ctx), &csapi.CreateDeviceKeysRequest{
		DeviceKeys: &csapi.DeviceKeys{DevEui: devEUI, NwkKey: appKey},
	})
	if err != nil {
		return fmt.Errorf("set device keys %s: %w", devEUI, err)
	}
	return nil
}

// DeleteDevice removes the device. Safe to call during rollback.
func (c *Client) DeleteDevice(ctx context.Context, devEUI string) error {
	_, err := c.device.Delete(c.auth(ctx), &csapi.DeleteDeviceRequest{DevEui: devEUI})
	return err
}

// Activation returns the device's assigned DevAddr once the OTAA join has
// completed, or an empty string while it has not. A non-empty DevAddr is proof
// the join-request reached ChirpStack and the join-accept was issued.
func (c *Client) Activation(ctx context.Context, devEUI string) (string, error) {
	resp, err := c.device.GetActivation(c.auth(ctx), &csapi.GetDeviceActivationRequest{DevEui: devEUI})
	if err != nil {
		return "", fmt.Errorf("get activation %s: %w", devEUI, err)
	}
	if resp.GetDeviceActivation() == nil {
		return "", nil
	}
	return resp.GetDeviceActivation().GetDevAddr(), nil
}

// EnqueueDownlink queues a downlink for the device on the given FPort. For a
// Class A device the LNS sends it in the RX window after the next uplink, where
// the simulator receives it on its inbound path.
func (c *Client) EnqueueDownlink(ctx context.Context, devEUI string, fPort uint32, data []byte) error {
	_, err := c.device.Enqueue(c.auth(ctx), &csapi.EnqueueDeviceQueueItemRequest{
		QueueItem: &csapi.DeviceQueueItem{DevEui: devEUI, FPort: fPort, Data: data},
	})
	if err != nil {
		return fmt.Errorf("enqueue downlink %s: %w", devEUI, err)
	}
	return nil
}

// LastSeen reports whether ChirpStack has recorded an uplink from the device
// (last_seen_at is set on the first received frame).
func (c *Client) LastSeen(ctx context.Context, devEUI string) (bool, error) {
	resp, err := c.device.Get(c.auth(ctx), &csapi.GetDeviceRequest{DevEui: devEUI})
	if err != nil {
		return false, fmt.Errorf("get device %s: %w", devEUI, err)
	}
	return resp.GetLastSeenAt() != nil, nil
}
