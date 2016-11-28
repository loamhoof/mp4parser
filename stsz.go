package mp4parser

import (
	"encoding/binary"
	"io"
)

type stszBox struct {
	baseBox
}

func (b *stszBox) Parse(r io.ReadSeeker, startOffset int64) error {
	size, offset, _, _, _, fields, err := parseFullBox(r, startOffset)
	if err != nil {
		return err
	}
	b.size = size
	b.fields = fields

	b4 := make([]byte, 4)

	if _, err := r.Read(b4); err != nil {
		return err
	}
	sampleSize := binary.BigEndian.Uint32(b4)
	b.fields = append(b.fields, &Field{"sample_size", sampleSize, offset, 32, 0})
	offset += 4

	if _, err := r.Read(b4); err != nil {
		return err
	}
	sampleCount := binary.BigEndian.Uint32(b4)
	b.fields = append(b.fields, &Field{"sample_count", sampleCount, offset, 32, 0})
	offset += 4

	if sampleSize == 0 {
		samplesOffset := offset
		samples := make([]Fields, sampleCount)
		for i := 0; uint32(i) < sampleCount; i++ {
			if _, err := r.Read(b4); err != nil {
				return err
			}
			samples[i] = Fields{{"entry_size", binary.BigEndian.Uint32(b4), offset, 32, 0}}
			offset += 4
		}
		b.fields = append(b.fields, &Field{"samples", samples, samplesOffset, uint64(offset-samplesOffset) * 8, 0})
	}

	return nil
}

func (b *stszBox) Type() string {
	return "stsz"
}
