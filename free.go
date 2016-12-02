package mp4parser

import (
	"io"
)

type freeBox struct {
	baseBox
	_type string
}

func (b *freeBox) Parse(r io.ReadSeeker, startOffset int64, pp ParsePlan) error {
	size, offset, _type, fields, err := parseBox(r, startOffset)
	if err != nil {
		return err
	}
	b.size = size
	b._type = _type
	b.fields = fields

	lData := size - uint64(offset-startOffset)
	bData := make([]byte, lData)
	if _, err := r.Read(bData); err != nil {
		return err
	}
	b.fields = append(b.fields, &Field{"data", string(bData), offset, lData * 8, 0})

	return nil
}

func (b *freeBox) Type() string {
	return b._type
}
