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
