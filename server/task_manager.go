package server

import (
	"encoding/json"
	"sync"
	"time"
)

type Task struct {
	Dur      time.Duration
	TaskDone chan struct{}
}

func (t Task) MarshalJSON() ([]byte, error) {
	return json.Marshal(struct {
		D string `json:"duration"`
	}{
		t.Dur.String(),
	})
}

func (t *Task) UnmarshalJSON(data []byte) error {
	tmp := struct {
		D string `json:"duration"`
	}{}

	if err := json.Unmarshal(data, &tmp); err != nil {
		return err
	}

	dur, err := time.ParseDuration(tmp.D)
	if err != nil {
		return err
	}

	*t = Task{
		Dur: dur,
	}

	return nil
}

type TaskManager struct {
	m       sync.RWMutex
	tasks   []Task
	newTask chan struct{}
}

func NewTaskManager() *TaskManager {
	return &TaskManager{
		m:       sync.RWMutex{},
		tasks:   make([]Task, 0),
		newTask: make(chan struct{}, 1),
	}
}

func (tm *TaskManager) AddTask(t Task) Task {
	tm.m.Lock()
	defer tm.m.Unlock()

	go func() {
		tm.newTask <- struct{}{}
	}()

	tm.tasks = append(tm.tasks, t)

	return t
}

func (tm *TaskManager) Shedule() []Task {
	tm.m.RLock()
	defer tm.m.RUnlock()

	return tm.tasks
}

func (tm *TaskManager) Time() time.Duration {
	tm.m.RLock()
	defer tm.m.RUnlock()

	totalDuration := time.Duration(0)
	for _, task := range tm.tasks {
		totalDuration += task.Dur
	}

	return totalDuration
}

func (tm *TaskManager) Empty() bool {
	tm.m.RLock()
	defer tm.m.RUnlock()

	return len(tm.tasks) == 0
}

func (tm *TaskManager) ProccessTask() {
	tm.m.Lock()
	task := tm.tasks[0]
	tm.tasks = tm.tasks[1:]
	tm.m.Unlock()

	time.Sleep(time.Duration(task.Dur))
	task.TaskDone <- struct{}{}
}
