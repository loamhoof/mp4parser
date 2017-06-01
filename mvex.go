package mp4parser

import (
	"io"
)

type MvexBox struct {
	baseBox
	children []Box
	Mehd     *MehdBox
	Trex     *TrexBox
}

func (b *MvexBox) Parse(r io.ReadSeeker, startOffset int64, pp ParsePlan, pc ParseContext) error {
	size, _, fields, children, err := parseContainerBox(r, startOffset, pp, pc)
	if err != nil {
		return err
	}
	b.size = size
	b.fields = fields
	b.children = children

	for _, child := range children {
		switch child := child.(type) {
		case *MehdBox:
			b.Mehd = child
		case *TrexBox:
			b.Trex = child
		default:
		}
	}

	return nil
}

func (b *MvexBox) Type() string {
	return "mvex"
}

func (b *MvexBox) Children() []Box {
	return b.children
}
