package mp4parser

import (
	"io"
)

type mdatBox struct {
	baseBox
}

func (b *mdatBox) Parse(r io.ReadSeeker, startOffset int64) error {
	size, _, _, fields, err := parseBox(r, startOffset)
	if err != nil {
		return err
	}
	b.size = size
	b.fields = fields

	return nil
}

func (b *mdatBox) Type() string {
	return "mdat"
}
