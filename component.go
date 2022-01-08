package ecs

type Component struct {
	entity Entity
}

func (component *Component) Entity() Entity {
	return component.entity
}

func (component *Component) setEntity(e Entity) {
	component.entity = e
}

type ComponentInterface interface {
	Entity() Entity
	setEntity(Entity)
}
