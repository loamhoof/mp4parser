package mp4parser

import (
	"encoding/binary"
	"io"
)

type StsdBox struct {
	baseBox
	children []Box
	// Soun []*AudioSampleEntry // TODO
	Vide []*VisualSampleEntry
	// Hint []*HintSampleEntry // TODO
	// Meta []*MetadataSampleEntry // TODO
}

func (b *StsdBox) Parse(r io.ReadSeeker, startOffset int64, pp ParsePlan, pc ParseContext) error {
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
	entryCount := binary.BigEndian.Uint32(b4)
	b.fields = append(b.fields, &Field{"entry_count", entryCount, offset, 32, 0})
	offset += 4

	b.children = make([]Box, 0, entryCount)
	_type := pc["handler_type"] + "_se"
	for i := 0; uint32(i) < entryCount; i++ {
		box := newBox(_type)
		if err := box.Parse(r, offset, pp[_type], pc); err != nil {
			return err
		}
		b.children = append(b.children, box)

		switch box := box.(type) {
		// TODO
		// case *AudioSampleEntry:
		// 	b.Soun = append(b.Soun, box)
		case *VisualSampleEntry:
			b.Vide = append(b.Vide, box)
		// TODO
		// case *HintSampleEntry:
		// 	b.Hint = append(b.Hint, box)
		// TODO
		// case *MetadataSampleEntry:
		// 	b.Meta = append(b.Meta, box)
		default:
		}
	}

	return nil
}

func (b *StsdBox) Type() string {
	return "stsd"
}

func (b *StsdBox) Children() []Box {
	return b.children
}

func parseSampleEntry(r io.ReadSeeker, startOffset int64) (size uint64, offset int64, _type string, fields Fields, err error) {
	size, offset, _type, fields, err = parseBox(r, startOffset)
	if err != nil {
		return
	}

	b2 := make([]byte, 2)

	if _, err = r.Seek(6, io.SeekCurrent); err != nil {
		return
	}
	offset += 6

	if _, err = r.Read(b2); err != nil {
		return
	}
	fields = append(fields, &Field{"data_reference_index", binary.BigEndian.Uint16(b2), offset, 16, 0})
	offset += 2

	return
}
