package mp4parser

import (
	"encoding/binary"
	"io"
)

type VisualSampleEntry struct {
	baseBox
	_type    string
	children []Box
	Clap     *ClapBox
	Pasp     *PaspBox
	AvcC     *AvcCBox
	Btrt     *BtrtBox
}

func (b *VisualSampleEntry) Parse(r io.ReadSeeker, startOffset int64, pp ParsePlan, pc ParseContext) error {
	size, offset, _type, fields, err := parseSampleEntry(r, startOffset)
	if err != nil {
		return err
	}
	b.size = size
	b._type = _type
	b.fields = fields

	b2 := make([]byte, 2)
	b4 := make([]byte, 4)
	b32 := make([]byte, 32)

	if _, err := r.Seek(16, io.SeekCurrent); err != nil {
		return err
	}
	offset += 16

	if _, err := r.Read(b2); err != nil {
		return err
	}
	fields = append(fields, &Field{"width", binary.BigEndian.Uint16(b2), offset, 16, 0})
	offset += 2

	if _, err := r.Read(b2); err != nil {
		return err
	}
	fields = append(fields, &Field{"height", binary.BigEndian.Uint16(b2), offset, 16, 0})
	offset += 2

	if _, err := r.Read(b4); err != nil {
		return err
	}
	fields = append(fields, &Field{"horizresolution", binary.BigEndian.Uint32(b4), offset, 32, 0})
	offset += 4

	if _, err := r.Read(b4); err != nil {
		return err
	}
	fields = append(fields, &Field{"vertresolution", binary.BigEndian.Uint32(b4), offset, 32, 0})
	offset += 4

	if _, err := r.Seek(4, io.SeekCurrent); err != nil {
		return err
	}
	offset += 4

	if _, err := r.Read(b2); err != nil {
		return err
	}
	fields = append(fields, &Field{"frame_count", binary.BigEndian.Uint16(b2), offset, 16, 0})
	offset += 2

	if _, err := r.Read(b32); err != nil {
		return err
	}
	fields = append(fields, &Field{"compressorname", string(b32), offset, 272, 0})
	offset += 32

	if _, err := r.Read(b2); err != nil {
		return err
	}
	fields = append(fields, &Field{"depth", binary.BigEndian.Uint16(b2), offset, 16, 0})
	offset += 2

	if _, err := r.Seek(2, io.SeekCurrent); err != nil {
		return err
	}
	offset += 2

	b.children = make([]Box, 0, 2)
	for offset-startOffset < int64(size) {
		if _, err := r.Seek(offset, io.SeekStart); err != nil {
			break
		}

		if _, err := r.Read(b4); err != nil {
			return err
		}
		l := binary.BigEndian.Uint32(b4)

		if _, err := r.Read(b4); err != nil {
			return err
		}
		boxType := string(b4)

		if _, ok := pp[boxType]; ok || pp == nil {
			box := newBox(boxType)
			if err := box.Parse(r, offset, pp[boxType], pc); err != nil {
				return err
			}
			b.children = append(b.children, box)

			switch box := box.(type) {
			case *ClapBox:
				b.Clap = box
			case *PaspBox:
				b.Pasp = box
			case *AvcCBox:
				b.AvcC = box
			case *BtrtBox:
				b.Btrt = box
			default:
			}
		}

		offset += int64(l)
	}

	return nil
}

func (b *VisualSampleEntry) Type() string {
	return "vide_se:" + b._type
}

func (b *VisualSampleEntry) Children() []Box {
	return b.children
}
