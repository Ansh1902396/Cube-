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

// Struct for TaskEvent
// user can control specific task by using this struct
type TaskEvent struct {
	ID        uuid.UUID
	State     State
	TimeStamp time.Time
	Task      Task
}

const (
	Pending State = iota
	Scheduled
	Running
	Completed
	Failed
)
