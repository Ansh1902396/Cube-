package manager

import (
	"bytes"
	"encoding/json"
	"fmt"
	"log"
	"net/http"

	"github.com/Ansh1902396/cube/task"
	"github.com/golang-collections/collections/queue"
	"github.com/google/uuid"
)

type Manager struct {
	Pending        queue.Queue
	TaskDb         map[string][]task.Task
	EventDb        map[string][]task.TaskEvent
	Workers        []string
	WorkersTaskMap map[string][]uuid.UUID
	TaskWorkerMap  map[uuid.UUID]string
	LastWorker     int
}

func (m *Manager) SelectWorker() string {
	var newWorker int
	if m.LastWorker+1 == len(m.Workers) {
		newWorker = m.LastWorker + 1
		m.LastWorker++
	} else {
		newWorker = 0
		m.LastWorker = 0
	}

	return m.Workers[newWorker]
}

func (m *Manager) UpdateTasks() {
	for _, worker := range m.Workers {
		log.Printf("Updating tasks for worker %v\n", worker)
		url := fmt.Sprintf("http://%s:%d/tasks", worker)

		resp, err := http.Get(url)

		if err != nil {
			log.Printf("Error getting tasks for worker %v\n", worker)
			continue
		}
		if resp.StatusCode != http.StatusOK {
			log.Printf("Error getting tasks for worker %v\n", worker)
			continue
		}

		d := json.NewDecoder(resp.Body)
		var tasks []task.Task

		err = d.Decode(&tasks)

		if err != nil {
			log.Printf("Error decoding tasks for worker %v\n", worker)
			continue

		}
	}

	for _, t := range tasks {
		log.Printf("Attempting to update task %v", t.ID)
		_, ok := m.TaskDb[t.ID]
		if !ok {
			log.Printf("Task with ID %s not found\n", t.ID)
			return
		}
		if m.TaskDb[t.ID].State != t.State {
			m.TaskDb[t.ID].State = t.State
		}
		m.TaskDb[t.ID].StartTime = t.StartTime
		m.TaskDb[t.ID].FinishTime = t.FinishTime
		m.TaskDb[t.ID].ContainerID = t.ContainerID
	}
}

func (m *Manager) SendWork() {
	if m.Pending.Len() > 0 {
		w := m.SelectWorker()

		e := m.Pending.Dequeue()
		te := e.(task.TaskEvent)

		t := te.Task
		log.Printf("Pulled %f off pending queue", t)
		m.EventDb[te.ID] = &te

		m.WorkersTaskMap[w] = append(m.WorkersTaskMap[w], t.ID)
		m.TaskWorkerMap[t.ID] = w
		t.State = task.Scheduled
		m.TaskDb[t.ID] = &t

		data, err := json.Marshal(te)
		if err != nil {
			log.Printf("Error marshalling task event: %v\n", te)
		}

		url := fmt.Sprintf("http://%s:%d/start", worker)
		resp, err := http.Post(url, "application/json", bytes.NewBuffer(data))

		if err != nil {
			log.Printf("Error sending task to worker: %v\n", w)
			m.Pending.Enqueue(t)
			return
		}

		d := json.NewDecoder(resp.Body)
		if resp.StatusCode != http.StatusCreated {
			e := w.ErrResponse{}

			err := d.Decode(&e)

			if err != nil {
				fmt.Println("Error decoding response")
				return
			}

			log.Printf("Error starting task: %v\n", e.Message)
			return
		}

		t = task.Task{}

		err = d.Decode(&t)
		if err != nil {
			fmt.Println("Error decoding response")
			return
		}
		log.Printf("%#v\n", t)

	} else {
		log.Println("No work in the queue")
	}
}
