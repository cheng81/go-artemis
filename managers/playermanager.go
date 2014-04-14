package managers

import (
	core "github.com/cheng81/go-artemis"
	util "github.com/cheng81/go-artemis/util"
)

func NewPlayerManager() *PlayerManager {
	return &PlayerManager{
		DefaultManager:   core.NewDefaultManager(),
		playerByEntity:   make(map[*core.Entity]string),
		entitiesByPlayer: make(map[string]*util.Bag),
	}
}

type PlayerManager struct {
	*core.DefaultManager
	playerByEntity   map[*core.Entity]string
	entitiesByPlayer map[string]*util.Bag
}

func (m *PlayerManager) TypeId() core.ManagerTypeId { return PlayerManagerTypeId }

func (m *PlayerManager) SetPlayer(e *core.Entity, player string) {
	m.playerByEntity[e] = player
	entities, ok := m.entitiesByPlayer[player]
	if !ok {
		entities = util.NewBag(32)
		m.entitiesByPlayer[player] = entities
	}
	entities.Add(e)
}

func (m *PlayerManager) EntitiesOfPlayer(player string) (util.ImmutableBag, bool) {
	bag, found := m.entitiesByPlayer[player]
	return bag, found
}

func (m *PlayerManager) RemoveFromPlayer(e *core.Entity) {
	player, ok := m.playerByEntity[e]
	if ok {
		entities, ok := m.entitiesByPlayer[player]
		if ok {
			entities.RemoveElem(e)
		}
	}
}

func (m *PlayerManager) Player(e *core.Entity) (string, bool) {
	player, found := m.playerByEntity[e]
	return player, found
}

func (m *PlayerManager) Deleted(e *core.Entity) {
	m.RemoveFromPlayer(e)
}
