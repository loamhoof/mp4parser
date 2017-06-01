package mp4parser

import (
	"io"
)

type TrafBox struct {
	baseBox
	children []Box
	Tfhd     *TfhdBox
	Trun     *TrunBox
}

func (b *TrafBox) Parse(r io.ReadSeeker, startOffset int64, pp ParsePlan, pc ParseContext) error {
	size, _, fields, children, err := parseContainerBox(r, startOffset, pp, pc)
	if err != nil {
		return err
	}
	b.size = size
	b.fields = fields
	b.children = children

	for _, child := range children {
		switch child := child.(type) {
		case *TfhdBox:
			b.Tfhd = child
		case *TrunBox:
			b.Trun = child
		default:
		}
	}

	return nil
}

func (b *TrafBox) Type() string {
	return "traf"
}

func (b *TrafBox) Children() []Box {
	return b.children
}
