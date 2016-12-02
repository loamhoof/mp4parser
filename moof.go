package mp4parser

import (
	"io"
)

type moofBox struct {
	baseBox
	children []Box
	Mfhd     *mfhdBox
	Traf     *trafBox
}

func (b *moofBox) Parse(r io.ReadSeeker, startOffset int64, pp ParsePlan) error {
	size, _, fields, children, err := parseContainerBox(r, startOffset, pp)
	if err != nil {
		return err
	}
	b.size = size
	b.fields = fields
	b.children = children

	for _, child := range children {
		switch child := child.(type) {
		case *mfhdBox:
			b.Mfhd = child
		case *trafBox:
			b.Traf = child
		default:
		}
	}

	return nil
}

func (b *moofBox) Type() string {
	return "moof"
}

func (b *moofBox) Children() []Box {
	return b.children
}
