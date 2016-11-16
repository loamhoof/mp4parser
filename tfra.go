package mp4parser

import (
	"encoding/binary"
	"io"
)

type tfraBox struct {
	offset int64
	length uint32
	data   Pairs
}

func (b *tfraBox) Parse(r io.ReadSeeker) error {
	if _, err := r.Seek(b.offset, io.SeekStart); err != nil {
		return err
	}

	bytes1 := make([]byte, 1)
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

	if _, err := r.Read(bytes4); err != nil {
		return err
	}
	trackID := binary.BigEndian.Uint32(bytes4)

	if _, err := r.Read(bytes4); err != nil {
		return err
	}
	n := binary.BigEndian.Uint32(bytes4)
	lengthSizeOfTrafNum := n & 48
	lengthSizeOfTrunNum := n & 12
	lengthSizeOfSampleNum := n & 3

	if _, err := r.Read(bytes4); err != nil {
		return err
	}
	numberOfEntry := binary.BigEndian.Uint32(bytes4)

	entries := make([]Pairs, numberOfEntry)
	for i := 0; uint32(i) < numberOfEntry; i++ {
		entry := make(Pairs, 0, 5)

		if version == 1 {
			if _, err := r.Read(bytes8); err != nil {
				return err
			}
			entry = append(entry, &Pair{"time", binary.BigEndian.Uint64(bytes8)})

			if _, err := r.Read(bytes8); err != nil {
				return err
			}
			entry = append(entry, &Pair{"moof_offset", binary.BigEndian.Uint64(bytes8)})
		} else {
			if _, err := r.Read(bytes4); err != nil {
				return err
			}
			entry = append(entry, &Pair{"time", binary.BigEndian.Uint32(bytes4)})

			if _, err := r.Read(bytes4); err != nil {
				return err
			}
			entry = append(entry, &Pair{"moof_offset", binary.BigEndian.Uint32(bytes4)})
		}

		bytesTraf := make([]byte, lengthSizeOfTrafNum+1)
		if _, err := r.Read(bytesTraf); err != nil {
			return err
		}
		trafNumber := binary.BigEndian.Uint64(append(make([]byte, 8-len(bytesTraf)), bytesTraf...))
		entry = append(entry, &Pair{"traf_number", trafNumber})

		bytesTrun := make([]byte, lengthSizeOfTrunNum+1)
		if _, err := r.Read(bytesTrun); err != nil {
			return err
		}
		trunNumber := binary.BigEndian.Uint64(append(make([]byte, 8-len(bytesTrun)), bytesTrun...))
		entry = append(entry, &Pair{"trun_number", trunNumber})

		bytesSample := make([]byte, lengthSizeOfSampleNum+1)
		if _, err := r.Read(bytesSample); err != nil {
			return err
		}
		sampleNumber := binary.BigEndian.Uint64(append(make([]byte, 8-len(bytesSample)), bytesSample...))
		entry = append(entry, &Pair{"sample_number", sampleNumber})

		entries[i] = entry
	}

	b.data = Pairs{
		{"track_ID", trackID},
		{"length_size_of_traf_num", lengthSizeOfTrafNum},
		{"length_size_of_trun_num", lengthSizeOfTrunNum},
		{"length_size_of_sample_num", lengthSizeOfSampleNum},
		{"number_of_entry", numberOfEntry},
		{"entries", entries},
	}

	return nil
}

func (b *tfraBox) Type() string {
	return "tfra"
}

func (b *tfraBox) Offset() int64 {
	return b.offset
}

func (b *tfraBox) Length() uint32 {
	return b.length
}

func (b *tfraBox) Children() []Box {
	return []Box{}
}

func (b *tfraBox) Data() Pairs {
	return b.data
}
