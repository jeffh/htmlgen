//go:build !js

package html

func Chan[X any](name string) chan X {
	return make(chan X, 1)
}

const Server = true
const Client = false
