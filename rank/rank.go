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
		set:         zset.NewZSet(),
		Locker:      sync.RWMutex{},
		PId:         actor_def.RankMgrSystemPID,
		rpcDispatch: actor.NewDispatcher(10),
	}
	r.register()
	return r
}

// ChangeScore 数据发生变化
func (r *Rank) ChangeScore(id string, score uint64, ts int64) *zset.ZSkipListNode {
	r.Locker.Lock()
	defer r.Locker.Unlock()
	if r.set.Length() >= maxNum && score <= r.set.MinScore() {
		return nil
	}
	if ts == 0 {
		ts = time.Now().UnixNano()
	}
	ele := r.set.Add(score, id, ts)
	if ele == nil {
		return nil
	}
	if r.set.Length() > maxNum {
		r.set.DeleteFirst()
	}
	return ele
}

func (r *Rank) GetRankById(id string) *RankSingleInfo {
	r.Locker.RLock()
	defer r.Locker.RUnlock()
	ranking, s := r.set.Rank(id, true)
	if s == nil {
		return nil
	}
	return &RankSingleInfo{Ranking: ranking, PlayerId: s.Key(), Score: s.Score()}
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

func (r *Rank) GetPlayerRankRange(playerId string, rangeNum uint32) ([]*zset.ZSkipListNode, uint32) {
	playerInfo := r.GetRankById(playerId)
	if playerInfo == nil {
		return nil, 0
	}
	start := uint32(1)
	if playerInfo.Ranking >= rangeNum {
		start = playerInfo.Ranking - rangeNum
	}
	end := playerInfo.Ranking + rangeNum

	list, _ := r.GetRange(start, end)
	return list, start
}
