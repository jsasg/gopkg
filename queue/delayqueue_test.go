package queue

import (
	"testing"
	"time"
)

func TestDelayQueue(t *testing.T) {
	dq := NewDelayQueue(64)
	dq.Start()
	defer dq.Stop()

	dq.AfterFunc(20*time.Second, func(v chan<- interface{}) {
		t.Log(20)
	})

	dq.AfterFunc(10*time.Second, func(v chan<- interface{}) {
		t.Log(10)
	})

	dq.AfterFunc(1*time.Second, func(v chan<- interface{}) {
		t.Log(1)
	})

	dq.AfterFunc(9*time.Second, func(v chan<- interface{}) {
		t.Log(9)
	})

	time.Sleep(20 * time.Second)

	dq.AfterFunc(5*time.Second, func(v chan<- interface{}) {
		t.Log(5)
	})

	dq.AfterFunc(2*time.Second, func(v chan<- interface{}) {
		t.Log(2)
	})

	select {}
}
