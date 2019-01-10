package routers

import (
	"awesomeProject/src/handler"
	"awesomeProject/src/routers/middleware"
	"net/http"

	"github.com/gin-gonic/gin"
)

func Load(g *gin.Engine, mw ...gin.HandlerFunc) *gin.Engine {
	g.Use(gin.Recovery())
	g.Use(middleware.NoCache)
	g.Use(middleware.NoCache)
	g.Use(middleware.NoCache)
	g.Use(middleware.NoCache)

	// 404Handler
	g.NoRoute(func(c *gin.Context) {
		c.String(http.StatusNotFound, "The incorrect API route.")
	})

	weChat := g.Group("weChat")
	{
		weChat.POST("/pushfile", handler.PushFile)
	}
	wechat := g.Group("wechat")
	{
		wechat.POST("/push", handler.PushController)
	}
	devices := g.Group("api/v1/devices/")
	{
		devices.POST("/control", handler.DeviceControl)
		devices.GET("/info/:id", handler.DeviceInfo)
		devices.GET("/list", handler.DeviceList)
		devices.GET("/online/:id", handler.DeviceOnlineStatus)
	}
	return g
}
