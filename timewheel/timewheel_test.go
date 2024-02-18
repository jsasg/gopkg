package timewheel

import (
	"fmt"
	"testing"
	"time"
)

func TestTimeWheel(t *testing.T) {
	tw := New(1*time.Second, 60)
	tw.Start()
	defer tw.Stop()

	fmt.Println(time.Now().Format("2006-01-02 15:04:05"))

	tw.AfterFunc(70*time.Second, func() {
		fmt.Println(time.Now().Format("2006-01-02 15:04:05"))
		t.Log("abcd", 70)
	})

	tw.AfterFunc(90*time.Second, func() {
		fmt.Println(time.Now().Format("2006-01-02 15:04:05"))
		t.Log("abcd", 90)
	})

	tw.AfterFunc(3600*time.Second, func() {
		fmt.Println(time.Now().Format("2006-01-02 15:04:05"))
		t.Log("abcd", 3600)
	})

	tw.AfterFunc(59*time.Second, func() {
		fmt.Println(time.Now().Format("2006-01-02 15:04:05"))
		t.Log("abcd", 59)
	})

	tw.AfterFunc(60*time.Second, func() {
		fmt.Println(time.Now().Format("2006-01-02 15:04:05"))
		t.Log("abcd", 60)
	})

	tw.AfterFunc(61*time.Second, func() {
		fmt.Println(time.Now().Format("2006-01-02 15:04:05"))
		t.Log("abcd", 61)
	})
	tw.AfterFunc(111*time.Second, func() {
		fmt.Println(time.Now().Format("2006-01-02 15:04:05"))
		t.Log("abcd", 111)
	})
	tw.AfterFunc(120*time.Second, func() {
		fmt.Println(time.Now().Format("2006-01-02 15:04:05"))
		t.Log("abcd", 120)
	})

	tw.AfterFunc(10*time.Second, func() {
		fmt.Println(time.Now().Format("2006-01-02 15:04:05"))
		t.Log("abcd", 10)
	})
	tw.AfterFunc(11*time.Second, func() {
		fmt.Println(time.Now().Format("2006-01-02 15:04:05"))
		t.Log("abcd", 11)
	})
	time.Sleep(10 * time.Second)
	tw.AfterFunc(1*time.Second, func() {
		fmt.Println(time.Now().Format("2006-01-02 15:04:05"))
		t.Log("abcd-2", 1)
	})
	tw.AfterFunc(19*time.Second, func() {
		fmt.Println(time.Now().Format("2006-01-02 15:04:05"))
		t.Log("abcd", 19)
	})
	tw.AfterFunc(20*time.Second, func() {
		fmt.Println(time.Now().Format("2006-01-02 15:04:05"))
		t.Log("abcd", 20)
	})
	tw.AfterFunc(1*time.Second, func() {
		fmt.Println(time.Now().Format("2006-01-02 15:04:05"))
		t.Log("abcd", 1)
	})
	tw.AfterFunc(1*time.Second, func() {
		fmt.Println(time.Now().Format("2006-01-02 15:04:05"))
		t.Log("abcd-1", 1)
	})
	tw.AfterFunc(2*time.Second, func() {
		fmt.Println(time.Now().Format("2006-01-02 15:04:05"))
		t.Log("abcd", 2)
	})
	tw.AfterFunc(180*time.Second, func() {
		fmt.Println(time.Now().Format("2006-01-02 15:04:05"))
		t.Log("abcd", 180)
	})
	tw.AfterFunc(119*time.Second, func() {
		fmt.Println(time.Now().Format("2006-01-02 15:04:05"))
		t.Log("abcd", 119)
	})
	tw.AfterFunc(120*time.Second, func() {
		fmt.Println(time.Now().Format("2006-01-02 15:04:05"))
		t.Log("abcd", 120)
	})
	tw.AfterFunc(121*time.Second, func() {
		fmt.Println(time.Now().Format("2006-01-02 15:04:05"))
		t.Log("abcd", 121)
	})
	tw.AfterFunc(122*time.Second, func() {
		fmt.Println(time.Now().Format("2006-01-02 15:04:05"))
		t.Log("abcd", 122)
	})
	tw.AfterFunc(125*time.Second, func() {
		fmt.Println(time.Now().Format("2006-01-02 15:04:05"))
		t.Log("abcd", 125)
	})
	tw.AfterFunc(128*time.Second, func() {
		fmt.Println(time.Now().Format("2006-01-02 15:04:05"))
		t.Log("abcd", 128)
	})
	tw.AfterFunc(130*time.Second, func() {
		fmt.Println(time.Now().Format("2006-01-02 15:04:05"))
		t.Log("abcd", 130)
	})

	select {}
}
