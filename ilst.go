package mp4parser

import (
	"io"
)

type IlstBox struct {
	baseBox
	children []Box
}

func (b *IlstBox) Parse(r io.ReadSeeker, startOffset int64, pp ParsePlan, pc ParseContext) error {
	size, _, fields, children, err := parseContainerBox(r, startOffset, pp, pc)
	if err != nil {
		return err
	}
	b.size = size
	b.fields = fields
	b.children = children

	return nil
}

func (b *IlstBox) Type() string {
	return "ilst"
}

func (b *IlstBox) Children() []Box {
	return b.children
}
