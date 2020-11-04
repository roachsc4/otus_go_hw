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
)

func Copy(fromPath string, toPath string, offset, limit int64) error {
	fromFileInfo, err := os.Stat(fromPath)
	if err != nil {
		return fmt.Errorf("failed to obtain source file info: %w", err)
	}

	if fromFileInfo.IsDir() {
		return ErrUnsupportedFile
	}

	if offset > fromFileInfo.Size() {
		return ErrOffsetExceedsFileSize
	}

	// If there is no limit, then we should write all source file's content excluding offset
	needToWrite := fromFileInfo.Size() - offset
	// But if there is limit and it is lesser than source file's content, we should write up to `limit` bytes
	if limit > 0 && needToWrite > limit {
		needToWrite = limit
	}

	fromFile, err := os.Open(fromPath)
	if err != nil {
		return fmt.Errorf("failed to open -from file: %w", err)
	}
	defer fromFile.Close()

	bar := pb.Full.Start(int(needToWrite))

	// Create source file reader
	fromFileReader := io.LimitReader(fromFile, needToWrite)
	// Set offset
	_, err = fromFile.Seek(offset, 0)
	if err != nil {
		return fmt.Errorf("failed to set offset to -from file: %w", err)
	}
	// Create target file
	targetFile, err := os.Create(toPath)
	if err != nil {
		return fmt.Errorf("failed to create target file: %w", err)
	}
	defer targetFile.Close()
	// Create target file writer connected to progress-bar
	toFileWriter := bar.NewProxyWriter(targetFile)
	// Copy content from source file to target file
	_, err = io.CopyN(toFileWriter, fromFileReader, needToWrite)
	if err != nil {
		return fmt.Errorf("failed to copy content from source to target file: %w", err)
	}
	bar.Finish()
	return nil
}
