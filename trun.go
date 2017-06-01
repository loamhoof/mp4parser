package mp4parser

import (
	"encoding/binary"
	"errors"
	"io"
)

type TrunBox struct {
	baseBox
}

func (b *TrunBox) Parse(r io.ReadSeeker, startOffset int64, pp ParsePlan, pc ParseContext) error {
	size, offset, _, version, flags, fields, err := parseFullBox(r, startOffset)
	if err != nil {
		return err
	}
	b.size = size
	b.fields = fields

	b4 := make([]byte, 4)

	if _, err := r.Read(b4); err != nil {
		return err
	}
	sampleCount := binary.BigEndian.Uint32(b4)
	b.fields = append(b.fields, &Field{"sample_count", sampleCount, offset, 32, 0})
	offset += 4

	if flags&0x01 == 0x01 {
		if _, err := r.Read(b4); err != nil {
			return err
		}
		dataOffset, read := binary.Varint(b4)
		if read <= 0 {
			return errors.New("TODO")
		}
		b.fields = append(b.fields, &Field{"data_offset", dataOffset, offset, 32, 0})
		offset += 4
	}

	if flags&0x04 == 0x04 {
		if _, err := r.Read(b4); err != nil {
			return err
		}
		b.fields = append(b.fields, &Field{"first_sample_flags", binary.BigEndian.Uint32(b4), offset, 32, 0})
		offset += 4
	}

	samplesOffset := offset
	samples := make([]Fields, sampleCount)
	for i := 0; uint32(i) < sampleCount; i++ {
		sample := make(Fields, 0, 4)

		if flags>>8&0x01 == 0x01 {
			if _, err := r.Read(b4); err != nil {
				return err
			}
			sample = append(sample, &Field{"sample_duration", binary.BigEndian.Uint32(b4), offset, 32, 0})
			offset += 4
		}

		if flags>>8&0x02 == 0x02 {
			if _, err := r.Read(b4); err != nil {
				return err
			}
			sample = append(sample, &Field{"sample_size", binary.BigEndian.Uint32(b4), offset, 32, 0})
			offset += 4
		}

		if flags>>8&0x04 == 0x04 {
			if _, err := r.Read(b4); err != nil {
				return err
			}
			sample = append(sample, &Field{"sample_flags", binary.BigEndian.Uint32(b4), offset, 32, 0})
			offset += 4
		}

		if flags>>8&0x08 == 0x08 {
			if _, err := r.Read(b4); err != nil {
				return err
			}
			if version == 0 {
				sample = append(sample, &Field{"sample_composition_time_offset", binary.BigEndian.Uint32(b4), offset, 32, 0})
				offset += 4
			} else {
				sampleCompositionTimeOffset, read := binary.Varint(b4)
				if read <= 0 {
					return errors.New("TODO")
				}
				sample = append(sample, &Field{"sample_composition_time_offset", sampleCompositionTimeOffset, offset, 32, 0})
				offset += 4
			}
		}

		samples[i] = sample
	}
	b.fields = append(b.fields, &Field{"samples", samples, samplesOffset, uint64(offset-samplesOffset) * 8, 0})

	return nil
}

func (b *TrunBox) Type() string {
	return "trun"
}
