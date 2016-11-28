package mp4parser

import (
	"io"
)

type mvexBox struct {
	baseBox
	children []Box
	Mehd     *mehdBox
	Trex     *trexBox
}

func (b *mvexBox) Parse(r io.ReadSeeker, startOffset int64) error {
	size, _, fields, children, err := parseContainerBox(r, startOffset)
	if err != nil {
		return err
	}
	b.size = size
	b.fields = fields
	b.children = children

	for _, child := range children {
		switch child := child.(type) {
		case *mehdBox:
			b.Mehd = child
		case *trexBox:
			b.Trex = child
		default:
		}
	}

	return nil
}

func (b *mvexBox) Type() string {
	return "mvex"
}

func (b *mvexBox) Children() []Box {
	return b.children
}
