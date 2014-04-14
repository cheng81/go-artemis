package managers

import (
	core "github.com/cheng81/go-artemis"
)

func NewTagManager() *TagManager {
	return &TagManager{
		DefaultManager: core.NewDefaultManager(),
		entitiesByTag:  make(map[string]*core.Entity),
		tagsByEntity:   make(map[*core.Entity]string),
	}
}

type TagManager struct {
	*core.DefaultManager
	entitiesByTag map[string]*core.Entity
	tagsByEntity  map[*core.Entity]string
}

func (m *TagManager) TypeId() core.ManagerTypeId { return TagManagerTypeId }

func (m *TagManager) Register(tag string, e *core.Entity) {
	m.entitiesByTag[tag] = e
	m.tagsByEntity[e] = tag
}
func (m *TagManager) Unregister(tag string) {
	e, found := m.entitiesByTag[tag]
	if found {
		delete(m.entitiesByTag, tag)
		delete(m.tagsByEntity, e)
	}
}
func (m *TagManager) Entity(tag string) (*core.Entity, bool) {
	e, found := m.entitiesByTag[tag]
	return e, found
}

func (m *TagManager) RegisteredTags() (out []string) {
	l := len(m.entitiesByTag)
	out = make([]string, l, l)
	i := 0
	for k, _ := range m.entitiesByTag {
		out[i] = k
		i++
	}
	return
}

func (m *TagManager) Deleted(e *core.Entity) {
	tag, found := m.tagsByEntity[e]
	if found {
		delete(m.tagsByEntity, e)
		delete(m.entitiesByTag, tag)
	}
}

// public Collection<String> getRegisteredTags() {
//   return tagsByEntity.values();
// }

// @Override
// public void deleted(Entity e) {
//   String removedTag = tagsByEntity.remove(e);
//   if(removedTag != null) {
//     entitiesByTag.remove(removedTag);
//   }
// }
