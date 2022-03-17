package main

import (
	"aproc/lib"
	"flag"
	"os"
	"os/signal"
	"syscall"

	"github.com/sirupsen/logrus"
)

var (
	settings []lib.Settings
	err      error
)

// 根据
func createCGroupForPID(pid uint64) error {
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
	// set log level
	level := flag.String("log-level", "Debug", "Set log level, value can be Trance, Debug, Info, Error, Warning, Fatal, Panic")
	flag.Parse()
	logrus.Print(*level)

	switch *level {
	case "Trace":
		logrus.SetReportCaller(true)
		logrus.SetLevel(logrus.TraceLevel)
	case "Debug":
		logrus.SetLevel(logrus.DebugLevel)
	case "Info":
		logrus.SetLevel(logrus.InfoLevel)
	case "Error":
		logrus.SetLevel(logrus.ErrorLevel)
	case "Warning":
		logrus.SetLevel(logrus.WarnLevel)
	case "Fatal":
		logrus.SetLevel(logrus.FatalLevel)
	case "Panic":
		logrus.SetLevel(logrus.PanicLevel)
	}
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
	settingWatcher := lib.NewSettingWatcher()
	settingWatcher.Watch()

	defer settingWatcher.Close()

	// 监控 /proc 目录的变动
	watcher := lib.NewProgressWatcher(2000)
	watcher.Watch()

	defer watcher.Close()

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
		case <-settingWatcher.Event:
			if settings, err = lib.GetSettings(); err != nil {
				watcher.Error <- err

				continue
			}

			logrus.Info("重载配置")

			if err := lib.ReloadManager(settings); err != nil {
				watcher.Error <- err
			}
		case err := <-settingWatcher.Error:
			watcher.Error <- err // 将 err 处理移交
		case res := <-sigInt:
			if res == syscall.SIGINT {
				return
			}
		}
	}
}
