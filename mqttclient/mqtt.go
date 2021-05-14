package mqttclient

import (
	"bufio"
	"encoding/json"
	"fmt"
	"log"
	"os"

	"github.com/Eldius/mqtt-listener-go/config"
	"github.com/Eldius/mqtt-listener-go/persistence"
	mqtt "github.com/eclipse/paho.mqtt.golang"
)

var messagePubHandler mqtt.MessageHandler = func(client mqtt.Client, msg mqtt.Message) {
	fmt.Printf("Received message: %s from topic: %s\n", msg.Payload(), msg.Topic())
}

var connectHandler mqtt.OnConnectHandler = func(client mqtt.Client) {
	fmt.Println("Connected")
}

var connectLostHandler mqtt.ConnectionLostHandler = func(client mqtt.Client, err error) {
	fmt.Printf("Connect lost: %v", err)
}

var msgReceivedHandler mqtt.MessageHandler = func(c mqtt.Client, m mqtt.Message) {
	log.Printf(`
---
topic: '%s'
msg:
%s
---
`, m.Topic(), string(m.Payload()))
	var values map[string]interface{}
	json.Unmarshal(m.Payload(), &values)
	entry := persistence.NewEntry(m.Topic(), values)
	_, err := persistence.Persist(entry)
	if err != nil {
		log.Printf("Failed to persist message:\n%s\n", err.Error())
	}
	m.Ack()
}

func Start() {
	broker := config.GetBrokerHost()
	port := config.GetBrokerPort()
	user := config.GetBrokerUser()
	pass := config.GetBrokerPass()
	topic := config.GetBrokerTopic()

	log.Printf("Broker host: '%s'\n", broker)

	opts := mqtt.NewClientOptions()
	opts.AddBroker(fmt.Sprintf("tcp://%s:%d", broker, port))
	opts.SetClientID("go_mqtt_client")
	if user != "" {
		opts.SetUsername(user)
		opts.SetPassword(pass)
	}
	opts.SetDefaultPublishHandler(messagePubHandler)
	opts.OnConnect = connectHandler
	opts.OnConnectionLost = connectLostHandler
	client := mqtt.NewClient(opts)
	if token := client.Connect(); token.Wait() && token.Error() != nil {
		panic(token.Error())
	}

	token := client.Subscribe(topic, 1, msgReceivedHandler)
	if token.Error() != nil {
		panic(token.Error().Error())
	}
	token.Wait()
	fmt.Printf("Subscribed to topic %s\n", topic)
	fmt.Println("Press 'Enter' to continue...")
	bufio.NewReader(os.Stdin).ReadBytes('\n')
}
