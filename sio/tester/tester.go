package main

import (
	"flag"
	"fmt"
	"github.com/jmore-reachtech/serial"
	"log"
	"os"
	"time"
)

func main() {
	ttyPtr := flag.String("tty", "/dev/ttymxc1", "Open tty device")
	baudPtr := flag.Int("baud", 115200, "Set baud rate")
	flag.Parse()

	c := &serial.Config{
		Name:        *ttyPtr,
		Baud:        *baudPtr,
		ReadTimeout: time.Millisecond * 500,
	}

	s, err := serial.OpenPort(c)
	if err != nil {
		log.Fatal("fatal error: ", err)
	}

	defer s.Close()

	for j := 0; j < 10; j++ {

		for i := 0; i < 12; i++ {
			msg := fmt.Sprintf("sl=%d\n", i*10)
			fmt.Println("sending: ", msg)

			n, err := s.Write([]byte(msg))
			if nil != err {
				fmt.Println(err)
				os.Exit(1)
			}
			fmt.Printf("wrote %d bytes \n", n)
			time.Sleep(50 * time.Millisecond)
		}

		for i := 12; i > 0; i-- {
			msg := fmt.Sprintf("sl=%d\n", i*10)
			fmt.Println("sending: ", msg)

			n, err := s.Write([]byte(msg))
			if nil != err {
				fmt.Println(err)
				os.Exit(1)
			}
			fmt.Printf("wrote %d bytes \n", n)
			time.Sleep(50 * time.Millisecond)
		}
	}
}
