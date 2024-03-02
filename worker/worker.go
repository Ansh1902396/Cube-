package worker

import (
	"fmt"

	"github.com/Ansh1902396/cube/task"
	"github.com/golang-collections/collections/queue"
	"github.com/google/uuid"
)

type Worker struct {
	Name      string
	Queue     queue.Queue
	Db        map[uuid.UUID]task.Task
	TaskCount int
}

func (w *Worker) CollectStats() {
	fmt.Println("Collecting stats")
}

func (w *Worker) RunTask() {
	fmt.Println("Running a task")
}

func (w *Worker) StopTask() {
	fmt.Println("Stopping a task")
}
