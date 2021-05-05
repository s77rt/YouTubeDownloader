package YouTubeDownloader

import (
	"os"
	"fmt"
	"context"
	"path/filepath"
)

type Task struct {
	Context context.Context
	CancelFunc context.CancelFunc
	Uuid string
	Tmpdir string
	PreCancelFuncs []func()
	PostCancelFuncs []func()
}

func (Task *Task) Done() {
	for _, f := range Task.PreCancelFuncs {
		f()
	}
	Task.CancelFunc()
	for _, f := range Task.PostCancelFuncs {
		f()
	}
}

func (Task *Task) Abort() {
	for _, f := range Task.PreCancelFuncs {
		f()
	}
	Task.CancelFunc()
	for _, f := range Task.PostCancelFuncs {
		f()
	}
}

func NewTask(uuid string) *Task {
	ctx, cancel := context.WithCancel(context.Background())

	task_tmpdir := filepath.Join(tmpdir, uuid)
	os.MkdirAll(task_tmpdir, 0o755)

	return &Task{
		ctx,
		cancel,
		uuid,
		task_tmpdir,
		[]func(){},
		[]func(){
			func(){os.RemoveAll(task_tmpdir)},
		},
	}
}

///////////////////////////////////////

type Tasks map[string]*Task

func (Tasks Tasks) New(uuid string) *Task {
	Tasks[uuid] = NewTask(uuid)
	return Tasks[uuid]
}

func (Tasks Tasks) Done(uuid string) error {
	task, ok := Tasks[uuid]
	if ok != true {
		return fmt.Errorf("task does not exist")
	}
	task.Done()
	delete(Tasks, uuid)
	return nil
}

func (Tasks Tasks) Abort(uuid string) error {
	task, ok := Tasks[uuid]
	if ok != true {
		return fmt.Errorf("task does not exist")
	}
	task.Abort()
	delete(Tasks, uuid)
	return nil
}

func NewTasks() Tasks {
	return make(Tasks)
}
