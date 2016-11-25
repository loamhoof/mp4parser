package mp4parser

import (
	"encoding/binary"
	"io"
)

func parseBox(r io.ReadSeeker, startOffset int64) (size uint64, offset int64, _type string, fields Fields, err error) {
	if _, err = r.Seek(startOffset, io.SeekStart); err != nil {
		return
	}

	fields = make(Fields, 0, 4)

	offset, err = r.Seek(0, io.SeekCurrent)
	if err != nil {
		return
	}

	b4 := make([]byte, 4)
	b8 := make([]byte, 8)
	b16 := make([]byte, 16)

	if _, err = r.Read(b4); err != nil {
		return
	}
	size = uint64(binary.BigEndian.Uint32(b4))
	fields = append(fields, &Field{"size", size, offset, 32})
	offset += 4

	if _, err = r.Read(b4); err != nil {
		return
	}
	_type = string(b4)
	fields = append(fields, &Field{"type", _type, offset, 32})
	offset += 4

	if size == 1 {
		if _, err = r.Read(b8); err != nil {
			return
		}
		size = binary.BigEndian.Uint64(b8)
		fields = append(fields, &Field{"largesize", size, offset, 64})
		offset += 8
	} else if size == 0 {
		var fileSize int64
		fileSize, err = r.Seek(0, io.SeekEnd)
		if err != nil {
			return
		}
		size = uint64(fileSize - offset - 8)
		_, err = r.Seek(offset, io.SeekStart)
		if err != nil {
			return
		}
	}

	if _type == "uuid" {
		if _, err = r.Read(b16); err != nil {
			return
		}
		fields = append(fields, &Field{"usertype", string(b16), offset, 128})
	}

	return

}

func parseFullBox(r io.ReadSeeker, startOffset int64) (size uint64, offset int64, _type string, version byte, flags uint32, fields Fields, err error) {
	size, offset, _type, fields, err = parseBox(r, startOffset)
	if err != nil {
		return
	}

	b4 := make([]byte, 4)

	if _, err = r.Read(b4); err != nil {
		return
	}
	version = b4[0]
	flags = binary.BigEndian.Uint32(b4) & 0xFFFFFF
	offset += 4

	fields = append(fields, &Field{"version", version, offset, 8})
	fields = append(fields, &Field{"flags", flags, offset, 24})

	return
}
