package model

import (
	"time"

	_ "github.com/go-sql-driver/mysql"
	_ "github.com/go-xorm/xorm"
	"github.com/jinzhu/gorm"
)

type Device struct {
	/* 设备信息表 */
	gorm.Model
	DeviceId   string `gorm:"not null;unique;size:14"`
	DeviceName string `gorm:"size:255"`
	IP         string `gorm:"size:100"`
}

type DevicesLifeCycle struct {
	/* 设备在线离线状态储存表 */
	gorm.Model
	DeviceId  string    `gorm:"not null;size:14"`
	OnlineAt  time.Time `gorm:"default:null"` // 上线时间
	OfflineAt time.Time `gorm:"default:null"` // 下线时间
}

func (d DevicesLifeCycle) TableName() string {
	return "devices_lifecycle"
}
