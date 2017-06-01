package mp4parser

import (
	"io"
)

type UdtaBox struct {
	baseBox
	children []Box
	Meta     *MetaBox
}

func (b *UdtaBox) Parse(r io.ReadSeeker, startOffset int64, pp ParsePlan, pc ParseContext) error {
	size, _, fields, children, err := parseContainerBox(r, startOffset, pp, pc)
	if err != nil {
		return err
	}
	b.size = size
	b.fields = fields
	b.children = children

	for _, child := range children {
		switch child := child.(type) {
		case *MetaBox:
			b.Meta = child
		default:
		}
	}

	return nil
}

func (b *UdtaBox) Type() string {
	return "udta"
}

func (b *UdtaBox) Children() []Box {
	return b.children
}
