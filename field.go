package mp4parser

import (
	"fmt"
)

type Field struct {
	Key    string
	Value  interface{}
	Offset int64
	Bits   uint64
}

func (f *Field) String() string {
	return fmt.Sprintf("%s: %v", f.Key, f.Value)
}

type Fields []*Field
