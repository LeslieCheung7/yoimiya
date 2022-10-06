package ds

import (
	"bytes"
	"github.com/stretchr/testify/assert"
	"reflect"
	"sort"
	"testing"
)

func TestAdaptiveRadixTreeIterator(t *testing.T) {
	tree := NewART()
	iter1 := tree.Iterator()
	assert.False(t, iter1.HasNext())

	var keys = [][]byte{[]byte("one"), []byte("two"), []byte("three"), []byte("four")}
	for v, k := range keys {
		tree.Put(k, v)
	}

	sort.Slice(keys, func(i, j int) bool {
		return bytes.Compare(keys[i], keys[j]) < 0
	})
	var targets [][]byte
	iter2 := tree.Iterator()
	for iter2.HasNext() {
		node, err := iter2.Next()
		assert.Nil(t, err)
		targets = append(targets, node.Key())
	}
	assert.Equal(t, keys, targets)
}

func TestAdaptiveRadixTreeGet(t *testing.T) {
	tree := NewART()
	tree.Put(nil, nil)
	tree.Put([]byte("0"), 0)

	type args struct {
		key []byte
	}
	tests := []struct {
		name string
		tree *AdaptiveRadixTree
		args args
		want interface{}
	}{
		{
			"nil", tree, args{key: nil}, nil,
		},
		{
			"zero", tree, args{key: []byte("0")}, 0,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if flag := tt.tree.Get(tt.args.key); !reflect.DeepEqual(flag, tt.want) {
				t.Errorf("Get() = %v, want %v", flag, tt.want)
			}
		})
	}
}

func TestAdaptiveRadixTreePut(t *testing.T) {
	tree := NewART()
	type args struct {
		key   []byte
		value interface{}
	}
	tests := []struct {
		name        string
		tree        *AdaptiveRadixTree
		args        args
		wantOldVal  interface{}
		wantUpdated bool
	}{
		{
			"nil", tree, args{key: nil, value: nil}, nil, false,
		},
		{
			"normal-1", tree, args{key: []byte("1"), value: 1}, nil, false,
		},
		{
			"normal-2", tree, args{key: []byte("1"), value: 1}, 1, true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			gotOldVal, gotUpdated := tt.tree.Put(tt.args.key, tt.args.value)
			if !reflect.DeepEqual(gotOldVal, tt.wantOldVal) {
				t.Errorf("Put() gotOldVal = %v, want %v", gotOldVal, tt.wantOldVal)
			}
			if gotUpdated != tt.wantUpdated {
				t.Errorf("Put) gotUpdated = %v, want %v", gotUpdated, tt.wantUpdated)
			}
		})
	}
}

func TestAdaptiveRadixTreeDelete(t *testing.T) {
	tree := NewART()
	tree.Put(nil, nil)
	tree.Put([]byte("0"), 0)
	tree.Put([]byte("1"), 1)
	tree.Put([]byte("1"), "rewrite-data")

	type args struct {
		key []byte
	}
	tests := []struct {
		name        string
		tree        *AdaptiveRadixTree
		args        args
		wantVal     interface{}
		wantUpdated bool
	}{
		{
			"nil", tree, args{key: nil}, nil, false,
		},
		{
			"zero", tree, args{key: []byte("0")}, 0, true,
		},
		{
			"rewrite-data", tree, args{key: []byte("1")}, "rewrite-data", true,
		},
		{
			"not-exist", tree, args{key: []byte("not-exist")}, nil, false,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			gotVal, gotUpdated := tt.tree.Delete(tt.args.key)
			if !reflect.DeepEqual(gotVal, tt.wantVal) {
				t.Errorf("Delete() gotOldVal = %v, want %v", gotVal, tt.wantVal)
			}
			if gotUpdated != tt.wantUpdated {
				t.Errorf("Delete() gotUpdated = %v, want %v", gotUpdated, tt.wantUpdated)
			}
		})
	}
}

func TestAdaptiveRadixTreePrefixScan(t *testing.T) {
	tree := NewART()
	tree.Put([]byte("aa"), 1)
	tree.Put([]byte("ab"), 2)
	tree.Put([]byte("ac"), 3)
	tree.Put([]byte("ad"), 4)
	tree.Put([]byte("ae"), 5)

	keys1 := tree.PrefixScan([]byte("a"), -1)
	assert.Equal(t, 0, len(keys1))

	keys2 := tree.PrefixScan([]byte("a"), 0)
	assert.Equal(t, 0, len(keys2))

	keys3 := tree.PrefixScan([]byte("a"), 1)
	assert.Equal(t, 1, len(keys3))

	keys4 := tree.PrefixScan([]byte("a"), 5)
	assert.Equal(t, 5, len(keys4))
}