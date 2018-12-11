package main

import (
	"awesomeProject/src/mqttbroker"
	"awesomeProject/src/routers"
	"errors"
	"github.com/eclipse/paho.mqtt.golang"
	"github.com/gin-gonic/gin"
	"log"
	"net/http"
	"os"
	"time"
)

func MqStart() {
	//mqtt.DEBUG = log.New(os.Stdout, "", 0)
	mqtt.ERROR = log.New(os.Stdout, "", 0)
	mq := new(mqttbroker.MQ)
	mqttbroker.MqConnect(mq, mqttbroker.HandlerFunc)
	mqttbroker.Sub(mq, "tf/Attendance/v1/notify", 1)
	mqttbroker.Sub(mq, "chat", 1)



	//c.Disconnect(250)
	select {}
}

func main() {
	// 初始化MQTT
	go MqStart()

	// Create the Gin engine.
	gin.SetMode(gin.DebugMode)
	g := gin.New()
	middlewares := []gin.HandlerFunc{}

	// Routes.
	routers.Load(
		// Cores.
		g,

		// Middlwares.
		middlewares...,
	)

	//Ping the server to make sure the router is working.
	go func() {
		if err := pingServer(); err != nil {
			//log.Fatal("The router has no response, or it might took too long to start up.", err)
		}
		log.Print("The router has been deployed successfully.")
	}()

	log.Printf("Start to listening the incoming requests on http address: %s", ":8080")
	log.Printf(http.ListenAndServe(":8080", g).Error())
}

// pingServer pings the http server to make sure the router is working.
func pingServer() error {
	for i := 0; i < 2; i++ {
		// Ping the server by sending a GET request to `/health`.
		resp, err := http.Get("http://127.0.0.1:8080" + "/wechat/push")
		if err == nil && resp.StatusCode == 200 {
			return nil
		}

		// Sleep for a second to continue the next ping.
		log.Print("Waiting for the router, retry in 1 second.")
		time.Sleep(time.Second)
	}
	return errors.New("Cannot connect to the router.")
}
