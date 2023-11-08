//go:build js && wasm

package main

import "testing"

func TestUserAgent(t *testing.T) {
	t.Logf("user-agent: %s\n", UserAgent())
}
