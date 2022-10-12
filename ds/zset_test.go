package ds

import (
	"github.com/stretchr/testify/assert"
	"testing"
)

func TestZAdd(t *testing.T) {
	zSet := initZSet()
	zSet.ZAdd("zset", 39, "aaa")

	c := zSet.ZCard("zset")
	assert.Equal(t, 7, c)
}

func TestZScore(t *testing.T) {
	zSet := initZSet()
	ok, s1 := zSet.ZScore("zset", "abc")
	assert.Equal(t, true, ok)
	assert.Equal(t, float64(19), s1)

	ok, s2 := zSet.ZScore("zset", "aaa")
	assert.Equal(t, false, ok)
	assert.Equal(t, float64(0), s2)
}

func TestZRank(t *testing.T) {
	key := "zset"
	zSet := initZSet()
	r1 := zSet.ZRank(key, "acb")
	assert.Equal(t, int64(0), r1)

	r2 := zSet.ZRank(key, "bac")
	assert.Equal(t, int64(1), r2)

	r3 := zSet.ZRank(key, "not exist")
	assert.Equal(t, int64(-1), r3)
}

func TestZRevRank(t *testing.T) {
	key := "zset"
	zSet := initZSet()
	r1 := zSet.ZRevRank(key, "acb")
	assert.Equal(t, int64(5), r1)

	r2 := zSet.ZRevRank(key, "bac")
	assert.Equal(t, int64(4), r2)

	r3 := zSet.ZRevRank(key, "not exist")
	assert.Equal(t, int64(-1), r3)
}

func TestZIncrBy(t *testing.T) {
	key := "zset"
	zSet := initZSet()

	incr1 := zSet.ZIncrBy(key, 300, "acb")
	assert.Equal(t, float64(312), incr1)
	r1 := zSet.ZRank(key, "acb")
	assert.Equal(t, int64(5), r1)

	incr2 := zSet.ZIncrBy(key, 300, "bac")
	assert.Equal(t, float64(317), incr2)
	r2 := zSet.ZRank(key, "bac")
	assert.Equal(t, int64(5), r2)
}

func TestZRange(t *testing.T) {
	key := "zset"
	zSet := initZSet()

	r := zSet.ZRange(key, 0, 5)
	assert.Equal(t, 6, len(r))

	for _, n := range r {
		assert.NotNil(t, n)
	}
}

func TestZRangeWithScores(t *testing.T) {
	key := "zset"
	zSet := initZSet()

	r := zSet.ZRangeWithScores(key, 0, 5)
	assert.Equal(t, 12, len(r))

	for _, n := range r {
		assert.NotNil(t, n)
	}
}

func TestZRevRange(t *testing.T) {
	key := "zset"
	zSet := initZSet()

	r := zSet.ZRevRange(key, 0, 5)
	assert.Equal(t, 6, len(r))

	for _, n := range r {
		assert.NotNil(t, n)
	}
}

func TestZRevRangeWithScores(t *testing.T) {
	key := "zset"
	zSet := initZSet()

	r := zSet.ZRevRangeWithScores(key, 0, 5)
	assert.Equal(t, 12, len(r))

	for _, n := range r {
		assert.NotNil(t, n)
	}
}

func TestZRem(t *testing.T) {
	key := "zset"
	zSet := initZSet()

	ok1 := zSet.ZRem(key, "acb")
	assert.Equal(t, true, ok1)
	ok1, _ = zSet.ZScore(key, "acb")
	assert.Equal(t, false, ok1)

	ok2 := zSet.ZRem(key, "not exists")
	assert.Equal(t, false, ok2)
}

func TestZGetByRank(t *testing.T) {
	key := "zset"
	zSet := initZSet()

	getRank := func(rank int) {
		val := zSet.ZGetByRank(key, rank)
		if val != nil {
			for _, n := range val {
				assert.NotNil(t, n)
			}
		}
	}

	getRank(0)
	getRank(4)
	getRank(5)
}

func TestZRevGetByRank(t *testing.T) {
	key := "zset"
	zSet := initZSet()

	getRevRank := func(rank int) {
		val := zSet.ZRevGetByRank(key, rank)
		if val != nil {
			for _, n := range val {
				assert.NotNil(t, n)
			}
		}
	}

	getRevRank(5)
	getRevRank(4)
	getRevRank(0)
}

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
