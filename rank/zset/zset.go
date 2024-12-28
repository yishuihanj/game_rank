package zset

import (
	"math/rand"
)

const (
	skipListMaxLevel = 8    // (1/p)^maxLevel >= maxNode
	skipListP        = 0.25 // SkipList P = 1/4
)

type zSkipListLevel struct {
	forward *ZSkipListNode
	span    uint32
}

// ZSkipListNode is an element of a skip list
type ZSkipListNode struct {
	backward *ZSkipListNode
	level    []zSkipListLevel
	score    uint64
	time     int64
	key      uint64
}

func zslCreateNode(level int, score uint64, key uint64, time int64) *ZSkipListNode {
	zn := &ZSkipListNode{
		time:  time,
		key:   key,
		score: score,
		level: make([]zSkipListLevel, level),
	}
	return zn
}

// Score return score
func (node *ZSkipListNode) Score() uint64 {
	return node.score
}

// Key return key
func (node *ZSkipListNode) Key() uint64 {
	return node.key
}

// Time 时间
func (node *ZSkipListNode) Time() int64 {
	return node.time
}

// positive if e > ele.
// negative if e < ele.
// 0 if e and ele are exactly the same.
func (node *ZSkipListNode) cmp(ele *ZSkipListNode) int {
	if node.time < ele.time {
		return 1
	}

	if node.time > ele.time {
		return -1
	}

	if node.key < ele.key {
		return 1
	}

	if node.key > ele.key {
		return -1
	}

	return 0
}

// Backward back
func (node *ZSkipListNode) Backward() *ZSkipListNode {
	return node.backward
}

// zSkipList represents a skip list
type zSkipList struct {
	header, tail *ZSkipListNode
	length       uint32
	level        int32 // current level count
}

// zslCreate creates a skip list
func zslCreate() *zSkipList {
	zsl := &zSkipList{
		level: 1,
	}
	zsl.header = zslCreateNode(skipListMaxLevel, 0, 0, 0, nil)
	return zsl
}

// insert element
func (list *zSkipList) insert(node *ZSkipListNode) *ZSkipListNode {
	var update [skipListMaxLevel]*ZSkipListNode
	var rank [skipListMaxLevel]uint32

	x := list.header
	for i := list.level - 1; i >= 0; i-- {
		if i == list.level-1 {
			rank[i] = 0
		} else {
			rank[i] = rank[i+1]
		}

		for x.level[i].forward != nil &&
			(x.level[i].forward.score < node.score ||
				x.level[i].forward.score == node.score &&
					x.level[i].forward.cmp(node) < 0) {
			rank[i] += x.level[i].span
			x = x.level[i].forward
		}
		update[i] = x
	}

	level := int32(len(node.level))
	if level > list.level {
		for i := list.level; i < level; i++ {
			rank[i] = 0
			update[i] = list.header
			update[i].level[i].span = list.length
		}
		list.level = level
	}

	x = node
	for i := int32(0); i < level; i++ {
		x.level[i].forward = update[i].level[i].forward
		update[i].level[i].forward = x
		x.level[i].span = update[i].level[i].span - (rank[0] - rank[i])
		update[i].level[i].span = (rank[0] - rank[i]) + 1
	}
	for i := level; i < list.level; i++ {
		update[i].level[i].span++
	}

	if update[0] == list.header {
		x.backward = nil
	} else {
		x.backward = update[0]
	}
	if x.level[0].forward == nil {
		list.tail = x
	} else {
		x.level[0].forward.backward = x
	}
	list.length++
	return x
}

