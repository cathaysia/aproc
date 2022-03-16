package main

import (
	"aproc/lib"
	"os"
	"os/signal"
	"syscall"

	"github.com/fsnotify/fsnotify"
	"github.com/sirupsen/logrus"
)

var (
	settings []lib.Settings
	err      error
)

// 根据
func createCGroupForPID(pid uint64) error {
	logrus.SetLevel(logrus.TraceLevel)
	logrus.SetReportCaller(true)

	procName, err := lib.GetProgressNameByPID(pid)
	if err != nil {
		return err
	}

	for _, v := range settings {
		if v.Proc == procName {
			// 这里要检查是否已经存在了
			manager, err := lib.CreateManager(pid, procName, &v.Resources)
			if err != nil {
				return err
			}

			err = manager.AddProc(pid)
			if err != nil {
				return err
			}
		}
	}

	return nil
}

func main() {
	// 检查是否是超级用户
	if os.Geteuid() != 0 {
		logrus.Fatalln("请使用超级用户运行此程序")
	}

	if settings, err = lib.GetSettings(); err != nil {
		logrus.Fatalln(err)
	}

	sigInt := make(chan os.Signal, 2)
	signal.Notify(sigInt, syscall.SIGINT)
	// 监控 /etc/aproc/settings.json 的变动
	settingWatcher, err := fsnotify.NewWatcher()
	if err != nil {
		logrus.Fatalln(err)
	}
	defer settingWatcher.Close()

	if err := settingWatcher.Add("/etc/aproc"); err != nil {
		logrus.Fatalln(err)
	}

	// 监控 /proc 目录的变动
	watcher := lib.NewProgressWatcher(2000)
	watcher.Watch()

	for {
		select {
		case event := <-watcher.Event:
			if event.IsCreate() {
				logrus.Debugf("%v is Created\n", event.PID)

				if err := createCGroupForPID(event.PID); err != nil {
					watcher.Error <- err
				}
			} else if event.IsDelete() {
				logrus.Debugf("%v is Deleted\n", event.PID)

				lib.CleanManager()
			}
		case err := <-watcher.Error:
			logrus.Errorln(err)
			signal.Notify(sigInt, syscall.SIGINT) // 请求退出进程
		case _, ok := <-settingWatcher.Events:
			if !ok {
				signal.Notify(sigInt, syscall.SIGINT)

				continue
			}

			if settings, err = lib.GetSettings(); err != nil {
				watcher.Error <- err

				continue
			}

			logrus.Traceln("重载配置")

			if err := lib.ReloadManager(settings); err != nil {
				watcher.Error <- err
			}
		case err, ok := <-settingWatcher.Errors:
			if !ok {
				signal.Notify(sigInt, syscall.SIGINT)

				continue
			}
			watcher.Error <- err
		case res := <-sigInt:
			if res == syscall.SIGINT {
				watcher.Exit()
				<-watcher.WaitForExit

				return
			}
		}
	}
}
