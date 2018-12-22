package mqttbroker

import (
	"awesomeProject/src/model"
	"encoding/json"
	"fmt"
	"time"

	mqtt "github.com/eclipse/paho.mqtt.golang"
)

/*
默认回调函数
*/
var HandlerFunc mqtt.MessageHandler = func(client mqtt.Client, msg mqtt.Message) {
	fmt.Printf("TOPIC: %s\n", msg.Topic())
	fmt.Printf("MSG: %s\n", msg.Payload())
	//var m Content
	//var m map[string]string
	//// TODO 添加责任链模式解析json
	//err := json.Unmarshal(msg.Payload(), &m)
	//if err != nil {
	//	log.Error("sub解析错误:", err)
	//}
	//
	//user := model.User{}
	//user.Name = ""
	//user.Balance = ""
	//rel, err := model.X.Insert(user)
	//_ = rel
	//
	//ChannelString <- m["content"]
	//fmt.Println(len(ChannelString))

}

var DeviceControlHandler mqtt.MessageHandler = func(client mqtt.Client, msg mqtt.Message) {
	DeviceControlChannel <- msg
}

var DeviceInfoHandler mqtt.MessageHandler = func(client mqtt.Client, msg mqtt.Message) {
	DeviceInfoChannel <- msg
}

type onlineMessage struct {
	DeviceId   string `json:"clientid" binding:"required"`
	DeviceName string `json:"username" binding:"required"`
	IP         string `json:"ipaddress" binding:"required"`
}

type offlineMessage struct {
	DeviceId string `json:"clientid" binding:"required"`
}

var DeviceOnlineHandler mqtt.MessageHandler = func(client mqtt.Client, msg mqtt.Message) {
	/* 设备上线时，将信息存入表中 */
	fmt.Printf("检测到设备上线，话题：%s，消息：%s\n", msg.Topic(), msg.Payload())
	//topic := msg.Topic()
	payload := msg.Payload()

	var json_data onlineMessage
	err := json.Unmarshal(payload, &json_data)
	if err != nil {
		fmt.Println("设备上线消息解码失败")
		return
	}
	info := model.DevicesOnlineOffLineStatus{DeviceId: json_data.DeviceId, OnlineAt: time.Now()}
	err = model.DB.Create(&info).Error
	if err != nil {
		fmt.Printf("设备%s上线消息添加失败：%s\n", json_data.DeviceId, err)
		return
	}
	fmt.Printf("设备%s上线信息已添加\n", json_data.DeviceId)

	var device_info model.Device
	model.DB.Where(&model.Device{DeviceId: json_data.DeviceId}).First(&device_info)
	if device_info.DeviceId == "" {
		// 新设备，添加到表
		new_device_info := model.Device{DeviceId: json_data.DeviceId, DeviceName: json_data.DeviceName, IP: json_data.IP}
		err = model.DB.Create(&new_device_info).Error
		if err != nil {
			fmt.Printf("发现新设备：%s，无法将其添加到设备表中：%s\n", json_data.DeviceId, err)
			return
		}
		fmt.Printf("发现新设备：%s，已将其添加到设备表中\n", json_data.DeviceId)
	} else {
		// 已在库中，更新信息
		err = model.DB.Model(&device_info).Update(model.Device{DeviceName: json_data.DeviceName, IP: json_data.IP}).Error
		if err != nil {
			fmt.Printf("设备：%s，资料更新失败：%s\n", json_data.DeviceId, err)
			return
		}
		fmt.Printf("设备：%s，资料已更新\n", json_data.DeviceId)
	}
}

var DeviceOfflineHandler mqtt.MessageHandler = func(client mqtt.Client, msg mqtt.Message) {
	/* 设备下线时，刷新表，在线周期结束 */
	fmt.Printf("检测到设备下线，话题：%s，消息：%s\n", msg.Topic(), msg.Payload())
	payload := msg.Payload()

	var json_data offlineMessage
	err := json.Unmarshal(payload, &json_data)
	if err != nil {
		fmt.Println("设备下线消息解码失败")
		return
	}
	var info model.DevicesOnlineOffLineStatus
	model.DB.Where(&model.DevicesOnlineOffLineStatus{DeviceId: json_data.DeviceId}).Last(&info)
	if info.DeviceId == "" {
		// 没有记录时，生成一条
		new_info := model.DevicesOnlineOffLineStatus{DeviceId: json_data.DeviceId, OfflineAt: time.Now()}
		err = model.DB.Create(&new_info).Error
		if err != nil {
			fmt.Printf("设备%s下线消息添加失败：%s\n", json_data.DeviceId, err)
			return
		}
		fmt.Printf("设备%s下线信息已添加\n", json_data.DeviceId)
	} else if info.DeviceId != "" && !info.OfflineAt.IsZero() {
		// 有记录且已填时，生成一条
		new_info := model.DevicesOnlineOffLineStatus{DeviceId: json_data.DeviceId, OfflineAt: time.Now()}
		err = model.DB.Create(&new_info).Error
		if err != nil {
			fmt.Printf("设备%s下消息添加失败：%s\n", json_data.DeviceId, err)
			return
		}
		fmt.Printf("设备%s下线信息已添加\n", json_data.DeviceId)
	} else if info.DeviceId != "" && info.OfflineAt.IsZero() {
		// 有记录且未填时，更新一条
		model.DB.Model(&info).Update(model.DevicesOnlineOffLineStatus{OfflineAt: time.Now()})
		fmt.Printf("设备%s下线信息已更新\n", json_data.DeviceId)
	}
}
