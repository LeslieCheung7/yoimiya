package logfile

import (
	"fmt"
	"github.com/stretchr/testify/assert"
	"reflect"
	"sync/atomic"
	"testing"
)

func TestOpenLogFile(t *testing.T) {
	t.Run("fileIo", func(t *testing.T) {
		testOpenLogFile(t, FileIo)
	})

	t.Run("mmap", func(t *testing.T) {
		testOpenLogFile(t, MMap)
	})
}

func testOpenLogFile(t *testing.T, ioType IOType) {
	type args struct {
		path   string
		fid    uint32
		fsize  int64
		ftype  FileType
		ioType IOType
	}
	tests := []struct {
		name    string
		args    args
		wantErr bool
	}{
		{
			"zero-size", args{path: "/tmp", fid: 0, fsize: 0, ftype: List, ioType: ioType}, true,
		},
		{
			"normal-size", args{path: "/tmp", fid: 1, fsize: 100, ftype: List, ioType: ioType}, false,
		},
		{
			"big-size", args{path: "/tmp", fid: 2, fsize: 1024 << 20, ftype: List, ioType: ioType}, false,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, err := OpenLogFile(tt.args.path, tt.args.fid, tt.args.fsize, tt.args.ftype, tt.args.ioType)
			defer func() {
				if got != nil && got.IoSelector != nil {
					_ = got.Delete()
				}
			}()

			if (err != nil) != tt.wantErr {
				t.Errorf("OpenLogFile() error = %v, wantErr %v", err, tt.wantErr)
			}
			if !tt.wantErr && got == nil {
				t.Errorf("OpenLogFile() got = nil, want not nil")
			}
		})
	}
}

func TestWrite(t *testing.T) {
	t.Run("fileIo", func(t *testing.T) {
		testWrite(t, FileIo)
	})

	t.Run("mmap", func(t *testing.T) {
		testWrite(t, MMap)
	})
}

func testWrite(t *testing.T, ioType IOType) {
	lf, err := OpenLogFile("/tmp", 1, 1<<20, List, ioType)
	assert.Nil(t, err)
	defer func() {
		if lf != nil {
			_ = lf.Close()
		}
	}()

	type fields struct {
		lf *LogFile
	}
	type args struct {
		buf []byte
	}
	tests := []struct {
		name    string
		fields  fields
		args    args
		wantErr bool
	}{
		{
			"nil", fields{lf: lf}, args{buf: nil}, false,
		},
		{
			"no-value", fields{lf: lf}, args{buf: []byte{}}, false,
		},
		{
			"normal-1", fields{lf: lf}, args{buf: []byte("normal-1")}, false,
		},
		{
			"normal-2", fields{lf: lf}, args{buf: []byte("normal-2")}, false,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if err := tt.fields.lf.Write(tt.args.buf); (err != nil) != tt.wantErr {
				t.Errorf("Write() error = %v, wantErr %v", err, tt.wantErr)
			}
		})
	}
}

func TestRead(t *testing.T) {
	t.Run("fieldIo", func(t *testing.T) {
		testRead(t, FileIo)
	})

	t.Run("mmap", func(t *testing.T) {
		testRead(t, MMap)
	})
}

func testRead(t *testing.T, ioType IOType) {
	lf, err := OpenLogFile("/tmp", 1, 1<<20, List, ioType)
	assert.Nil(t, err)
	defer func() {
		if lf != nil {
			_ = lf.Close()
		}
	}()

	data := [][]byte{
		[]byte("some data 0"),
		[]byte("some data 1"),
		[]byte("some data 2"),
		[]byte("some data 3"),
		[]byte("some data 4"),
	}
	offset := writeSomeData(lf, data)

	type fields struct {
		lf *LogFile
	}
	type args struct {
		offset int64
		size   uint32
	}
	tests := []struct {
		name    string
		fields  fields
		args    args
		want    []byte
		wantErr bool
	}{
		{
			"read-0", fields{lf: lf}, args{offset: offset[0], size: uint32(len(data[0]))}, data[0], false,
		},
		{
			"read-1", fields{lf: lf}, args{offset: offset[1], size: uint32(len(data[1]))}, data[1], false,
		},
		{
			"read-2", fields{lf: lf}, args{offset: offset[2], size: uint32(len(data[2]))}, data[2], false,
		},
		{
			"read-3", fields{lf: lf}, args{offset: offset[3], size: uint32(len(data[3]))}, data[3], false,
		},
		{
			"read-4", fields{lf: lf}, args{offset: offset[4], size: uint32(len(data[4]))}, data[4], false,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, err := tt.fields.lf.Read(tt.args.offset, tt.args.size)
			if (err != nil) != tt.wantErr {
				t.Errorf("Read() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if !reflect.DeepEqual(got, tt.want) {
				t.Errorf("Read() got = %v, want %v", got, tt.want)
			}
		})
	}
}

func writeSomeData(lf *LogFile, data [][]byte) []int64 {
	var offset []int64
	for _, v := range data {
		off := atomic.LoadInt64(&lf.WriteAt)
		offset = append(offset, off)
		if err := lf.Write(v); err != nil {
			panic(fmt.Sprintf("write data err. err = %v", err))
		}
	}
	return offset
}

