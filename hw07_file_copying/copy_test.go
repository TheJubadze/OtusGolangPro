package main

import (
	"bytes"
	"fmt"
	"os"
	"testing"

	"github.com/stretchr/testify/require"
)

func TestRandom(t *testing.T) {
	const fromName = "/dev/urandom"
	to, err := os.CreateTemp(".", "tmp.txt")
	if err != nil {
		t.Fatal(err)
	}
	defer func() {
		if err := to.Close(); err != nil {
			t.Fatal(err)
		}
		if err := os.Remove(to.Name()); err != nil {
			t.Fatal(err)
		}
	}()

	err = Copy(fromName, to.Name(), offset, limit)

	require.Equal(t, ErrUnsupportedFile, err, "wrong or no error returned")
}

func TestCopy(t *testing.T) {
	tests := []struct {
		offset int64
		limit  int64
		err    error
	}{
		{0, 0, nil},
		{0, 10, nil},
		{0, 1000, nil},
		{0, 10000, nil},
		{100, 1000, nil},
		{6000, 1000, nil},
		{60000, 1000, ErrOffsetExceedsFileSize},
	}
	for _, tt := range tests {
		t.Run("", func(t *testing.T) {
			testCopy(t, tt.offset, tt.limit, tt.err)
		})
	}
}

func testCopy(t *testing.T, offset, limit int64, expectedErr error) {
	t.Helper()
	const fromName = "./testdata/input.txt"
	to, err := os.CreateTemp(".", "tmp.txt")
	if err != nil {
		t.Fatal(err)
	}
	defer func() {
		if err := to.Close(); err != nil {
			t.Fatal(err)
		}
		if err := os.Remove(to.Name()); err != nil {
			t.Fatal(err)
		}
	}()
	expectedFilePath := fmt.Sprintf("./testdata/out_offset%d_limit%d.txt", offset, limit)

	err = Copy(fromName, to.Name(), offset, limit)

	f1, err1 := os.ReadFile(to.Name())
	f2, err2 := os.ReadFile(expectedFilePath)
	require.Equal(t, expectedErr, err, "copy failed")
	require.Nil(t, err1, "actual file read failed")
	require.Nil(t, err2, "expected file read failed")
	require.True(t, bytes.Equal(f1, f2), fmt.Sprintf("expected and actual files are not equal: %s", expectedFilePath))
}
