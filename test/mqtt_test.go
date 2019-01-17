package main

import (
	"fmt"

	MQTT "github.com/eclipse/paho.mqtt.golang"
)

var onConnect MQTT.OnConnectHandler = func(c MQTT.Client) {
	if token := c.Subscribe("test", 1, onMessage); token.Wait() && token.Error() != nil {
		panic(token.Error())
	}

}

var count = 0
var onMessage MQTT.MessageHandler = func(c MQTT.Client, msg MQTT.Message) {
	count++
	fmt.Printf("Count:%d,MSG: %s\n", count, msg.Payload())
}

func MqttStart() {
	//MQTT.ERROR = log.New(os.Stdout, "", 0)
	opts := MQTT.NewClientOptions().AddBroker("tcp://106.12.130.179:1883")
	opts.SetClientID("E470-B8A3-1CA9")
	opts.SetOnConnectHandler(onConnect)
	opts.SetDefaultPublishHandler(onMessage)

	client := MQTT.NewClient(opts)
	if token := client.Connect(); token.Wait() && token.Error() != nil {
		panic(token.Error())
	} else {
		fmt.Println("已连接到MQTT服务器")
	}
	select {}
}

func main() {
	MqttStart()
}
