package proximitor

import (
	"fmt"
	"time"

	mqtt "github.com/eclipse/paho.mqtt.golang"
	"github.com/stianeikeland/go-rpio/v4"
)

const gpioTrigger = 18
const gpioEcho = 24
const topic = "home/serina/serina-rpi/distance_cm"

func Start() {
	opts := mqtt.NewClientOptions().AddBroker("tcp://localhost:1883")
	opts.AutoReconnect = true
	opts.SetKeepAlive(2 * time.Second)
	opts.SetPingTimeout(1 * time.Second)
	client := mqtt.NewClient(opts)

	if token := client.Connect(); token.Wait() && token.Error() != nil {
		panic(token.Error())
	}

	doEvery(2*time.Second, publishDistance(client))
}

func publishDistance(client mqtt.Client) func() {
	return func() {
		err := rpio.Open()
		if err != nil {
			fmt.Print(err)
		}

		distanceSensor := NewHCSR04(gpioTrigger, gpioEcho)
		distanceCm := distanceSensor.Measure()

		if token := client.Publish(topic, 0, false, fmt.Sprintf("%f", distanceCm)); token.Wait() && token.Error() != nil {
			fmt.Print(token.Error())
		}
	}
}

func doEvery(d time.Duration, f func()) {
	for range time.Tick(d) {
		f()
	}
}
