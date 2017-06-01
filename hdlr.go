package mp4parser

import (
	"io"
)

type HdlrBox struct {
	baseBox
}

func (b *HdlrBox) Parse(r io.ReadSeeker, startOffset int64, pp ParsePlan, pc ParseContext) error {
	size, offset, _, _, _, fields, err := parseFullBox(r, startOffset)
	if err != nil {
		return err
	}
	b.size = size
	b.fields = fields

	b4 := make([]byte, 4)

	if _, err := r.Seek(4, io.SeekCurrent); err != nil {
		return err
	}
	offset += 4

	if _, err := r.Read(b4); err != nil {
		return err
	}
	b.fields = append(b.fields, &Field{"handler_type", string(b4), offset, 32, 0})
	offset += 4

	// Parsing the stsd box requires the handler_type
	pc["handler_type"] = string(b4)

	if _, err := r.Seek(12, io.SeekCurrent); err != nil {
		return err
	}
	offset += 12

	lName := size - uint64(offset-startOffset)
	bName := make([]byte, lName)
	if _, err := r.Read(bName); err != nil {
		return err
	}
	b.fields = append(b.fields, &Field{"name", string(bName), offset, lName * 8, 0})

	return nil
}

func (b *HdlrBox) Type() string {
	return "hdlr"
}
