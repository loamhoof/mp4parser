package mp4parser

import (
	"io"
)

type udtaBox struct {
	baseBox
	children []Box
	Meta     *metaBox
}

func (b *udtaBox) Parse(r io.ReadSeeker, startOffset int64) error {
	size, _, fields, children, err := parseContainerBox(r, startOffset)
	if err != nil {
		return err
	}
	b.size = size
	b.fields = fields
	b.children = children

	for _, child := range children {
		switch child := child.(type) {
		case *metaBox:
			b.Meta = child
		default:
		}
	}

	return nil
}

func (b *udtaBox) Type() string {
	return "udta"
}

func (b *udtaBox) Children() []Box {
	return b.children
}
