package main

import (
	"fmt"
	"game_rank/pkg/actor"
	rank "game_rank/rank"
	"time"
)

// 该排行榜使用跳表结构，使用actor模型来保证了请求的顺序
// todo 该排行榜没有实现数据的持久化，可作为单区单服使用，在玩家数量不多时，比如20000人

func main() {
	fmt.Println("hello game_rank")
	gameInit()
	fmt.Println("开始增加排行榜")
	rank.UpdateScore("1111", 100, time.Now().UnixNano())
	rank.UpdateScore("2222", 99, time.Now().UnixNano())
	rank.UpdateScore("3333", 98, time.Now().UnixNano())
	err, resp := rank.GetPlayerRank("1111")
	if err != nil {
		fmt.Println("1111 玩家的排行榜获取失败", err.Error())
		return
	}
	fmt.Println("1111 玩家的排行是", resp.Ranking)

	fmt.Println("获取前2的排行榜数据")
	err, lists := rank.GetTopN(2)
	if err != nil {
		fmt.Println("获取前2的排行榜数据失败", err.Error())
		return
	}
	for _, list := range lists {
		fmt.Println("玩家:", list.PlayerId, "排名：", list.Ranking)
	}
	fmt.Println("获取3333玩家的周边左右1玩家的排名")
	err, lists1 := rank.GetPlayerRankRange("3333", 1)
	if err != nil {
		fmt.Println("获取3333玩家的周边排名失败", err.Error())
		return
	}
	for _, list := range lists1 {
		fmt.Println("玩家:", list.PlayerId, "排名：", list.Ranking)
	}

}

// 注册
func gameInit() {
	r := rank.NewRank()
	actor.RegisterActor(r, 20000)
}
