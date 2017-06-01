package mp4parser

import (
	"io"
)

type MdatBox struct {
	baseBox
}

func (b *MdatBox) Parse(r io.ReadSeeker, startOffset int64, pp ParsePlan, pc ParseContext) error {
	size, _, _, fields, err := parseBox(r, startOffset)
	if err != nil {
		return err
	}
	b.size = size
	b.fields = fields

	return nil
}

func (b *MdatBox) Type() string {
	return "mdat"
}
