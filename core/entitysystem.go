package core

import (
	util "github.com/cheng81/go-artemis/util"
)

type EntitySystemTypeId uint

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

// var sysIdx = uint(0)

func NewEntitySystem(aspect *Aspect, typeId EntitySystemTypeId, processor EntitySystemProcessor) (out *EntitySystem) {
	out = &EntitySystem{
		// EntityObserver: EmptyEntityObserver(),
		processor:   processor,
		typeId:      typeId,
		systemIndex: uint(typeId),
		world:       nil,
		actives:     util.NewBag(64),
		aspect:      aspect,
		allSet:      aspect.allSet,
		exclSet:     aspect.exclSet,
		oneSet:      aspect.oneSet,
		passive:     false,
		dummy:       aspect.allSet.Empty() && aspect.oneSet.Empty(),
	}
	processor.SetEntitySystem(out)
	return
}

type EntitySystem struct {
	// EntityObserver
	systemIndex uint
	typeId      EntitySystemTypeId
	world       *World
	actives     *util.Bag //*Entity
	aspect      *Aspect
	allSet      *util.BitSet
	exclSet     *util.BitSet
	oneSet      *util.BitSet

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

	contains := e.systemBits.Get(es.systemIndex)
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

	// if es.typeId == 18 {
	// fmt.Println("EntitySystem.Check - ", es.typeId, e.Id(), interested, contains, es.actives.Size(), es.systemIndex, e.systemBits.ToStr())
	// }

	if interested && !contains {
		// fmt.Println("EntitySystem.Check - insert", es.TypeId(), e.Id())
		es.insert(e)
	} else if !interested && contains {
		// fmt.Println("EntitySystem.Check - remove", es.TypeId(), e.Id())
		es.remove(e)
	}
}

func (es *EntitySystem) insert(e *Entity) {
	es.actives.Add(e)
	e.systemBits.Set(es.systemIndex)
	// if es.typeId == 18 {
	// 	fmt.Println("EntitySystem.insert - inserted", es.typeId, e.Id(), es.actives.Size())
	// }
	es.processor.Inserted(e)
}
func (es *EntitySystem) remove(e *Entity) {
	es.actives.RemoveElem(e)
	// fmt.Println("EntitySystem.remove", e.Id(), l, es.actives.Size())
	// if es.typeId == 18 {
	// 	fmt.Println("EntitySystem.remove - removed", es.typeId, e.Id(), es.actives.Size())
	// }
	e.systemBits.Unset(es.systemIndex)
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
	// fmt.Println("EntitySystem.Deleted - entity deleted")
	if e.systemBits.Get(es.systemIndex) {
		es.remove(e)
	}
}

func (es *EntitySystem) SetWorld(w *World) {
	es.world = w
}
func (es *EntitySystem) World() *World { return es.world }

// type SystemIndexManager struct {
// 	index   uint
// 	indices map[reflect.Type]uint
// }

// func (i *SystemIndexManager) getIndexFor(gType reflect.Type) uint {
// 	index, ok := i.indices[gType]
// 	if !ok {
// 		index = i.index
// 		i.indices[gType] = index
// 		i.index++
// 	}
// 	return index
// }

// var systemIndices = &SystemIndexManager{uint(0), make(map[reflect.Type]uint)}

// func GetIndexFor(gType reflect.Type) uint {
// 	return systemIndices.getIndexFor(gType)
// }
