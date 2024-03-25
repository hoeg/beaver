package app

import (
	"log"

	"github.com/g4s8/go-lifecycle/pkg/lifecycle"
	event "github.com/hoeg/beaver/internal/events"
)

func Start() {
	api := &API{}
	svc := NewAPI(api)

	lf := lifecycle.New(lifecycle.DefaultConfig)
	svc.RegisterLifecycle("web", lf)
	ec, err := event.NewEventController("hoeg.com/artifact-id")
	if err != nil {
		log.Fatalf("failed to create event controller: %v", err)
	}
	lf.RegisterService(NewEventController(ec))

	lf.Start()
	sig := lifecycle.NewSignalHandler(lf, nil)
	sig.Start(lifecycle.DefaultShutdownConfig)
	if err := sig.Wait(); err != nil {
		log.Fatalf("shutdown error: %v", err)
	}
	log.Print("shutdown complete")
}
