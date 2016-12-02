package mp4parser

import (
	"io"
)

type moovBox struct {
	baseBox
	children []Box
	Mvhd     *mvhdBox
	Trak     *trakBox
	Mvex     *mvexBox
}

func (b *moovBox) Parse(r io.ReadSeeker, startOffset int64, pp ParsePlan) error {
	size, _, fields, children, err := parseContainerBox(r, startOffset, pp)
	if err != nil {
		return err
	}
	b.size = size
	b.fields = fields
	b.children = children

	for _, child := range children {
		switch child := child.(type) {
		case *mvhdBox:
			b.Mvhd = child
		case *trakBox:
			b.Trak = child
		case *mvexBox:
			b.Mvex = child
		default:
		}
	}

	return nil
}

func (b *moovBox) Type() string {
	return "moov"
}

func (b *moovBox) Children() []Box {
	return b.children
}
