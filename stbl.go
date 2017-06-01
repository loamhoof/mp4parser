package mp4parser

import (
	"io"
)

type StblBox struct {
	baseBox
	children []Box
	Stts     *SttsBox
	Stsc     *StscBox
	Stsz     *StszBox
	Stco     *StcoBox
	Stsd     *StsdBox
}

func (b *StblBox) Parse(r io.ReadSeeker, startOffset int64, pp ParsePlan, pc ParseContext) error {
	size, _, fields, children, err := parseContainerBox(r, startOffset, pp, pc)
	if err != nil {
		return err
	}
	b.size = size
	b.fields = fields
	b.children = children

	for _, child := range children {
		switch child := child.(type) {
		case *SttsBox:
			b.Stts = child
		case *StscBox:
			b.Stsc = child
		case *StszBox:
			b.Stsz = child
		case *StcoBox:
			b.Stco = child
		case *StsdBox:
			b.Stsd = child
		default:
		}
	}

	return nil
}

func (b *StblBox) Type() string {
	return "stbl"
}

func (b *StblBox) Children() []Box {
	return b.children
}
