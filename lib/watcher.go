package lib

import (
	"errors"
	"os"
	"time"
)

type ProgressWatcher struct {
	duration    int64       // 多久检查一次进程变动
	timer       *time.Timer // 定时器定时器会在每次检查结束后开始计时
	procs       []uint64    // 保存的上次进程列表，用来和当前进程进行 diff
	exit        chan bool   // 布尔标志位，用来通知 Watch() 退出
	WaitForExit chan bool   // 布尔标志位，用来
	hadBeenWait bool        // 布尔标志位，用来查看 Watch 时候已经被调用过了

	Event chan *ProgressEvent
	Error chan error
}

func NewProgressWatcher(second int64) *ProgressWatcher {
	return &ProgressWatcher{
		duration:    second,
		timer:       time.NewTimer(time.Second * time.Duration(second)),
		procs:       make([]uint64, 0),
		exit:        make(chan bool),
		WaitForExit: make(chan bool),
		hadBeenWait: false,

		Event: make(chan *ProgressEvent),
		Error: make(chan error),
	}
}

// 启动 watch 逻辑
func (watcher *ProgressWatcher) Watch() error {
	if watcher.hadBeenWait {
		return nil
	}

	watcher.hadBeenWait = true

	if _, err := os.Stat("/proc"); err != nil {
		return errors.New("不存在 /proc 路径")
	}

	go func() {
		for {
			select {
			case <-watcher.timer.C:
				curProcs, err := GetCurrentProgressList()
				if err != nil {
					watcher.Error <- err

					continue
				}

				// 被删除的进程列表
				for _, pid := range Difference(watcher.procs, curProcs) {
					watcher.Event <- NewProgressEvent(pid, EventDelete)
				}
				// 新建的进程列表
				for _, pid := range Difference(curProcs, watcher.procs) {
					watcher.Event <- NewProgressEvent(pid, EventCreate)
				}

				watcher.procs = curProcs
				// 重置定时器
				watcher.timer.Reset(time.Second * time.Duration(watcher.duration))
			case <-watcher.exit:
				watcher.timer.Stop()
				watcher.WaitForExit <- true

				break
			}
		}
	}()

	return nil
}

// 请求退出
func (watcher *ProgressWatcher) Exit() {
	watcher.exit <- true
}

type EventType int

const (
	EventCreate EventType = 0
	EventDelete EventType = 1
)

type ProgressEvent struct {
	PID       uint64
	eventType EventType
}

func NewProgressEvent(pid uint64, eType EventType) *ProgressEvent {
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
