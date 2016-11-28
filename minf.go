package mp4parser

import (
	"io"
)

type minfBox struct {
	baseBox
	children []Box
	Vmhd     *vmhdBox
	Dinf     *dinfBox
}

func (b *minfBox) Parse(r io.ReadSeeker, startOffset int64) error {
	size, _, fields, children, err := parseContainerBox(r, startOffset)
	if err != nil {
		return err
	}
	b.size = size
	b.fields = fields
	b.children = children

	for _, child := range children {
		switch child := child.(type) {
		case *vmhdBox:
			b.Vmhd = child
		case *dinfBox:
			b.Dinf = child
		default:
		}
	}

	return nil
}

func (b *minfBox) Type() string {
	return "minf"
}

func (b *minfBox) Children() []Box {
	return b.children
}
