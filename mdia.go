package mp4parser

import (
	"io"
)

type MdiaBox struct {
	baseBox
	children []Box
	Mdhd     *MdhdBox
	Hdlr     *HdlrBox
	Minf     *MinfBox
}

func (b *MdiaBox) Parse(r io.ReadSeeker, startOffset int64, pp ParsePlan, pc ParseContext) error {
	size, _, fields, children, err := parseContainerBox(r, startOffset, pp, pc)
	if err != nil {
		return err
	}
	b.size = size
	b.fields = fields
	b.children = children

	for _, child := range children {
		switch child := child.(type) {
		case *MdhdBox:
			b.Mdhd = child
		case *HdlrBox:
			b.Hdlr = child
		case *MinfBox:
			b.Minf = child
		default:
		}
	}

	return nil
}

func (b *MdiaBox) Type() string {
	return "mdia"
}

func (b *MdiaBox) Children() []Box {
	return b.children
}
