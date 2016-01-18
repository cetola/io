package tio_test

import (
	"bufio"
	"fmt"
	"github.com/everlook/io/tio"
	"os"
	"strings"
	"testing"
)

const testKey = "key"
const testValue = "value"

// test we can add a new item
func TestAdd(t *testing.T) {
	tio := tio.NewTio("test")
	tio.ItemAdd(testKey, testValue)
}

// test we can find an item
func TestFind(t *testing.T) {
	tio := tio.NewTio("test")
	tio.ItemAdd(testKey, testValue)

	v, _ := tio.ItemFind(testKey)
	if v != testValue {
		t.Error("Expected '%s', got %s", testValue, v)
	}
}

// test we can remove an item
func TestRemove(t *testing.T) {
	tio := tio.NewTio("test")
	tio.ItemAdd(testKey, testValue)

	v, _ := tio.ItemFind(testKey)
	if v != testValue {
		t.Error("Expected '%s', got %s", testValue, v)
	}

	tio.ItemRemove(testKey)
	_, exists := tio.ItemFind(testKey)
	if exists != false {
		t.Error("Expected false, got true")
	}
}

// test adding items form a text file
func TestAddFromFile(t *testing.T) {
	f, err := os.Open("./sample/sample.txt")
	defer f.Close()

	if err != nil {
		t.Error("Error opening sample.txt")
	}

	gui := tio.NewTio("gui")
	micro := tio.NewTio("micro")

	scanner := bufio.NewScanner(f)

	for scanner.Scan() {
		var line = scanner.Text()
		if string(line[0]) == "#" || string(line[0]) == "/" {
			fmt.Println("skipping comment")
			continue
		}

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
			gui.ItemAdd(key, value)
		case "M":
			micro.ItemAdd(key, value)
		default:
		}
	}

	mapping := gui.ItemTranslate("meter.value=5")
	if mapping != "x=5" {
		t.Error("expected x=5 got", mapping)
	}

	mapping = gui.ItemTranslate("notfound")
	if mapping != "notfound" {
		t.Error("expected notfound got", mapping)
	}

}
