package core

import (
	util "github.com/cheng81/go-artemis/util"
)

type Aspect struct {
	allSet  *util.BitSet
	exclSet *util.BitSet
	oneSet  *util.BitSet
}

func newAspect() *Aspect {
	return &Aspect{
		allSet:  util.NewBitSet(64),
		exclSet: util.NewBitSet(64),
		oneSet:  util.NewBitSet(64),
	}
}

func (a *Aspect) All(cTypes ...ComponentTypeId) *Aspect {
	for _, t := range cTypes {
		a.allSet.Set(t.Uint())
	}
	return a
}
func (a *Aspect) Exclude(cTypes ...ComponentTypeId) *Aspect {
	for _, t := range cTypes {
		a.exclSet.Set(t.Uint())
	}
	return a
}
func (a *Aspect) One(cTypes ...ComponentTypeId) *Aspect {
	for _, t := range cTypes {
		a.oneSet.Set(t.Uint())
	}
	return a
}

func AspectFor(cTypes ...ComponentTypeId) *Aspect {
	return AspectForAll(cTypes...)
}
func AspectForAll(cTypes ...ComponentTypeId) *Aspect {
	return newAspect().All(cTypes...)
}
func AspectForOne(cTypes ...ComponentTypeId) *Aspect {
	return newAspect().One(cTypes...)
}
func AspectEmpty() *Aspect {
	return newAspect()
}
