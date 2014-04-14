package systems

func newIntervalProcessor(interval float64) *IntervalProcessor {
	return &IntervalProcessor{NewBaseProcessor(), 0, interval}
}

type IntervalProcessor struct {
	*BaseProcessor
	acc, interval float64
}

func (iep *IntervalProcessor) CheckProcessing() bool {
	iep.acc += iep.World().Delta()
	if iep.acc > iep.interval {
		iep.acc -= iep.interval
		return true
	}
	return false
}
