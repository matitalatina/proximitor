package proximitor

import (
	"time"

	"github.com/stianeikeland/go-rpio/v4"
)

const timeOutDuration = time.Second

// A type for the HC-SR04 ultrasonic distance meter
type HCSR04 struct {
	triggerPin rpio.Pin
	echoPin    rpio.Pin
}

func NewHCSR04(triggerNum, echoNum int) *HCSR04 {
	h := new(HCSR04)
	h.triggerPin = rpio.Pin(triggerNum)
	h.echoPin = rpio.Pin(echoNum)

	h.triggerPin.Output()
	h.echoPin.Input()

	return h
}

// returns the distance in centimeters
func (h *HCSR04) Measure() float64 {

	h.triggerPin.Low()
	time.Sleep(time.Microsecond * 30)
	h.triggerPin.High()
	time.Sleep(time.Microsecond * 30)
	h.triggerPin.Low()
	time.Sleep(time.Microsecond * 30)

	// sometimes the HC-SR04 stalls, if so we just break on a set timeout
	// todo check where it stalls (in which loop)
loopHigh:
	for timeout := time.After(timeOutDuration); ; {
		select {
		case <-timeout:
			break loopHigh
		default:
		}
		status := h.echoPin.Read()
		if status == rpio.High {
			break
		}
	}

	begin := time.Now()

loopLow:
	for timeout := time.After(timeOutDuration); ; {
		select {
		case <-timeout:
			break loopLow
		default:
		}
		status := h.echoPin.Read()
		if status == rpio.Low {
			break
		}
	}

	end := time.Now()
	diff := end.Sub(begin)
	cm := (float64(diff.Nanoseconds()) * 17150) / 1000000000

	return cm
}
