package model

import (
	"github.com/go-xorm/xorm"
	"github.com/go-xorm/core"
	_ "github.com/mattn/go-sqlite3"
)

// User 用户数据表
type User struct {
	ID       int    `xorm:"pk autoincr notnull"`
	UserName     string `xorm:"varchar(32) notnull unique index"`
	Password string `xorm:"varchar(128) notnull"`
}

// Assets 资产数据表
type Assets struct {
	ID         int    `xorm:"pk autoincr notnull"`
	IP         string `xorm:"notnull unique"`
	HostName   string `xorm:"notnull unique"`
}

// Monitor 系统状态
type Monitor struct {
	ID int `xorm:"pk autoincr notnull"`
	CreateTime int64 `xorm:"index"`
	HostName string	`xorm:"varchar(64) notnull index"`
	IP	string	`xorm:"varchar(32) notnull index"`
	CPUPercent float64
	MemoryPercent float64
	DisksPercent string
	Load1 float64
	Load5 float64
	Load15 float64
	TCPEstablished int
	TCPTimeWait int
}

// Connect 连接数据库
func Connect(driveName, dataSourceName string, showSQL bool) (*xorm.Engine, error) {
	orm, err := xorm.NewEngine(driveName, dataSourceName)
	if err != nil {
		return nil, err
	}
	orm.SetMapper(core.GonicMapper{})
	orm.ShowSQL(showSQL)
	return orm, nil
}

// Migrate 同步表结构
func Migrate(orm *xorm.Engine,) error {
	err := orm.Sync2(&User{}, &Assets{}, &Monitor{})
	if err != nil {
		return err
	}
	return nil
}