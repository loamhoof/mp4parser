package mp4parser

import (
	"encoding/binary"
	"errors"
	"io"
	"math"
	"strconv"
	"time"
)

var PaddingErr = errors.New("Is just padding, not a box")

type GPMFBox struct {
	baseBox
	children []Box
	Meta     *MetaBox
}

func (b *GPMFBox) Parse(r io.ReadSeeker, startOffset int64, pp ParsePlan, pc ParseContext) error {
	gpmfSize, offset, _, fields, err := parseBox(r, startOffset)
	if err != nil {
		return err
	}
	b.size = gpmfSize
	b.fields = fields

	endOffset := startOffset + int64(gpmfSize)

	children := make([]Box, 0, 2)
	for {
		// less than 4 bytes remaining in the GPMF box (only 0 should happen)
		if endOffset-offset < 4 {
			break
		}

		box := &GenericGPMFBox{}

		if err := box.Parse(r, offset, pp, pc); err != nil {
			if errors.Is(err, PaddingErr) {
				offset += 4

				break
			}

			return err
		}

		offset += int64(box.size)
		children = append(children, box)
	}

	if offset != endOffset {
		if _, err := r.Seek(endOffset, io.SeekStart); err != nil {
			return err
		}
	}

	b.children = children

	return nil
}

func (b *GPMFBox) Type() string {
	return "GPMF"
}

func (b *GPMFBox) Children() []Box {
	return b.children
}

type GenericGPMFBox struct {
	baseBox
	_type string
}

func (b *GenericGPMFBox) Parse(r io.ReadSeeker, startOffset int64, pp ParsePlan, pc ParseContext) error {
	b1 := make([]byte, 1)
	b2 := make([]byte, 2)
	b4 := make([]byte, 4)
	offset := startOffset
	fields := make(Fields, 0, 5)

	if _, err := r.Read(b4); err != nil {
		return err
	}
	_type := string(b4)
	b._type = _type
	fields = append(fields, &Field{"type", _type, offset, 32, 0})
	offset += 4

	if _type == "\x00\x00\x00\x00" {
		return PaddingErr
	}

	if _, err := r.Read(b1); err != nil {
		return err
	}
	typeChar := string(b1)
	fields = append(fields, &Field{"typeChar", typeChar, offset, 8, 0})
	offset += 1

	if _, err := r.Read(b1); err != nil {
		return err
	}
	size := b1[0]
	fields = append(fields, &Field{"size", size, offset, 8, 0})
	offset += 1

	if _, err := r.Read(b2); err != nil {
		return err
	}
	repeat := binary.BigEndian.Uint16(b2)
	fields = append(fields, &Field{"repeat", repeat, offset, 8, 0})
	offset += 2

	if typeChar == "\x00" {
		b.size = uint64(offset - startOffset)
		b.fields = fields

		return nil
	}

	values := make([][]byte, repeat)
	for i := 0; i < int(repeat); i++ {
		bvalue := make([]byte, size)
		if _, err := r.Read(bvalue); err != nil {
			return err
		}

		values[i] = bvalue
	}
	castValues, err := b.castSliceToTypeChar(typeChar, values, pc)
	if err != nil {
		return err
	}
	totalSize := uint16(size) * repeat
	fields = append(fields, &Field{"values", castValues, offset, uint64(totalSize), 0})
	offset += int64(totalSize)

	if totalSize%4 != 0 {
		padding := int64(4 - (totalSize % 4))
		if _, err := r.Seek(padding, io.SeekCurrent); err != nil {
			return err
		}
		fields = append(fields, &Field{"padding", nil, offset, uint64(padding), 0})
		offset += padding
	}

	b.size = uint64(offset - startOffset)
	b.fields = fields

	return nil
}

func (b *GenericGPMFBox) castSliceToTypeChar(typeChar string, values [][]byte, pc ParseContext) ([]interface{}, error) {
	nValues := len(values)
	castValues := make([]interface{}, nValues)

	if typeChar != "?" {
		for i, v := range values {
			castValue, err := b.castValueToTypeChar(typeChar, v)
			if err != nil {
				return nil, err
			}

			if len(castValue) == 1 {
				castValues[i] = castValue[0]
			} else {
				castValues[i] = castValue
			}
		}

		return castValues, nil
	}

	typ, ok := pc["gpmf_type"]
	if !ok {
		return nil, errors.New("Complex GPMF box, but the type has not been declared previously")
	}

	for i, v := range values {
		castValue, err := b.castValueToTypeChar(typ, v)
		if err != nil {
			return nil, err
		}
		castValues[i] = castValue
	}

	return castValues, nil
}

