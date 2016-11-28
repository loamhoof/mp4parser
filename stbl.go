package mp4parser

import (
	"io"
)

type stblBox struct {
	baseBox
	children []Box
	Stts     *sttsBox
	Stsc     *stscBox
	Stsz     *stszBox
	Stco     *stcoBox
}

func (b *stblBox) Parse(r io.ReadSeeker, startOffset int64) error {
	size, _, fields, children, err := parseContainerBox(r, startOffset)
	if err != nil {
		return err
	}
	b.size = size
	b.fields = fields
	b.children = children

	for _, child := range children {
		switch child := child.(type) {
		case *sttsBox:
			b.Stts = child
		case *stscBox:
			b.Stsc = child
		case *stszBox:
			b.Stsz = child
		case *stcoBox:
			b.Stco = child
		default:
		}
	}

	return nil
}

func (b *stblBox) Type() string {
	return "stbl"
}

func (b *stblBox) Children() []Box {
	return b.children
}
