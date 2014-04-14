package artemis

import (
	util "github.com/cheng81/go-artemis/util"
)

type Performer func(EntityObserver, *Entity)

func NewWorld() (out *World) {
	em := NewEntityManager()
	cm := NewComponentManager()
	out = &World{
		entityManager:    em,
		componentManager: cm,

		added:   util.NewBag(64),
		changed: util.NewBag(64),
		deleted: util.NewBag(64),
		enable:  util.NewBag(64),
		disable: util.NewBag(64),

		managers:    make(map[ManagerTypeId]Manager),
		managersBag: util.NewBag(64),

		systems:    make(map[EntitySystemTypeId]*EntitySystem),
		systemsBag: util.NewBag(64),
	}

	out.AddManager(em)
	out.AddManager(cm)

	return
}

type World struct {
	entityManager    *EntityManager
	componentManager *ComponentManager

	delta float64

	added, changed, deleted, enable, disable *util.Bag

	managers    map[ManagerTypeId]Manager
	managersBag *util.Bag

	systems    map[EntitySystemTypeId]*EntitySystem
	systemsBag *util.Bag
}

func (w *World) Initialize() {
	w.managersBag.ForEach(func(_ int, m interface{}) {
		m.(Manager).Initialize()
	})
	w.systemsBag.ForEach(func(_ int, s interface{}) {
		s.(*EntitySystem).Initialize()
	})
}

func (w *World) EntityManager() *EntityManager       { return w.entityManager }
func (w *World) ComponentManager() *ComponentManager { return w.componentManager }

func (w *World) AddManager(m Manager) Manager {
	w.managers[m.TypeId()] = m
	w.managersBag.Add(m)
	m.SetWorld(w)
	return m
}

func (w *World) ManagerOfType(t ManagerTypeId) Manager {
	return w.managers[t]
}

func (w *World) DeleteManager(m Manager) {
	delete(w.managers, m.TypeId())
	w.managersBag.RemoveElem(m)
}

func (w *World) Delta() float64     { return w.delta }
func (w *World) SetDelta(d float64) { w.delta = d }

func (w *World) AddEntity(e *Entity)     { w.added.Add(e) }
func (w *World) ChangedEntity(e *Entity) { w.changed.Add(e) }
func (w *World) EnableEntity(e *Entity)  { w.enable.Add(e) }
func (w *World) DisableEntity(e *Entity) { w.disable.Add(e) }

func (w *World) DeleteEntity(e *Entity) {
	if !w.deleted.Contains(e) {
		w.deleted.Add(e)
	}
}

func (w *World) CreateEntity() *Entity {
	return w.entityManager.createEntityInstance()
}

func (w *World) EntityById(id uint) *Entity {
	return w.entityManager.Entity(id)
}

func (w *World) Systems() *util.Bag {
	return w.systemsBag
}

func (w *World) AddActiveSystem(s *EntitySystem) *EntitySystem {
	return w.AddSystem(s, false)
}
func (w *World) AddSystem(s *EntitySystem, passive bool) *EntitySystem {
	s.SetWorld(w)
	s.passive = passive

	w.systems[s.TypeId()] = s
	w.systemsBag.Add(s)
	return s
}

func (w *World) DeleteSystem(s *EntitySystem) {
	delete(w.systems, s.TypeId())
	w.systemsBag.RemoveElem(s)
}

func (w *World) notify(p Performer, e *Entity, bag *util.Bag) {
	bag.ForEach(func(_ int, eo interface{}) {
		p(eo.(EntityObserver), e)
	})
}

func (w *World) notifySystems(p Performer, e *Entity) {
	w.notify(p, e, w.systemsBag)
}
func (w *World) notifyManagers(p Performer, e *Entity) {
	w.notify(p, e, w.managersBag)
}

func (w *World) SystemOfType(t EntitySystemTypeId) *EntitySystem {
	return w.systems[t]
}

func (w *World) check(entities *util.Bag, p Performer) {
	if !entities.Empty() {
		entities.ForEach(func(_ int, e interface{}) {
			w.notifyManagers(p, e.(*Entity))
			w.notifySystems(p, e.(*Entity))
		})
		entities.Clear()
	}
}

func (w *World) Process() {
	w.check(w.added, func(eo EntityObserver, e *Entity) { eo.Added(e) })
	w.check(w.changed, func(eo EntityObserver, e *Entity) { eo.Changed(e) })
	w.check(w.disable, func(eo EntityObserver, e *Entity) { eo.Disabled(e) })
	w.check(w.enable, func(eo EntityObserver, e *Entity) { eo.Enabled(e) })
	w.check(w.deleted, func(eo EntityObserver, e *Entity) { eo.Deleted(e) })

	w.componentManager.clean()

	w.systemsBag.ForEach(func(_ int, si interface{}) {
		s := si.(*EntitySystem)
		if !s.passive {
			s.Process()
		}
	})
}
