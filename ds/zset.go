package ds

import "math/rand"

// zset is the implementation of sorted set.

const (
	maxLevel    = 32
	probability = 0.25
)

type (
	// SortedSet sorted set struct.
	SortedSet struct {
		record map[string]*SortedSetNode
	}

	// SortedSetNode node of sorted set.
	SortedSetNode struct {
		dict map[string]*sklNode
		skl  *skipList
	}

	sklLevel struct {
		forward *sklNode
		span    uint64
	}

	sklNode struct {
		member   string
		score    float64
		backward *sklNode
		level    []*sklLevel
	}

	skipList struct {
		head   *sklNode
		tail   *sklNode
		length int64
		level  int16
	}
)

// New create a new sorted set.
func New() *SortedSet {
	return &SortedSet{record: make(map[string]*SortedSetNode)}
}

// ZAdd adds the specified member with the specified score to the sorted set stored at key.
func (z *SortedSet) ZAdd(key string, score float64, member string) {
	if !z.exist(key) {
		node := &SortedSetNode{
			dict: make(map[string]*sklNode),
			skl:  newSkipList(),
		}
		z.record[key] = node
	}

	item := z.record[key]
	v, exist := item.dict[member]

	var node *sklNode
	if exist {
		if score != v.score {
			item.skl.sklDelete(v.score, member)
			node = item.skl.sklInsert(score, member)
		}
	} else {
		node = item.skl.sklInsert(score, member)
	}

	if node != nil {
		item.dict[member] = node
	}
}

func (z *SortedSet) exist(key string) bool {
	_, exist := z.record[key]
	return exist
}

func newSkipList() *skipList {
	return &skipList{
		level: 1,
		head:  sklNewNode(maxLevel, 0, ""),
	}
}

func sklNewNode(level int16, score float64, member string) *sklNode {
	node := &sklNode{
		score:  score,
		member: member,
		level:  make([]*sklLevel, level),
	}

	for i := range node.level {
		node.level[i] = new(sklLevel)
	}
	return node
}

func (skl *skipList) sklInsert(score float64, member string) *sklNode {
	updates := make([]*sklNode, maxLevel)
	rank := make([]uint64, maxLevel)

	p := skl.head
	for i := skl.level - 1; i >= 0; i-- {
		if i == skl.level-1 {
			rank[i] = 0
		} else {
			rank[i] = rank[i+1]
		}

		if p.level[i] != nil {
			for p.level[i].forward != nil &&
				(p.level[i].forward.score < score ||
					(p.level[i].forward.score == score && p.level[i].forward.member < member)) {
				rank[i] += p.level[i].span
				p = p.level[i].forward
			}
		}
		updates[i] = p
	}

	// update skipList's head node if current level is greater than skl.level
	level := randomLevel()
	if level > skl.level {
		for i := skl.level; i < level; i++ {
			rank[i] = 0
			updates[i] = skl.head
			updates[i].level[i].span = uint64(skl.length)
		}
		skl.level = level
	}

	p = sklNewNode(level, score, member)
	for i := int16(0); i < level; i++ {
		p.level[i].forward = updates[i].level[i].forward
		updates[i].level[i].forward = p

		p.level[i].span = updates[i].level[i].span - (rank[0] - rank[i])
		updates[i].level[i].span = (rank[0] - rank[i]) + 1
	}

	for i := level; i < skl.level; i++ {
		updates[i].level[i].span++
	}

	if updates[0] == skl.head {
		p.backward = nil
	} else {
		p.backward = updates[0]
	}

	if p.level[0].forward != nil {
		p.level[0].forward.backward = p
	} else {
		skl.tail = p
	}

	skl.length++
	return p
}

func (skl *skipList) sklDelete(score float64, member string) {
	update := make([]*sklNode, maxLevel)
	p := skl.head

	for i := skl.level - 1; i >= 0; i-- {
		for p.level[i].forward != nil &&
			(p.level[i].forward.score < score ||
				(p.level[i].forward.score == score && p.level[i].forward.member < member)) {
			p = p.level[i].forward
		}
		update[i] = p
	}

	p = p.level[0].forward
	if p != nil && score == p.score && p.member == member {
		skl.sklDeleteNode(p, update)
		return
	}
}

func (skl *skipList) sklDeleteNode(p *sklNode, updates []*sklNode) {
	for i := int16(0); i < skl.level; i++ {
		if updates[i].level[i].forward == p {
			updates[i].level[i].span += p.level[i].span - 1
			updates[i].level[i].forward = p.level[i].forward
		} else {
			updates[i].level[i].span--
		}
	}

	if p.level[0].forward != nil {
		p.level[0].forward.backward = p.backward
	} else {
		skl.tail = p.backward
	}

	for skl.level > 1 && skl.head.level[skl.level-1].forward == nil {
		skl.level--
	}

	skl.length--
}

func randomLevel() int16 {
	var level int16 = 1
	for level < maxLevel {
		if rand.Float64() < probability {
			break
		}
		level++
	}
	return level
}
