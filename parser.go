package mp4parser

import (
	"fmt"
	"os"
	"strings"
)

func New(f *os.File) (Box, error) {
	box := &mp4Box{}
	if err := box.Parse(f); err != nil {
		return nil, err
	}

	return box, nil
}

func Fmt(b Box, data bool) string {
	return fmtBox(b, data, 0)
}

func fmtBox(b Box, data bool, offset int) string {
	str := fmt.Sprintf("%s%s (%v, %v)\n", strings.Repeat("-", offset*2), b.Type(), b.Offset(), b.Length())
	if data && b.Data() != nil {
		str += fmt.Sprintf("%s↳%s\n", strings.Repeat(" ", offset*2+5), b.Data())
	}

	for _, child := range b.Children() {
		str += fmtBox(child, data, offset+1)
	}

	return str
}
