package mp4parser

import (
	"encoding/binary"
	"fmt"
	"io"
	"strings"
)

type MP4 struct {
	children []Box
	Ftyp     *FtypBox
	Free     *FreeBox
	Moov     *MoovBox
	Moof     []*MoofBox
	Mdat     []*MdatBox
	Mfra     *MfraBox
}

func (m *MP4) Parse(r io.ReadSeeker, offset int64, pp ParsePlan, pc ParseContext) error {
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

		if _, ok := pp[boxType]; ok || pp == nil {
			box := newBox(boxType)
			if err := box.Parse(r, offset, pp[boxType], pc); err != nil {
				return err
			}
			children = append(children, box)

			switch box := box.(type) {
			case *FtypBox:
				m.Ftyp = box
			case *FreeBox:
				m.Free = box
			case *MoovBox:
				m.Moov = box
			case *MoofBox:
				m.Moof = append(m.Moof, box)
			case *MdatBox:
				m.Mdat = append(m.Mdat, box)
			case *MfraBox:
				m.Mfra = box
			default:
			}
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
