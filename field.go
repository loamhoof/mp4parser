package mp4parser

import (
	"fmt"
	"strconv"
)

type Field struct {
	Key        string
	Value      interface{}
	Offset     int64
	Bits       uint64
	BitsOffset uint8
}

func (f *Field) String() string {
	quotedValue := strconv.Quote(fmt.Sprintf("%v", f.Value))

	return fmt.Sprintf("%s: %v", f.Key, quotedValue[1:len(quotedValue)-1])
}

type Fields []*Field

func (fields Fields) Get(key string) *Field {
	for _, f := range fields {
		if f.Key == key {
			return f
		}
	}

	return nil
}
