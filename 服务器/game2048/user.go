package game2048

import (
	"encoding/json"
	"errors"
	"fmt"
	"math/rand/v2"
	"strconv"
	"time"

	"2024challenge/dynamicflag"
)

const (
	ItemIDRemovePiece = 1
	ItemIDDoubleMoney = 2
	ItemIDV50         = 3
	ItemIDFlag8       = 4
	ItemIDFlagB       = 5
)

type UserData struct {
	GameData         *Game `json:"game_data,omitempty"`
	MoneyCount       int64 `json:"money_count,omitempty"`
	RemovePieceCount int64 `json:"remove_piece_count,omitempty"`
	DoubleMoneyCount int64 `json:"double_money_count,omitempty"`
	V50Count         int64 `json:"v50_count,omitempty"`
	Flag8Count       int64 `json:"flag_8_count,omitempty"`
	FlagBCount       int64 `json:"flag_b_count,omitempty"`
	DoubleMoneyBuff  bool  `json:"double_money_buff,omitempty"`
}

type UserDataResponseGameData struct {
	Tiles    []int64 `json:"tiles"`
	Score    int64   `json:"score,omitempty"`
	GameOver bool    `json:"game_over,omitempty"`
}
type UserDataResponse struct {
	GameData         *UserDataResponseGameData `json:"game_data,omitempty"`
	MoneyCount       int64                     `json:"money_count,omitempty"`
	RemovePieceCount int64                     `json:"remove_piece_count,omitempty"`
	DoubleMoneyCount int64                     `json:"double_money_count,omitempty"`
	V50Count         int64                     `json:"v50_count,omitempty"`
	Flag8Count       int64                     `json:"flag_8_count,omitempty"`
	FlagBCount       int64                     `json:"flag_b_count,omitempty"`
	DoubleMoneyBuff  bool                      `json:"double_money_buff,omitempty"`
}

func ParseUserData(s string) *UserData {
	var u UserData
	err := json.Unmarshal([]byte(s), &u)
	if err != nil {
		u = UserData{}
		u.NewGame()
	}
	return &u
}
func SerializeUserData(u *UserData) string {
	b, err := json.Marshal(u)
	if err != nil {
		return ""
	}
	return string(b)
}
func ToUserDataResponse(u *UserData) *UserDataResponse {
	var responseGameData *UserDataResponseGameData
	if u.GameData != nil {
		responseGameData = &UserDataResponseGameData{
			Tiles:    u.GameData.Tiles,
			Score:    u.GameData.Score,
			GameOver: u.GameData.IsGameOver(),
		}
	}
	return &UserDataResponse{
		GameData:         responseGameData,
		MoneyCount:       u.MoneyCount,
		RemovePieceCount: u.RemovePieceCount,
		DoubleMoneyCount: u.DoubleMoneyCount,
		V50Count:         u.V50Count,
		Flag8Count:       u.Flag8Count,
		FlagBCount:       u.FlagBCount,
		DoubleMoneyBuff:  u.DoubleMoneyBuff,
	}
}
func (u *UserData) BuyItem(shopItemID int, buyCount int64) error {
	if buyCount < 0 {
		return errors.New("buyCount 不能小于 0")
	}
	if buyCount == 0 {
		return errors.New("buyCount 不能为 0")
	}
	var shopItem *ShopItem
	for _, shopItem1 := range ShopItemList {
		if shopItem1.ShopItemID == shopItemID {
			shopItem = shopItem1
			break
		}
	}
	if shopItem == nil {
		return errors.New("shop_item_id 商品 id 不存在")
	}

	newMoneyCount := u.MoneyCount - buyCount*shopItem.Price
	if newMoneyCount < 0 {
		return errors.New("钱不够")
	}
	if newMoneyCount > u.MoneyCount {
		return errors.New("购买商品之后钱怎么还变多了？不知道出什么 bug 了，暂时先拦一下 ^_^")
	}
	u.MoneyCount = newMoneyCount
	u.AddItem(shopItem.BagItemID, buyCount*shopItem.Count)
	return nil
}
func (u *UserData) AddItem(itemID int, count int64) {
	switch itemID {
	case ItemIDRemovePiece:
		u.RemovePieceCount += count
	case ItemIDDoubleMoney:
		u.DoubleMoneyCount += count
	case ItemIDV50:
		u.V50Count += count
	case ItemIDFlag8:
		u.Flag8Count += count
	case ItemIDFlagB:
		u.FlagBCount += count
	}
}
func (u *UserData) UseItem(itemID int, uid uint64) (result string, err error) {
	switch itemID {
	case ItemIDRemovePiece:
		if u.GameData == nil {
			return "", errors.New("请先开始游戏，再使用道具")
		}
		if u.RemovePieceCount <= 0 {
			return "", errors.New("道具数量不足")
		}
		var tileIndexList []int
		for i, tile := range u.GameData.Tiles {
			if tile == 2 || tile == 4 {
				tileIndexList = append(tileIndexList, i)
			}
		}
		if len(tileIndexList) == 0 {
			return "", errors.New("没有2或4的格子")
		}
		u.RemovePieceCount -= 1
		u.GameData.Tiles[tileIndexList[rand.IntN(len(tileIndexList))]] = 0
		return "", nil
	case ItemIDDoubleMoney:
		if u.GameData == nil {
			return "", errors.New("请先开始游戏，再使用道具")
		}
		if u.DoubleMoneyCount <= 0 {
			return "", errors.New("道具数量不足")
		}
		if u.DoubleMoneyBuff {
			return "", errors.New("当前已经处于双倍状态了，请消耗掉当前状态再使用道具")
		}
		u.DoubleMoneyCount -= 1
		u.DoubleMoneyBuff = true
		return "", nil
	case ItemIDV50:
		if u.V50Count <= 0 {
			return "", errors.New("道具数量不足")
		}
		u.V50Count -= 1
		return "竟然真的有人v我50，真的太感动了。作为奖励呢，我就提示你一下吧，关键词是“溢出”。", nil
	case ItemIDFlag8:
		if u.Flag8Count <= 0 {
			return "", errors.New("道具数量不足")
		}
		u.Flag8Count -= 1
		return "flag8{OaOjIK}", nil
	case ItemIDFlagB:
		if u.FlagBCount <= 0 {
			return "", errors.New("道具数量不足")
		}
		u.FlagBCount -= 1
		flagContent, expiredAt := dynamicflag.CalcFlag(strconv.FormatUint(uid, 10), "_2024_52pojie_flagB_", time.Now())
		return fmt.Sprintf("flagB{%s} 过期时间: %s", flagContent, expiredAt.Format(time.DateTime)), nil
	default:
		return "", errors.New("未知 item_id")
	}
}
func (u *UserData) Move(direction int) (err error) {
	if _, ok := MoveDirectionData[direction]; !ok {
		return errors.New("direction 错误")
	}
	if u.GameData == nil {
		return errors.New("请先开始游戏")
	}
	prevScore := u.GameData.Score
	if u.GameData.Move(direction) {
		deltaScore := u.GameData.Score - prevScore
		u.MoneyCount += deltaScore
		if u.DoubleMoneyBuff {
			u.DoubleMoneyBuff = false
			u.MoneyCount += deltaScore
		}
		u.GameData.AddRandomTile()
	}
	return nil
}
func (u *UserData) NewGame() {
	u.MoneyCount -= 200
	if u.MoneyCount < 0 {
		u.MoneyCount = 0
	}
	u.GameData = NewGame()
	u.GameData.AddRandomTile()
	u.GameData.AddRandomTile()
}
