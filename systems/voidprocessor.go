package systems

import (
	core "github.com/cheng81/go-artemis"
	util "github.com/cheng81/go-artemis/util"
)

func NewBaseProcessor() *BaseProcessor {
	return new(BaseProcessor)
}

type BaseProcessor struct {
	*core.EntitySystem
}

func (p *BaseProcessor) SetEntitySystem(es *core.EntitySystem) { p.EntitySystem = es }
func (_ *BaseProcessor) CheckProcessing() bool                 { return true }
func (_ *BaseProcessor) Begin()                                {}
func (_ *BaseProcessor) End()                                  {}
func (_ *BaseProcessor) Initialize()                           {}
func (_ *BaseProcessor) Inserted(_ *core.Entity)               {}
func (_ *BaseProcessor) Removed(_ *core.Entity)                {}
func (_ *BaseProcessor) ProcessEntities(_ util.ImmutableBag)   {}

// this is something different:
// VoidProcessor.java has a processSystem abstract method,
// which is used by systems not really interested into
// entities (a bit dumb, IMHO, but ok)
// so, rename this to BaseProcessor or emptyProcessor

type SystemProcessor func(*core.EntitySystem)

func NewVoidProcessor(proc SystemProcessor) *VoidProcessor {
	return &VoidProcessor{NewBaseProcessor(), proc}
}

type VoidProcessor struct {
	*BaseProcessor
	processor SystemProcessor
}

func (p VoidProcessor) ProcessEntities(_ util.ImmutableBag) {
	p.processor(p.EntitySystem)
}
