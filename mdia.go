package mp4parser

import (
	"io"
)

type mdiaBox struct {
	baseBox
	children []Box
	Mdhd     *mdhdBox
	Hdlr     *hdlrBox
	Minf     *minfBox
}

func (b *mdiaBox) Parse(r io.ReadSeeker, startOffset int64) error {
	size, _, fields, children, err := parseContainerBox(r, startOffset)
	if err != nil {
		return err
	}
	b.size = size
	b.fields = fields
	b.children = children

	for _, child := range children {
		switch child := child.(type) {
		case *mdhdBox:
			b.Mdhd = child
		case *hdlrBox:
			b.Hdlr = child
		case *minfBox:
			b.Minf = child
		default:
		}
	}

	return nil
}

func (b *mdiaBox) Type() string {
	return "mdia"
}

func (b *mdiaBox) Children() []Box {
	return b.children
}
