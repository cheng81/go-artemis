package core

type ComponentTypeId uint

func (i ComponentTypeId) Uint() uint { return uint(i) }

type Component interface {
	TypeId() ComponentTypeId
}
