package main

import "testing"
import "net"
import "fmt"
import "bufio"
import "time"
import "strings"

const (
	MESSAGE_OK = "OK"
	MESSAGE_EMPTY = ""
)

func Sleep() {
	time.Sleep(100 * time.Millisecond)
}

func Test(t *testing.T) {

	go main()
	Sleep()
	fmt.Println("Starting tests")
	conn, _ := net.Dial("tcp", "127.0.0.1:8888")

	// Test insert START
	fmt.Println("Testing insert method")
	fmt.Fprint(conn, "set Izya eq 19\n")
	Sleep()

	message, _ := bufio.NewReader(conn).ReadString('\n')
	Sleep()

	fmt.Println("Response:" + message)
	message = strings.Replace(strings.Replace(message, " ", "", -1), "\n", "", -1)
	fmt.Println("Cleaned response:" + message)

	if message != MESSAGE_OK {
		t.Error("message = " + message)
		t.Error("expected = " + MESSAGE_OK)
		t.Error("Insert not passed")
	}
	// Test insert END

	// Test get START
	fmt.Fprint(conn, "get Izya")
	fmt.Println("listen")
	message, _ = bufio.NewReader(conn).ReadString('\n')
	Sleep()

	fmt.Println("Response:" + message)
	message = strings.Replace(strings.Replace(message, " ", "", -1), "\n", "", -1)
	fmt.Println("Cleaned response:" + message)

	expected := strings.Replace(Person{Name:"Izya", Age: 19}.String(), " ", "", -1)
	if message != expected {
		t.Error("message = " + message)
		t.Error("expected = " + expected)
		t.Error("Select not passed")
	}
	Sleep()
	// Test get END

	// Test update START
	fmt.Fprint(conn, "update Izya set 23\n")
	Sleep()

	message, _ = bufio.NewReader(conn).ReadString('\n')
	Sleep()

	fmt.Println("Response:" + message)
	message = strings.Replace(strings.Replace(message, " ", "", -1), "\n", "", -1)
	fmt.Println("Cleaned response:" + message)

	if message != MESSAGE_OK {
		t.Error("message = " + message)
		t.Error("expected = " + MESSAGE_OK)
		t.Error("Update not passed")
	}
	// Test update END

	// Test get START
	fmt.Fprint(conn, "get Izya\n")
	Sleep()

	message, _ = bufio.NewReader(conn).ReadString('\n')
	Sleep()

	fmt.Println("Response:" + message)
	message = strings.Replace(strings.Replace(message, " ", "", -1), "\n", "", -1)
	fmt.Println("Cleaned response:" + message)

	expected = strings.Replace(Person{Name:"Izya", Age: 23}.String(), " ", "", -1)
	if message != expected {
		t.Error("message = " + message)
		t.Error("expected = " + expected)
		t.Error("Get after update not passed")
	}
	Sleep()
	// Test get END

	// Test delete START
	fmt.Fprint(conn, "delete Izya")
	Sleep()

	message, _ = bufio.NewReader(conn).ReadString('\n')
	message = strings.Replace(strings.Replace(message, " ", "", -1), "\n", "", -1)
	if message != MESSAGE_OK {
		t.Error("message = " + message)
		t.Error("expected = " + MESSAGE_OK)
		t.Error("Detele not passed")
	}
	// Test delete END

	// Test get after delete START
	fmt.Fprint(conn, "get Izya")

	message, _ = bufio.NewReader(conn).ReadString('\n')
	message = strings.Replace(strings.Replace(message, " ", "", -1), "\n", "", -1)
	if message != MESSAGE_EMPTY {
		t.Error("message = " + message)
		t.Error("expected = " + MESSAGE_EMPTY)
		t.Error("Detele not passed")
	}
	// Test get after delete END


	fmt.Fprint(conn, "exit")

	fmt.Println("Test passed successfully")
}
