package app

import (
	"time"

	"github.com/g4s8/go-lifecycle/pkg/types"
	event "github.com/hoeg/beaver/internal/events"
)

func NewEventController(e *event.Controller) types.ServiceConfig {
	return types.ServiceConfig{
		Name:         "k8s-event-controller",
		StartupHook:  e.Start,
		ShutdownHook: e.Stop,
		RestartPolicy: types.ServiceRestartPolicy{
			RestartOnFailure: true,
			RestartCount:     3,
			RestartDelay:     time.Millisecond * 200,
		},
	}
}
