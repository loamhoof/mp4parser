package mp4parser

import (
	"encoding/binary"
	"io"
)

type tkhdBox struct {
	offset int64
	length uint32
	data   Pairs
}

func (b *tkhdBox) Parse(r io.ReadSeeker) error {
	if _, err := r.Seek(b.offset, io.SeekStart); err != nil {
		return err
	}

	bytes1 := make([]byte, 1)
	bytes2 := make([]byte, 2)
	bytes4 := make([]byte, 4)
	bytes8 := make([]byte, 8)

	if _, err := r.Read(bytes4); err != nil {
		return err
	}
	l := binary.BigEndian.Uint32(bytes4)

	b.length = l

	if _, err := r.Seek(4, io.SeekCurrent); err != nil {
		return err
	}

	if _, err := r.Read(bytes1); err != nil {
		return err
	}
	version := bytes1[0]

	if _, err := r.Seek(3, io.SeekCurrent); err != nil {
		return err
	}

	b.data = make(Pairs, 0, 10)

	if version == 1 {
		if _, err := r.Read(bytes8); err != nil {
			return err
		}
		b.data = append(b.data, &Pair{"creation_time", binary.BigEndian.Uint64(bytes8)})

		if _, err := r.Read(bytes8); err != nil {
			return err
		}
		b.data = append(b.data, &Pair{"modification_time", binary.BigEndian.Uint64(bytes8)})

		if _, err := r.Read(bytes4); err != nil {
			return err
		}
		b.data = append(b.data, &Pair{"timescale", binary.BigEndian.Uint32(bytes4)})

		if _, err := r.Seek(4, io.SeekCurrent); err != nil {
			return err
		}

		if _, err := r.Read(bytes8); err != nil {
			return err
		}
		b.data = append(b.data, &Pair{"duration", binary.BigEndian.Uint64(bytes8)})
	} else {
		if _, err := r.Read(bytes4); err != nil {
			return err
		}
		b.data = append(b.data, &Pair{"creation_time", binary.BigEndian.Uint32(bytes4)})

		if _, err := r.Read(bytes4); err != nil {
			return err
		}
		b.data = append(b.data, &Pair{"modification_time", binary.BigEndian.Uint32(bytes4)})

		if _, err := r.Read(bytes4); err != nil {
			return err
		}
		b.data = append(b.data, &Pair{"timescale", binary.BigEndian.Uint32(bytes4)})

		if _, err := r.Seek(4, io.SeekCurrent); err != nil {
			return err
		}

		if _, err := r.Read(bytes4); err != nil {
			return err
		}
		b.data = append(b.data, &Pair{"duration", binary.BigEndian.Uint32(bytes4)})
	}

	if _, err := r.Seek(8, io.SeekCurrent); err != nil {
		return err
	}

	if _, err := r.Read(bytes2); err != nil {
		return err
	}
	b.data = append(b.data, &Pair{"layer", binary.BigEndian.Uint16(bytes2)})

	if _, err := r.Read(bytes2); err != nil {
		return err
	}
	b.data = append(b.data, &Pair{"alternate_group", binary.BigEndian.Uint16(bytes2)})

	if _, err := r.Read(bytes2); err != nil {
		return err
	}
	b.data = append(b.data, &Pair{"volume", binary.BigEndian.Uint16(bytes2)})

	if _, err := r.Seek(2, io.SeekCurrent); err != nil {
		return err
	}

	var matrix [9]uint32
	for i := 0; i < 9; i++ {
		if _, err := r.Read(bytes4); err != nil {
			return err
		}
		matrix[i] = binary.BigEndian.Uint32(bytes4)
	}
	b.data = append(b.data, &Pair{"matrix", matrix})

	if _, err := r.Read(bytes4); err != nil {
		return err
	}
	b.data = append(b.data, &Pair{"width", binary.BigEndian.Uint32(bytes4)})

	if _, err := r.Read(bytes4); err != nil {
		return err
	}
	b.data = append(b.data, &Pair{"height", binary.BigEndian.Uint32(bytes4)})

	return nil
}

func (b *tkhdBox) Type() string {
	return "tkhd"
}

func (b *tkhdBox) Offset() int64 {
	return b.offset
}

func (b *tkhdBox) Length() uint32 {
	return b.length
}

func (b *tkhdBox) Children() []Box {
	return []Box{}
}

func (b *tkhdBox) Data() Pairs {
	return b.data
}
