package rank

type RankSingleInfo struct {
	Ranking  uint32
	MemberId uint64
	Score    uint64
}

type RankSetSingleReq struct {
	MemberId uint64
	Score    uint64
}

type RankGetByMemberIdReq struct {
	MemberId uint64
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
