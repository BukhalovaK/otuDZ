package main

import (
	"errors"
	"io"
	"os"

	pb "github.com/cheggaaa/pb/v3"
)

var (
	ErrUnsupportedFile       = errors.New("unsupported file")
	ErrOffsetExceedsFileSize = errors.New("offset exceeds file size")
)

func Copy(fromPath, toPath string, offset, limit int64) error {
	fromFileInfo, err := os.Stat(fromPath)
	if err != nil {
		return err
	}

	fromFileSize := fromFileInfo.Size()

	if fromFileSize == 0 {
		return ErrUnsupportedFile
	}

	if offset > fromFileSize {
		return ErrOffsetExceedsFileSize
	}

	fileFrom, err := os.OpenFile(fromPath, os.O_RDONLY, 0o644)
	if err != nil {
		return err
	}
	defer fileFrom.Close()
	_, err = fileFrom.Seek(offset, 0)
	if err != nil {
		return err
	}

	fileTo, err := os.Create(toPath)
	if err != nil {
		return err
	}
	defer fileTo.Close()

	limitBar := fromFileSize - offset
	if limit != 0 && limit < limitBar {
		limitBar = limit
	}
	bar := pb.Full.Start64(limitBar)
	barReader := bar.NewProxyReader(fileFrom)

	_, err = io.CopyN(fileTo, barReader, limitBar)

	bar.Finish()

	if err != nil && errors.Is(err, io.EOF) {
		return err
	}

	return nil
}
