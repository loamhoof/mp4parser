package mp4parser

import (
	"io"
)

type MinfBox struct {
	baseBox
	children []Box
	Vmhd     *VmhdBox
	Dinf     *DinfBox
	Stbl     *StblBox
}

func (b *MinfBox) Parse(r io.ReadSeeker, startOffset int64, pp ParsePlan, pc ParseContext) error {
	size, _, fields, children, err := parseContainerBox(r, startOffset, pp, pc)
	if err != nil {
		return err
	}
	b.size = size
	b.fields = fields
	b.children = children

	for _, child := range children {
		switch child := child.(type) {
		case *VmhdBox:
			b.Vmhd = child
		case *DinfBox:
			b.Dinf = child
		case *StblBox:
			b.Stbl = child
		default:
		}
	}

	return nil
}

func (b *MinfBox) Type() string {
	return "minf"
}

func (b *MinfBox) Children() []Box {
	return b.children
}
