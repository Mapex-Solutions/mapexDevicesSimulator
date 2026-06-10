package hub

// clientBuffer bounds a client's pending frames; a client that falls this far
// behind drops frames instead of blocking the broadcaster.
const clientBuffer = 256
