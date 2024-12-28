package rank

import (
	"game_rank/actor_def"
	"game_rank/pkg/actor"
	"github.com/pkg/errors"
	"time"
)

// UpdateScore
//
//	@Description: 更新玩家分数
//	@param pid
//	@param playerId 玩家ID
//	@param score 分数
//	@param ts 时间戳 纳秒
func UpdateScore(playerId string, score uint64, ts int64) {
	if ts == 0 {
		ts = time.Now().UnixNano()
	}
	_ = actor.Send(actor_def.GamePid, actor_def.RankMgrSystemPID, &actor.Message{
		Id:  actor_def.MsgId_Rank_SetRankData,
		Uid: 0,
		Data: &RankSetSingleReq{
			PlayerId: playerId,
			Score:    score,
			Ts:       ts,
		},
	})
}

// GetPlayerRank
//
//	@Description: 通过玩家ID获取玩家排行榜信息
//	@param memberId
//	@return error
//	@return *RankSingleInfo
func GetPlayerRank(playerId string) (error, *RankSingleInfo) {
	rsp := actor.SyncRequest(actor_def.GamePid, actor_def.RankMgrSystemPID, &actor.Message{
		Id:   actor_def.MsgId_Rank_GetRankByMemberId,
		Uid:  0,
		Data: playerId,
	})
	if rsp.Err != nil {
		return rsp.Err, nil
	}
	resp := rsp.Data.(*RankGetByMemberIdRsp)
	if resp == nil {
		return errors.New("data is nil"), nil
	}
	return nil, resp.Data
}

// GetTopN
//
//	@Description: 获取排行榜前N名，N必须 >0
//	@param n
//	@return error
//	@return []*RankSingleInfo
func GetTopN(n uint32) (error, []*RankSingleInfo) {
	if n == 0 {
		return errors.New("n is 0"), nil
	}
	rsp := actor.SyncRequest(actor_def.GamePid, actor_def.RankMgrSystemPID, &actor.Message{
		Id:  actor_def.MsgId_Rank_GetRankRange,
		Uid: 0,
		Data: &RankGetRangeReq{
			Start: 1,
			End:   n,
		},
	})
	if rsp.Err != nil {
		return rsp.Err, nil
	}
	resp := rsp.Data.(*RankGetRangeRsp)
	if resp == nil {
		return errors.New("data is nil"), nil
	}
	return nil, resp.List
}

// GetPlayerRankRange
//
//	@Description:
//	@param playerId 获取玩家周边排名
//	@param rangeNum 周围
//	@return error
//	@return []*RankSingleInfo
func GetPlayerRankRange(playerId string, rangeNum uint32) (error, []*RankSingleInfo) {
	rsp := actor.SyncRequest(actor_def.GamePid, actor_def.RankMgrSystemPID, &actor.Message{
		Id:  actor_def.MsgId_Rank_GetPlayerRankRange,
		Uid: 0,
		Data: &GetPlayerRankRangeReq{
			PlayerId: playerId,
			RangeNum: rangeNum,
		},
	})
	if rsp.Err != nil {
		return rsp.Err, nil
	}
	resp := rsp.Data.(*GetPlayerRankRangeResp)
	if resp == nil {
		return errors.New("data is nil"), nil
	}
	return nil, resp.List
}
