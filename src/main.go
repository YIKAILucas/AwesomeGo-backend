package main

import (
	"awesomeProject/src/conf"
	"awesomeProject/src/handler"
	"awesomeProject/src/model"
	"awesomeProject/src/middleware/mqttbroker"
	"awesomeProject/src/routers"
	"errors"
	"fmt"
	"log"
	"net/http"
	"os"
	"time"

	"github.com/eclipse/paho.mqtt.golang"
	"github.com/gin-gonic/gin"
	"github.com/spf13/pflag"
	"github.com/spf13/viper"
)

func MqStart() {
	//mqtt.DEBUG = log.New(os.Stdout, "", 0)
	mqtt.ERROR = log.New(os.Stdout, "", 0)
	mq := &mqttbroker.BrokerInfo{}
	mqttbroker.InfoInit(mq)
	mqttbroker.MqConnect(mq, mqttbroker.HandlerFunc)
	mqttbroker.Sub(mq, "tf/Attendance/v1/notify", 1, nil)
	mqttbroker.Sub(mq, "chat", 1, nil)
	mqttbroker.Sub(mq, "tf/Attendance/v1/devices/+/control", 1, mqttbroker.DeviceControlHandler)
	mqttbroker.Sub(mq, "tf/Attendance/v1/devices/+/info", 1, mqttbroker.DeviceInfoHandler)
	mqttbroker.Sub(mq, "$SYS/brokers/+/clients/+/connected", 1, mqttbroker.DeviceOnlineHandler)
	mqttbroker.Sub(mq, "$SYS/brokers/+/clients/+/disconnected", 1, mqttbroker.DeviceOfflineHandler)
	//c.Disconnect(250)

	for data := range handler.HTTPPubChannel {
		err := mqttbroker.Pub(mq, data.Topic, 1, data.Payload, false)
		if err != nil {
			fmt.Printf("消息发布失败：%s,Payload：%s", data.Topic, data.Payload)
		}
	}
}

var (
	cfg = pflag.StringP("conf", "c", "", "apiserver conf file path.")
)

func main() {
	pflag.Parse()
	// init conf
	if err := conf.Init(*cfg); err != nil {
		panic(err)
	}
	model.DBInit()

	// 初始化MQTT
	go MqStart()
	go mqttbroker.MqSaveDeviceControlLoop()
	go mqttbroker.MqSaveDeviceInfoLoop()

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
	log.Printf(http.ListenAndServe(viper.GetString("addr"), g).Error())
}

// pingServer pings the http server to make sure the router is working.
func pingServer() error {
	for i := 0; i < viper.GetInt("max_ping_count"); i++ {
		// Ping the server by sending a GET request to `/health`.
		resp, err := http.Get(viper.GetString("url") + "/sd/health")
		if err == nil && resp.StatusCode == 200 {
			return nil
		}

		// Sleep for a second to continue the next ping.
		log.Print("Waiting for the router, retry in 1 second.")
		time.Sleep(time.Second)
	}
	return errors.New("Cannot connect to the router.")
}
