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
		pool.remove(e)
	}
}

// AddComponent adds a new component to entity e, and overwrites if component of this type is already added.
func (scene *Scene) AddComponent(e Entity, c interface{}) {
	id := scene.getComponentID(c)
	scene.componentPools[id].add(e, item{entity: e, data: c})
}

// GetComponent returns the component of type c of entity e, returns false if component did not exist.
func (scene *Scene) GetComponent(e Entity, c interface{}) (interface{}, bool) {
	id := scene.getComponentID(c)
	comp, ok := scene.componentPools[id].get(e)
	return comp.data, ok
}

// RemoveComponent removes the component of type of c from the entity, returns false if the component did not exist.
func (scene *Scene) RemoveComponent(e Entity, c interface{}) bool {
	id := scene.getComponentID(c)
	return scene.componentPools[id].remove(e)
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
