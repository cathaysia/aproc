package lib

import (
	"errors"
	"fmt"
	"io/ioutil"
	"log"
	"os"
	"strconv"
	"strings"
)

// 获取当前进程列表
func GetCurrentProgressList() ([]uint64, error) {
	dirs, err := ioutil.ReadDir("/proc")
	if err != nil {
		return nil, errors.New("不存在 /proc 路径")
	}

	result := make([]uint64, 0)

	for _, dir := range dirs {
		if !dir.IsDir() {
			continue
		}

		commFile := fmt.Sprintf("/proc/%v/comm", dir.Name())

		// 判断出来文件夹是不是进程代表的文件夹
		if _, err := os.Stat(commFile); err != nil {
			continue
		}

		if pid, err := strconv.Atoi(dir.Name()); err == nil {
			result = append(result, uint64(pid))
		}
	}

	return result, nil
}

func GetProgressNameByPID(pid uint64) (string, error) {
	proc, err := ioutil.ReadFile(fmt.Sprintf("/proc/%v/comm", pid))
	if err != nil {
		return "", err
	}

	res := strings.ReplaceAll(string(proc), "\n", "")

	return res, nil
}

func HasProgress(name string) []int {
	// 检查是否存在名为 name 的进程
	dirs, err := ioutil.ReadDir("/proc")
	if err != nil {
		log.Fatalln(err)
	}

	result := make([]int, 0)

	for _, dir := range dirs {
		if !dir.IsDir() {
			continue
		}

		commFile := fmt.Sprintf("/proc/%v/comm", dir.Name())

		if _, err := os.Stat(commFile); err != nil {
			continue
		}

		if proc, err := ioutil.ReadFile(commFile); err == nil {
			proc := strings.ReplaceAll(string(proc), "\n", "")

			if proc != name {
				continue
			}

			if pid, err := strconv.Atoi(dir.Name()); err == nil {
				result = append(result, pid)
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
