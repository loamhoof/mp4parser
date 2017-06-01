package mp4parser

import (
	"io"
)

type EdtsBox struct {
	baseBox
	children []Box
	Elst     *ElstBox
}

func (b *EdtsBox) Parse(r io.ReadSeeker, startOffset int64, pp ParsePlan, pc ParseContext) error {
	size, _, fields, children, err := parseContainerBox(r, startOffset, pp, pc)
	if err != nil {
		return err
	}
	b.size = size
	b.fields = fields
	b.children = children

	for _, child := range children {
		switch child := child.(type) {
		case *ElstBox:
			b.Elst = child
		default:
		}
	}

	return nil
}

func (b *EdtsBox) Type() string {
	return "edts"
}

func (b *EdtsBox) Children() []Box {
	return b.children
}
