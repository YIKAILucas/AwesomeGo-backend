package handler

import (
	"net/http"

	"github.com/gin-gonic/gin"
)

// 发送控制指令到设备(即：发布一条消息到MQTT)
func DeviceControl(c *gin.Context) {
	c.String(http.StatusOK, "Hello %s %s", "a", "b")
}

// 获取设备指令结果(即：从数据库中查询结果)
func DeviceInfo(c *gin.Context) {
	c.String(http.StatusOK, "Hello %s %s", "a", "b")
}
