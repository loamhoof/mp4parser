package mp4parser

import (
	"encoding/binary"
	"io"
)

type ftypBox struct {
	offset int64
	length uint32
	data   Pairs
}

func (b *ftypBox) Parse(r io.ReadSeeker) error {
	if _, err := r.Seek(b.offset, io.SeekStart); err != nil {
		return err
	}

	bytes := make([]byte, 4)

	if _, err := r.Read(bytes); err != nil {
		return err
	}
	l := binary.BigEndian.Uint32(bytes)

	b.length = l

	if _, err := r.Seek(4, io.SeekCurrent); err != nil {
		return err
	}

	b.data = make(Pairs, 0, 3)

	if _, err := r.Read(bytes); err != nil {
		return err
	}
	b.data = append(b.data, &Pair{"major_brand", string(bytes)})

	if _, err := r.Read(bytes); err != nil {
		return err
	}
	b.data = append(b.data, &Pair{"minor_version", binary.BigEndian.Uint32(bytes)})

	nBrands := (l - 16) / 4
	compatibleBrands := make([]string, nBrands)
	for i := 0; uint32(i) < nBrands; i++ {
		if _, err := r.Read(bytes); err != nil {
			return err
		}
		compatibleBrands[i] = string(bytes)
	}
	b.data = append(b.data, &Pair{"compatible_brands", compatibleBrands})

	return nil
}

func (b *ftypBox) Type() string {
	return "ftyp"
}

func (b *ftypBox) Offset() int64 {
	return b.offset
}

func (b *ftypBox) Length() uint32 {
	return b.length
}

func (b *ftypBox) Children() []Box {
	return []Box{}
}

func (b *ftypBox) Data() Pairs {
	return b.data
}
