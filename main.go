package main

import (
	"uc_framework/internal/core"
)

func main() {
	// Initialize core config and auth first
	core.Start()
	// Keep main running forever (web server runs in a goroutine from core.Start())
	select {}
}

