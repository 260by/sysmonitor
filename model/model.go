package model

// Stats 系统状态
type Monitor struct {
	CreateTime int64
	HostName string
	IP	string
	CPUPercent float64
	MemoryPercent float64
	DisksPercent string
}