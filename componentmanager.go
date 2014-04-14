package artemis

import (
	util "github.com/cheng81/go-artemis/util"
	// "reflect"
)

const (
	componentManagerTypeId = ManagerTypeId(0)
)

func NewComponentManager() (out *ComponentManager) {
	out = &ComponentManager{
		DefaultManager:   NewDefaultManager(),
		componentsByType: util.NewBag(64),
		deleted:          util.NewBag(64),
	}
	return
}

type ComponentManager struct {
	*DefaultManager
	componentsByType *util.Bag
	deleted          *util.Bag
}

func (_ *ComponentManager) TypeId() ManagerTypeId {
	return componentManagerTypeId
}

func (t *ComponentManager) removeComponentsOfEntity(e *Entity) {
	componentBits := e.componentBits
	componentBits.WithSetBitsFrom(0, func(index uint) {
		t.componentsByType.GetAt(index).(*util.Bag).SetAt(e.id, nil)
	})
	componentBits.Clear()
}

func (t *ComponentManager) addComponent(e *Entity, c Component) {
	typeId := c.TypeId()
	t.componentsByType.EnsureCapacity(typeId.Uint())

	components := t.componentsByType.GetAt(typeId.Uint())
	if components == nil {
		components = util.NewBag(64)
		t.componentsByType.SetAt(typeId.Uint(), components)
	}

	components.(*util.Bag).SetAt(e.id, c)
	e.componentBits.Set(typeId.Uint())
}

func (t *ComponentManager) removeComponent(e *Entity, typeId ComponentTypeId) {
	if e.componentBits.Get(typeId.Uint()) {
		t.componentsByType.GetAt(typeId.Uint()).(*util.Bag).SetAt(e.id, nil)
		e.componentBits.Unset(typeId.Uint())
	}
}

func (t *ComponentManager) getComponentsByType(typeId ComponentTypeId) *util.Bag {
	out := t.componentsByType.GetAt(typeId.Uint())
	if out == nil {
		out = util.NewBag(64)
		t.componentsByType.SetAt(typeId.Uint(), out.(*util.Bag))
	}
	return out.(*util.Bag)
}

func (t *ComponentManager) getComponent(e *Entity, typeId ComponentTypeId) (out Component) {
	components := t.componentsByType.GetAt(typeId.Uint())
	if components != nil {
		c, ok := components.(*util.Bag).GetAtSafe(e.id)
		if ok {
			return c.(Component)
		}
		// return components.(*util.Bag).GetAt(e.id).(Component)
	}
	return nil
}

func (t *ComponentManager) GetComponentsFor(e *Entity, fill *util.Bag) *util.Bag {
	componentBits := e.componentBits

	componentBits.WithSetBitsFrom(0, func(index uint) {
		fill.Add(t.componentsByType.GetAt(index).(*util.Bag).GetAt(e.id))
	})

	return fill
}

func (t *ComponentManager) clean() {
	for i := uint(0); i < t.deleted.Size(); i++ {
		el := t.deleted.GetAt(i)
		if el != nil {
			t.removeComponentsOfEntity(el.(*Entity))
		}
	}
	t.deleted.Clear()
}

func (t *ComponentManager) Deleted(e *Entity) {
	t.deleted.Add(e)
}
