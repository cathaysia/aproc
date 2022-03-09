package lib

import (
	"os"
	"strconv"
	"time"
)

type ProgressWatcher struct {
	duration    int64       // 多久检查一次进程变动
	timer       *time.Timer // 定时器定时器会在每次检查结束后开始计时
	procs       *[]string   // 保存的上次进程列表，用来和当前进程进行 diff
	exit        chan bool   // 布尔标志位，用来通知 Watch() 退出
	waitForExit chan bool   // 布尔标志位，用来
	hadBeenWait bool        // 布尔标志位，用来查看 Watch 时候已经被调用过了

	Event chan *ProgressEvent
	Error chan error
}

func NewProgressWatcher(second int64) *ProgressWatcher {
	return &ProgressWatcher{
		duration:    second,
		timer:       time.NewTimer(time.Second * time.Duration(second)),
		procs:       nil,
		exit:        make(chan bool),
		waitForExit: make(chan bool),
		hadBeenWait: false,

		Event: make(chan *ProgressEvent),
		Error: make(chan error),
	}
}

// 监控 /proc 路径
func (watcher *ProgressWatcher) Watch() error {
	if watcher.hadBeenWait {
		return nil
	}

	watcher.hadBeenWait = true

	if _, err := os.Stat("/proc"); err != nil {
		return err
	}

	go func() {
		for {
			select {
			case <-watcher.timer.C:
				if watcher.procs == nil {
					for _, proc := range *GetProgressList() {
						if pid, err := strconv.Atoi(proc); err == nil {
							watcher.Event <- NewProgressEvent(pid, EventCreate)
						}
					}

					continue
				}

				curProcs := GetProgressList()

				// 被删除的进程列表
				for _, proc := range Difference(*watcher.procs, *curProcs) {
					if pid, err := strconv.Atoi(proc); err == nil {
						watcher.Event <- NewProgressEvent(pid, EventCreate)
					}
				}
				// 新建的进程列表
				for _, proc := range Difference(*curProcs, *watcher.procs) {
					if pid, err := strconv.Atoi(proc); err == nil {
						watcher.Event <- NewProgressEvent(pid, EventDelete)
					}
				}
				// 重置定时器
				watcher.timer.Reset(time.Second * time.Duration(watcher.duration))
			case <-watcher.exit:
				watcher.timer.Stop()
				watcher.waitForExit <- true

				break
			}
		}
	}()
	return nil
}

func (watcher *ProgressWatcher) Exit() {
	watcher.exit <- true
	<-watcher.waitForExit
}

type EventType int

const (
	EventCreate EventType = 0
	EventDelete EventType = 1
)

type ProgressEvent struct {
	PID       int
	eventType EventType
}

func NewProgressEvent(pid int, eType EventType) *ProgressEvent {
	return &ProgressEvent{
		PID:       pid,
		eventType: eType,
	}
}

func (event *ProgressEvent) IsCreate() bool {
	return event.eventType == EventCreate
}

func (event *ProgressEvent) IsDelete() bool {
	return event.eventType == EventDelete
}
