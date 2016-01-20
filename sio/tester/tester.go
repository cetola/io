package main

import (
	"flag"
	"fmt"
	"log"
	"math/rand"
	"os"
	"time"

	"github.com/jmore-reachtech/serial"
)

const dur int = 10

func random(min, max int) int {
	rand.Seed(time.Now().UnixNano())
	return rand.Intn(max-min) + min
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

	s, err := serial.OpenPort(c)
	if err != nil {
		log.Fatal("fatal error: ", err)
	}
	defer s.Close()

	for j := 0; j < 100; j++ {

		for i := 0; i < 18; i++ {
			w := fmt.Sprintf("w=%d\n", random(0, i+1))
			f := fmt.Sprintf("f=%d\n", random(0, i+1))
			t := fmt.Sprintf("t=%d\n", random(0, i+1))
			p := fmt.Sprintf("s=%d\n", random(0, (i+1)*10))
			o := fmt.Sprintf("o=%d\n", i%2)

			_, err := s.Write([]byte(w))
			checkError(err)
			time.Sleep((time.Duration(dur) * time.Millisecond))
			_, err = s.Write([]byte(f))
			checkError(err)
			time.Sleep((time.Duration(dur) * time.Millisecond))
			_, err = s.Write([]byte(t))
			checkError(err)
			time.Sleep((time.Duration(dur) * time.Millisecond))
			_, err = s.Write([]byte(o))
			checkError(err)
			time.Sleep((time.Duration(dur) * time.Millisecond))
			_, err = s.Write([]byte(p))
			checkError(err)

		}

		for i := 18; i > 0; i-- {
			w := fmt.Sprintf("w=%d\n", random(0, i+1))
			f := fmt.Sprintf("f=%d\n", random(0, i+1))
			t := fmt.Sprintf("t=%d\n", random(0, i+1))
			p := fmt.Sprintf("s=%d\n", random(0, (i+1)*10))
			o := fmt.Sprintf("o=%d\n", i%2)

			_, err := s.Write([]byte(w))
			checkError(err)
			time.Sleep((time.Duration(dur) * time.Millisecond))
			_, err = s.Write([]byte(f))
			checkError(err)
			time.Sleep((time.Duration(dur) * time.Millisecond))
			_, err = s.Write([]byte(t))
			checkError(err)
			time.Sleep((time.Duration(dur) * time.Millisecond))
			_, err = s.Write([]byte(o))
			checkError(err)
			time.Sleep((time.Duration(dur) * time.Millisecond))
			_, err = s.Write([]byte(p))
			checkError(err)
			time.Sleep((time.Duration(dur) * time.Millisecond))
		}
	}
}

// If err is non-nil, print it out and halt.
func checkError(err error) {
	if err != nil {
		fmt.Fprintf(os.Stderr, "Error: %s\n", err.Error())
		os.Exit(1)
	}
}
