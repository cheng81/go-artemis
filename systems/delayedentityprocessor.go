package systems

import (
	core "github.com/cheng81/go-artemis"
	util "github.com/cheng81/go-artemis/util"
)

func NewDelayedEntitySystem(
	typeId core.EntitySystemTypeId,
	aspect *core.Aspect,
	processor DelayedProcessor) *core.EntitySystem {
	dep := newDelayedEntityProcessor(processor)
	return core.NewEntitySystem(aspect, typeId, dep)
}

func newDelayedEntityProcessor(processor DelayedProcessor) (out *DelayedEntityProcessor) {
	out = &DelayedEntityProcessor{
		EntitySystem: nil,
		delay:        0,
		acc:          0,
		running:      true,
		processor:    processor,
	}
	return
}

type DelayedProcessor interface {
	RemainingDelay(*DelayedEntityProcessor, *core.Entity) float64
	ProcessDelta(*DelayedEntityProcessor, *core.Entity, float64)
	ProcessExpired(*DelayedEntityProcessor, *core.Entity)
}

type DelayedEntityProcessor struct {
	*core.EntitySystem

	delay, acc float64

	running bool

	processor DelayedProcessor
}

func (_ *DelayedEntityProcessor) Initialize() {}

func (p *DelayedEntityProcessor) SetEntitySystem(es *core.EntitySystem) {
	p.EntitySystem = es
}

func (p *DelayedEntityProcessor) ProcessEntities(entities util.ImmutableBag) {
	entities.ForEach(func(_ int, ei interface{}) {
		e := ei.(*core.Entity)
		p.processor.ProcessDelta(p, e, p.acc)
		remaining := p.processor.RemainingDelay(p, e)
		if remaining <= 0 {
			p.processor.ProcessExpired(p, e)
		} else {
			p.OfferDelay(remaining)
		}
	})
	p.Stop()
}

func (p *DelayedEntityProcessor) Inserted(e *core.Entity) {
	delay := p.processor.RemainingDelay(p, e)
	if delay > 0 {
		p.OfferDelay(delay)
	}
}

func (_ *DelayedEntityProcessor) Removed(_ *core.Entity) {}

func (_ *DelayedEntityProcessor) Begin() {}
func (_ *DelayedEntityProcessor) End()   {}

func (p *DelayedEntityProcessor) CheckProcessing() bool {
	if p.running {
		p.acc += p.World().Delta()
		if p.acc > p.delay {
			return true
		}
	}
	return false
}

func (p *DelayedEntityProcessor) Restart(delay float64) {
	p.delay = delay
	p.acc = 0
	p.running = true
}

func (p *DelayedEntityProcessor) OfferDelay(delay float64) {
	if !p.running || delay < p.RemainingTimeUntilProcessing() {
		p.Restart(delay)
	}
}

func (p *DelayedEntityProcessor) InitialDelay() float64 { return p.delay }

func (p *DelayedEntityProcessor) RemainingTimeUntilProcessing() float64 {
	if p.running {
		return p.delay - p.acc
	}
	return 0
}

func (p *DelayedEntityProcessor) Running() bool { return p.running }

func (p *DelayedEntityProcessor) Stop() {
	p.running = false
	p.acc = 0
}
