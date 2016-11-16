package mp4parser

import (
	"encoding/binary"
	"io"
)

type subsBox struct {
	size   uint64
	fields Fields
}

func (b *subsBox) Parse(r io.ReadSeeker, startOffset int64) error {
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
	b.fields = append(b.fields, &Field{"entry_count", entryCount, offset, 32})
	offset += 4

	entries := make([]Fields, entryCount)
	for i := 0; uint32(i) < entryCount; i++ {
		entry := make(Fields, 0, 3)

		if _, err := r.Read(b4); err != nil {
			return err
		}
		entry = append(entry, &Field{"sample_delta", binary.BigEndian.Uint32(b4), offset, 32})
		offset += 4

		if _, err := r.Read(b2); err != nil {
			return err
		}
		subsampleCount := binary.BigEndian.Uint16(b2)
		entry = append(entry, &Field{"subsample_count", subsampleCount, offset, 16})
		offset += 2

		subsamples := make([]Fields, subsampleCount)
		for j := 0; uint16(j) < subsampleCount; j++ {
			subsample := make(Fields, 0, subsampleCount)

			if version == 1 {
				if _, err := r.Read(b4); err != nil {
					return err
				}
				subsample = append(subsample, &Field{"subsample_size", binary.BigEndian.Uint32(b4), offset, 32})
				offset += 4
			} else {
				if _, err := r.Read(b2); err != nil {
					return err
				}
				subsample = append(subsample, &Field{"subsample_size", binary.BigEndian.Uint16(b2), offset, 16})
				offset += 2
			}

			if _, err := r.Read(b1); err != nil {
				return err
			}
			subsample = append(subsample, &Field{"subsample_priority", b1[0], offset, 8})
			offset += 1

			if _, err := r.Read(b1); err != nil {
				return err
			}
			subsample = append(subsample, &Field{"discardable", b1[0], offset, 8})
			offset += 1

			if _, err := r.Seek(4, io.SeekCurrent); err != nil {
				return err
			}
			offset += 4

			subsamples[j] = subsample
		}

		entry = append(entry, &Field{"subsamples", subsamples, offset, 0}) // TODO

		entries[i] = entry
	}
	b.fields = append(b.fields, &Field{"entries", entries, offset, 0}) // TODO

	return nil
}

func (b *subsBox) Type() string {
	return "subs"
}

func (b *subsBox) Offset() int64 {
	return b.fields[0].Offset
}

func (b *subsBox) Size() uint64 {
	return b.size
}

func (b *subsBox) Children() []Box {
	return []Box{}
}

func (b *subsBox) Data() Fields {
	return b.fields
}
