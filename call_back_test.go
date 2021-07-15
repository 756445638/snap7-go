package snap7go

import (
	"fmt"
	"testing"
	"time"
)

func TestSomeCallBack(t *testing.T) {
	server := NewS7Server()
	server.SetEventsCallback(justPrintEvent)
	server.SetReadEventsCallback(justPrintEvent)

	var data [1024]byte
	go func() {
		for ; ; time.Sleep(time.Millisecond * 50) {
			for k, _ := range data {
				data[k]++
			}
		}
	}()
	err := server.RegisterArea(SrvAreaPA, 0, data[:])
	if err != nil {
		t.Fatal(err)
	}
	err = server.Start()
	if err != nil {
		t.Fatal(err)
	}
	defer func() {
		err = server.Stop()
		if err != nil {
			t.Fatal(err)
			return
		}
		server.Destroy()
	}()
	client := NewS7Client()
	err = client.ConnectTo("127.0.0.1", 0, 1)
	if err != nil {
		t.Fatal(err)
		return
	}
	for i := 0; i < 10; func() {
		i++
		time.Sleep(time.Millisecond * 50)
	}() {
		data, err := client.ReadArea(S7AreaPA, 0, 1, 10, S7WLWord)
		if err != nil {
			t.Fatal(err)
			return
		}
		fmt.Println("read data:", data)
	}
}
