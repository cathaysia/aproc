package internal_test

import (
	"aproc/internal"
	"testing"
	"time"
)

func TestProgressWatcher(t *testing.T) {
	t.Parallel()

	watcher := internal.NewProgressWatcher(1)
	result := make(chan bool)

	go func() {
		for {
			select {
			case event := <-watcher.Event:
				if event.IsCreate() {
					if event.PID == 1 {
						result <- true
					}
				}
			}
		}
	}()

	watcher.Watch()

	select {
	case <-result:
		return
	case <-time.NewTimer(time.Second * 3).C:
		t.Fatalf("获取 PID == 1 的进程超时")
	}
}

func TestProgressEvent(t *testing.T) {
	t.Parallel()

	createEvent := internal.NewProgressEvent(10, internal.EventCreate)

	if !createEvent.IsCreate() {
		t.Fatalf("ProgressEvent CreateEvent 失败")
	}

	if createEvent.PID != 10 {
		t.Fatalf("ProgressEvent CreateEvent PID 失败")
	}

	deleteEvent := internal.NewProgressEvent(20, internal.EventDelete)

	if !deleteEvent.IsDelete() {
		t.Fatalf("ProgressEvent DeleteEvent 失败")
	}

	if deleteEvent.PID != 20 {
		t.Fatalf("ProgressEvent DeleteEvent PID 失败")
	}
}
