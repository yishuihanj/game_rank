package rank

type RankSingleInfo struct {
	Ranking  uint32
	PlayerId string
	Score    uint64
}

type RankSetSingleReq struct {
	PlayerId string
	Score    uint64
	Ts       int64
}

type RankGetByMemberIdReq struct {
	PlayerId string
}

type RankGetByMemberIdRsp struct {
	Data *RankSingleInfo
}

type RankGetRangeReq struct {
	Start uint32
	End   uint32
}

type RankGetRangeRsp struct {
	List []*RankSingleInfo
}

type GetPlayerRankRangeReq struct {
	PlayerId string
	RangeNum uint32
}

type GetPlayerRankRangeResp struct {
	List []*RankSingleInfo
}
