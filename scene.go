package ecs3

import (
	"reflect"
)

// Entity is an enitity created by a scene. An entity should only be created from Scene.NewEntity.
type Entity uint32

// Scene contains all entities an their components.
type Scene struct {
	componentPools     []pool
	componentIDs       map[reflect.Type]uint32
	currentComponentID uint32
	entityCounter      Entity
}

// NewEntity creates a new entity, and returns it.
func (scene *Scene) NewEntity() Entity {
	scene.entityCounter++
	return scene.entityCounter - 1
}

// RemoveEntity removes an entity from the scene.
func (scene *Scene) RemoveEntity(e Entity) {
	for _, pool := range scene.componentPools {
		pool.Remove(e)
	}
}

// AddComponent adds a new component to entity e.
func (scene *Scene) AddComponent(e Entity, c interface{}) {
	id := scene.getComponentID(c)
	scene.componentPools[id].Add(e, item{entity: e, data: c})
}

// RemoveComponent removes the component of type of c from the entity, returns false if the component did not exist
func (scene *Scene) RemoveComponent(e Entity, c interface{}) bool {
	id := scene.getComponentID(c)
	return scene.componentPools[id].Remove(e)
}

func (scene *Scene) getComponentID(component interface{}) uint32 {
	componentType := reflect.TypeOf(component)
	if scene.componentIDs == nil {
		scene.componentIDs = make(map[reflect.Type]uint32)
	}
	id, ok := scene.componentIDs[componentType]
	if !ok {
		id = scene.currentComponentID
		scene.componentIDs[componentType] = id
		scene.componentPools = append(scene.componentPools, pool{nil, make(map[Entity]uint32)})
		scene.currentComponentID++
	}
	return id
}
