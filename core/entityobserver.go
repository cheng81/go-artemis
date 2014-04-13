package core

func EmptyEntityObserver() EntityObserver {
	return emptyObserver
}

type emptyEntityObserver struct{}

var emptyObserver = emptyEntityObserver(struct{}{})

func (_ emptyEntityObserver) Added(_ *Entity)    {}
func (_ emptyEntityObserver) Changed(_ *Entity)  {}
func (_ emptyEntityObserver) Deleted(_ *Entity)  {}
func (_ emptyEntityObserver) Enabled(_ *Entity)  {}
func (_ emptyEntityObserver) Disabled(_ *Entity) {}

type EntityObserver interface {
	Added(*Entity)
	Changed(*Entity)
	Deleted(*Entity)
	Enabled(*Entity)
	Disabled(*Entity)
}
