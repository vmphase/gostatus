package gateway

import "sync/atomic"

var connected atomic.Bool

// Connected reports whether the gateway is currently connected.
func Connected() bool {
	return connected.Load()
}
