package app

import (
	"time"

	"github.com/g4s8/go-lifecycle/pkg/types"
	"github.com/hoeg/beaver/internal/events"
)

func NewEventController() types.ServiceConfig {
	return types.ServiceConfig{
		Name:         "k8s-event-controller",
		StartupHook:  events.Start,
		ShutdownHook: events.Stop,
		RestartPolicy: types.ServiceRestartPolicy{
			RestartOnFailure: true,
			RestartCount:     3,
			RestartDelay:     time.Millisecond * 200,
		},
	}
}
