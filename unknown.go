package mp4parser

import (
	"io"
	"strconv"
)

type UnknownBox struct {
	baseBox
	_type string
}

func (b *UnknownBox) Parse(r io.ReadSeeker, startOffset int64, pp ParsePlan, pc ParseContext) error {
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
	b.fields = append(b.fields, &Field{"unknown", bData, offset, lData * 8, 0})

	return nil
}

func (b *UnknownBox) Type() string {
	return strconv.Quote("#" + b._type + "#")
}
