package lib_test

import (
	"aproc/lib"
	"testing"
	"time"
)

func TestProgressWatcher(t *testing.T) {
	t.Parallel()

	watcher := lib.NewProgressWatcher(1)
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

	if err := watcher.Watch(); err != nil {
		t.Fatal(err)
	}

	select {
	case <-result:
		return
	case <-time.NewTimer(time.Second * 3).C:
		t.Fatalf("获取 PID == 1 的进程超时")
	}
}

func TestProgressEvent(t *testing.T) {
	t.Parallel()

	createEvent := lib.NewProgressEvent(10, lib.EventCreate)

	if !createEvent.IsCreate() {
		t.Fatalf("ProgressEvent CreateEvent 失败")
	}

	if createEvent.PID != 10 {
		t.Fatalf("ProgressEvent CreateEvent PID 失败")
	}

	deleteEvent := lib.NewProgressEvent(20, lib.EventDelete)

	if !deleteEvent.IsDelete() {
		t.Fatalf("ProgressEvent DeleteEvent 失败")
	}

	if deleteEvent.PID != 20 {
		t.Fatalf("ProgressEvent DeleteEvent PID 失败")
	}
}
