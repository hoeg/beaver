package app

import (
	"log"

	"github.com/g4s8/go-lifecycle/pkg/lifecycle"
)

func Start() {
	api := &API{}
	svc := NewAPI(api)

	lf := lifecycle.New(lifecycle.DefaultConfig)
	svc.RegisterLifecycle("web", lf)

	lf.Start()
	sig := lifecycle.NewSignalHandler(lf, nil)
	sig.Start(lifecycle.DefaultShutdownConfig)
	if err := sig.Wait(); err != nil {
		log.Fatalf("shutdown error: %v", err)
	}
	log.Print("shutdown complete")
}
