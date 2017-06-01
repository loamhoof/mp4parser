package mp4parser

import (
	"errors"
	"io"
)

type ParsePlan map[string]ParsePlan

type ParseContext map[string]string

func Parse(r io.ReadSeeker) (*MP4, error) {
	mp4 := &MP4{}

	return mp4, mp4.Parse(r, 0, nil, Background())
}

func ParseOnly(r io.ReadSeeker, pp ParsePlan) (*MP4, error) {
	mp4 := &MP4{}

	return mp4, mp4.Parse(r, 0, pp, Background())
}

func Find(r io.ReadSeeker, boxType string) (Box, error) {
	b4 := make([]byte, 4)
	b1 := make([]byte, 1)

	if _, err := r.Read(b4); err != nil {
		return nil, err
	}

	for {
		if string(b4) == boxType {
			box := newBox(boxType)

			offset, err := r.Seek(-8, io.SeekCurrent)
			if err != nil {
				return nil, err
			}

			if err := box.Parse(r, offset, nil, Background()); err != nil {
				return nil, err
			}

			return box, nil
		}

		if _, err := r.Read(b1); err != nil {
			return nil, err
		}

		b4 = append(b4[1:], b1...)
	}
}

func Replace(w io.ReadWriteSeeker, f *Field, v interface{}) error {
	var wb []byte
	switch v := v.(type) {
	case []byte:
		wb = v
	case string:
		wb = []byte(v)
	case uint8:
		wb = []byte{v}
	case int8:
		wb = []byte{byte(v)}
	case uint16:
		wb = []byte{
			byte(v >> 8),
			byte(v & 0xFF),
		}
	case int16:
		wb = []byte{
			byte(v >> 8),
			byte(v & 0xFF),
		}
	case uint32:
		wb = []byte{
			byte(v >> 24),
			byte(v >> 16 & 0xFF),
			byte(v >> 8 & 0xFF),
			byte(v & 0xFF),
		}
	case int32:
		wb = []byte{
			byte(v >> 24),
			byte(v >> 16 & 0xFF),
			byte(v >> 8 & 0xFF),
			byte(v & 0xFF),
		}
	case uint64:
		wb = []byte{
			byte(v >> 56),
			byte(v >> 48 & 0xFF),
			byte(v >> 40 & 0xFF),
			byte(v >> 32 & 0xFF),
			byte(v >> 24 & 0xFF),
			byte(v >> 16 & 0xFF),
			byte(v >> 8 & 0xFF),
			byte(v & 0xFF),
		}
	case int64:
		wb = []byte{
			byte(v >> 56),
			byte(v >> 48 & 0xFF),
			byte(v >> 40 & 0xFF),
			byte(v >> 32 & 0xFF),
			byte(v >> 24 & 0xFF),
			byte(v >> 16 & 0xFF),
			byte(v >> 8 & 0xFF),
			byte(v & 0xFF),
		}
	case bool:
		if v {
			wb = []byte{1}
		} else {
			wb = []byte{0}
		}
	default:
		return errors.New("TODO")
	}
	wbl := int((f.Bits + 7) / 8)
	wb = wb[len(wb)-wbl : len(wb)]

	if _, err := w.Seek(f.Offset, io.SeekStart); err != nil {
		return err
	}

	if f.BitsOffset == 0 && f.Bits%8 == 0 {
		if _, err := w.Write(wb); err != nil {
			return err
		}

		return nil
	}

	nBits := uint64(f.BitsOffset) + f.Bits
	restBits := byte(8 - nBits%8)
	rbl := int((nBits + uint64(restBits)) / 8)
	rb := make([]byte, rbl)

	if _, err := w.Read(rb); err != nil {
		return err
	}

	var offMask, restMask, i byte
	for i = 0; i < f.BitsOffset; i++ {
		offMask += 1 << (7 - i)
	}
	for i = 0; i < restBits; i++ {
		restMask += 1 << i
	}

	result := make([]byte, rbl)

	if rbl == 1 {
		result[0] = rb[0]&offMask | wb[0]&((offMask^restMask^0xFF)>>restBits)<<restBits | rb[0]&restMask
	} else {
		result[rbl-1] = (wb[wbl-1]&((restMask<<(8-restBits))^0xFF))<<restBits | rb[rbl-1]&restMask

		for i := 2; i < rbl; i++ {
			result[rbl-i] = (wb[wbl-i+1]&(restMask^0xFF))>>(8-restBits) | (wb[wbl-i]&restMask)<<restBits
		}

		result[0] = rb[0]&offMask | (((wb[0] & (restMask ^ 0xFF)) >> (8 - restBits)) & (offMask ^ 0xFF))
	}

	if _, err := w.Seek(f.Offset, io.SeekStart); err != nil {
		return err
	}

	if _, err := w.Write(result); err != nil {
		return err
	}

	return nil
}
