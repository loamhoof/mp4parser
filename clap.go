package mp4parser

import (
	"encoding/binary"
	"io"
)

type ClapBox struct {
	baseBox
}

func (b *ClapBox) Parse(r io.ReadSeeker, startOffset int64, pp ParsePlan, pc ParseContext) error {
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
	fields = append(fields, &Field{"cleanApertureWidthN", binary.BigEndian.Uint32(b4), offset, 32, 0})
	offset += 4

	if _, err := r.Read(b4); err != nil {
		return err
	}
	fields = append(fields, &Field{"cleanApertureWidthD", binary.BigEndian.Uint32(b4), offset, 32, 0})
	offset += 4

	if _, err := r.Read(b4); err != nil {
		return err
	}
	fields = append(fields, &Field{"cleanApertureHeightN", binary.BigEndian.Uint32(b4), offset, 32, 0})
	offset += 4

	if _, err := r.Read(b4); err != nil {
		return err
	}
	fields = append(fields, &Field{"cleanApertureHeightD", binary.BigEndian.Uint32(b4), offset, 32, 0})
	offset += 4

	if _, err := r.Read(b4); err != nil {
		return err
	}
	fields = append(fields, &Field{"horizOffN", binary.BigEndian.Uint32(b4), offset, 32, 0})
	offset += 4

	if _, err := r.Read(b4); err != nil {
		return err
	}
	fields = append(fields, &Field{"horizOffD", binary.BigEndian.Uint32(b4), offset, 32, 0})
	offset += 4

	if _, err := r.Read(b4); err != nil {
		return err
	}
	fields = append(fields, &Field{"vertOffN", binary.BigEndian.Uint32(b4), offset, 32, 0})
	offset += 4

	if _, err := r.Read(b4); err != nil {
		return err
	}
	fields = append(fields, &Field{"vertOffD", binary.BigEndian.Uint32(b4), offset, 32, 0})

	return nil
}

func (b *ClapBox) Type() string {
	return "clap"
}
