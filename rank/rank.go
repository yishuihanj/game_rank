package rank

import (
	"game_rank/actor_def"
	"game_rank/pkg/actor"
	"game_rank/rank/zset"
	"sync"
	"time"
)

const (
	//todo 默认先写20000个
	maxNum = 20000
)

type Rank struct {
	rpcDispatch *actor.Dispatcher
	set         *zset.ZSet
	Locker      sync.RWMutex
	PId         actor.PID
}

func NewRank() *Rank {
	r := &Rank{
		set:    zset.NewZSet(),
		Locker: sync.RWMutex{},
		PId:    actor.NewPID(1, actor_def.Actor_Name_Single_Rank),
		rpcDispatch: actor.NewDispatcher(10, func(uid uint64, id uint32, ms int64) {

		}),
	}
	r.rpcRegister()
	return r
}

// ChangeScore 数据发生变化
func (r *Rank) ChangeScore(id uint64, score uint64) *zset.ZSkipListNode {
	r.Locker.Lock()
	defer r.Locker.Unlock()
	if r.set.Length() >= maxNum && score <= r.set.MinScore() {
		return nil
	}

	t := time.Now().UnixNano()
	ele := r.set.Add(score, id, t)
	if ele == nil {
		return nil
	}
	if r.set.Length() > maxNum {
		r.set.DeleteFirst()
	}
	return ele
}

func (r *Rank) GetRankById(id uint64) *RankSingleInfo {
	r.Locker.RLock()
	defer r.Locker.RUnlock()
	ranking, s := r.set.Rank(id, true)
	if s == nil {
		return nil
	}
	return &RankSingleInfo{Ranking: ranking, MemberId: s.Key(), Score: s.Score()}
}

func (r *Rank) Length() uint32 {
	return r.set.Length()
}

func (r *Rank) GetRange(rankBegin uint32, rankEnd uint32) ([]*zset.ZSkipListNode, uint32) {
	r.Locker.RLock()
	defer r.Locker.RUnlock()
	rangeNodes := []*zset.ZSkipListNode{}
	start := r.set.Range(rankBegin, rankEnd, true, &rangeNodes)
	return rangeNodes, start
}

// TotalNum 总人数
func (r *Rank) TotalNum() uint32 {
	return r.set.Length()
}

// MinScore 最小积分
func (r *Rank) MinScore() uint64 {
	r.Locker.RLock()
	defer r.Locker.RUnlock()
	return r.set.MinScore()
}
