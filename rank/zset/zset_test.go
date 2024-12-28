package zset

import (
	"math/rand"
	"testing"
)

func init() {
}

func BenchmarkSkipList_RandomLevel(b *testing.B) {
	l := zslCreate()
	for i := 0; i < b.N; i++ {
		l.randomLevel()
	}
}

func TestAdd(t *testing.T) {
	zs := NewZSet()
	// key相同时，覆盖积分
	zs.Add(1, 1, 0)
	zs.Add(2, 1, 1)
	if zs.Length() != 1 || zs.Tail().time != 1 {
		t.Error("add")
	}
	// 时间相同、积分相同  key小的排前面
	zs.Add(2, 2, 1)
	rank, score := zs.Rank(2, true)
	if rank != 2 {
		t.Error("add")
	}

	// 积分不同
	zs.Add(100, 3, 1)
	rank, _ = zs.Rank(3, true)
	if rank != 1 {
		t.Error("add")
	}

	// 积分相同 时间不同
	zs.Add(100, 4, 0)
	rank, _ = zs.Rank(4, true)
	if rank != 1 {
		t.Error("add")
	}

	zs.Add(9, 3, 1)
	rank, score = zs.Rank(3, true)
	if rank != 2 || score != 9 || zs.Length() != 4 {
		t.Error("add")
	}
}

func TestZSkipList_GetSimilar(t *testing.T) {
	zs := NewZSet()
	zs.zsl.GetSimilar(0)
	zs.Add(1, 12, 0)
	zs.Add(4, 14, 0)
	zs.Add(5, 15, 0)
	zs.Add(7, 17, 0)
	zs.Add(9, 19, 0)

	var testCase = []struct {
		score uint64
		name  uint64
	}{
		{0, 12},
		{1, 12},
		{2, 12},
		{3, 13},
		{4, 14},
		{5, 15},
		{6, 16},
		{7, 17},
		{8, 18},
		{9, 19},
		{10, 20},
		{11, 21},
		{18, 22},
		{123, 23},
		{1234, 24},
	}

	for _, v := range testCase {
		node := zs.zsl.GetSimilar(v.score)
		if node.Key() != v.name {
			t.Error("GetSimilar", v, node.Key())
		}
	}
}

func TestDelete(t *testing.T) {
	zs := NewZSet()
	zs.Add(1, 3, 0)
	zs.Delete(3)
	zs.Delete(3)
	if zs.Length() != 0 {
		t.Error("delete")
	}

	zs.Add(1, 1, 0)
	zs.Add(1, 2, 0)
	zs.Add(1, 3, 0)
	zs.Add(1, 4, 0)
	zs.Add(1, 5, 0)
	zs.Delete(3)
	zs.Delete(88)
	if zs.Length() != 4 {
		t.Error("delete")
	}
	rank, score := zs.Rank(4, true)
	if rank != 3 || score != 1 {
		t.Error("delete")
	}
}

func TestGetRank(t *testing.T) {
	zs := NewZSet()
	zs.Add(1, 1, 0)
	zs.Add(1, 2, 0)
	zs.Add(1, 3, 0)
	zs.Add(1, 4, 0)
	zs.Add(1, 5, 0)
	rank, score := zs.Rank(4, true)
	if rank != 4 || score != 1 {
		t.Error("TestGetRank")
	}

	rank, score = zs.Rank(1, false)
	if rank != 5 || score != 1 {
		t.Error("TestGetRank", rank, score)
	}

	rank, score = zs.Rank(13, true)
	if rank != 0 || score != 0 {
		t.Error("get rank with not exist id ")
	}
}

func TestGetElementByRank(t *testing.T) {
	zs := NewZSet()
	for i := 1; i <= 10; i++ {
		zs.Add(uint64(i), uint64(i), 0)
	}
	node := zs.zsl.getElementByRank(6)
	if node.Key() != 6 {
		t.Error("TestGetElementByRank")
	}

	node = zs.zsl.getElementByRank(0)
	if node != zs.zsl.header {
		t.Error("TestGetElementByRank")
	}

	node = zs.zsl.getElementByRank(1)
	if node.Key() != 1 {
		t.Error("TestGetElementByRank")
	}

	node = zs.zsl.getElementByRank(10)
	if node.Key() != zs.Tail().Key() {
		t.Error("TestGetElementByRank")
	}

	node = zs.zsl.getElementByRank(11)
	if node != nil {
		t.Error("TestGetElementByRank")
	}
}

func TestRange(t *testing.T) {
	zs := NewZSet()
	for i := 0; i < 1000; i++ {
		zs.Add(uint64(i), 10000+uint64(i), 0)
	}
	ids := make([]*ZSkipListNode, 0, 5000)
	zs.Range(1, 20, true, &ids)
	if len(ids) != 20 || ids[0].key != 10999 {
		t.Error("zset TestRange")
	}

	ids = ids[:0]
	zs.Range(0, 0, true, &ids)
	if len(ids) != 1000 || ids[0].key != 10999 {
		t.Error("zset TestRange")
	}

	ids = ids[:0]
	zs.Range(0, 121212, false, &ids)
	if len(ids) != 1000 || ids[0].key != 10000 {
		t.Error("zset TestRange")
	}

	ids = ids[:0]
	zs.Range(12, 1, true, &ids)
	if len(ids) != 0 {
		t.Error("zset TestRange")
	}
}

func BenchmarkAdd(b *testing.B) {
	r := NewZSet()
	for i := 0; i < b.N; i++ {
		r.Add(uint64(i), uint64(i)%20000, 0)
	}
}

func BenchmarkChange(b *testing.B) {
	total := 5000

	r := NewZSet()
	for i := 0; i < total; i++ {
		r.Add(rand.Uint64()%100000000, uint64(i), rand.Int63())
	}

	var q = make([]uint64, 10000000)
	for i := 0; i < len(q); i++ {
		q[i] = rand.Uint64() % 100000000
	}

	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		id := uint64(i) % uint64(total)
		score := q[i%10000000]
		r.Add(score, id, int64(score))
		if r.Length() > uint32(total) {
			r.DeleteFirst()
		}
	}

	if r.Length() != uint32(total) {
		b.Error("ll")
	}
}

func BenchmarkRange(b *testing.B) {
	zs := NewZSet()
	for i := 0; i < 5000; i++ {
		zs.Add(uint64(i), 10000+uint64(i), 10)
	}

	ids := make([]*ZSkipListNode, 0, 5000)
	zs.Range(1, 20, true, &ids)
	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		ids = ids[:0]
		zs.Range(850, 870, true, &ids)
	}
}
