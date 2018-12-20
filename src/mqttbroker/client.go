package mqttbroker

import (
	"fmt"
	"os"
	"time"

	mqtt "github.com/eclipse/paho.mqtt.golang"
)

type MQCallback interface {
	handlerFunc(client mqtt.Client, msg mqtt.Message)
}

type Content struct {
	Content string `json:"content"`
}

var ChannelString chan string = make(chan string, 5)
var DeviceControlChannel = make(chan mqtt.Message, 2000)
var DeviceInfoChannel = make(chan mqtt.Message, 2000)

type BrokerInfo struct {
	client    mqtt.Client
	brokerURL string
	clientId  string
	userName  string
}

func InfoInit(mq *BrokerInfo) {
	mq.brokerURL = "tcp://106.12.130.179:1883"
	//mq.clientId = string(rand.Int())
	mq.clientId = "E470-B8A3-1"
	mq.userName = "golang-server"
}

func MqConnect(mq *BrokerInfo, handler mqtt.MessageHandler) bool {
	// 连接broker
	opts := mqtt.NewClientOptions().AddBroker(mq.brokerURL).SetClientID(mq.clientId)
	opts.SetUsername(mq.userName)
	opts.SetKeepAlive(30 * time.Second)
	opts.SetAutoReconnect(true)
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

/*
定义回调函数
*/
var HandlerFunc mqtt.MessageHandler = func(client mqtt.Client, msg mqtt.Message) {
	fmt.Printf("TOPIC: %s\n", msg.Topic())
	fmt.Printf("MSG: %s\n", msg.Payload())
	//var m Content
	//var m map[string]string
	//// TODO 添加责任链模式解析json
	//err := json.Unmarshal(msg.Payload(), &m)
	//if err != nil {
	//	log.Error("sub解析错误:", err)
	//}
	//
	//user := model.User{}
	//user.Name = ""
	//user.Balance = ""
	//rel, err := model.X.Insert(user)
	//_ = rel
	//
	//ChannelString <- m["content"]
	//fmt.Println(len(ChannelString))

}

var DeviceControlHandler mqtt.MessageHandler = func(client mqtt.Client, msg mqtt.Message) {
	DeviceControlChannel <- msg
}

var DeviceInfoHandler mqtt.MessageHandler = func(client mqtt.Client, msg mqtt.Message) {
	DeviceInfoChannel <- msg
}

func Sub(mq *BrokerInfo, topic string, qos byte, callback mqtt.MessageHandler) {
	if token := mq.client.Subscribe(topic, qos, callback); token.Wait() && token.Error() != nil {
		fmt.Println(token.Error())
		os.Exit(1)
	}
}

func Pub(mq *BrokerInfo, topic string, qos byte, payload interface{}, retain bool) error {
	token := mq.client.Publish(topic, qos, retain, payload)
	if token.Wait() && token.Error() != nil {
		return token.Error()
	}
	return nil
}

func UnSub(mq *BrokerInfo, topic string) {
	if token := mq.client.Unsubscribe(topic); token.Wait() && token.Error() != nil {
		fmt.Println(token.Error())
		os.Exit(1)
	}
}
