package main

import (
	"github.com/hoeg/beaver/internal/app"
	"github.com/hoeg/beaver/internal/events"
)

func main() {
	//user the Merge API as the server interface
	app.Start()
	events.ListenForEvents()
}
