package mp4parser

import (
	"io"
)

type trakBox struct {
	baseBox
	children []Box
	Tkhd     *tkhdBox
	Mdia     *mdiaBox
}

func (b *trakBox) Parse(r io.ReadSeeker, startOffset int64) error {
	size, _, fields, children, err := parseContainerBox(r, startOffset)
	if err != nil {
		return err
	}
	b.size = size
	b.fields = fields
	b.children = children

	for _, child := range children {
		switch child := child.(type) {
		case *tkhdBox:
			b.Tkhd = child
		case *mdiaBox:
			b.Mdia = child
		default:
		}
	}

	return nil
}

func (b *trakBox) Type() string {
	return "trak"
}

func (b *trakBox) Children() []Box {
	return b.children
}
