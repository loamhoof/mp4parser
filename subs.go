package mp4parser

import (
	"encoding/binary"
	"io"
)

type SubsBox struct {
	baseBox
}

func (b *SubsBox) Parse(r io.ReadSeeker, startOffset int64, pp ParsePlan, pc ParseContext) error {
	size, offset, _, version, _, fields, err := parseFullBox(r, startOffset)
	if err != nil {
		return err
	}
	b.size = size
	b.fields = fields

	b1 := make([]byte, 1)
	b2 := make([]byte, 2)
	b4 := make([]byte, 4)

	if _, err := r.Read(b4); err != nil {
		return err
	}
	entryCount := binary.BigEndian.Uint32(b4)
	b.fields = append(b.fields, &Field{"entry_count", entryCount, offset, 32, 0})
	offset += 4

	entriesOffset := offset
	entries := make([]Fields, entryCount)
	for i := 0; uint32(i) < entryCount; i++ {
		entry := make(Fields, 0, 3)

		if _, err := r.Read(b4); err != nil {
			return err
		}
		entry = append(entry, &Field{"sample_delta", binary.BigEndian.Uint32(b4), offset, 32, 0})
		offset += 4

		if _, err := r.Read(b2); err != nil {
			return err
		}
		subsampleCount := binary.BigEndian.Uint16(b2)
		entry = append(entry, &Field{"subsample_count", subsampleCount, offset, 16, 0})
		offset += 2

		subsamplesOffset := offset
		subsamples := make([]Fields, subsampleCount)
		for j := 0; uint16(j) < subsampleCount; j++ {
			subsample := make(Fields, 0, subsampleCount)

			if version == 1 {
				if _, err := r.Read(b4); err != nil {
					return err
				}
				subsample = append(subsample, &Field{"subsample_size", binary.BigEndian.Uint32(b4), offset, 32, 0})
				offset += 4
			} else {
				if _, err := r.Read(b2); err != nil {
					return err
				}
				subsample = append(subsample, &Field{"subsample_size", binary.BigEndian.Uint16(b2), offset, 16, 0})
				offset += 2
			}

			if _, err := r.Read(b1); err != nil {
				return err
			}
			subsample = append(subsample, &Field{"subsample_priority", b1[0], offset, 8, 0})
			offset += 1

			if _, err := r.Read(b1); err != nil {
				return err
			}
			subsample = append(subsample, &Field{"discardable", b1[0], offset, 8, 0})
			offset += 1

			if _, err := r.Seek(4, io.SeekCurrent); err != nil {
				return err
			}
			offset += 4

			subsamples[j] = subsample
		}

		entry = append(entry, &Field{"subsamples", subsamples, subsamplesOffset, uint64(offset-subsamplesOffset) * 8, 0})

		entries[i] = entry
	}
	b.fields = append(b.fields, &Field{"entries", entries, entriesOffset, uint64(offset-entriesOffset) * 8, 0})

	return nil
}

func (b *SubsBox) Type() string {
	return "subs"
}
