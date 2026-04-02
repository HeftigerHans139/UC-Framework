package main

import (
	"uc_framework/internal/core"
)

func main() {
	core.Start()
	// Keep process alive: core starts long-running services in goroutines.
	select {}
}
