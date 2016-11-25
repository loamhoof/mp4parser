package mp4parser

import (
	"io"
)

func Parse(r io.ReadSeeker) (*MP4, error) {
	mp4 := &MP4{}

	return mp4, mp4.Parse(r, 0)
}
