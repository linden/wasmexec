//go:build js && wasm

package main

import (
	"fmt"
	"syscall/js"
)

func UserAgent() string {
	return js.Global().Get("navigator").Get("userAgent").String()
}

func main() {
	fmt.Printf("user-agent: %v\n", UserAgent())
}
