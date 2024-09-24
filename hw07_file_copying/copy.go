package main

import (
	"errors"
	"fmt"
	"io"
	"os"

	"github.com/cheggaaa/pb/v3"
)

var (
	ErrUnsupportedFile       = errors.New("unsupported file")
	ErrOffsetExceedsFileSize = errors.New("offset exceeds file size")
	ErrFromPathRequired      = errors.New("fromPath required")
	ErrToPathRequired        = errors.New("toPath required")
	ErrSameFile              = errors.New("can copy to source file")
	ErrLimitLessZero         = errors.New("limit cannot be less than zero")
	ErrOffsetLessZero        = errors.New("offset cannot be less than zero")
	ErrFileNotFound          = errors.New("file not found")
	ErrOpenFile              = errors.New("cant open file")
	ErrReadFileState         = errors.New("cant read file stat")
	ErrFileSeek              = errors.New("cant seek file")
	ErrCreateFile            = errors.New("cant create file")
	ErrCopyFile              = errors.New("cant copy to file")
)

func Copy(fromPath, toPath string, offset, limit int64) error {
	if fromPath == "" {
		return ErrFromPathRequired
	}
	if toPath == "" {
		return ErrToPathRequired
	}
	if fromPath == toPath {
		return ErrSameFile
	}
	if offset < 0 {
		return ErrOffsetLessZero
	}
	if limit < 0 {
		return ErrLimitLessZero
	}

	file, err := os.OpenFile(fromPath, os.O_RDONLY, 0o644)
	if err != nil {
		if os.IsNotExist(err) {
			return ErrFileNotFound
		}
		return ErrOpenFile
	}
	defer file.Close()

	fileStat, err := file.Stat()
	if err != nil {
		return ErrReadFileState
	}

	fileSize := fileStat.Size()
	if offset > fileSize {
		return ErrOffsetExceedsFileSize
	}

	if limit == 0 {
		limit = fileSize - offset
	}

	bar := pb.Full.Start64(limit)
	defer bar.Finish()

	_, err = file.Seek(offset, 0)
	if err != nil {
		return ErrFileSeek
	}

	outFile, err := os.Create(toPath)
	if err != nil {
		fmt.Println(toPath)
		return ErrCreateFile
	}
	defer outFile.Close()

	barReader := bar.NewProxyReader(file)
	_, err = io.CopyN(outFile, barReader, limit)
	if err != nil {
		if errors.Is(err, io.EOF) {
			return nil
		}
		return ErrCopyFile
	}

	return nil
}
