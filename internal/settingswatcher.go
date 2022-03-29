package internal

import (
	"sync"

	"github.com/fsnotify/fsnotify"
)

type SettingWatcherEvent int

const (
	SettingEventChange SettingWatcherEvent = 0
)

type SettingWatcher struct {
	Event       chan SettingWatcherEvent
	Error       chan error
	exit        chan bool // 布尔标志位，用来通知 Watch() 退出
	WaitForExit chan bool // 布尔标志位，用来
	once        sync.Once
}

func NewSettingWatcher() *SettingWatcher {
	return &SettingWatcher{
		Event:       make(chan SettingWatcherEvent),
		Error:       make(chan error),
		exit:        make(chan bool),
		WaitForExit: make(chan bool),
		once:        sync.Once{},
	}
}

func (watcher *SettingWatcher) Watch() {
	watcher.once.Do(watcher.watchImpl)
}

func (watcher *SettingWatcher) watchImpl() {
	settingWatcher, err := fsnotify.NewWatcher()
	if err != nil {
		watcher.Error <- err

		return
	}

	if err := settingWatcher.Add("/etc/aproc"); err != nil {
		watcher.Error <- err

		return
	}

	go func() {
		defer settingWatcher.Close()

		for {
			select {
			case _, ok := <-settingWatcher.Events:
				if !ok {
					continue
				}
				watcher.Event <- SettingEventChange
			case err, ok := <-settingWatcher.Errors:
				if !ok {
					continue
				}
				watcher.Error <- err
			case <-watcher.exit:
				watcher.WaitForExit <- true

				return
			}
		}
	}()
}

func (watcher *SettingWatcher) Close() {
	watcher.exit <- true
	<-watcher.WaitForExit
}
