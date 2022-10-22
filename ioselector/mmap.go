package ioselector

import (
	"io"
	"os"
	"yoimiya/mmap"
)

// MMapSelector represents using memory-mapped file I/O
type MMapSelector struct {
	fd     *os.File
	buf    []byte
	bufLen int64
}

func NewMMapSelector(fname string, fsize int64) (IOSelector, error) {
	if fsize <= 0 {
		return nil, ErrInvalidFsize
	}
	file, err := openFile(fname, fsize)
	if err != nil {
		return nil, err
	}
	buf, err := mmap.Mmap(file, true, fsize)
	if err != nil {
		return nil, err
	}
	return &MMapSelector{fd: file, buf: buf, bufLen: int64(len(buf))}, nil
}

// Write copy slice b into mapped region(buf) at offset.
func (ms *MMapSelector) Write(b []byte, offset int64) (int, error) {
	length := int64(len(b))
	if length <= 0 {
		return 0, nil
	}
	if offset < 0 || length+offset > ms.bufLen {
		return 0, io.EOF
	}
	return copy(ms.buf[offset:], b), nil
}

// Read copy data from mapped region(buf) into slice b at offset.
func (ms *MMapSelector) Read(b []byte, offset int64) (int, error) {
	if offset < 0 || offset >= ms.bufLen {
		return 0, io.EOF
	}
	if offset+int64(len(b)) >= ms.bufLen {
		return 0, io.EOF
	}
	return copy(b, ms.buf[offset:]), nil
}

// Sync synchronize the mapped buffer to the file's contents on disk.
func (ms *MMapSelector) Sync() error {
	return mmap.Msync(ms.buf)
}

// Close sync/unmap mapped buffer and close fd.
func (ms *MMapSelector) Close() error {
	if err := mmap.Msync(ms.buf); err != nil {
		return nil
	}
	if err := mmap.Munmap(ms.buf); err != nil {
		return nil
	}
	return ms.fd.Close()
}

func (ms *MMapSelector) Delete() error {
	if err := mmap.Munmap(ms.buf); err != nil {
		return nil
	}
	ms.buf = nil

	if err := ms.fd.Truncate(0); err != nil {
		return nil
	}
	if err := ms.fd.Close(); err != nil {
		return nil
	}
	return os.Remove(ms.fd.Name())
}
