// Copyright 2022 Øystein Berntzen

// Licensed under the Apache License, Version 2.0 (the "License");
// you may not use this file except in compliance with the License.
// You may obtain a copy of the License at

// 	http://www.apache.org/licenses/LICENSE-2.0

// Unless required by applicable law or agreed to in writing, software
// distributed under the License is distributed on an "AS IS" BASIS,
// WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
// See the License for the specific language governing permissions and
// limitations under the License.

package ecs

import (
	"reflect"
)

// Entity is an enitity created by a scene. An entity should only be created from Scene.NewEntity.
type Entity struct {
	id    uint32
	scene *Scene
}

// Scene contains all entities an their components.
type Scene struct {
	componentPools     []pool
	componentIDs       map[reflect.Type]uint32
	currentComponentID uint32
	entityCounter      uint32
}

// NewEntity creates a new entity, and returns it.
func (scene *Scene) NewEntity() Entity {
	scene.entityCounter++
	return Entity{scene.entityCounter - 1, scene}
}

func (scene *Scene) removeEntity(e *Entity) {
	for _, pool := range scene.componentPools {
		pool.remove(e)
	}
}

func (scene *Scene) addComponent(e *Entity, c ComponentInterface) {
	id := scene.getComponentID(c)
	scene.componentPools[id].add(e, c)
}

func (scene *Scene) getComponent(e *Entity, c ComponentInterface) ComponentInterface {
	id := scene.getComponentID(c)
	return scene.componentPools[id].get(e)
}

func (scene *Scene) removeComponent(e *Entity, c ComponentInterface) bool {
	id := scene.getComponentID(c)
	return scene.componentPools[id].remove(e)
}

func (scene *Scene) AllComponents(c ComponentInterface) []ComponentInterface {
	id := scene.getComponentID(c)
	return scene.componentPools[id].components
}

func (scene *Scene) getComponentID(component ComponentInterface) uint32 {
	componentType := reflect.TypeOf(component)
	if scene.componentIDs == nil {
		scene.componentIDs = make(map[reflect.Type]uint32)
	}
	id, ok := scene.componentIDs[componentType]
	if !ok {
		id = scene.currentComponentID
		scene.componentIDs[componentType] = id
		scene.componentPools = append(scene.componentPools, pool{nil, make(map[uint32]uint32)})
		scene.currentComponentID++
	}
	return id
}

// AddComponent adds a new component to the entity, and overwrites if component of this type is already added.
func (entity *Entity) AddComponent(component ComponentInterface) {
	entity.scene.addComponent(entity, component)
}

// GetComponent returns the component of type c of entity e, returns false if component did not exist.
func (entity *Entity) GetComponent(component ComponentInterface) ComponentInterface {
	return entity.scene.getComponent(entity, component)
}

// RemoveComponent removes the component of type of c from the entity, returns false if the component did not exist.
func (entity *Entity) RemoveComponent(component ComponentInterface) bool {
	return entity.scene.removeComponent(entity, component)
}

// RemoveEntity the entity from the scene.
func (entity *Entity) Remove() {
	entity.scene.removeEntity(entity)
}
