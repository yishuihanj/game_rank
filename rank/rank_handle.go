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
	r.ChangeScore(req.PlayerId, req.Score, req.Ts)
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
			PlayerId: s.Key(),
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
	memberId := msg.Data.(string)
	info := r.GetRankById(memberId)
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

// 获取玩家周边排名
func (r *Rank) RankGetPlayerRankRange(ctx context.Context, msg *actor.Message) error {
	req := msg.Data.(*GetPlayerRankRangeReq)
	var err error
	list, start := r.GetPlayerRankRange(req.PlayerId, req.RangeNum)
	if list == nil {
		err = errors.New("failed to get information")
	}
	var rankList = make([]*RankSingleInfo, 0)
	if len(list) > 0 {
		rankList = make([]*RankSingleInfo, 0, len(list))
		for i, s := range list {
			ex := &RankSingleInfo{
				Ranking:  start + uint32(i),
				PlayerId: s.Key(),
				Score:    s.Score(),
			}
			rankList = append(rankList, ex)
		}
	}
	resp := &GetPlayerRankRangeResp{List: rankList}
	msg.Response(actor.RespMessage{
		Err:  err,
		Data: resp,
	})
	return nil
}
