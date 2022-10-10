package ds

import "testing"

// nil ------------------- cab --- abc ---------------
// nil ------------------- cab --- abc ---------------
// nil ------------------- cab --- abc ---------------
// nil ------------------- cab --- abc ----------- bca
// nil ------------------- cab --- abc --- cba --- bca
// nil ------------------- cab --- abc --- cba --- bca
// nil --- acb --- bac --- cab --- abc --- cba --- bca
func initZSet() *SortedSet {
	zSet := New()
	zSet.ZAdd("zset", 19, "abc")
	zSet.ZAdd("zset", 12, "acb")
	zSet.ZAdd("zset", 17, "bac")
	zSet.ZAdd("zset", 32, "bca")
	zSet.ZAdd("zset", 17, "cab")
	zSet.ZAdd("zset", 21, "cba")
	return zSet
}

func TestZAdd(t *testing.T) {
	zSet := initZSet()
	zSet.ZAdd("zset", 39, "aaa")
}
