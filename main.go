package main

import (
	"encoding/xml"
	"fmt"
	"net"
	"os"
	"sync"
	"strings"
	"regexp"
	"bytes"
	"path/filepath"
	"log"
	"strconv"
)

type Person struct {
	XMLName xml.Name `xml:"person"`
	Name    string   `xml:"name"`
	Age     int      `xml:"age"`
}

func (p Person) String() string {
	return fmt.Sprintf("{Name: %s, Age: %d}", p.Name, p.Age)
}

type People struct {
	XMLName xml.Name `xml:"people"`
	People  []Person `xml:"person"`
}

type FileSystem struct {
	lock             sync.Mutex
	DatabaseFilePath string
}

func (f FileSystem) ReadEntries() ([]Person, error) {

	file, _ := os.Open(f.DatabaseFilePath)

	var people People
	if err := xml.NewDecoder(file).Decode(&people); err != nil {
		handleError(err)
	}

	file.Close()

	return people.People, nil
}

func (f FileSystem)FlushDatabase(people People) error {

	os.Remove(f.DatabaseFilePath)
	file, _ := os.Create(f.DatabaseFilePath)

	fmt.Println("Locking writable database")
	f.lock.Lock()
	fmt.Println("Locked writable database")
	if err := xml.NewEncoder(file).Encode(people); err != nil {
		return err
	}
	fmt.Println("Unlocking writable database")
	f.lock.Unlock()
	fmt.Println("Unlocked writable database")

	file.Close()

	return nil
}

const (
	CONN_HOST = "localhost"
	CONN_PORT = "8888"
	CONN_TYPE = "tcp"
)

func main() {
	databaseFilePath, err := filepath.Abs("people.xml")
	handleError(err)

	fmt.Println("Init file system")
	fileSystem := FileSystem{DatabaseFilePath: databaseFilePath}
	fmt.Println("Reading saved entries")
	entries, _ := fileSystem.ReadEntries()
	p := People{People: entries}

	l, err := net.Listen(CONN_TYPE, CONN_HOST + ":" + CONN_PORT)
	handleError(err)
	// Close the listener when the application closes.
	defer l.Close()
	fmt.Println("Listening on " + CONN_HOST + ":" + CONN_PORT)
	for {
		// Listen for an incoming connection.
		conn, err := l.Accept()
		fmt.Println("accept")
		handleError(err)
		// Handle connections in a new goroutine.
		go handleConnection(conn, &fileSystem, &p)
	}
}

// Handles incoming requests.
func handleConnection(conn net.Conn, c *FileSystem, p*People) {
	// Make a buffer to hold incoming data.
	var buf [512]byte
	empty_string := ""
	for {
		copy(buf[:], empty_string) //make  buffer empty
		_, err := conn.Read(buf[0:])
		handleError(err)
		n := bytes.Index(buf[:], []byte{0})
		s := string(buf[:n])
		go handleRequest(s, c, conn, p)
	}
}

func handleError(e error) {
	if e != nil {
		log.Fatal(e)
		os.Exit(1)
	}
}

func (p *People) selectByKey(name string) []Person {
	selecting_data := []Person{}

	for _, item := range p.People {
		if item.Name == name {
			selecting_data = append(selecting_data, item)
		}
	}
	return selecting_data
}

func (p *People) deleteByKey(name string) {
	for index := 0; index < len(p.People); index++ {
		if p.People[index].Name == name {
			p.People = append((p.People)[:index], (p.People)[index + 1:]...)
		}
	}
}

func (p *People) updateByKey(name string, age int) {
	for index := 0; index < len(p.People); index++ {
		if p.People[index].Name == name {
			p.People[index].Age = age
		}
	}
}

func (p *People) addItem(name string, age int) {
	human := Person{Name: name, Age: age}
	p.People = append(p.People, human)
}

const (
	GET_REGEX = `^get\s+\w+`
	SET_REGEX = `^set\s+\w+\s+eq\s+-?\d+\s+.*`
	DELETE_REGEX = `^delete\s+\w+\s+`
	UPDATE_REGEX = `^update\s+\w+\s+set\s+-?\d+\s+`
)

func handleRequest(command string, fs *FileSystem, conn net.Conn, p *People) {
	fmt.Println("handling request: " + command)
	params := strings.Fields(command)
	// remove empty entries and remove whitespaces
	fmt.Println(params)
	fmt.Println(params[0])
	// remove empty entries and remove whitespaces
	switch strings.ToLower(params[0]) {
	case "exit":
		os.Exit(0)
	case "set":
		fmt.Println("Got 'set' request:")
		insert_regex := regexp.MustCompile(SET_REGEX)
		if insert_regex.MatchString(command) {
			age, err := strconv.Atoi(params[3])
			handleError(err)
			p.addItem(params[1], age)
			conn.Write([]byte("OK\n"))
			go fs.FlushDatabase(*p)
		} else {
			conn.Write([]byte(string("Your command didn't match the pattern\n")))
		}
	case "get":
		fmt.Println("Got 'get' request:")
		select_regex := regexp.MustCompile(GET_REGEX)
		if select_regex.MatchString(command) {
			for _, elem := range p.selectByKey(params[1]) {
				//data += elem + " "
				conn.Write([]byte(elem.String() + " "))
			}
			conn.Write([]byte("\n"))
		} else {
			conn.Write([]byte(string("Your command didn't match the pattern\n")))
		}
	case "delete":
		fmt.Println("Got 'delete' request:")
		delete_regex := regexp.MustCompile(DELETE_REGEX)
		if delete_regex.MatchString(command) {
			p.deleteByKey(params[1])
			conn.Write([]byte("OK\n"))
			go fs.FlushDatabase(*p)
		} else {
			conn.Write([]byte(string("Your command didn't match the pattern\n")))
		}

	case "update":
		fmt.Println("Got 'update' request:")
		update_regex := regexp.MustCompile(UPDATE_REGEX)
		if update_regex.MatchString(command) {
			age, err := strconv.Atoi(params[3])
			handleError(err)
			p.updateByKey(params[1], age)
			conn.Write([]byte("OK\n"))
			go fs.FlushDatabase(*p)
		} else {
			conn.Write([]byte(string("Your command didn't match the pattern\n")))
		}

	case "help":
		fmt.Println("Got 'help' request:")
		conn.Write([]byte(string("set [name] eq [age]\n")))
		conn.Write([]byte(string("update [name] set [age]\n")))
		conn.Write([]byte(string("get [name]\n")))
		conn.Write([]byte(string("delete [name]\n")))
	default:
		fmt.Println("Got wront request:")
		conn.Write([]byte(string("Your command didn't match the pattern\n")))
	}
}