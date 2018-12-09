package routers

import (
	"awesomeProject/src/handler"
	"awesomeProject/src/routers/middleware"
	"github.com/gin-gonic/gin"
	"net/http"
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
		weChat.POST("/push", handler.PushController)
		weChat.POST("/pushfile", handler.PushFile)
	}

	return g
}
