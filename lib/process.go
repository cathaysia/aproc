package lib

import (
	"errors"
	"fmt"
	"io/ioutil"
	"os"
	"strconv"
	"strings"
)

var (
	ErrSystemNotSupport = errors.New("不支持不包含 /proc 路径的系统")
	errPocessInvalid    = errors.New("进程无效")
)

func ProcessInvalidError(pid uint64) error {
	return fmt.Errorf("ProcessInvalidError %w : %v", errPocessInvalid, pid)
}

// 获取当前进程列表
func GetCurrentProgressList() ([]uint64, error) {
	dirs, err := ioutil.ReadDir("/proc")
	if err != nil {
		return nil, ErrSystemNotSupport
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
		return "", ProcessInvalidError(pid)
	}

	res := strings.ReplaceAll(string(proc), "\n", "")

	return res, nil
}
