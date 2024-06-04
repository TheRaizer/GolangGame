package util

import "time"

func Sleep(ms uint64, c chan bool) {
	time.Sleep(time.Duration(ms))
	c <- true
}
