package managers

import (
	core "github.com/cheng81/go-artemis/core"
	util "github.com/cheng81/go-artemis/util"
)

func NewTeamManager() *TeamManager {
	return &TeamManager{
		DefaultManager: core.NewDefaultManager(),
		playersByTeam:  make(map[string]*util.Bag),
		teamByPlayer:   make(map[string]string),
	}
}

type TeamManager struct {
	*core.DefaultManager
	playersByTeam map[string]*util.Bag
	teamByPlayer  map[string]string
}

func (m *TeamManager) TypeId() core.ManagerTypeId { return TeamManagerTypeId }

func (m *TeamManager) Team(player string) string {
	return m.teamByPlayer[player]
}

func (m *TeamManager) SetTeam(player, team string) {
	m.RemoveFromTeam(player)
	m.teamByPlayer[player] = team
	players, ok := m.playersByTeam[team]
	if !ok {
		players = util.NewBag(64)
		m.playersByTeam[team] = players
	}
	players.Add(player)
}

func (m *TeamManager) RemoveFromTeam(player string) {
	team, ok := m.teamByPlayer[player]
	if ok {
		players, ok := m.playersByTeam[team]
		if ok {
			players.RemoveElem(player)
		}
	}
}

func (m *TeamManager) Players(team string) *util.Bag {
	return m.playersByTeam[team]
}
