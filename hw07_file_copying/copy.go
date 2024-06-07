package main

import (
	"errors"
	"io"
	"log"
	"os"

	"github.com/cheggaaa/pb/v3"
)

var (
	ErrUnsupportedFile       = errors.New("unsupported file")
	ErrOffsetExceedsFileSize = errors.New("offset exceeds file size")
)

func Copy(fromPath, toPath string, offset, limit int64) error {
	srcFile, err := os.Open(fromPath)
	if err != nil {
		return err
	}
	defer closeFile(srcFile)

	fi, err := srcFile.Stat()
	if err != nil || fi.Size() <= 0 {
		return ErrUnsupportedFile
	}

	if offset > fi.Size() {
		return ErrOffsetExceedsFileSize
	}

	dstFile, err := os.Create(toPath)
	if err != nil {
		return err
	}
	defer closeFile(dstFile)

	if limit == 0 || limit > fi.Size() {
		limit = fi.Size()
	}
	if limit+offset > fi.Size() {
		limit = fi.Size() - offset
	}

	err = copyFile(srcFile, dstFile, offset, limit)
	if err != nil {
		return err
	}

	return nil
}

func copyFile(srcFile *os.File, dstFile *os.File, offset int64, limit int64) error {
	if limit <= 0 {
		return nil
	}

	_, err := srcFile.Seek(offset, 0)
	if err != nil {
		return err
	}

	bar := pb.Full.Start64(limit)
	reader := bar.NewProxyReader(io.LimitReader(srcFile, limit))

	_, err = io.Copy(dstFile, reader)
	if err != nil && !errors.Is(err, io.EOF) {
		return err
	}

	bar.Finish()

	return nil
}

func closeFile(file *os.File) {
	err := file.Close()
	if err != nil {
		log.Panic(err)
	}
}
