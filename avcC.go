package mp4parser

import (
	"encoding/binary"
	"errors"
	"io"
)

type AvcCBox struct {
	baseBox
}

func (b *AvcCBox) Parse(r io.ReadSeeker, startOffset int64, pp ParsePlan, pc ParseContext) error {
	size, offset, _, fields, err := parseBox(r, startOffset)
	if err != nil {
		return err
	}
	b.size = size
	b.fields = fields

	b1 := make([]byte, 1)
	b2 := make([]byte, 2)

	if _, err := r.Read(b1); err != nil {
		return err
	}
	configurationVersion := b1[0]
	fields = append(fields, &Field{"configurationVersion", configurationVersion, offset, 8, 0})
	offset += 1

	if configurationVersion != 1 {
		return errors.New("Unsupported configuration version")
	}

	if _, err := r.Read(b1); err != nil {
		return err
	}
	fields = append(fields, &Field{"AVCProfileIndication", b1[0], offset, 8, 0})
	offset += 1

	if _, err := r.Read(b1); err != nil {
		return err
	}
	fields = append(fields, &Field{"profile_compatibility", b1[0], offset, 8, 0})
	offset += 1

	if _, err := r.Read(b1); err != nil {
		return err
	}
	fields = append(fields, &Field{"AVCLevelIndication", b1[0], offset, 8, 0})
	offset += 1

	if _, err := r.Read(b1); err != nil {
		return err
	}
	fields = append(fields, &Field{"lengthSizeMinusOne", b1[0] & 0x3, offset, 2, 6})
	offset += 1

	if _, err := r.Read(b1); err != nil {
		return err
	}
	numOfSequenceParameterSets := b1[0] & 0x1F
	fields = append(fields, &Field{"numOfSequenceParameterSets", numOfSequenceParameterSets, offset, 5, 3})
	offset += 1

	sequenceParameterSetsOffset := offset
	sequenceParameterSets := make([]Fields, numOfSequenceParameterSets)
	for i := 0; i < int(numOfSequenceParameterSets); i++ {
		sequenceParameterSet := make(Fields, 0, 2)

		if _, err := r.Read(b2); err != nil {
			return err
		}
		sequenceParameterSetLength := binary.BigEndian.Uint16(b2)
		sequenceParameterSet = append(sequenceParameterSet, &Field{"sequenceParameterSetLength", sequenceParameterSetLength, offset, 16, 0})
		offset += 2

		bLength := make([]byte, sequenceParameterSetLength)

		if _, err := r.Read(bLength); err != nil {
			return err
		}
		sequenceParameterSet = append(sequenceParameterSet, &Field{"sequenceParameterSetNALUnit", bLength, offset, 8 * uint64(sequenceParameterSetLength), 0})
		offset += int64(sequenceParameterSetLength)

		sequenceParameterSets[i] = sequenceParameterSet
	}
	b.fields = append(b.fields, &Field{"sequenceParameterSets", sequenceParameterSets, sequenceParameterSetsOffset, uint64(offset-sequenceParameterSetsOffset) * 8, 0})

	if _, err := r.Read(b1); err != nil {
		return err
	}
	numOfPictureParameterSets := b1[0]
	fields = append(fields, &Field{"numOfPictureParameterSets", numOfPictureParameterSets, offset, 8, 0})
	offset += 1

	pictureParameterSetsOffset := offset
	pictureParameterSets := make([]Fields, numOfPictureParameterSets)
	for i := 0; i < int(numOfPictureParameterSets); i++ {
		pictureParameterSet := make(Fields, 0, 2)

		if _, err := r.Read(b2); err != nil {
			return err
		}
		pictureParameterSetLength := binary.BigEndian.Uint16(b2)
		pictureParameterSet = append(pictureParameterSet, &Field{"pictureParameterSetLength", pictureParameterSetLength, offset, 16, 0})
		offset += 2

		bLength := make([]byte, pictureParameterSetLength)

		if _, err := r.Read(bLength); err != nil {
			return err
		}
		pictureParameterSet = append(pictureParameterSet, &Field{"pictureParameterSetNALUnit", bLength, offset, 8 * uint64(pictureParameterSetLength), 0})
		offset += int64(pictureParameterSetLength)

		pictureParameterSets[i] = pictureParameterSet
	}
	b.fields = append(b.fields, &Field{"pictureParameterSets", pictureParameterSets, pictureParameterSetsOffset, uint64(offset-pictureParameterSetsOffset) * 8, 0})

	return nil
}

func (b *AvcCBox) Type() string {
	return "avcC"
}
