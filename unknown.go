package mp4parser

import (
	"io"
)

type unknownBox struct {
	baseBox
	_type string
}

func (b *unknownBox) Parse(r io.ReadSeeker, startOffset int64, pp ParsePlan) error {
	size, offset, _, fields, err := parseBox(r, startOffset)
	if err != nil {
		return err
	}
	b.size = size
	b.fields = fields

	lData := size - uint64(offset-startOffset)
	bData := make([]byte, lData)
	if _, err := r.Read(bData); err != nil {
		return err
	}
	b.fields = append(b.fields, &Field{"unknown", string(bData), offset, lData * 8, 0})

	return nil
}

func (b *unknownBox) Type() string {
	return "#" + b._type + "#"
}
