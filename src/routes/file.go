package routes

import (
	"fmt"
	"github.com/gin-gonic/gin"
	"io"
	"log"
	"net/http"
	"os"
)

func MyController() *gin.Engine{

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
文件上传
 */
	router.POST("/upload", func(c *gin.Context) {
		name := c.PostForm("name")
		fmt.Println(name)
		file, header, err := c.Request.FormFile("upload")
		if err != nil {
			c.String(http.StatusBadRequest, "Bad request")
			return
		}
		filename := header.Filename

		fmt.Println(file, err, filename)

		out, err := os.Create(filename)
		if err != nil {
			log.Fatal(err)
		}
		defer out.Close()
		_, err = io.Copy(out, file)
		if err != nil {
			log.Fatal(err)
		}
		c.String(http.StatusCreated, "upload successful")
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

	/**
	分组路由
	 */
	v1 := router.Group("/v1")

	v1.GET("/login", func(c *gin.Context) {
		c.String(http.StatusOK, "v1 login")
	})

	v2 := router.Group("/v2")

	v2.GET("/login", func(c *gin.Context) {
		c.String(http.StatusOK, "v2 login")
	})

	return router
}
