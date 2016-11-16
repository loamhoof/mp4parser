package mp4parser

import (
	"encoding/binary"
	"io"
)

type mp4Box struct {
	children []Box
}

func (b *mp4Box) Parse(r io.ReadSeeker) error {
	if closer, ok := r.(io.Closer); ok {
		defer closer.Close()
	}

	var offset int64
	children := make([]Box, 0, 1)
	for {
		if _, err := r.Seek(offset, io.SeekStart); err != nil {
			break
		}

		bytes := make([]byte, 4)

		if _, err := r.Read(bytes); err != nil {
			if err == io.EOF {
				break
			}
			return err
		}
		l := binary.BigEndian.Uint32(bytes)

		if _, err := r.Read(bytes); err != nil {
			return err
		}
		boxType := string(bytes)

		box := newBox(boxType, offset)
		if err := box.Parse(r); err != nil {
			return err
		}
		children = append(children, box)

		offset += int64(l)
	}

	b.children = children

	return nil
}

func (b *mp4Box) Type() string {
	return "mp4"
}

func (b *mp4Box) Offset() int64 {
	return 0
}

func (b *mp4Box) Length() uint32 {
	return 0
}

func (b *mp4Box) Children() []Box {
	return b.children
}

func (b *mp4Box) Data() Pairs {
	return nil
}
