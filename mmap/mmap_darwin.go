package mmap

import (
	"golang.org/x/sys/unix"
	"os"
	"syscall"
	"unsafe"
)

func mmap(fd *os.File, writable bool, size int64) ([]byte, error) {
	typ := unix.PROT_READ
	if writable {
		typ |= unix.PROT_WRITE
	}
	return unix.Mmap(int(fd.Fd()), 0, int(size), typ, unix.MAP_SHARED)
}

func munmap(b []byte) error {
	return unix.Munmap(b)
}

// madvise is required because the unix package does not support the madvise system call on OS X.
func madvise(b []byte, readAhead bool) error {
	advice := unix.MADV_NORMAL
	if !readAhead {
		advice = unix.MADV_RANDOM
	}

	_, _, err := syscall.Syscall(syscall.SYS_MADVISE, uintptr(unsafe.Pointer(&b[0])), uintptr(len(b)), uintptr(advice))
	if err != 0 {
		return err
	}
	return nil
}

func msync(b []byte) error {
	return unix.Msync(b, unix.MS_SYNC)
}
