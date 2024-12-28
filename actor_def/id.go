package actor_def

import "game_rank/pkg/actor"

const (
	Actor_System_Begin = iota // 系统 actor id 从1-999
	Actor_System_Game
	Actor_System_RankMgr
)

var (
	GamePid                = actor.NewPID(Actor_System_Game, "GamePid")
	RankMgrSystemPID       = actor.NewPID(Actor_System_RankMgr, "Actor_System_RankMgr")
	Actor_Name_Single_Rank = "Actor_Name_Single_Rank"
)
