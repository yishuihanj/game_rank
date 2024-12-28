package rank

import (
	"context"
	"errors"
	"game_rank/pkg/actor"
	"game_rank/rank/zset"
)

// set排行榜积分
func (r *Rank) RankSetRankData(ctx context.Context, msg *actor.Message) error {
	req := msg.Data.(*RankSetSingleReq)
	r.ChangeScore(req.MemberId, uint64(req.Score))
	return nil
}

// 获取排行榜list
func (r *Rank) RankGetRankRange(ctx context.Context, msg *actor.Message) error {
	resp := actor.RespMessage{}
	defer func() {
		msg.Response(resp)
	}()
	req := msg.Data.(*RankGetRangeReq)
	var start uint32
	var list []*zset.ZSkipListNode
	if req.Start == 0 && req.End == 0 {
		//获取全部
		list = r.set.GetAll()
		start = 1
	} else {
		list, start = r.GetRange(req.Start, req.End)
	}

	var rankList []*RankSingleInfo
	for i, s := range list {
		ex := &RankSingleInfo{
			Ranking:  start + uint32(i),
			MemberId: s.Key(),
			Score:    s.Score(),
		}
		rankList = append(rankList, ex)
	}
	data := &RankGetRangeRsp{
		List: rankList,
	}
	resp.Data = data
	return nil
}

// 获取成员排名信息
func (r *Rank) RankGetRankByMemberId(ctx context.Context, msg *actor.Message) error {
	req := msg.Data.(*RankGetByMemberIdReq)
	info := r.GetRankById(req.MemberId)
	var err error
	if info == nil {
		err = errors.New("failed to get information")
	}
	resp := &RankGetByMemberIdRsp{Data: info}
	msg.Response(actor.RespMessage{
		Err:  err,
		Data: resp,
	})
	return nil
}
