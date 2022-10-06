package logfile

import (
	"encoding/binary"
	"hash/crc32"
)

// MaxHeaderSize max entry header size.
// crc32      type     kSize     vSize    expiredAt
//   4    +    1    +    5    +    5    +    10    =    25
const MaxHeaderSize = 25

// EntryType type of entry
type EntryType byte

const (
	// TypeDelete represents entry type is delete.
	TypeDelete EntryType = iota + 1

	// TypeListMeta represents entry is list meta.
	TypeListMeta
)

// LogEntry is the data will be appended in log file.
type LogEntry struct {
	Key       []byte
	Value     []byte
	ExpiredAt int64
	Type      EntryType
}

type entryHeader struct {
	crc32     uint32
	typ       EntryType
	kSize     uint32
	vSize     uint32
	expiredAt int64
}

// EncodeEntry will encode entry into a byte slice.
// The encoded Entry look like:
// +-------+--------+------------+--------------+------------+-------+---------+
// |  crc  |  type  |  key size  |  value size  | expiredAt  |  key  |  value  |
// +-------+--------+------------+--------------+------------+-------+---------+
// |-------------------------HEADER--------------------------|
//         |-------------------------------crc check---------------------------|
func EncodeEntry(e *LogEntry) ([]byte, int) {
	if e == nil {
		return nil, 0
	}
	header := make([]byte, MaxHeaderSize)
	header[4] = byte(e.Type)
	var index = 5
	index += binary.PutVarint(header[index:], int64(len(e.Key)))
	index += binary.PutVarint(header[index:], int64(len(e.Value)))
	index += binary.PutVarint(header[index:], e.ExpiredAt)

	var size = index + len(e.Key) + len(e.Value)
	buf := make([]byte, size)
	// header.
	copy(buf[:index], header[:])
	// key.
	copy(buf[index:], e.Key)
	// value.
	copy(buf[index+len(e.Key):], e.Value)

	crc := crc32.ChecksumIEEE(buf[4:])
	binary.LittleEndian.PutUint32(buf[:4], crc)
	return buf, size
}

// decodeHeader returns the entry and entry's size.
func decodeHeader(buf []byte) (*entryHeader, int64) {
	if len(buf) <= 4 {
		return nil, 0
	}
	header := &entryHeader{
		crc32: binary.LittleEndian.Uint32(buf[:4]),
		typ:   EntryType(buf[4]),
	}
	var index = 5
	keySize, n := binary.Varint(buf[index:])
	header.kSize = uint32(keySize)
	index += n

	valueSize, n := binary.Varint(buf[index:])
	header.vSize = uint32(valueSize)
	index += n

	expiredAt, n := binary.Varint(buf[index:])
	header.expiredAt = expiredAt
	return header, int64(index + n)
}

func getEntryCrc(e *LogEntry, header []byte) uint32 {
	if e == nil {
		return 0
	}
	crc := crc32.ChecksumIEEE(header[:])
	crc = crc32.Update(crc, crc32.IEEETable, e.Key)
	crc = crc32.Update(crc, crc32.IEEETable, e.Value)
	return crc
}
