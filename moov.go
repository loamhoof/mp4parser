package mp4parser

import (
	"io"
)

type MoovBox struct {
	baseBox
	children []Box
	Mvhd     *MvhdBox
	Trak     []*TrakBox
	Mvex     *MvexBox
	Udta     *UdtaBox
}

func (b *MoovBox) Parse(r io.ReadSeeker, startOffset int64, pp ParsePlan, pc ParseContext) error {
	size, _, fields, children, err := parseContainerBox(r, startOffset, pp, pc)
	if err != nil {
		return err
	}
	b.size = size
	b.fields = fields
	b.children = children

	for _, child := range children {
		switch child := child.(type) {
		case *MvhdBox:
			b.Mvhd = child
		case *TrakBox:
			b.Trak = append(b.Trak, child)
		case *MvexBox:
			b.Mvex = child
		case *UdtaBox:
			b.Udta = child
		default:
		}
	}

	return nil
}

func (b *MoovBox) Type() string {
	return "moov"
}

func (b *MoovBox) Children() []Box {
	return b.children
}
