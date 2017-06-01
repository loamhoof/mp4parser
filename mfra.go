package mp4parser

import (
	"io"
)

type MfraBox struct {
	baseBox
	children []Box
	Tfra     *TfraBox
	Mfro     *MfroBox
}

func (b *MfraBox) Parse(r io.ReadSeeker, startOffset int64, pp ParsePlan, pc ParseContext) error {
	size, _, fields, children, err := parseContainerBox(r, startOffset, pp, pc)
	if err != nil {
		return err
	}
	b.size = size
	b.fields = fields
	b.children = children

	for _, child := range children {
		switch child := child.(type) {
		case *TfraBox:
			b.Tfra = child
		case *MfroBox:
			b.Mfro = child
		default:
		}
	}

	return nil
}

func (b *MfraBox) Type() string {
	return "mfra"
}

func (b *MfraBox) Children() []Box {
	return b.children
}
