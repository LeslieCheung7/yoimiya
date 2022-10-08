package logfile

import (
	"reflect"
	"testing"
)

func TestEncodeEntry(t *testing.T) {
	type args struct {
		e *LogEntry
	}
	tests := []struct {
		name  string
		args  args
		want1 []byte
		want2 int
	}{
		{
			"nil", args{e: nil}, nil, 0,
		},
		{
			"no-fields", args{e: &LogEntry{}}, []byte{28, 223, 68, 33, 0, 0, 0, 0}, 8,
		},
		{
			"no-key-value", args{e: &LogEntry{ExpiredAt: 443434211}},
			[]byte{51, 97, 150, 123, 0, 0, 0, 198, 147, 242, 166, 3}, 12,
		},
		{
			"with-key-value", args{e: &LogEntry{Key: []byte("kv"), Value: []byte("kv"), ExpiredAt: 443434211}},
			[]byte{17, 125, 5, 80, 0, 4, 4, 198, 147, 242, 166, 3, 107, 118, 107, 118},
			16,
		},
		{
			"type-delete", args{e: &LogEntry{Key: []byte("kv"), Value: []byte("kv"), ExpiredAt: 443434211, Type: TypeDelete}},
			[]byte{126, 49, 160, 203, 1, 4, 4, 198, 147, 242, 166, 3, 107, 118, 107, 118},
			16,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got1, got2 := EncodeEntry(tt.args.e)
			if !reflect.DeepEqual(got1, tt.want1) {
				t.Errorf("EncodeEntry() got1 = %v, want1 %v", got1, tt.want1)
			}
			if got2 != tt.want2 {
				t.Errorf("EncodeEntry() got2 = %v, want2 %v", got2, tt.want2)
			}
		})
	}
}
