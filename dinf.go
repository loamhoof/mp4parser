package mp4parser

import (
	"io"
)

type DinfBox struct {
	baseBox
	children []Box
	Dref     *DrefBox
}

func (b *DinfBox) Parse(r io.ReadSeeker, startOffset int64, pp ParsePlan, pc ParseContext) error {
	size, _, fields, children, err := parseContainerBox(r, startOffset, pp, pc)
	if err != nil {
		return err
	}
	b.size = size
	b.fields = fields
	b.children = children

	for _, child := range children {
		switch child := child.(type) {
		case *DrefBox:
			b.Dref = child
		default:
		}
	}

	return nil
}

func (b *DinfBox) Type() string {
	return "dinf"
}

func (b *DinfBox) Children() []Box {
	return b.children
}
