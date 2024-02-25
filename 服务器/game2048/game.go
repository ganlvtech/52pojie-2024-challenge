package game2048

import (
	"math/rand/v2"
)

const MoveUp = 1
const MoveDown = 2
const MoveLeft = 3
const MoveRight = 4

var MoveDirectionData = map[int][][]int{
	MoveUp: {
		{0, 4, 8, 12},
		{1, 5, 9, 13},
		{2, 6, 10, 14},
		{3, 7, 11, 15},
	},
	MoveDown: {
		{12, 8, 4, 0},
		{13, 9, 5, 1},
		{14, 10, 6, 2},
		{15, 11, 7, 3},
	},
	MoveLeft: {
		{0, 1, 2, 3},
		{4, 5, 6, 7},
		{8, 9, 10, 11},
		{12, 13, 14, 15},
	},
	MoveRight: {
		{3, 2, 1, 0},
		{7, 6, 5, 4},
		{11, 10, 9, 8},
		{15, 14, 13, 12},
	},
}

type Game struct {
	Tiles []int64 `json:"tiles"` // 第一行从左到右下标是 0 1 2 3
	Score int64   `json:"score,omitempty"`
}

func NewGame() *Game {
	return &Game{
		Tiles: make([]int64, 16),
		Score: 0,
	}
}

func (g *Game) Move(direction int) (moved bool) {
	for _, line := range MoveDirectionData[direction] {
		fallTileIndex := 0
		mergeableTileIndex := -1
		for i := 0; i < len(line); i++ {
			if g.Tiles[line[i]] != 0 { // 不是空的方块
				// 先尝试合并，再尝试下落
				if mergeableTileIndex >= 0 && g.Tiles[line[i]] == g.Tiles[line[mergeableTileIndex]] {
					// 合并
					g.Tiles[line[mergeableTileIndex]] += g.Tiles[line[i]]
					g.Tiles[line[i]] = 0
					g.Score += g.Tiles[line[mergeableTileIndex]]
					moved = true
					mergeableTileIndex = -1 // 只能连续合并一次
				} else {
					if i != fallTileIndex {
						g.Tiles[line[fallTileIndex]] = g.Tiles[line[i]]
						g.Tiles[line[i]] = 0
						moved = true
					}
					mergeableTileIndex = fallTileIndex
					fallTileIndex += 1
				}
			}
		}
	}
	return
}

// AddRandomTile 随机生成 1 块方块，返回刚刚添加的方块的下标
func (g *Game) AddRandomTile() int {
	var emptyIndexList []int
	for i, v := range g.Tiles {
		if v == 0 {
			emptyIndexList = append(emptyIndexList, i)
		}
	}
	if len(emptyIndexList) == 0 {
		return -1
	}
	index := emptyIndexList[rand.IntN(len(emptyIndexList))]
	g.Tiles[index] = GetRandomTileNumber()
	return index
}

// IsGameOver 判断游戏是否结束
func (g *Game) IsGameOver() bool {
	// 有空格，则没有结束
	for _, v := range g.Tiles {
		if v == 0 {
			return false
		}
	}
	// 没有空格，但是有相邻相等，则没有结束
	for _, line := range MoveDirectionData[MoveUp] {
		for i := 0; i < len(line)-1; i++ {
			if g.Tiles[line[i]] == g.Tiles[line[i+1]] {
				return false
			}
		}
	}
	for _, line := range MoveDirectionData[MoveLeft] {
		for i := 0; i < len(line)-1; i++ {
			if g.Tiles[line[i]] == g.Tiles[line[i+1]] {
				return false
			}
		}
	}
	// 没有空格，所有相邻都不相等，则游戏结束
	return true
}

// GetRandomTileNumber 90% 的概率生成 2，10% 的概率生成 4
func GetRandomTileNumber() int64 {
	if rand.IntN(10) == 0 {
		return 4
	} else {
		return 2
	}
}
