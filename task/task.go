package task

import (
	"context"
	"io"
	"log"
	"os"
	"time"

	"github.com/docker/docker/api/types"
	"github.com/docker/docker/client"
	"github.com/docker/go-connections/nat"
	"github.com/google/uuid"
)

type State int

type Docker struct {
	Client      *client.Client
	Config      Config
	ContainerId string
}

type DockerResult struct {
	Error       error
	Action      string
	ContainerId string
	Result      string
}

type Config struct {
	Name          string
	AttachStdin   bool
	AttachStdout  bool
	AttachStderr  bool
	Cmd           []string
	Image         string
	Memory        int64
	Disk          int64
	Env           []string
	RestartPolicy string
}

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

// The Run Function
func (d *Docker) Run() DockerResult {
	ctx := context.Background()
	reader, err := d.Client.ImagePull(ctx, d.Config.Image, types.ImagePullOptions{})
	if err != nil {
		log.Printf("Error pulling image: %v\n", d.Config)
		return DockerResult{Error: err}
	}

	io.Copy(os.Stdout, reader)
}

func (cli *Client) ContainerCreate(ctx context.Context, config *container.Config, hostConfig *container.HostConfig, networkingConfig *network.NetworkingConfig, platform *specs.Platform, containerName string)
