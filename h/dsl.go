//go:build !js

package h

// Chan creates a buffered channel with capacity 1.
// The name parameter is for identification purposes.
func Chan[X any](name string) chan X {
	return make(chan X, 1)
}

// Server indicates whether the code is running on the server side.
// This is true when not compiled for WebAssembly (js build tag).
const Server = true

// Client indicates whether the code is running on the client side.
// This is false when not compiled for WebAssembly (js build tag).
const Client = false
