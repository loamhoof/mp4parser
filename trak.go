package mp4parser

import (
	"io"
)

type TrakBox struct {
	baseBox
	children []Box
	Tkhd     *TkhdBox
	Mdia     *MdiaBox
}

func (b *TrakBox) Parse(r io.ReadSeeker, startOffset int64, pp ParsePlan, pc ParseContext) error {
	size, _, fields, children, err := parseContainerBox(r, startOffset, pp, pc)
	if err != nil {
		return err
	}
	b.size = size
	b.fields = fields
	b.children = children

	for _, child := range children {
		switch child := child.(type) {
		case *TkhdBox:
			b.Tkhd = child
		case *MdiaBox:
			b.Mdia = child
		default:
		}
	}

	return nil
}

func (b *TrakBox) Type() string {
	return "trak"
}

func (b *TrakBox) Children() []Box {
	return b.children
}
