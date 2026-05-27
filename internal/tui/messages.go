package tui

import (
	"runpod-top/internal/api"
	"time"
)

type TickMsg time.Time

type PodsUpdatedMsg struct {
	Pods []api.Pod
}

type APIErrorMsg struct {
	Err error
}

type ActionResultMsg struct {
	PodID  string
	Action string
	Err    error
}
