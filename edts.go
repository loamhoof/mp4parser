package mp4parser

import (
	"io"
)

type edtsBox struct {
	baseBox
	children []Box
	Elst     *elstBox
}

func (b *edtsBox) Parse(r io.ReadSeeker, startOffset int64) error {
	size, _, fields, children, err := parseContainerBox(r, startOffset)
	if err != nil {
		return err
	}
	b.size = size
	b.fields = fields
	b.children = children

	for _, child := range children {
		switch child := child.(type) {
		case *elstBox:
			b.Elst = child
		default:
		}
	}

	return nil
}

func (b *edtsBox) Type() string {
	return "edts"
}

func (b *edtsBox) Children() []Box {
	return b.children
}
