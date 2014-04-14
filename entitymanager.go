package artemis

import (
	util "github.com/cheng81/go-artemis/util"
)

const (
	entityManagerTypeId = ManagerTypeId(1)
)

func NewEntityManager() *EntityManager {
	return &EntityManager{
		DefaultManager: NewDefaultManager(),
		entities:       util.NewBag(64),
		disabled:       util.NewBitSet(64),
		idPool:         newIdentifierPool(),
	}
}

type EntityManager struct {
	*DefaultManager
	entities *util.Bag
	disabled *util.BitSet

	active                  uint
	added, created, deleted uint64

	idPool *identifierPool
}

func (_ *EntityManager) TypeId() ManagerTypeId { return entityManagerTypeId }

func (e *EntityManager) createEntityInstance() (out *Entity) {
	out = NewEntity(e.world, e.idPool.checkOut())
	e.created += 1
	return
}

func (em *EntityManager) Added(e *Entity) {
	em.active++
	em.added++
	em.entities.SetAt(e.id, e)
}
func (em *EntityManager) Enabled(e *Entity) {
	em.disabled.Unset(e.id)
}
func (em *EntityManager) Disabled(e *Entity) {
	em.disabled.Set(e.id)
}
func (em *EntityManager) Deleted(e *Entity) {
	em.entities.SetAt(e.id, nil)
	em.disabled.Unset(e.id)
	em.idPool.checkIn(e.id)
	em.active--
	em.deleted++
}

func (e *EntityManager) EntityActive(id uint) bool  { return nil != e.entities.GetAt(id) }
func (e *EntityManager) EntityEnabled(id uint) bool { return !e.disabled.Get(id) }
func (e *EntityManager) Entity(id uint) *Entity     { return e.entities.GetAt(id).(*Entity) }

func (e *EntityManager) ActiveEntitiesCount() uint { return e.active }
func (e *EntityManager) TotalCreated() uint64      { return e.created }
func (e *EntityManager) TotalAdded() uint64        { return e.added }
func (e *EntityManager) TotalDeleted() uint64      { return e.deleted }

func newIdentifierPool() (out *identifierPool) {
	out = &identifierPool{
		ids:             util.NewBag(64),
		nextAvailableId: 0,
	}
	return
}

type identifierPool struct {
	ids             *util.Bag
	nextAvailableId uint
}

func (p *identifierPool) checkOut() uint {
	if p.ids.Size() > 0 {
		return p.ids.RemoveLast().(uint)
	}
	out := p.nextAvailableId
	p.nextAvailableId += 1
	return out
}

func (p *identifierPool) checkIn(id uint) {
	p.ids.Add(id)
}
