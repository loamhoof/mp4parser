package mp4parser

import (
	"encoding/binary"
	"io"
)

type ftypBox struct {
	size   uint64
	fields Fields
}

func (b *ftypBox) Parse(r io.ReadSeeker, startOffset int64) error {
	size, offset, _, fields, err := parseBox(r, startOffset)
	if err != nil {
		return err
	}
	b.size = size
	b.fields = fields

	b4 := make([]byte, 4)

	if _, err := r.Read(b4); err != nil {
		return err
	}
	b.fields = append(b.fields, &Field{"major_brand", string(b4), offset, 32, 0})
	offset += 4

	if _, err := r.Read(b4); err != nil {
		return err
	}
	b.fields = append(b.fields, &Field{"minor_version", binary.BigEndian.Uint32(b4), offset, 32, 0})
	offset += 4

	nBrands := (size - 16) / 4
	compatibleBrands := make([]string, nBrands)
	for i := 0; uint64(i) < nBrands; i++ {
		if _, err := r.Read(b4); err != nil {
			return err
		}
		compatibleBrands[i] = string(b4)
	}
	b.fields = append(b.fields, &Field{"compatible_brands", compatibleBrands, offset, 32 * nBrands, 0})

	return nil
}

func (b *ftypBox) Type() string {
	return "ftyp"
}

func (b *ftypBox) Offset() int64 {
	return b.fields[0].Offset
}

func (b *ftypBox) Size() uint64 {
	return b.size
}

func (b *ftypBox) Children() []Box {
	return []Box{}
}

func (b *ftypBox) Data() Fields {
	return b.fields
}