// delete element
func (list *zSkipList) delete(node *ZSkipListNode) *ZSkipListNode {
	var update [skipListMaxLevel]*ZSkipListNode
	x := list.header
	for i := list.level - 1; i >= 0; i-- {
		for next := x.level[i].forward; next != nil &&
			(next.score < node.score ||
				next.score == node.score &&
					next.cmp(node) < 0); next = x.level[i].forward {
			x = next
		}
		update[i] = x
	}
	x = x.level[0].forward
	if x != nil && x.score == node.score && x.key == node.key {
		for i := int32(0); i < list.level; i++ {
			if update[i].level[i].forward == x {
				update[i].level[i].span += x.level[i].span - 1
				update[i].level[i].forward = x.level[i].forward
			} else {
				update[i].level[i].span--
			}
		}
		if x.level[0].forward == nil {
			list.tail = x.backward
		} else {
			x.level[0].forward.backward = x.backward
		}
		for list.level > 1 && list.header.level[list.level-1].forward == nil {
			list.level--
		}

		list.length--
		return x
	}
	return nil
}

// Find the rank for an element.
// Returns 0 when the element cannot be found, rank otherwise.
// Note that the rank is 1-based
func (list *zSkipList) zslGetRank(node *ZSkipListNode) uint32 {
	var rank uint32
	x := list.header
	for i := list.level - 1; i >= 0; i-- {
		for next := x.level[i].forward; next != nil &&
			(next.score < node.score ||
				next.score == node.score &&
					(node.time < next.time ||
						next.cmp(node) <= 0)); next = x.level[i].forward {
			rank += x.level[i].span
			x = next
		}
		if x.key != 0 && x.key == node.key {
			return rank
		}
	}
	return 0
}

func (list *zSkipList) randomLevel() int {
	lvl := 1
	for lvl < skipListMaxLevel && rand.Float64() < skipListP {
		lvl++
	}
	return lvl
}

// Finds an element by its rank. The rank argument needs to be 1-based.
func (list *zSkipList) getElementByRank(rank uint32) *ZSkipListNode {
	if rank == list.length {
		return list.tail
	}

	if rank == 1 {
		return list.header.level[0].forward
	}

	var traversed uint32
	x := list.header
	for i := list.level - 1; i >= 0; i-- {
		for x.level[i].forward != nil && traversed+x.level[i].span <= rank {
			traversed += x.level[i].span
			x = x.level[i].forward
		}
		if traversed == rank {
			return x
		}
	}
	return nil
}

// GetSimilar 获得跟score最相近的节点
func (list *zSkipList) GetSimilar(score uint64) *ZSkipListNode {
	x := list.header

	for i := list.level - 1; i >= 0; i-- {
		for n := x.level[i].forward; n != nil && n.Score() < score; n = x.level[i].forward {
			x = x.level[i].forward
		}
	}

	if x == list.header {
		return x.level[0].forward
	}

	var delta = func(x1, x2 uint64) uint64 {
		if x1 <= x2 {
			return x2 - x1
		}
		return x1 - x2
	}

	if x.level[0].forward != nil && delta(x.level[0].forward.Score(), score) < delta(score, x.Score()) {
		return x.level[0].forward
	}
	return x
}

// ZSet set
type ZSet struct {
	dict map[uint64]*ZSkipListNode
	zsl  *zSkipList
}

// NewZSet create ZSet
func NewZSet() *ZSet {
	zs := &ZSet{
		dict: make(map[uint64]*ZSkipListNode),
		zsl:  zslCreate(),
	}
	return zs
}

// Add a new element or update the score of an existing element
func (zs *ZSet) Add(score uint64, key uint64, t int64) *ZSkipListNode {
	if node := zs.dict[key]; node != nil {
		oldScore := node.score
		if score == oldScore {
			return nil
		}
		if next := node.level[0].forward; score > oldScore && (next == nil || score < next.score) {
			node.score = score
			node.time = t
		} else if score < oldScore && (node.backward == nil || score > node.backward.score) {
			node.score = score
			node.time = t
		} else {
			zs.zsl.delete(node)
			node.score = score
			node.time = t
			zs.zsl.insert(node)
		}
		return node
	}

	lvl := zs.zsl.randomLevel()
	node := zslCreateNode(lvl, score, key, t)
	zs.zsl.insert(node)
	zs.dict[key] = node
	return node
}

