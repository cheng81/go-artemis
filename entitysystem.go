package artemis

import (
	util "github.com/cheng81/go-artemis/util"
)

type EntitySystemTypeId uint

func (t EntitySystemTypeId) Uint() uint { return uint(t) }

type EntitySystemProcessor interface {
	CheckProcessing() bool
	Begin() //Called before processing of entities
	End()
	Initialize()
	Inserted(*Entity)
	Removed(*Entity)
	ProcessEntities(util.ImmutableBag)
	SetEntitySystem(*EntitySystem)
}

func NewEntitySystem(aspect *Aspect, typeId EntitySystemTypeId, processor EntitySystemProcessor) (out *EntitySystem) {
	out = &EntitySystem{
		processor: processor,
		typeId:    typeId,
		world:     nil,
		actives:   util.NewBag(64),
		aspect:    aspect,
		allSet:    aspect.allSet,
		exclSet:   aspect.exclSet,
		oneSet:    aspect.oneSet,
		passive:   false,
		dummy:     aspect.allSet.Empty() && aspect.oneSet.Empty(),
	}
	processor.SetEntitySystem(out)
	return
}

type EntitySystem struct {
	typeId  EntitySystemTypeId
	world   *World
	actives *util.Bag //*Entity
	aspect  *Aspect
	allSet  *util.BitSet
	exclSet *util.BitSet
	oneSet  *util.BitSet

	passive, dummy bool
	processor      EntitySystemProcessor
}

func (es *EntitySystem) Initialize() {
	es.processor.Initialize()
}

func (es *EntitySystem) TypeId() EntitySystemTypeId {
	return es.typeId
}

func (es *EntitySystem) Process() {
	if es.processor.CheckProcessing() {
		es.processor.Begin()
		es.processor.ProcessEntities(es.actives)
		es.processor.End()
	}
}

func (es *EntitySystem) Check(e *Entity) {
	if es.dummy {
		return
	}

	contains := e.systemBits.Get(es.typeId.Uint())
	interested := true

	cBits := e.componentBits
	if !es.allSet.Empty() {
		es.allSet.WithSetBitsFrom(0, func(index uint) {
			if !cBits.Get(index) {
				interested = false
			}
		})
	}

	if interested && !es.exclSet.Empty() {
		interested = !es.exclSet.Intersects(cBits)
	}

	if interested && !es.oneSet.Empty() {
		interested = es.oneSet.Intersects(cBits)
	}

	if interested && !contains {
		es.insert(e)
	} else if !interested && contains {
		es.remove(e)
	}
}

func (es *EntitySystem) insert(e *Entity) {
	es.actives.Add(e)
	e.systemBits.Set(es.typeId.Uint())
	es.processor.Inserted(e)
}
func (es *EntitySystem) remove(e *Entity) {
	es.actives.RemoveElem(e)
	e.systemBits.Unset(es.typeId.Uint())
	es.processor.Removed(e)
}

func (es *EntitySystem) Added(e *Entity)   { es.Check(e) }
func (es *EntitySystem) Changed(e *Entity) { es.Check(e) }
func (es *EntitySystem) Enabled(e *Entity) { es.Check(e) }

func (es *EntitySystem) Disabled(e *Entity) {
	if e.systemBits.Get(e.id) {
		es.remove(e)
	}
}
func (es *EntitySystem) Deleted(e *Entity) {
	if e.systemBits.Get(es.typeId.Uint()) {
		es.remove(e)
	}
}

func (es *EntitySystem) SetWorld(w *World) {
	es.world = w
}
func (es *EntitySystem) World() *World { return es.world }
