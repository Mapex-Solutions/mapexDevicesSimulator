// Package band holds the minimal per-region LoRaWAN parameters the simulator needs:
// the RX2 downlink window (frequency + data rate), the standard join/rx delays, and
// a representative uplink frequency + spreading factor used to fill simulated radio
// metadata (freq/SF) on console frames. These are the well-known regional defaults;
// the tables are intentionally small (not a full band plan).
package band

import "time"

// Region is the subset of band parameters used by the device session.
type Region struct {
	Name             string
	UplinkFrequency  uint64        // Hz, a representative default channel
	UplinkSF         int           // representative spreading factor (DR-derived)
	UplinkDR         int           // data-rate index for UplinkSF (LoRaWAN 1.1 uplink MIC)
	UplinkChannel    int           // channel index of UplinkFrequency (LoRaWAN 1.1 uplink MIC)
	RX2Frequency     uint64        // Hz, fixed RX2 downlink frequency
	RX2SF            int           // RX2 spreading factor
	RX1Delay         time.Duration // delay before the RX1 window after an uplink
	JoinAcceptDelay1 time.Duration // delay before RX1 for a join accept
}

// RX2Delay is RX1Delay + 1s by spec.
func (r Region) RX2Delay() time.Duration { return r.RX1Delay + time.Second }

// JoinAcceptDelay2 is JoinAcceptDelay1 + 1s by spec.
func (r Region) JoinAcceptDelay2() time.Duration { return r.JoinAcceptDelay1 + time.Second }

// regions maps a region id (matching the GatewayRegion enum on the wire) to its
// defaults.
var regions = map[string]Region{
	"EU868": {Name: "EU868", UplinkFrequency: 868_100_000, UplinkSF: 7, UplinkDR: 5, UplinkChannel: 0, RX2Frequency: 869_525_000, RX2SF: 12, RX1Delay: time.Second, JoinAcceptDelay1: 5 * time.Second},
	"US915": {Name: "US915", UplinkFrequency: 902_300_000, UplinkSF: 7, UplinkDR: 3, UplinkChannel: 0, RX2Frequency: 923_300_000, RX2SF: 12, RX1Delay: time.Second, JoinAcceptDelay1: 5 * time.Second},
	"AU915": {Name: "AU915", UplinkFrequency: 915_200_000, UplinkSF: 7, UplinkDR: 5, UplinkChannel: 0, RX2Frequency: 923_300_000, RX2SF: 12, RX1Delay: time.Second, JoinAcceptDelay1: 5 * time.Second},
	"AS923": {Name: "AS923", UplinkFrequency: 923_200_000, UplinkSF: 7, UplinkDR: 5, UplinkChannel: 0, RX2Frequency: 923_200_000, RX2SF: 10, RX1Delay: time.Second, JoinAcceptDelay1: 5 * time.Second},
	"CN470": {Name: "CN470", UplinkFrequency: 470_300_000, UplinkSF: 7, UplinkDR: 5, UplinkChannel: 0, RX2Frequency: 505_300_000, RX2SF: 12, RX1Delay: time.Second, JoinAcceptDelay1: 5 * time.Second},
	"IN865": {Name: "IN865", UplinkFrequency: 865_062_500, UplinkSF: 7, UplinkDR: 5, UplinkChannel: 0, RX2Frequency: 866_550_000, RX2SF: 10, RX1Delay: time.Second, JoinAcceptDelay1: 5 * time.Second},
	"KR920": {Name: "KR920", UplinkFrequency: 922_100_000, UplinkSF: 7, UplinkDR: 5, UplinkChannel: 0, RX2Frequency: 921_900_000, RX2SF: 12, RX1Delay: time.Second, JoinAcceptDelay1: 5 * time.Second},
	"RU864": {Name: "RU864", UplinkFrequency: 868_900_000, UplinkSF: 7, UplinkDR: 5, UplinkChannel: 0, RX2Frequency: 869_100_000, RX2SF: 12, RX1Delay: time.Second, JoinAcceptDelay1: 5 * time.Second},
}

// Get returns the region parameters for the given id, defaulting to EU868 when the
// id is unknown so a misconfigured device still has sane radio metadata.
func Get(region string) Region {
	if r, ok := regions[region]; ok {
		return r
	}
	return regions["EU868"]
}
