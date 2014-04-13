package systems

import (
	// "fmt"
	core "github.com/cheng81/go-artemis/core"
	util "github.com/cheng81/go-artemis/util"
)

func NewEntitySystem(typeId core.EntitySystemTypeId, aspect *core.Aspect, processor EntityProcessor) *core.EntitySystem {
	return core.NewEntitySystem(aspect, typeId, newEntitySystemProcessor(processor))
}

func newEntitySystemProcessor(processor EntityProcessor) core.EntitySystemProcessor {
	return &EntitySystemProcessor{NewBaseProcessor(), processor}
}

// type EntityProcessor func(*core.Entity)

type EntityProcessor interface {
	Process(*core.Entity)
}

type EntitySystemProcessor struct {
	*BaseProcessor
	processor EntityProcessor
}

func (esp *EntitySystemProcessor) ProcessEntities(es util.ImmutableBag) {
	// if uint(esp.TypeId()) == 18 {
	// 	fmt.Println("EntitySystemProcessor.ProcessEntities", es.Size())
	// }
	es.ForEach(func(_ int, ei interface{}) {
		// if uint(esp.TypeId()) == 18 {
		// 	fmt.Println("EntitySystemProcessor.ProcessEntities - entity", ei.(*core.Entity).Id())
		// }
		// esp.processor(ei.(*core.Entity))
		esp.processor.Process(ei.(*core.Entity))
	})
}
