package main

import (
	"github.com/chrisenytc/skynet/config"
	"github.com/chrisenytc/skynet/proxy"
)

func main() {
	// Load configs
	config.Load()

	// Load proxy server
	proxy.Load()
}
