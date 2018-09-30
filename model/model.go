package model

// Stats 系统状态
type Stats struct {
	CreateTime int64
	Memory
	Disks []Disk
}


type Memory struct {
	Total uint64
	Free uint64
	UsePercent float64
}

type Disk struct {
	MountPoint string
	Total uint64
	Free uint64
	Used uint64
	UsePercent float64
}