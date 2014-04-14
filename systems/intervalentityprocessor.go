package systems

import (
	core "github.com/cheng81/go-artemis"
	util "github.com/cheng81/go-artemis/util"
)

func NewIntervalEntitySystem(
	typeId core.EntitySystemTypeId,
	aspect *core.Aspect,
	interval float64,
	processor EntityProcessor) *core.EntitySystem {
	iep := newIntervalEntityProcessor(interval, processor)
	return core.NewEntitySystem(aspect, typeId, iep)
}

func newIntervalEntityProcessor(interval float64, processor EntityProcessor) *IntervalEntityProcessor {
	return &IntervalEntityProcessor{newIntervalProcessor(interval), processor}
}

type IntervalEntityProcessor struct {
	*IntervalProcessor
	processor EntityProcessor
}

func (iep *IntervalEntityProcessor) ProcessEntities(es util.ImmutableBag) {
	es.ForEach(func(_ int, ei interface{}) {
		iep.processor.Process(ei.(*core.Entity))
	})

}