// Delete the element 'ele' from the sorted set,
// return 1 if the element existed and was deleted, 0 otherwise
func (zs *ZSet) Delete(id uint64) int {
	node := zs.dict[id]
	if node == nil {
		return 0
	}
	zs.zsl.delete(node)
	delete(zs.dict, id)
	return 1
}

// Rank return 1-based rank or 0 if not exist
func (zs *ZSet) Rank(id uint64, reverse bool) (uint32, *ZSkipListNode) {
	node := zs.dict[id]
	if node != nil {
		rank := zs.zsl.zslGetRank(node)
		if rank > 0 {
			if reverse {
				//return zs.zsl.length - rank + 1, node.score
				return zs.zsl.length - rank + 1, node
			}
			//return rank, node.score
			return rank, node
		}
	}
	return 0, nil
}

// Score return score
func (zs *ZSet) Score(id uint64) uint64 {
	node := zs.dict[id]
	if node != nil {
		return node.score
	}
	return 0
}

// Range return 1-based elements in [start, end]
func (zs *ZSet) Range(start uint32, end uint32, reverse bool, retNode *[]*ZSkipListNode) uint32 {
	if start == 0 {
		start = 1
	}
	if end == 0 {
		end = zs.zsl.length
	}
	if start > end || start > zs.zsl.length {
		return 0
	}
	if end > zs.zsl.length {
		end = zs.zsl.length
	}
	rangeLen := end - start + 1
	if reverse {
		node := zs.zsl.getElementByRank(zs.zsl.length - start + 1)
		for i := uint32(0); i < rangeLen; i++ {
			*retNode = append(*retNode, node)
			node = node.backward
		}
	} else {
		node := zs.zsl.getElementByRank(start)
		for i := uint32(0); i < rangeLen; i++ {
			*retNode = append(*retNode, node)
			node = node.level[0].forward
		}
	}
	return start
}

// Range return 1-based elements in [start, end]
func (zs *ZSet) MemberNode(rank uint32, reverse bool) *ZSkipListNode {
	if rank == 0 {
		rank = 1
	}
	if rank > zs.zsl.length {
		return nil
	}
	if reverse {
		return zs.zsl.getElementByRank(zs.zsl.length - rank + 1)
	} else {
		return zs.zsl.getElementByRank(rank)
	}
}

// Length return the element count
func (zs *ZSet) Length() uint32 {
	return zs.zsl.length
}

// MinScore return min score
func (zs *ZSet) MinScore() uint64 {
	first := zs.zsl.header.level[0].forward
	if first != nil {
		return first.score
	}
	return 0
}

// Tail return the last element
func (zs *ZSet) Tail() *ZSkipListNode {
	if zs.zsl.tail != nil {
		return zs.zsl.tail
	}
	return nil
}

// Head return the head
func (zs *ZSet) Head() *ZSkipListNode {
	if zs.zsl.tail != nil {
		return zs.zsl.header
	}
	return nil
}

func (zs *ZSet) GetAll() []*ZSkipListNode {
	if zs.zsl.tail == nil {
		return nil
	}
	list := []*ZSkipListNode{}
	tmp := zs.zsl.tail
	for tmp != nil {
		list = append(list, tmp)
		tmp = tmp.backward
	}
	return list
}

// DeleteFirst the first element
func (zs *ZSet) DeleteFirst() *ZSkipListNode {
	node := zs.zsl.header.level[0].forward
	zs.zsl.delete(node)
	delete(zs.dict, node.key)
	return node
}

// GetSimilar 获得跟score最相近的节点
func (zs *ZSet) GetSimilar(score uint64) *ZSkipListNode {
	return zs.zsl.GetSimilar(score)
}

func (zs *ZSet) GetRankByNode(node *ZSkipListNode) uint32 {
	return zs.zsl.zslGetRank(node)
}
