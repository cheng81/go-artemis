package systems

import (
	core "github.com/cheng81/go-artemis"
	util "github.com/cheng81/go-artemis/util"
)

func NewEntitySystem(typeId core.EntitySystemTypeId, aspect *core.Aspect, processor EntityProcessor) *core.EntitySystem {
	return core.NewEntitySystem(aspect, typeId, newEntitySystemProcessor(processor))
}

func newEntitySystemProcessor(processor EntityProcessor) core.EntitySystemProcessor {
	return &EntitySystemProcessor{NewBaseProcessor(), processor}
}

type EntityProcessor interface {
	Process(*core.Entity)
}

type EntitySystemProcessor struct {
	*BaseProcessor
	processor EntityProcessor
}

func (esp *EntitySystemProcessor) ProcessEntities(es util.ImmutableBag) {
	es.ForEach(func(_ int, ei interface{}) {
		esp.processor.Process(ei.(*core.Entity))
	})
}
