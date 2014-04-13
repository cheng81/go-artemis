package core

type ManagerTypeId uint

type Manager interface {
	EntityObserver
	TypeId() ManagerTypeId
	Initialize()
	SetWorld(*World)
	World() *World
}

func NewDefaultManager() *DefaultManager {
	return &DefaultManager{
		EmptyEntityObserver(),
		nil,
	}
}

type DefaultManager struct {
	EntityObserver
	world *World
}

func (m *DefaultManager) Initialize()       {}
func (m *DefaultManager) SetWorld(w *World) { m.world = w }
func (m *DefaultManager) World() *World     { return m.world }
