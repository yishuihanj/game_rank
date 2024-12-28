package rank

import (
	"context"
	"game_rank/actor_def"
	"game_rank/pkg/actor"
)

func (r *Rank) PID() actor.PID {
	return r.PId
}

func (r *Rank) OnStop() {
}

func (r *Rank) Process(msg *actor.Message) {
	r.rpcDispatch.Dispatch(context.Background(), msg)
}

func (r *Rank) Register() {

}

func (r *Rank) register() {
	r.rpcDispatch.Register(actor_def.MsgId_Rank_SetRankData, r.RankSetRankData)
	r.rpcDispatch.Register(actor_def.MsgId_Rank_GetRankRange, r.RankGetRankRange)
	r.rpcDispatch.Register(actor_def.MsgId_Rank_GetRankByMemberId, r.RankGetRankByMemberId)
	r.rpcDispatch.Register(actor_def.MsgId_Rank_GetPlayerRankRange, r.RankGetPlayerRankRange)
}
