package handler

import (
	"awesomeProject/src/middleware/mongo/mongoDao"
	"github.com/gin-gonic/gin"
)

func ActiveDevice(c *gin.Context) {
	mongoDao.GetDAU()
	//c.JSON(http.StatusOK, gin.H{"success": true, "msg": "", "result": doc})
}
