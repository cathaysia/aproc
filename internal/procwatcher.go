package internal

import (
	"os"
	"sync"
	"time"
)

type ProgressWatcher struct {
	duration    int64 // 多久检查一次进程变动，单位为毫秒
	done        chan void
	waitForExit chan void
	once        sync.Once

	Event chan ProgressEvent
	Error chan error
}

func NewProgressWatcher(msecond int64) *ProgressWatcher {
	return &ProgressWatcher{
		duration:    msecond,
		done:        make(chan void),
		waitForExit: make(chan void),

		Event: make(chan ProgressEvent),
		Error: make(chan error),
	}
}

func (watcher *ProgressWatcher) Watch() {
	watcher.once.Do(watcher.watchImpl)
}

// 启动 watch 逻辑
func (watcher *ProgressWatcher) watchImpl() {
	if _, err := os.Stat("/proc"); err != nil {
		watcher.Error <- ErrSystemNotSupport
		return
	}

	go func() {
		timer := time.NewTimer(time.Millisecond * time.Duration(watcher.duration))
		procs := make([]uint64, 0)

		for {
			select {
			case <-timer.C:
				curProcs, err := GetCurrentProgressList()
				if err != nil {
					watcher.Error <- err

					continue
				}

				// 被删除的进程列表
				for _, pid := range Difference(procs, curProcs) {
					watcher.Event <- *NewProgressEvent(pid, EventDelete)
				}
				// 新建的进程列表
				for _, pid := range Difference(curProcs, procs) {
					watcher.Event <- *NewProgressEvent(pid, EventCreate)
				}

				procs = curProcs
				// 重置定时器
				timer.Reset(time.Second * time.Duration(watcher.duration))
			case <-watcher.done:
				timer.Stop()
				close(watcher.waitForExit)

				break
			}
		}
	}()
}

// 请求退出
func (watcher *ProgressWatcher) Close() {
	close(watcher.done)
	<-watcher.waitForExit // 阻塞至退出
}

type ProcessWatcherEventType int

const (
	EventCreate ProcessWatcherEventType = 0
	EventDelete ProcessWatcherEventType = 1
)

type ProgressEvent struct {
	PID       uint64
	eventType ProcessWatcherEventType
}

func NewProgressEvent(pid uint64, eType ProcessWatcherEventType) *ProgressEvent {
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
