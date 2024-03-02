package task

import (
	"time"

	"github.com/docker/go-connections/nat"
	"github.com/google/uuid"
)

type State int

type Task struct {
	ID            uuid.UUID
	Name          string
	State         State
	Image         string
	Memory        string
	Disk          string
	ExposedPorts  nat.PortSet
	PortBindings  map[string]string
	RestartPolicy string
	StartTime     time.Time
	FinishTime    time.Time
}

const (
	Pending State = iota
	Scheduled
	Running
	Completed
	Failed
)