func TestReadLogEntry(t *testing.T) {
	t.Run("fileIo", func(t *testing.T) {
		testLogFileReadLogEntry(t, FileIo)
	})

	t.Run("mmap", func(t *testing.T) {
		testLogFileReadLogEntry(t, MMap)
	})
}

func testLogFileReadLogEntry(t *testing.T, ioType IOType) {
	lf, err := OpenLogFile("/tmp", 1, 1<<20, Strs, ioType)
	assert.Nil(t, err)
	defer func() {
		if lf != nil {
			_ = lf.Delete()
		}
	}()

	// write some entries.
	entries := []*LogEntry{
		{
			ExpiredAt: 123332, Type: 0,
		},
		{
			ExpiredAt: 123332, Type: TypeDelete,
		},
		{
			Key: []byte(""), Value: []byte(""), ExpiredAt: 994332343, Type: TypeDelete,
		},
		{
			Key: []byte("k1"), Value: nil, ExpiredAt: 7844332343,
		},
		{
			Key: nil, Value: []byte("some data"), ExpiredAt: 99400542343,
		},
		{
			Key: []byte("k2"), Value: []byte("some data"), ExpiredAt: 8847333912,
		},
		{
			Key: []byte("k3"), Value: []byte("some data"), ExpiredAt: 8847333912, Type: TypeDelete,
		},
	}
	var values [][]byte
	for _, e := range entries {
		v, _ := EncodeEntry(e)
		values = append(values, v)
	}
	offsets := writeSomeData(lf, values)

	type fields struct {
		lf *LogFile
	}
	type args struct {
		offset int64
	}
	tests := []struct {
		name    string
		fields  fields
		args    args
		want    *LogEntry
		want1   int64
		wantErr bool
	}{
		{
			"read-entry-0", fields{lf: lf}, args{offset: offsets[0]}, entries[0], int64(len(values[0])), false,
		},
		{
			"read-entry-1", fields{lf: lf}, args{offset: offsets[1]}, entries[1], int64(len(values[1])), false,
		},
		{
			"read-entry-2", fields{lf: lf}, args{offset: offsets[2]}, &LogEntry{ExpiredAt: 994332343, Type: TypeDelete}, int64(len(values[2])), false,
		},
		{
			"read-entry-3", fields{lf: lf}, args{offset: offsets[3]}, &LogEntry{Key: []byte("k1"), Value: []byte{}, ExpiredAt: 7844332343}, int64(len(values[3])), false,
		},
		{
			"read-entry-4", fields{lf: lf}, args{offset: offsets[4]}, &LogEntry{Key: []byte{}, Value: []byte("some data"), ExpiredAt: 99400542343}, int64(len(values[4])), false,
		},
		{
			"read-entry-5", fields{lf: lf}, args{offset: offsets[5]}, entries[5], int64(len(values[5])), false,
		},
		{
			"read-entry-6", fields{lf: lf}, args{offset: offsets[6]}, entries[6], int64(len(values[6])), false,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, got1, err := tt.fields.lf.ReadLogEntry(tt.args.offset)
			if (err != nil) != tt.wantErr {
				t.Errorf("ReadLogEntry() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if !reflect.DeepEqual(got, tt.want) {
				t.Errorf("ReadLogEntry() got = %v, want %v", got, tt.want)
			}
			if got1 != tt.want1 {
				t.Errorf("ReadLogEntry() got1 = %v, want %v", got1, tt.want1)
			}
		})
	}
}

func TestSync(t *testing.T) {
	sync := func(ioType IOType) {
		lf, err := OpenLogFile("/tmp", 0, 100, Hash, ioType)
		assert.Nil(t, err)
		defer func() {
			if lf != nil {
				_ = lf.Close()
			}
		}()
		err = lf.Sync()
		assert.Nil(t, err)
	}

	t.Run("fileIo", func(t *testing.T) {
		sync(FileIo)
	})

	t.Run("mmap", func(t *testing.T) {
		sync(MMap)
	})
}

func TestClose(t *testing.T) {
	closeLf := func(ioType IOType) {
		lf, err := OpenLogFile("/tmp", 0, 100, Sets, ioType)
		assert.Nil(t, err)
		defer func() {
			if lf != nil {
				_ = lf.Delete()
			}
		}()
		err = lf.Close()
		assert.Nil(t, err)
	}

	t.Run("fileIo", func(t *testing.T) {
		closeLf(FileIo)
	})

	t.Run("mmap", func(t *testing.T) {
		closeLf(MMap)
	})
}

func TestDelete(t *testing.T) {
	deleteLf := func(ioType IOType) {
		lf, err := OpenLogFile("/tmp", 0, 100, ZSet, ioType)
		assert.Nil(t, err)
		err = lf.Delete()
		assert.Nil(t, err)
	}

	t.Run("fileIo", func(t *testing.T) {
		deleteLf(FileIo)
	})

	t.Run("mmap", func(t *testing.T) {
		deleteLf(MMap)
	})
}
