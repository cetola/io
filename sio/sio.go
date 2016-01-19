package main

import (
	"bufio"
	"flag"
	"fmt"
	"github.com/jmore-reachtech/io/tio"
	"github.com/jmore-reachtech/serial"
	"io"
	"log"
	"net"
	"os"
	"os/signal"
	"strings"
	"syscall"
	"time"
)

// channel timeout used with select
const timeout = (5 * time.Millisecond)

// message map for gui -> micro
var mapGui = tio.NewTio("gui")

// message map for micro -> gui
var mapMicro = tio.NewTio("micro")

// serial read, we run this in a goroutine so we can block on read
// on data read send through the channel
func serialRead(p *serial.Port, ch chan string) {
	for {
		buf := make([]byte, 128)
		n, err := p.Read(buf)
		if err != nil && err != io.EOF {
			log.Fatal(err)
		}
		if n > 0 {
			ch <- fmt.Sprint(string(buf[:n]))
		}
	}
}

// handle the serial port communication. When data comes in through
// the serial port write it to the channel for the socket to see. When
// data comes in from the socket over the channel translate the message
// and write to serial port
func handlePort(p *serial.Port, ch chan string) {
	readCh := make(chan string)

	go serialRead(p, readCh)

	for {
		select {
		case s := <-ch:
			{
				// map gui -> micro
				trans := mapGui.ItemTranslate(s)
				_, err := p.Write([]byte(trans))
				if err != nil {
					fmt.Println(err)
				}
			}
		case r := <-readCh:
			{
				ch <- fmt.Sprint(r)
			}
		case <-time.After(timeout):
			continue
		}
	}
}

// listen for gui client to connect to our socket
func accept(listener *net.UnixListener, ch chan string) {
	for {
		// we are going to eat the serial data until
		// we get a socket connection so we don't block the channel
		select {
		case <-ch:
			log.Println("eating serial data")
		default:
		}

		// set timeout to fall through so we can check the channel for
		// serial data
		listener.SetDeadline(time.Now().Add(100 * time.Millisecond))
		conn, err := listener.AcceptUnix()
		if nil != err {
			if opErr, ok := err.(*net.OpError); ok && opErr.Timeout() {
				continue
			}
			log.Println(err)
		} else {
			// we have connection, call handle, we only handle one connection
			// so no goroutine here
			handleSocket(conn, ch)
		}
	}
}

// socket read, we run this in a goroutine so we can block on read
// on data read send through the channel
func socketRead(conn *net.UnixConn, ch chan string) {
	for {
		buf := make([]byte, 512)
		n, err := conn.Read(buf)
		if nil != err {
			log.Fatal(err)
		}

		if n > 0 {
			ch <- fmt.Sprint(string(buf[:n]))
		}
	}
}

// handle the socket connection. When data comes in on the socket write it
// to the channel so the serial port can see it. When data comes in over the
// channel translate the message and write it to the socket
func handleSocket(conn *net.UnixConn, ch chan string) {
	defer conn.Close()

	readCh := make(chan string)

	go socketRead(conn, readCh)

	for {
		select {
		case s := <-ch:
			{
				// map micro -> gui
				trans := mapMicro.ItemTranslate(s)
				_, err := conn.Write([]byte(trans))
				if err != nil {
					fmt.Println(err)
				}
			}
		case r := <-readCh:
			{
				ch <- fmt.Sprint(r)
			}
		case <-time.After(timeout):
			continue
		}
	}
}

// create the translate map using the translate.txt file.
func initMapping() {
	f, err := os.Open("translate.txt")
	defer f.Close()

	// if there is no translate.txt file die
	if err != nil {
		fmt.Println("Error opening translate.txt")
		os.Exit(1)
	}

	scanner := bufio.NewScanner(f)

	// read each line and parse out the mapping
	for scanner.Scan() {
		var line = scanner.Text()
		if len(line) == 0 {
			break
		}
		if string(line[0]) == "#" || string(line[0]) == "/" {
			continue
		}

		// the mapping is split by a ,
		var trans = strings.SplitN(line, ",", 2)

		// look for = in message
		eq := strings.Index(trans[0], "=")
		var key string
		if eq == -1 {
			key = trans[0][2:]
		} else {
			key = trans[0][2:eq]
		}

		// look for = in trans
		eq = strings.Index(trans[1], "=")
		var value string
		if eq == -1 {
			value = trans[1][2:]
		} else {
			value = trans[1][2:eq]
		}

		switch string(trans[0][0]) {
		case "G":
			mapGui.ItemAdd(key, value)
		case "M":
			mapMicro.ItemAdd(key, value)
		default:
		}
	}

}

// init is called on start, here we create the translate mapping
func init() {
	initMapping()
}

func main() {
	ttyPtr := flag.String("tty", "/dev/ttymxc1", "Open tty device")
	baudPtr := flag.Int("baud", 115200, "Set baud rate")
	flag.Parse()

	c := &serial.Config{
		Name:        *ttyPtr,
		Baud:        *baudPtr,
		ReadTimeout: time.Millisecond * 500,
	}

	// this is the channel that moves data from the serial and socket
	// goroutines
	ch := make(chan string)

	s, err := serial.OpenPort(c)
	if err != nil {
		log.Fatal(err)
	}

	_, err = s.Write([]byte("Hello!"))
	if err != nil {
		log.Fatal(err)
	}

	go handlePort(s, ch)

	laddr, err := net.ResolveUnixAddr("unix", "/tmp/tioSocket")
	if nil != err {
		log.Fatalln(err)
	}

	listener, err := net.ListenUnix("unix", laddr)
	if nil != err {
		log.Fatalln(err)
	}
	log.Println("listening on", listener.Addr())
	go accept(listener, ch)

	// Handle SIGINT and SIGTERM.
	ex := make(chan os.Signal)
	signal.Notify(ex, syscall.SIGINT, syscall.SIGTERM)
	log.Println(<-ex)

	log.Println("Closing up...")
	s.Close()
}
