package mp4parser

import (
	"io"
)

type trafBox struct {
	baseBox
	children []Box
	Tfhd     *tfhdBox
	Trun     *trunBox
}

func (b *trafBox) Parse(r io.ReadSeeker, startOffset int64, pp ParsePlan) error {
	size, _, fields, children, err := parseContainerBox(r, startOffset, pp)
	if err != nil {
		return err
	}
	b.size = size
	b.fields = fields
	b.children = children

	for _, child := range children {
		switch child := child.(type) {
		case *tfhdBox:
			b.Tfhd = child
		case *trunBox:
			b.Trun = child
		default:
		}
	}

	return nil
}

func (b *trafBox) Type() string {
	return "traf"
}

func (b *trafBox) Children() []Box {
	return b.children
}
