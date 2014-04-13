package managers

import (
	core "github.com/cheng81/go-artemis/core"
	util "github.com/cheng81/go-artemis/util"
)

func NewGroupManager() *GroupManager {
	return &GroupManager{
		DefaultManager:  core.NewDefaultManager(),
		entitiesByGroup: make(map[string]*util.Bag),
		groupsByEntity:  make(map[*core.Entity]*util.Bag),
	}
}

type GroupManager struct {
	*core.DefaultManager
	entitiesByGroup map[string]*util.Bag
	groupsByEntity  map[*core.Entity]*util.Bag
}

func (_ *GroupManager) TypeId() core.ManagerTypeId { return GroupManagerTypeId }

func (m *GroupManager) Add(e *core.Entity, group string) {
	entities, ok := m.entitiesByGroup[group]
	if !ok {
		entities = util.NewBag(64)
		m.entitiesByGroup[group] = entities
	}
	entities.Add(e)

	groups, ok := m.groupsByEntity[e]
	if !ok {
		groups = util.NewBag(64)
		m.groupsByEntity[e] = groups
	}
	groups.Add(group)
}

func (m *GroupManager) Remove(e *core.Entity, group string) {
	entities, ok := m.entitiesByGroup[group]
	if ok {
		entities.RemoveElem(e)
	}
	groups, ok := m.groupsByEntity[e]
	if ok {
		groups.RemoveElem(group)
	}
}

func (m *GroupManager) RemoveFromAllGroups(e *core.Entity) {
	groups, ok := m.groupsByEntity[e]
	if ok {
		groups.ForEach(func(_ int, gi interface{}) {
			entities, ok := m.entitiesByGroup[gi.(string)]
			if ok {
				entities.RemoveElem(e)
			}
		})
		groups.Clear()
	}
}

func (m *GroupManager) Entities(group string) (util.ImmutableBag, bool) {
	out, found := m.entitiesByGroup[group]
	return out, found
}

func (m *GroupManager) Groups(e *core.Entity) (util.ImmutableBag, bool) {
	out, found := m.groupsByEntity[e]
	return out, found
}

func (m *GroupManager) InGroup(e *core.Entity, group string) bool {
	groups, ok := m.groupsByEntity[e]
	if ok {
		for i := uint(0); i < groups.Size(); i++ {
			if group == groups.GetAt(i).(string) {
				return true
			}
		}
	}
	return false
}

func (m *GroupManager) Deleted(e *core.Entity) {
	m.RemoveFromAllGroups(e)
}
