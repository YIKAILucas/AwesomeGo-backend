package mqtt

import (
	"github.com/eclipse/paho.mqtt.golang"
	"log"
	"os"
)

func MqStart() {
	//mqtt.DEBUG = log.New(os.Stdout, "", 0)
	mqtt.ERROR = log.New(os.Stdout, "", 0)
	mq:=new(MQ)
	mqConnect(mq, handlerFunc)
	sub(mq, "chat", 1)

	//c.Disconnect(250)
	select {}
}

func main() {

	MqStart()

}
