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
	text := "insert 222 9 into test_table\n"
	// send to socket
	fmt.Fprintf(conn, text)
	fmt.Println("send select")
	time.Sleep(100 * time.Millisecond)

	text = "select 222 from test_table\n"
	// send to socket
	fmt.Fprintf(conn, text)
	fmt.Println("listen")
	message, _ := bufio.NewReader(conn).ReadString('\n')
	time.Sleep(100 * time.Millisecond)

	fmt.Print("Message from server:" + message + "end")
	message = strings.Replace(strings.Replace(message, " ", "", -1), "\n", "", -1)
	fmt.Print("Message from server:" + message + "end")

	if message != "9" {
		t.Error("select return wrong string")

	}
	time.Sleep(100 * time.Millisecond)

	fmt.Fprintf(conn, "update 222 to 10 in test_table\n")
	time.Sleep(100 * time.Millisecond)
	fmt.Fprintf(conn, "select 222 from test_table\n")

	message, _ = bufio.NewReader(conn).ReadString('\n')
	message = strings.Replace(strings.Replace(message, " ", "", -1), "\n", "", -1)
	if message != "10" {
		t.Error("update works wrong")

	}
	time.Sleep(100 * time.Millisecond)

	fmt.Fprintf(conn, "delete 222 in test_table\n")
	time.Sleep(100 * time.Millisecond)
	fmt.Fprintf(conn, "select 222 from test_table\n")

	message, _ = bufio.NewReader(conn).ReadString('\n')
	message = strings.Replace(strings.Replace(message, " ", "", -1), "\n", "", -1)
	if message != "" {
		t.Error("delete works wrong")

	}

	fmt.Fprintf(conn, "exit 1")

	fmt.Println("OK")
}
