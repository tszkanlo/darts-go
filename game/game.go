package game

import (
	"dartsvader-go/model"
	"dartsvader-go/websocket"
	"encoding/json"

	"github.com/google/uuid"
)

var (
	game     *model.Game
	gameView string
)

const (
	WebsocketStartGame    = "start"
	WebsocketInsertThrow  = "insert_throw"
	WebsocketInsertDelete = "delete_throw"
)

func GetGame() *model.Game {
	if game == nil {
		game = model.NewGame()
	}

	return game
}

func SetPlayer(player *model.Player) {
	GetGame().SetPlayer(player)
}

func Throw(c *model.CamCommand) {
	player := GetGame().GetCurrentPlayer()
	if player.HasMoreThrow() == false {
		GetGame().NextPlayer()

		player = GetGame().GetCurrentPlayer()
	}

	if player.GetCurrentRoundID() == "" {
		player.IncRound()
	}

	thr := &model.Throw{
		ID:       uuid.New().String(),
		Score:    c.Score,
		Modifier: c.Modifier,
	}
	player.SetThrow(thr)

	jsonThr, _ := json.Marshal(struct {
		Command  string       `json:"command"`
		ID       string       `json:"id"`
		PlayerID string       `json:"playerId"`
		RoundID  string       `json:"roundId"`
		Thr      *model.Throw `json:"throw"`
	}{
		Command:  WebsocketInsertThrow,
		ID:       thr.ID,
		PlayerID: player.ID,
		RoundID:  player.GetCurrentRoundID(),
		Thr:      thr,
	})

	websocket.BroadcastMsg(jsonThr)
}

func WebsocketOnConnectMsg() []byte {
	g, _ := json.Marshal(struct {
		Command string      `json:"command"`
		Game    *model.Game `json:"game"`
	}{
		Command: WebsocketStartGame,
		Game:    GetGame(),
	})

	return g
}
