package mp4parser

import (
	"encoding/binary"
	"fmt"
	"io"
	"strings"
)

type MP4 struct {
	children []Box
	Ftyp     Box
	Free     Box
	Moov     Box
	Moof     []Box
	Mdat     []Box
	Mfra     Box
}

func (m *MP4) Parse(r io.ReadSeeker, offset int64) error {
	bytes := make([]byte, 4)

	children := make([]Box, 0, 1)
	for {
		if _, err := r.Seek(offset, io.SeekStart); err != nil {
			break
		}

		if _, err := r.Read(bytes); err != nil {
			if err == io.EOF {
				break
			}
			return err
		}
		l := binary.BigEndian.Uint32(bytes)

		if _, err := r.Read(bytes); err != nil {
			return err
		}
		boxType := string(bytes)

		box := newBox(boxType)
		if err := box.Parse(r, offset); err != nil {
			return err
		}
		children = append(children, box)

		switch box.Type() {
		case "ftyp":
			m.Ftyp = box
		case "free":
			m.Free = box
		case "moov":
			m.Moov = box
		case "moof":
			m.Moof = append(m.Moof, box)
		case "mdat":
			m.Mdat = append(m.Mdat, box)
		case "mfra":
			m.Mfra = box
		}

		offset += int64(l)
	}

	m.children = children

	return nil
}

func (m *MP4) Type() string {
	return "mp4"
}

func (m *MP4) Offset() int64 {
	return 0
}

func (m *MP4) Size() uint64 {
	var size uint64 = 0
	for _, b := range m.Children() {
		size += b.Size()
	}

	return size
}

func (m *MP4) Children() []Box {
	return m.children
}

func (m *MP4) Data() Fields {
	return nil
}

func (m *MP4) String() string {
	return fmtChildren(m, 0)
}

func fmtChildren(b Box, offset int) string {
	str := ""

	for _, child := range b.Children() {
		str += fmtBox(child, offset)
	}

	return str
}

func fmtBox(b Box, offset int) string {
	str := fmt.Sprintf("%s%s (%v, %v)\n", strings.Repeat("-", offset*2), b.Type(), b.Offset(), b.Size())

	return str + fmtChildren(b, offset+1)
}
