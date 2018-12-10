package handler

import (
	"github.com/gin-gonic/gin"
	"io"
	"log"
	"net/http"
	"os"
)

func MyController() *gin.Engine {

	router := gin.Default()
	f, _ := os.Create("gin.log")
	gin.DefaultWriter = io.MultiWriter(f)

	router.GET("/", func(c *gin.Context) {
		c.String(http.StatusOK, "Hello World")
	})
	router.GET("/user/:name", func(c *gin.Context) {
		// 有默认值的参数
		_ = c.DefaultQuery("firstname", "Guest")
		// 默认返回空字符串的queryParam
		_ = c.Query("lastname")
		// pathParam
		name := c.Param("name")
		c.String(http.StatusOK, "Hello %s", name)
	})

	router.POST("/form_post", func(c *gin.Context) {
		message := c.PostForm("message")
		nick := c.DefaultPostForm("nick", "anonymous")

		c.JSON(http.StatusOK, gin.H{
			"status": gin.H{
				"status_code": http.StatusOK,
				"status":      "ok",
			},
			"message": message,
			"nick":    nick,
		})
	})

	/**
多文件上传
 */
	router.POST("/multi/upload", func(c *gin.Context) {
		err := c.Request.ParseMultipartForm(200000)
		if err != nil {
			log.Fatal(err)
		}

		formdata := c.Request.MultipartForm

		files := formdata.File["upload"]
		for i, _ := range files {
			file, err := files[i].Open()
			defer file.Close()
			if err != nil {
				log.Fatal(err)
			}

			out, err := os.Create(files[i].Filename)

			defer out.Close()

			if err != nil {
				log.Fatal(err)
			}

			_, err = io.Copy(out, file)

			if err != nil {
				log.Fatal(err)
			}

			c.String(http.StatusCreated, "upload successful")

		}

	})

	return router
}
