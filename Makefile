
.EXPORT_ALL_VARIABLES:

MQTT_BROKER_HOST = 192.168.100.195
MQTT_BROKER_PORT = 1883
MQTT_BROKER_USER = ""
MQTT_BROKER_PASS = ""
MQTT_BROKER_TOPIC = \#
MQTT_BROKER_RECONNECT = true
MQTT_CORS_ALLOW_ORIGINS = http://localhost:3000

dockerbuild:
	docker buildx build \
		--push \
		--platform linux/arm/v7,linux/arm64/v8,linux/amd64 \
		--tag eldius/mqtt-listener-go:latest \
		--tag eldius/mqtt-listener-go:$(shell git rev-parse --short HEAD) \
		.

startback:
	go run main.go start

startfront:
	cd static ; yarn start


listen:
	mosquitto_sub -h 192.168.100.195 -t '#'
