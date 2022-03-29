package main

import (
	"aproc/internal"
	"flag"
	"os"
	"os/signal"
	"syscall"

	"github.com/sirupsen/logrus"
)

var (
	settings []internal.Settings
	err      error
	level    string
)

func main() {
	parseFlags()
	setLogLevel(level)

	checkPermission()

	if settings, err = internal.GetSettings(); err != nil {
		logrus.Fatalln(err)
	}

	sigInt := make(chan os.Signal, 2)
	signal.Notify(sigInt, syscall.SIGINT)
	// 监控 /etc/aproc/settings.json 的变动
	settingWatcher := internal.NewSettingWatcher()
	settingWatcher.Watch()

	defer settingWatcher.Close()

	// 监控 /proc 目录的变动
	procWatcher := internal.NewProgressWatcher(2000)
	procWatcher.Watch()

	defer procWatcher.Close()

	for {
		select {
		case event := <-procWatcher.Event:
			if err := handleProcressEvent(event); err != nil {
				procWatcher.Error <- err
			}
		case err := <-procWatcher.Error:
			logrus.Errorln(err)
			signal.Notify(sigInt, syscall.SIGINT) // 请求退出进程
		case event := <-settingWatcher.Event:
			if err := handleSettingsEvent(&event); err != nil {
				procWatcher.Error <- err
			}
		case err := <-settingWatcher.Error:
			procWatcher.Error <- err // 将 err 处理移交
		case res := <-sigInt:
			if res == syscall.SIGINT {
				return
			}
		}
	}
}

func setLogLevel(level string) {
	switch level {
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
}

func parseFlags() {
	flag.StringVar(&level, "log-level", "Debug", "Set log level, value can be Trance, Debug, Info, Error, Warning, Fatal, Panic")
	flag.Parse()
}

func checkPermission() {
	// 检查是否是超级用户
	if os.Geteuid() != 0 {
		logrus.Fatalln("run this program as a superuser")
	}
}

func handleProcressEvent(event *internal.ProgressEvent) error {
	if event.IsCreate() {
		logrus.Debugf("%v is Created\n", event.PID)

		if err := createCGroupForPID(event.PID); err != nil {
			return err
		}
	} else if event.IsDelete() {
		logrus.Debugf("%v is Deleted\n", event.PID)

		internal.CleanManager()
	}
	return err
}

func handleSettingsEvent(event *internal.SettingWatcherEvent) error {
	if settings, err = internal.GetSettings(); err != nil {
		return err
	}

	logrus.Info("Reload settings")

	if err := internal.ReloadManager(settings); err != nil {
		return err
	}
	return nil
}

// 根据
func createCGroupForPID(pid uint64) error {
	procName, err := internal.GetProgressNameByPID(pid)
	if err != nil {
		return err
	}

	for _, v := range settings {
		if v.Proc == procName {
			// 这里要检查是否已经存在了
			manager, err := internal.CreateManager(pid, procName, &v.Resources)
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
