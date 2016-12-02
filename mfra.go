package mp4parser

import (
	"io"
)

type mfraBox struct {
	baseBox
	children []Box
	Tfra     *tfraBox
	Mfro     *mfroBox
}

func (b *mfraBox) Parse(r io.ReadSeeker, startOffset int64, pp ParsePlan) error {
	size, _, fields, children, err := parseContainerBox(r, startOffset, pp)
	if err != nil {
		return err
	}
	b.size = size
	b.fields = fields
	b.children = children

	for _, child := range children {
		switch child := child.(type) {
		case *tfraBox:
			b.Tfra = child
		case *mfroBox:
			b.Mfro = child
		default:
		}
	}

	return nil
}

func (b *mfraBox) Type() string {
	return "mfra"
}

func (b *mfraBox) Children() []Box {
	return b.children
}
