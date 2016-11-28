package mp4parser

import (
	"io"
)

type dinfBox struct {
	baseBox
	children []Box
	Dref     *drefBox
}

func (b *dinfBox) Parse(r io.ReadSeeker, startOffset int64) error {
	size, _, fields, children, err := parseContainerBox(r, startOffset)
	if err != nil {
		return err
	}
	b.size = size
	b.fields = fields
	b.children = children

	for _, child := range children {
		switch child := child.(type) {
		case *drefBox:
			b.Dref = child
		default:
		}
	}

	return nil
}

func (b *dinfBox) Type() string {
	return "dinf"
}

func (b *dinfBox) Children() []Box {
	return b.children
}
