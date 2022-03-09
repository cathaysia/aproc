package lib

import (
	"fmt"
	"io/ioutil"
	"log"
	"os"
	"strings"
)

// 获取当前进程列表
func GetProgressList() *[]string {
	// 检查时候存在名为 name 的进程
	dirs, err := ioutil.ReadDir("/proc")
	if err != nil {
		log.Fatalln(err)
	}

	result := make([]string, 0)

	for _, dir := range dirs {
		if dir.IsDir() {
			commFile := fmt.Sprintf("/proc/%v/comm", dir.Name())
			if _, err := os.Stat(commFile); err == nil {
				// 这里就能判断出来是进程代表的文件夹了
				result = append(result, dir.Name())
			}
		}
	}

	return &result
}

func HasProgress(name string) []string {
	// 检查时候存在名为 name 的进程
	dirs, err := ioutil.ReadDir("/proc")
	if err != nil {
		log.Fatalln(err)
	}

	result := make([]string, 0)

	for _, dir := range dirs {
		if !dir.IsDir() {
			continue
		}

		commFile := fmt.Sprintf("/proc/%v/comm", dir.Name())

		if _, err := os.Stat(commFile); err == nil {
			if proc, err := ioutil.ReadFile(commFile); err == nil {
				proc := strings.ReplaceAll(string(proc), "\n", "")
				if proc == name {
					result = append(result, dir.Name())
				}
			}
		}
	}

	return result
}

func ProgressCGroup(pid string) string {
	filePath := fmt.Sprintf("/proc/%v/cgroup", pid)

	var result string

	if content, err := ioutil.ReadFile(filePath); err == nil {
		result = strings.ReplaceAll(string(content)[3:], "\n", "")
	} else {
		log.Fatalln(err)
	}

	return result
}

func CreateCGroup(name string) {
	filePath := fmt.Sprintf("/sys/fs/cgroup/%v", name)
	if err := os.Mkdir(filePath, os.ModePerm); err != nil {
		log.Println(err)
	}
}

func AddToCGroup(cgroup string, pid string) {
	filePath := fmt.Sprintf("/sys/fs/cgroup/%v/cgroup.procs", cgroup)

	if _, err := os.Stat(filePath); err != nil {
		log.Fatalln(err)
	}

	if err := os.WriteFile(filePath, []byte(pid), os.ModePerm); err != nil {
		log.Fatalln(err)
	}
}

func SetCPULimit(cgroup string, max string) {
	filePath := fmt.Sprintf("/sys/fs/cgroup/%v/cpu.max", cgroup)

	if _, err := os.Stat(filePath); err != nil {
		log.Fatalln(err)
	}

	if err := os.WriteFile(filePath, []byte(max), os.ModePerm); err != nil {
		log.Fatalln(err)
	}
}
