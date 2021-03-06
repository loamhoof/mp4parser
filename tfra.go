package mp4parser

import (
	"encoding/binary"
	"io"
)

type TfraBox struct {
	baseBox
}

func (b *TfraBox) Parse(r io.ReadSeeker, startOffset int64, pp ParsePlan, pc ParseContext) error {
	size, offset, _, version, _, fields, err := parseFullBox(r, startOffset)
	if err != nil {
		return err
	}
	b.size = size
	b.fields = fields

	b1 := make([]byte, 1)
	b4 := make([]byte, 4)
	b8 := make([]byte, 8)

	if _, err := r.Read(b4); err != nil {
		return err
	}
	b.fields = append(b.fields, &Field{"track_ID", binary.BigEndian.Uint32(b4), offset, 32, 0})
	offset += 4

	if _, err := r.Seek(3, io.SeekCurrent); err != nil {
		return err
	}
	offset += 3

	if _, err := r.Read(b1); err != nil {
		return err
	}
	lengthSizeOfTrafNum := b1[0] & 48
	b.fields = append(b.fields, &Field{"length_size_of_traf_num", lengthSizeOfTrafNum, offset, 2, 2})
	lengthSizeOfTrunNum := b1[0] & 12
	b.fields = append(b.fields, &Field{"length_size_of_trun_num", lengthSizeOfTrunNum, offset, 2, 4})
	lengthSizeOfSampleNum := b1[0] & 3
	b.fields = append(b.fields, &Field{"length_size_of_sample_num", lengthSizeOfSampleNum, offset, 2, 6})
	offset += 1

	if _, err := r.Read(b4); err != nil {
		return err
	}
	numberOfEntry := binary.BigEndian.Uint32(b4)
	b.fields = append(b.fields, &Field{"number_of_entry", numberOfEntry, offset, 32, 0})
	offset += 4

	entriesOffset := offset
	entries := make([]Fields, numberOfEntry)
	for i := 0; uint32(i) < numberOfEntry; i++ {
		entry := make(Fields, 0, 5)

		if version == 1 {
			if _, err := r.Read(b8); err != nil {
				return err
			}
			entry = append(entry, &Field{"time", binary.BigEndian.Uint64(b8), offset, 64, 0})
			offset += 8

			if _, err := r.Read(b8); err != nil {
				return err
			}
			entry = append(entry, &Field{"moof_offset", binary.BigEndian.Uint64(b8), offset, 64, 0})
			offset += 8
		} else {
			if _, err := r.Read(b4); err != nil {
				return err
			}
			entry = append(entry, &Field{"time", binary.BigEndian.Uint32(b4), offset, 32, 0})
			offset += 4

			if _, err := r.Read(b4); err != nil {
				return err
			}
			entry = append(entry, &Field{"moof_offset", binary.BigEndian.Uint32(b4), offset, 32, 0})
			offset += 4
		}

		lTraf := lengthSizeOfTrafNum + 1
		bTraf := make([]byte, lTraf)
		if _, err := r.Read(bTraf); err != nil {
			return err
		}
		trafNumber := binary.BigEndian.Uint64(append(make([]byte, 8-lTraf), bTraf...))
		entry = append(entry, &Field{"traf_number", trafNumber, offset, 8 * uint64(lTraf), 0})
		offset += int64(lTraf)

		lTrun := lengthSizeOfTrunNum + 1
		bTrun := make([]byte, lTrun)
		if _, err := r.Read(bTrun); err != nil {
			return err
		}
		trunNumber := binary.BigEndian.Uint64(append(make([]byte, 8-lTrun), bTrun...))
		entry = append(entry, &Field{"trun_number", trunNumber, offset, 8 * uint64(lTrun), 0})
		offset += int64(lTrun)

		lSample := lengthSizeOfSampleNum + 1
		bSample := make([]byte, lSample)
		if _, err := r.Read(bSample); err != nil {
			return err
		}
		sampleNumber := binary.BigEndian.Uint64(append(make([]byte, 8-lSample), bSample...))
		entry = append(entry, &Field{"sample_number", sampleNumber, offset, 8 * uint64(lSample), 0})
		offset += int64(lSample)

		entries[i] = entry
	}

	b.fields = append(b.fields, &Field{"entries", entries, entriesOffset, uint64(offset-entriesOffset) * 8, 0})

	return nil
}

func (b *TfraBox) Type() string {
	return "tfra"
}
