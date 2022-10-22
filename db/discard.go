package db

import (
	"errors"
	"sync"
	"yoimiya/ioselector"
)

const (
	discardRecordSize       = 12
	discardFileSize   int64 = 2 << 12 // 8kb, contains mostly 682 records in file.
	discardFileName         = "discard"
)

// ErrDiscardNoSpace no enough space for discard file.
var ErrDiscardNoSpace = errors.New("not enough space can be allocated for the discard file")

type discard struct {
	sync.Mutex
	once     *sync.Once
	valChan  chan *indexNode
	file     ioselector.IOSelector
	freeList []int64          // contains file offset that can be allocated.
	location map[uint32]int64 // offset of each fid.
}
