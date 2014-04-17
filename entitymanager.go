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
	}
}

type EntityManager struct {
	*DefaultManager
	entities *util.Bag
	disabled *util.BitSet

	active                  uint
	added, created, deleted uint64

	entPool *entityPool
}

func (_ *EntityManager) TypeId() ManagerTypeId { return entityManagerTypeId }

func (e *EntityManager) SetWorld(w *World) {
	e.entPool = newEntityPool(w)
	e.DefaultManager.SetWorld(w)
}

func (e *EntityManager) createEntityInstance() (out *Entity) {
	out = e.entPool.checkOut()
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
	em.entPool.checkIn(e)
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

func newEntityPool(w *World) *entityPool {
	return &entityPool{
		world:       w,
		ents:        util.NewBag(256),
		nextAvailId: 0,
	}
}

type entityPool struct {
	world       *World
	ents        *util.Bag
	nextAvailId uint
}

func (e *entityPool) checkOut() (out *Entity) {
	if e.ents.Size() > 0 {
		out = e.ents.RemoveLast().(*Entity)
		out.reset()
		return
	}
	out = newEntity(e.world, e.nextAvailId)
	e.nextAvailId += 1
	return
}
func (ep *entityPool) checkIn(e *Entity) {
	ep.ents.Add(e)
}
