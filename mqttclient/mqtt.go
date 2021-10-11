package mqttclient

import (
	"bufio"
	"crypto/tls"
	"encoding/json"
	"fmt"
	"log"
	"net/url"
	"os"

	"github.com/Eldius/mqtt-listener-go/config"
	"github.com/Eldius/mqtt-listener-go/model"
	"github.com/Eldius/mqtt-listener-go/persistence"
	mqtt "github.com/eclipse/paho.mqtt.golang"
	"github.com/google/uuid"
)

var repo persistence.Repository

var messagePubHandler mqtt.MessageHandler = func(_ mqtt.Client, msg mqtt.Message) {
	log.Printf("Received message: %s from topic: %s\n", msg.Payload(), msg.Topic())
}

var connectHandler mqtt.OnConnectHandler = func(client mqtt.Client) {
	log.Printf("Connected (%v/%v)", client.IsConnected(), client.IsConnectionOpen())
	topic := config.GetBrokerTopic()
	token := client.Subscribe(topic, 1, msgReceivedHandler)
	if err := token.Error(); err != nil {
		log.Printf("Failed to subscribe to topic: %s", err.Error())
	}
	token.Wait()
	log.Printf("Subscribed to topic %s\n", topic)
}

var connectionLostHandler mqtt.ConnectionLostHandler = func(client mqtt.Client, err error) {
	autoReconnect := config.GetBrokerAutoreconnect()
	log.Printf("Connection lost: %v", err.Error())
	log.Printf("Auto reconnect is active: %v", autoReconnect)
	if !autoReconnect {
		log.Println("Sutting down...")
		os.Exit(1)
	}
	//client.Connect()
}

var reconnectingHandler = func(c mqtt.Client, co *mqtt.ClientOptions) {
	log.Println("Reconecting...")
}

var connectionAttemptHandler = func(broker *url.URL, tlsCfg *tls.Config) *tls.Config {
	log.Printf("Attempting to connect to '%s'", broker.String())
	return tlsCfg
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
	entry := model.NewEntry(m.Topic(), values)
	_, err := repo.Persist(entry)
	if err != nil {
		log.Printf("Failed to persist message:\n%s\n", err.Error())
	}
	m.Ack()
}

func buildOpts() *mqtt.ClientOptions {
	broker := config.GetBrokerHost()
	port := config.GetBrokerPort()
	user := config.GetBrokerUser()
	pass := config.GetBrokerPass()
	autoReconnect := config.GetBrokerAutoreconnect()

	opts := mqtt.NewClientOptions()
	opts.AddBroker(fmt.Sprintf("tcp://%s:%d", broker, port))
	opts.SetClientID("go_mqtt_client_" + uuid.NewString())
	if user != "" {
		opts.SetUsername(user)
		opts.SetPassword(pass)
	}
	//opts.SetDefaultPublishHandler(messagePubHandler)
	opts.OnConnect = connectHandler
	opts.OnConnectionLost = connectionLostHandler
	opts.OnReconnecting = reconnectingHandler
	opts.OnConnectAttempt = connectionAttemptHandler
	//opts.SetOnConnectHandler(connectHandler)
	opts.AutoReconnect = autoReconnect
	log.Printf("Broker host: '%s'\n", broker)
	log.Printf("Auto reconnect: '%v'\n", autoReconnect)

	return opts
}

func Connect(_repo persistence.Repository) {
	repo = _repo
	client := mqtt.NewClient(buildOpts())
	if token := client.Connect(); token.Wait() && token.Error() != nil {
		panic(token.Error())
	}
}

func Start() {

	client := mqtt.NewClient(buildOpts())
	if token := client.Connect(); token.Wait() && token.Error() != nil {
		panic(token.Error())
	}

	//topic := config.GetBrokerTopic()
	//token := client.Subscribe(topic, 1, msgReceivedHandler)
	//if token.Error() != nil {
	//	panic(token.Error().Error())
	//}
	//token.Wait()
	//log.Printf("Subscribed to topic %s\n", topic)

	fmt.Println("Press 'Enter' to continue...")
	bufio.NewReader(os.Stdin).ReadBytes('\n')
}
