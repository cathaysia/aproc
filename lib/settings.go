package lib

import (
	"encoding/json"
	"io/ioutil"
	"os"

	v2 "github.com/containerd/cgroups/v2"
)

// import types of v2 to json
// type CPU struct { // v2/cpu.go
// 	Weight uint64 `json:"weight"`
// 	Max    string `json:"max"`
// 	Cpus   string `json:"cpus"`
// 	Mems   string `json:"mems"`
// }

// type Memory struct {
// 	Swap int64 `json:"swap"`
// 	Max  int64 `json:"max"`
// 	Low  int64 `json:"low"`
// 	High int64 `json:"high"`
// }
// type Pids struct {
// 	Max int64 `json:"max"`
// }

// type BFQ struct {
// 	Weight uint16 `json:"weight"`
// }

// type Entry struct {
// 	Type  string `json:"type"` // can be rbps, wbps, riops, wiops. it is IOType in cgroup
// 	Major int64  `json:"major"`
// 	Minor int64  `json:"minor"`
// 	Rate  uint64 `json:"rate"`
// }

// type IO struct {
// 	BFQ BFQ     `json:"bfq"`
// 	Max []Entry `json:"max"`
// }

// type RDMA struct {
// 	Limit []RDMAEntry `json:"limit"`
// }

// type RDMAEntry struct {
// 	Device     string `json:"device"`
// 	HcaHandles uint32 `json:"hcahandles"`
// 	HcaObjects uint32 `json:"hcaobjects"`
// }

// type HugeTlb []HugeTlbEntry

// type HugeTlbEntry struct {
// 	HugePageSize string `json:"hugepagesize"`
// 	Limit        uint64 `json:"hugepagesize"`
// }

// // Resources for a cgroups v2 unified hierarchy
// type Resources struct {
// 	CPU     CPU                       `json:"cpu"`
// 	Memory  Memory                    `json:"memory"`
// 	Pids    Pids                      `json:"pids"`
// 	IO      IO                        `json:"io"`
// 	RDMA    RDMA                      `json:"rdma"`
// 	HugeTlb HugeTlb                   `json:"hugeTlb"`
// 	Devices []specs.LinuxDeviceCgroup `json:"devices"`
// }

type Settings struct {
	Proc      string
	Resources v2.Resources
}

// 从 /etc/aproc/settings.json 获取数据
func GetSettings() ([]Settings, error) {
	if _, err := os.Stat("/etc/aproc"); err != nil {
		if err = os.MkdirAll("/etc/aproc", os.ModePerm); err != nil {
			return nil, err
		}
	}
	// 读取 settingsFile 文件
	settingsFile, err := os.OpenFile("/etc/aproc/settings.json", os.O_RDWR|os.O_CREATE, os.ModePerm)
	if err != nil {
		return nil, err
	}

	defer settingsFile.Close()

	settingsData, err := ioutil.ReadAll(settingsFile)
	if err != nil {
		return nil, err
	}

	var settings []Settings

	if err = json.Unmarshal(settingsData, &settings); err != nil {
		return nil, err
	}

	return settings, nil
}

func CreateEmptySettings() error {
	// create file
	if _, err := os.Stat("/etc/aproc"); err != nil {
		if err = os.MkdirAll("/etc/aproc", os.ModePerm); err != nil {
			return err
		}
	}

	settingsFile, err := os.OpenFile("/etc/aproc/settings.json", os.O_RDWR|os.O_CREATE, os.ModePerm)
	if err != nil {
		return err
	}

	settingsFile.Close()

	settings := []Settings{
		{
			Proc: "systemd",
		},
	}

	data, err := json.Marshal(settings)
	if err != nil {
		return err
	}

	return ioutil.WriteFile("/etc/aproc/settings.json", data, os.ModePerm)
}
