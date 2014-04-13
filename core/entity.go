package core

import (
	util "github.com/cheng81/go-artemis/util"
)

func NewEntity(world *World, id uint) (out *Entity) {
	out = &Entity{
		id:            id,
		componentBits: util.NewBitSet(64),
		systemBits:    util.NewBitSet(64),
		world:         world,
		eManager:      world.entityManager,
		cManager:      world.componentManager,
	}
	out.reset()
	return
}

type Entity struct {
	uuid          util.UUID
	id            uint
	componentBits *util.BitSet
	systemBits    *util.BitSet
	world         *World
	eManager      *EntityManager
	cManager      *ComponentManager
}

func (e *Entity) Id() uint { return e.id }

func (e *Entity) AddComponent(c Component) *Entity {
	e.cManager.addComponent(e, c)
	return e
}

func (e *Entity) RemoveComponent(c Component) *Entity {
	return e.RemoveComponentOfType(c.TypeId())
}

func (e *Entity) RemoveComponentOfType(typeId ComponentTypeId) *Entity {
	e.cManager.removeComponent(e, typeId)
	return e
}

func (e *Entity) Active() bool  { return e.eManager.EntityActive(e.id) }
func (e *Entity) Enabled() bool { return e.eManager.EntityEnabled(e.id) }

func (e *Entity) Component(typeId ComponentTypeId) Component {
	return e.cManager.getComponent(e, typeId)
}

func (e *Entity) Components(fill *util.Bag) *util.Bag {
	return e.cManager.GetComponentsFor(e, fill)
}

func (e *Entity) AddToWorld()      { e.world.AddEntity(e) }
func (e *Entity) ChangedInWorld()  { e.world.ChangedEntity(e) }
func (e *Entity) DeleteFromWorld() { e.world.DeleteEntity(e) }
func (e *Entity) Enable()          { e.world.EnableEntity(e) }
func (e *Entity) Disable()         { e.world.DisableEntity(e) }

func (e *Entity) Uuid() util.UUID { return e.uuid }
func (e *Entity) World() *World   { return e.world }

func (e *Entity) HasSystem(es *EntitySystem) bool {
	return e.systemBits.Get(es.systemIndex)
}

func (e *Entity) HasComponent(ct ComponentTypeId) bool {
	return e.componentBits.Get(ct.Uint())
}

func (e *Entity) reset() {
	uuid, err := util.NewUUID()
	if err != nil {
		panic(err)
	}
	e.componentBits.Clear()
	e.systemBits.Clear()
	e.uuid = uuid
}
