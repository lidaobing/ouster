package scene

import (
	"github.com/tiancaiamao/ouster"
	"github.com/tiancaiamao/ouster/aoi"
	"github.com/tiancaiamao/ouster/data"
	"github.com/tiancaiamao/ouster/player"
	"math"
	"time"
)

type Map struct {
	data.Map
	players []*player.Player
	aoi     *aoi.Aoi

	quit      chan struct{}
	event     chan interface{}
	heartbeat <-chan time.Time
}

func New(m *data.Map) *Map {
	ret := new(Map)
	ret.players = make([]*player.Player, 0, 200)

	ret.quit = make(chan struct{})
	ret.event = make(chan interface{})
	ret.heartbeat = time.Tick(50 * time.Microsecond)

	return ret
}

func (m *Map) Player(playerId uint32) *player.Player {
	if playerId >= uint32(len(m.players)) {
		return nil
	}
	return m.players[playerId]
}

func (m *Map) HeartBeat() {
	for _, pc := range m.players {
		if pc.State == player.MOVE {
			v := pc.Speed()
			if ouster.Distance(pc.Pos, pc.To) < v {
				pc.State = player.STAND
				pc.Pos.X = pc.To.X
				pc.Pos.Y = pc.To.Y
			} else {
				dx := pc.To.X - pc.Pos.X
				dy := pc.To.Y - pc.Pos.Y
				angle := math.Atan2(float64(dy), float64(dx))
				vx := v * float32(math.Cos(angle))
				vy := v * float32(math.Sin(angle))

				pc.Pos.X += vx
				pc.Pos.Y += vy
			}
		}
	}
}

func (m *Map) Login(player *player.Player) error {
	idx := len(m.players)
	m.players = append(m.players, player)
	player.Id = uint32(idx)
	player.Scene = m.Title

	return nil
}
