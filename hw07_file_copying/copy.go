package main

import (
	"errors"
	"io"
	"os"

	"github.com/cheggaaa/pb/v3"
)

var (
	ErrUnsupportedFile       = errors.New("unsupported file")
	ErrOffsetExceedsFileSize = errors.New("offset exceeds file size")
	ErrFromPathIsUndefined   = errors.New("source file name is undefined")
	ErrToPathIsUndefined     = errors.New("destination file name is undefined")
	ErrOffsetIsNegative      = errors.New("offset is negative")
	ErrLimitIsNegative       = errors.New("limit is negative")
)

func ValidateParameters(fromPath, toPath string, offset, limit, size int64) error {
	if fromPath == "" {
		return ErrFromPathIsUndefined
	}
	if toPath == "" {
		return ErrToPathIsUndefined
	}
	if offset < 0 {
		return ErrOffsetIsNegative
	}
	if limit < 0 {
		return ErrLimitIsNegative
	}
	if offset > size {
		return ErrOffsetExceedsFileSize
	}

	return nil
}

func Copy(fromPath, toPath string, offset, limit int64) error {
	fileInfo, err := os.Stat(fromPath)
	if err != nil {
		return err
	}
	size := fileInfo.Size()
	if size == 0 {
		return ErrUnsupportedFile
	}

	if err = ValidateParameters(fromPath, toPath, offset, limit, size); err != nil {
		return err
	}

	source, err := os.Open(fromPath)
	if err != nil {
		return err
	}
	defer source.Close()

	if _, err = source.Seek(offset, io.SeekStart); err != nil {
		return err
	}

	count := limit
	if count == 0 || limit > size-offset {
		count = size - offset
	}

	destination, err := os.Create(toPath)
	if err != nil {
		return err
	}
	defer destination.Close()

	bar := pb.Full.Start64(count)
	defer bar.Finish()

	barReader := bar.NewProxyReader(source)

	_, err = io.CopyN(destination, barReader, count)

	if err != nil {
		return err
	}

	return nil
}
