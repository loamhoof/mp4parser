package mp4parser

import (
	"io"
)

type MoofBox struct {
	baseBox
	children []Box
	Mfhd     *MfhdBox
	Traf     *TrafBox
}

func (b *MoofBox) Parse(r io.ReadSeeker, startOffset int64, pp ParsePlan, pc ParseContext) error {
	size, _, fields, children, err := parseContainerBox(r, startOffset, pp, pc)
	if err != nil {
		return err
	}
	b.size = size
	b.fields = fields
	b.children = children

	for _, child := range children {
		switch child := child.(type) {
		case *MfhdBox:
			b.Mfhd = child
		case *TrafBox:
			b.Traf = child
		default:
		}
	}

	return nil
}

func (b *MoofBox) Type() string {
	return "moof"
}

func (b *MoofBox) Children() []Box {
	return b.children
}
