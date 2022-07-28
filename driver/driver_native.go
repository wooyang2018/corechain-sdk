//go:build !wasm
// +build !wasm

package driver

import (
	"github.com/wooyang2018/corechain-sdk/code"
	"github.com/wooyang2018/corechain-sdk/driver/native"
)

// Serve run contract in native environment
func Serve(contract code.Contract) {
	driver := native.New()
	driver.Serve(contract)
}
