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
	"errors"
	"fmt"
	"reflect"
)

// Scene contains all entities, components and systems.
type Scene struct {
	entityCounter uint32

	componentPools     []poolInterface
	componentIDs       map[reflect.Type]uint32
	currentComponentID uint32

	systems []SystemInterface
}

// NewEntity creates a new entity, and returns it.
func (scene *Scene) NewEntity() Entity {
	scene.entityCounter++
	return Entity{scene.entityCounter, scene}
}

// AddSystem adds the system to the scene.
func (scene *Scene) AddSystem(system SystemInterface) {
	system.setScene(scene)
	scene.systems = append(scene.systems, system)
}

// Init calls Init functions on all systems.
func (scene *Scene) Init() {
	for _, system := range scene.systems {
		if initSystem, ok := system.(InitListener); ok {
			initSystem.Init()
		}
	}
}

// Update calls Update functions on all systems.
func (scene *Scene) Update(dt float64) {
	for _, system := range scene.systems {
		system.Update(dt)
	}
}

// Delete calls Delete functions on all systems.
func (scene *Scene) Delete() {
	for _, system := range scene.systems {
		if deleteSystem, ok := system.(DeleteListener); ok {
			deleteSystem.Delete()
		}
	}
}

func (scene *Scene) removeEntity(entity *Entity) {
	entity.id = 0
	entity.scene = nil
	for _, pool := range scene.componentPools {
		pool.remove(entity)
	}
}

func AllComponents[T ComponentInterface](scene *Scene) []T {
	id := getComponentID[T](scene)
	componentPool, ok := scene.componentPools[id].(*pool[T])
	if !ok {
		panic("Not working")
	}
	return componentPool.components
}

// AddComponent adds a new component to the entity, and overwrites if component of this type is already added.
func AddComponent[T ComponentInterface](component *T) error {
	if (*component).Entity().scene == nil || (*component).Entity().id == 0 {
		return errors.New("ecs: entity not registered to a scene (or has been deleted)")
	}

	id := getComponentID[T]((*component).Entity().scene)
	componentPool, ok := (*component).Entity().scene.componentPools[id].(*pool[T])
	if !ok {
		panic("Not working")
	}
	componentPool.add((*component).Entity(), component)

	return nil
}

// GetComponent returns the component of type c of entity e, returns false if component did not exist.
func GetComponent[T ComponentInterface](entity *Entity) (*T, error) {
	if entity.scene == nil || entity.id == 0 {
		return nil, errors.New("ecs: entity not registered to a scene (or has been deleted)")
	}

	id := getComponentID[T](entity.scene)
	componentPool, ok := entity.scene.componentPools[id].(*pool[T])
	if !ok {
		panic("Not working")
	}

	result := componentPool.get(entity)
	if result == nil {
		return nil, fmt.Errorf("ecs: no component of type %s added to entity", reflect.TypeOf(new(T)))
	}
	return result, nil
}

// RemoveComponent removes the component of type of c from the entity, returns false if the component did not exist.
func RemoveComponent[T ComponentInterface](entity *Entity) error {
	if entity.scene == nil || entity.id == 0 {
		return errors.New("ecs: entity not registered to a scene (or has been deleted)")
	}
	id := getComponentID[T](entity.scene)
	if !entity.scene.componentPools[id].remove(entity) {
		return fmt.Errorf("ecs: no component of type %s added to entity", reflect.TypeOf(new(T)))
	}
	return nil
}

func getComponentID[T ComponentInterface](scene *Scene) uint32 {
	componentType := reflect.TypeOf(new(T))
	if scene.componentIDs == nil {
		scene.componentIDs = make(map[reflect.Type]uint32)
	}
	id, ok := scene.componentIDs[componentType]
	if !ok {
		id = scene.currentComponentID
		scene.componentIDs[componentType] = id
		scene.componentPools = append(scene.componentPools, &pool[T]{make([]T, 0, 1), make(map[uint32]uint32)})
		scene.currentComponentID++
	}
	return id
}
