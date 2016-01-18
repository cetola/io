package tio

import (
	"fmt"
	"strings"
)

type Tio struct {
	id      string
	mapping map[string]string
}

func NewTio(id string) *Tio {
	tio := new(Tio)
	tio.id = id
	tio.mapping = make(map[string]string)

	return tio
}

func (t *Tio) ItemAdd(key string, value string) {
	fmt.Printf("adding item: '%s':'%s'\n", key, value)

	t.mapping[key] = value
}

func (t *Tio) ItemFind(key string) (string, bool) {
	fmt.Println("lookging up ", key)
	v, exists := t.mapping[key]
	return v, exists
}

func (t *Tio) ItemRemove(key string) {
	delete(t.mapping, key)
}

func (t *Tio) ItemTranslate(msg string) string {
	var (
		value  string
		exists bool
	)

	eq := strings.Index(msg, "=")

	if eq == -1 {
		value, exists = t.ItemFind(msg)
		if exists {
			return value
		} else {
			return msg
		}
	} else {
		value, exists = t.ItemFind(msg[:eq])
		if exists {
			return value + msg[eq:]
		} else {
			return msg
		}

	}

}
