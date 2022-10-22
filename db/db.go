package db

import (
	"errors"
	"math"
	"yoimiya/logfile"
)

var (
	// ErrKeyNotFound key not found.
	ErrKeyNotFound = errors.New("key not found")

	// ErrLogFileNotFound log file not found.
	ErrLogFileNotFound = errors.New("log file not found")

	// ErrWrongNumberOfArgs doesn't match key-value pair numbers.
	ErrWrongNumberOfArgs = errors.New("wrong number of arguments")

	// ErrIntegerOverflow overflows int64 limitations.
	ErrIntegerOverflow = errors.New("increment of decrement overflow")

	// ErrWrongValueType value is not a number.
	ErrWrongValueType = errors.New("value is not an integer")

	// ErrWrongIndex index is out of range.
	ErrWrongIndex = errors.New("index is out of range")

	// ErrGCRunning log file gc is running.
	ErrGCRunning = errors.New("log file gc is running, retry later")
)

const (
	logFileTypeNum   = 5
	encodeHeaderSize = 10
	initialListSeq   = math.MaxUint32 / 2
	discardFilePath  = "DISCARD"
	lockFileName     = "FLOCK"
)

type (
	YoimiyaDB struct {
		activeLogFiles   map[DataType]*logfile.LogFile
		archivedLogFiles map[DataType]archivesFiles
		fidMap           map[DataType][]uint32
		discards         map[DataType]*discard
	}

	archivesFiles map[uint32]*logfile.LogFile

	indexNode struct {
		value     []byte
		fid       uint32
		offset    int64
		entrySize int
		expiredAt int64
	}
)
