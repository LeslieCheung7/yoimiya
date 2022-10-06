package ds

import art "github.com/plar/go-adaptive-radix-tree"

type AdaptiveRadixTree struct {
	tree art.Tree
}

func NewART() *AdaptiveRadixTree {
	return &AdaptiveRadixTree{tree: art.New()}
}

func (ds *AdaptiveRadixTree) Size() int {
	return ds.tree.Size()
}

func (ds *AdaptiveRadixTree) Get(key []byte) interface{} {
	value, _ := ds.tree.Search(key)
	return value
}

func (ds *AdaptiveRadixTree) Put(key []byte, value interface{}) (oldVal interface{}, updated bool) {
	return ds.tree.Insert(key, value)
}

func (ds *AdaptiveRadixTree) Delete(key []byte) (val interface{}, updated bool) {
	return ds.tree.Delete(key)
}

func (ds *AdaptiveRadixTree) Iterator() art.Iterator {
	return ds.tree.Iterator()
}

func (ds *AdaptiveRadixTree) PrefixScan(prefix []byte, count int) (keys [][]byte) {
	cb := func(node art.Node) bool {
		if node.Kind() != art.Leaf {
			return true
		}
		if count <= 0 {
			return false
		}
		keys = append(keys, node.Key())
		count--
		return true
	}

	if len(prefix) == 0 {
		ds.tree.ForEach(cb)
	} else {
		ds.tree.ForEachPrefix(prefix, cb)
	}
	return
}
