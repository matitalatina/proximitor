package proximitor

import (
	"fmt"
	"time"

	mqtt "github.com/eclipse/paho.mqtt.golang"
	"github.com/stianeikeland/go-rpio/v4"
)

const gpioTrigger = 18
const gpioEcho = 24

func Start() {
	doEvery(2*time.Second, publishDistance)
}

func publishDistance() {
	err := rpio.Open()
	if err != nil {
		fmt.Print(err)
	}
	distance := NewHCSR04(gpioTrigger, gpioEcho)
	distanceCm := distance.Measure()
	fmt.Printf("%f cm", distanceCm)

	const TOPIC = "home/serina/serina-rpi/distance_cm"

	opts := mqtt.NewClientOptions().AddBroker("tcp://localhost:1883")
	client := mqtt.NewClient(opts)
	if token := client.Connect(); token.Wait() && token.Error() != nil {
		fmt.Print(token.Error())
	}

	if token := client.Publish(TOPIC, 0, false, fmt.Sprintf("%f", distanceCm)); token.Wait() && token.Error() != nil {
		fmt.Print(token.Error())
	}
}

func doEvery(d time.Duration, f func()) {
	for range time.Tick(d) {
		f()
	}
}
