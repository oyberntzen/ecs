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

// Scene contains all entities, components and systems.
type Scene struct {
	entityCounter uint32

	componentPools     []pool
	componentIDs       map[reflect.Type]uint32
	currentComponentID uint32

	systems []SystemInterface
}

// NewEntity creates a new entity, and returns it.
func (scene *Scene) NewEntity() Entity {
	scene.entityCounter++
	return Entity{scene.entityCounter - 1, scene}
}

func (scene *Scene) AddSystem(system SystemInterface) {
	system.setScene(scene)
	scene.systems = append(scene.systems, system)
}

func (scene *Scene) Init() {
	for _, system := range scene.systems {
		if initSystem, ok := system.(InitListener); ok {
			initSystem.Init()
		}
	}
}

func (scene *Scene) Update(dt float64) {
	for _, system := range scene.systems {
		system.Update(dt)
	}
}

func (scene *Scene) Delete() {
	for _, system := range scene.systems {
		if deleteSystem, ok := system.(DeleteListener); ok {
			deleteSystem.Delete()
		}
	}
}

func (scene *Scene) removeEntity(entity *Entity) {
	for _, pool := range scene.componentPools {
		pool.remove(entity)
	}
}

func (scene *Scene) addComponent(entity *Entity, component ComponentInterface) {
	id := scene.getComponentID(component)
	scene.componentPools[id].add(entity, component)
}

func (scene *Scene) getComponent(entity *Entity, component ComponentInterface) ComponentInterface {
	id := scene.getComponentID(component)
	return scene.componentPools[id].get(entity)
}

func (scene *Scene) removeComponent(entity *Entity, component ComponentInterface) bool {
	id := scene.getComponentID(component)
	return scene.componentPools[id].remove(entity)
}

func (scene *Scene) allComponents(component ComponentInterface) []ComponentInterface {
	id := scene.getComponentID(component)
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
