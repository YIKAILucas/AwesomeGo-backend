package mqttbroker

import (
	"fmt"
	"math/rand"
	"os"
	"time"

	"github.com/eclipse/paho.mqtt.golang"
)

type MQ struct {
	client    mqtt.Client
	brokerURL string
	clientId  string
	userName  string
}
type MQCallback interface {
	handlerFunc(client mqtt.Client, msg mqtt.Message)
}

/**
定义回调函数
 */
var HandlerFunc mqtt.MessageHandler = func(client mqtt.Client, msg mqtt.Message) {
	fmt.Printf("TOPIC: %s\n", msg.Topic())
	fmt.Printf("MSG: %s\n", msg.Payload())
}

func MqConnect(mq *MQ, handler mqtt.MessageHandler) bool {
	mq.brokerURL = "tcp://106.12.130.179:1883"
	mq.clientId = string(rand.Int())
	mq.userName = "golang-server"

	// 连接broker
	opts := mqtt.NewClientOptions().AddBroker(mq.brokerURL).SetClientID(mq.clientId)
	opts.SetUsername(mq.userName)
	opts.SetKeepAlive(30 * time.Second)
	// 设置回调
	opts.SetDefaultPublishHandler(handler)
	//opts.SetPingTimeout(1 * time.Second)

	//create client
	mq.client = mqtt.NewClient(opts)
	if token := mq.client.Connect(); token.Wait() && token.Error() != nil {
		panic(token.Error())
	}
	return mq.client.IsConnected()
}

func Sub(mq *MQ, topic string, qos byte) {
	if token := mq.client.Subscribe(topic, qos, nil); token.Wait() && token.Error() != nil {
		fmt.Println(token.Error())
		os.Exit(1)
	}
}

func Pub(mq *MQ, topic string, qos byte, payload interface{}, retain bool) {
	token := mq.client.Publish("chat", 0, retain, payload)
	token.Wait()
}

func UnSub(mq *MQ, topic string) {
	if token := mq.client.Unsubscribe(topic); token.Wait() && token.Error() != nil {
		fmt.Println(token.Error())
		os.Exit(1)
	}
}
