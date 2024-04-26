package main

import (
	"fmt"
	"os"
	"strconv"

	"github.com/Ansh1902396/cube/task"
	"github.com/Ansh1902396/cube/worker"
	"github.com/docker/docker/client"

	"github.com/golang-collections/collections/queue"
	"github.com/google/uuid"
)

func createContainer() (*task.Docker, *task.DockerResult) {
	c := task.Config{
		Name:  "test-container-1",
		Image: "postgres:13",
		Env: []string{
			"POSTGRES_USER=ansh",
			"POSTGRES_PASSWORD=secret",
		},
	}
	dc, _ := client.NewClientWithOpts(client.FromEnv)
	d := task.Docker{
		Client: dc,
		Config: c,
	}

	result := d.Run()

	if result.Error != nil {
		fmt.Printf("%v\n", result.Error)
		return nil, nil
	}

	fmt.Printf("Container is started")
	return &d, &result

}

func stopContainer(d *task.Docker) *task.DockerResult {
	result := d.Stop(d.ContainerId)
	if result.Error != nil {
		fmt.Printf("%v\n", result.Error)
		return nil
	}
	fmt.Printf(
		"Container %s has been stopped and removed\n", result.ContainerId)
	return &result
}

func main() {
	host := os.Getenv("TESSERACT_HOST")
	port, _ := strconv.Atoi(os.Getenv("TESSERACT_PORT"))

	fmt.Println("starting cube worker")

	w := worker.Worker{
		Queue: *queue.New(),
		Db:    make(map[uuid.UUID]*task.Task),
	}

	api := worker.Api{
		Address: host,
		Port:    port,
		Worker:  &w,
	}

	go runTasks(&w)
	go w.CollectStats()
	api.Start()
}
