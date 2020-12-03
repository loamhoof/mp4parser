package mp4parser

import (
	"encoding/binary"
	"fmt"
	"io"
	"strconv"
)

type UdtaBox struct {
	baseBox
	children []Box
	Meta     *MetaBox
}

func (b *UdtaBox) Parse(r io.ReadSeeker, startOffset int64, pp ParsePlan, pc ParseContext) error {
	size, _, fields, children, err := parseContainerBoxWithFactory(r, startOffset, pp, pc, newUdtaBox)
	if err != nil {
		return err
	}
	b.size = size
	b.fields = fields
	b.children = children

	for _, child := range children {
		switch child := child.(type) {
		case *MetaBox:
			b.Meta = child
		default:
		}
	}

	return nil
}

func (b *UdtaBox) Type() string {
	return "udta"
}

func (b *UdtaBox) Children() []Box {
	return b.children
}

func newUdtaBox(_type string) Box {
	box := newBox(_type)

	switch box.(type) {
	default:
		return box
	case *UnknownBox:
		return &GenericUdtaBox{_type: _type}
	}
}

type GenericUdtaBox struct {
	baseBox
	_type string
}

func (b *GenericUdtaBox) Parse(r io.ReadSeeker, startOffset int64, pp ParsePlan, pc ParseContext) error {
	size, offset, _, fields, err := parseBox(r, startOffset)
	if err != nil {
		return err
	}
	b.size = size
	b.fields = fields

	// if not international text
	if b._type[0] != '\xa9' {
		lData := size - uint64(offset-startOffset)
		bData := make([]byte, lData)
		if _, err := r.Read(bData); err != nil {
			return err
		}
		b.fields = append(b.fields, &Field{"unknown", string(bData), offset, lData * 8, 0})

		return nil
	}

	for size-uint64(offset-startOffset) > 0 {
		b2 := make([]byte, 2)

		if _, err := r.Read(b2); err != nil {
			return err
		}
		textSize := binary.BigEndian.Uint16(b2)
		b.fields = append(b.fields, &Field{"text_size", textSize, offset, 16, 0})
		offset += 2

		if _, err := r.Read(b2); err != nil {
			return err
		}
		n := binary.BigEndian.Uint16(b2)

		b.fields = append(b.fields, &Field{"pad", n >> 15, offset, 1, 0})
		language := fmt.Sprintf("%c%c%c", n>>10&0x1F+0x60, n>>5&0x1F+0x60, n&0x1F+0x60)
		b.fields = append(b.fields, &Field{"language", language, offset, 15, 1})

		offset += 2

		bData := make([]byte, textSize)
		if _, err := r.Read(bData); err != nil {
			return err
		}
		// the encoding of the text is actually a bit more complicated
		// https://developer.apple.com/library/archive/documentation/QuickTime/QTFF/QTFFChap2/qtff2.html#//apple_ref/doc/uid/TP40000939-CH204-SW1
		b.fields = append(b.fields, &Field{"text", string(bData), offset, uint64(textSize) * 8, 0})

		offset += int64(textSize)
	}

	return nil
}

func (b *GenericUdtaBox) Type() string {
	quoted := strconv.Quote(b._type)

	return quoted[1 : len(quoted)-1]
}
