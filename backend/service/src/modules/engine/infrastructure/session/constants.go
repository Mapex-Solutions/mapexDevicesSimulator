package session

import "time"

// mqttQuiesceMillis is the graceful-flush window when closing an MQTT session.
const mqttQuiesceMillis = 250

// joinTimeout bounds how long Open waits for the OTAA join accept before failing,
// so the session manager retries with backoff rather than blocking forever.
const joinTimeout = 6 * time.Second

// statInterval is how often the gateway reports statistics so the LNS keeps it
// marked online (a real packet forwarder sends a stat roughly every 30s).
const statInterval = 30 * time.Second