func (b *GenericGPMFBox) castValueToTypeChar(typeChar string, value []byte) ([]interface{}, error) {
	castValue := make([]interface{}, 0, len(typeChar))
	lenValue := len(value)

	offset := 0
	for {
		for _, typ := range typeChar {
			switch typ {
			case 'b':
				castValue = append(castValue, int8(value[offset]))
				offset += 1
			case 'B':
				castValue = append(castValue, value[offset])
				offset += 1
			case 'c': // this type is probably broken in a complex structure
				if len(typeChar) == 1 {
					castValue = append(castValue, string(value))
					offset += len(value)
				} else {
					castValue = append(castValue, string(value[offset]))
					offset += 1
				}
			case 'd':
				castValue = append(castValue, math.Float64frombits(binary.BigEndian.Uint64(value[offset:offset+8])))
				offset += 8
			case 'f':
				castValue = append(castValue, math.Float32frombits(binary.BigEndian.Uint32(value[offset:offset+4])))
				offset += 4
			case 'F':
				castValue = append(castValue, string(value[offset:offset+4]))
				offset += 4
			case 'G': // `G` is guid[16], string should work
				castValue = append(castValue, string(value[offset:offset+16]))
				offset += 16
			case 'j':
				castValue = append(castValue, int64(binary.BigEndian.Uint64(value[offset:offset+8])))
				offset += 8
			case 'J':
				castValue = append(castValue, binary.BigEndian.Uint64(value[offset:offset+8]))
				offset += 8
			case 'l':
				castValue = append(castValue, int32(binary.BigEndian.Uint32(value[offset:offset+4])))
				offset += 4
			case 'L':
				castValue = append(castValue, binary.BigEndian.Uint32(value[offset:offset+4]))
				offset += 4
			case 'q':
				castValue = append(castValue, Fixed1516{int16(binary.BigEndian.Uint16(value[offset : offset+2])), binary.BigEndian.Uint16(value[offset+2 : offset+4])})
				offset += 4
			case 'Q':
				castValue = append(castValue, Fixed3132{int32(binary.BigEndian.Uint32(value[offset : offset+4])), binary.BigEndian.Uint32(value[offset+4 : offset+8])})
				offset += 8
			case 's':
				castValue = append(castValue, int16(binary.BigEndian.Uint16(value[offset:offset+2])))
				offset += 2
			case 'S':
				castValue = append(castValue, binary.BigEndian.Uint16(value[offset:offset+2]))
				offset += 2
			case 'U':
				datetime, err := time.Parse("060102150405.999", string(value[offset:offset+16]))
				if err != nil {
					return nil, err
				}
				castValue = append(castValue, datetime)
				offset += 16
			}
		}

		if offset == lenValue {
			break
		}
	}

	return castValue, nil
}

func (b *GenericGPMFBox) Type() string {
	quoted := strconv.Quote(b._type)

	return quoted[1 : len(quoted)-1]
}

func ParseGPMFSample(r io.ReadSeeker, startOffset, sampleSize uint32) ([]*GenericGPMFBox, error) {
	endOffset := startOffset + sampleSize
	offset := startOffset
	pc := Background()

	children := make([]*GenericGPMFBox, 0, 2)
	for {
		// less than 4 bytes remaining in the sample (only 0 should happen)
		if endOffset-offset < 4 {
			break
		}

		box := &GenericGPMFBox{}

		if err := box.Parse(r, int64(offset), nil, pc); err != nil {
			if errors.Is(err, PaddingErr) {
				offset += 4

				break
			}

			return nil, err
		}

		if box.Type() == "TYPE" {
			pc["gpmf_type"] = box.Data().Get("values").Value.([]interface{})[0].(string)
		}

		offset += uint32(box.size)
		children = append(children, box)
	}

	if offset != endOffset {
		if _, err := r.Seek(int64(endOffset), io.SeekStart); err != nil {
			return nil, err
		}
	}

	return children, nil
}

type Fixed1516 struct {
	Int  int16
	Frac uint16
}

type Fixed3132 struct {
	Int  int32
	Frac uint32
}
