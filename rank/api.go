package rank

import (
	"game_rank/actor_def"
	"game_rank/pkg/actor"
)

// SetRankScore 设置排行榜
func SetRankScore(pid actor.PID, req *RankSetSingleReq) {
	actor.Send(pid, actor_def.RankMgrSystemPID, &actor.Message{
		Id:   actor_def.MsgId_Rank_SetRankData,
		Uid:  0,
		Data: req,
	})
}

// GetRankByMemberId 通过memberId获取排行
func GetRankByMemberId(pid actor.PID, req *RankGetByMemberIdReq) (error, *RankGetByMemberIdRsp) {
	rsp := actor.SyncRequest(pid, actor_def.RankMgrSystemPID, &actor.Message{
		Id:   actor_def.MsgId_Rank_GetRankByMemberId,
		Uid:  0,
		Data: req,
	})
	if rsp.Err != nil {
		return rsp.Err, nil
	}
	return nil, rsp.Data.(*RankGetByMemberIdRsp)
}
