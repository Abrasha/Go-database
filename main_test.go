package main

import "testing"
import "net"
import "fmt"
import "bufio"
import "time"
import "strings"

func Test(t *testing.T) {

	go main()
	time.Sleep(100 * time.Millisecond)
	fmt.Println("test starting")
	conn, _ := net.Dial("tcp", "127.0.0.1:8888")
	fmt.Println("send insert")
	// read in input from stdin
	text := "set Izya eq 19"
	// send to socket
	fmt.Fprintf(conn, text)
	fmt.Println("send select")
	time.Sleep(100 * time.Millisecond)

	text = "get Izya"

	// send to socket
	fmt.Fprintf(conn, text)
	fmt.Println("listen")
	message, _ := bufio.NewReader(conn).ReadString('\n')
	time.Sleep(100 * time.Millisecond)

	fmt.Print("Message from server:" + message + "end")
	message = strings.Replace(strings.Replace(message, " ", "", -1), "\n", "", -1)
	fmt.Print("Message from server:" + message + "end")

	if message != "19" {
		t.Error("select return wrong string")

	}
	time.Sleep(100 * time.Millisecond)

	fmt.Fprintf(conn, "update Izya set 23")
	time.Sleep(100 * time.Millisecond)
	fmt.Fprintf(conn, "get Izya")

	message, _ = bufio.NewReader(conn).ReadString('\n')
	message = strings.Replace(strings.Replace(message, " ", "", -1), "\n", "", -1)
	if message != "23" {
		t.Error("update works wrong")

	}
	time.Sleep(100 * time.Millisecond)

	fmt.Fprintf(conn, "delete Izya")
	time.Sleep(100 * time.Millisecond)
	fmt.Fprintf(conn, "get Izya")

	message, _ = bufio.NewReader(conn).ReadString('\n')
	message = strings.Replace(strings.Replace(message, " ", "", -1), "\n", "", -1)
	if message != "" {
		t.Error("delete works wrong")

	}

	fmt.Fprintf(conn, "exit")

	fmt.Println("OK")
}
