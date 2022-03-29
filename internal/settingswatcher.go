package internal

import (
	"sync"

	"github.com/fsnotify/fsnotify"
)

type SettingWatcherEvent int
type void struct{}

const (
	SettingEventChange SettingWatcherEvent = 0
)

type SettingWatcher struct {
	done        chan void
	waitForExit chan void
	once        sync.Once

	Event chan SettingWatcherEvent
	Error chan error
}

func NewSettingWatcher() *SettingWatcher {
	return &SettingWatcher{
		done:        make(chan void),
		waitForExit: make(chan void),
		once:        sync.Once{},

		Event: make(chan SettingWatcherEvent),
		Error: make(chan error),
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
			case <-watcher.done:
				close(watcher.waitForExit)

				return
			}
		}
	}()
}

func (watcher *SettingWatcher) Close() {
	close(watcher.done)
	<-watcher.waitForExit
}
