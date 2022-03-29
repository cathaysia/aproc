package internal

import (
	"github.com/fsnotify/fsnotify"
)

type SettingEvent int

const (
	SettingEventChange SettingEvent = 0
)

type SettingWatcher struct {
	Event       chan SettingEvent
	Error       chan error
	hadBeenWait bool
	exit        chan bool // 布尔标志位，用来通知 Watch() 退出
	WaitForExit chan bool // 布尔标志位，用来
}

func NewSettingWatcher() *SettingWatcher {
	return &SettingWatcher{
		Event:       make(chan SettingEvent),
		Error:       make(chan error),
		hadBeenWait: false,
		exit:        make(chan bool),
		WaitForExit: make(chan bool),
	}
}

func (watcher *SettingWatcher) Watch() {
	if watcher.hadBeenWait {
		return
	}

	watcher.hadBeenWait = true

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
