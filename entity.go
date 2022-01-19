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

// Entity is an enitity created by a scene. An entity should only be created from Scene.NewEntity.
type Entity struct {
	id    uint32
	scene *Scene
}

// AddComponent adds a new component to the entity, and overwrites if component of this type is already added.
func (entity *Entity) AddComponent(component ComponentInterface) error {
	if entity.scene == nil || entity.id == 0 {
		return errors.New("ecs: entity not registered to a scene (or has been deleted)")
	}
	entity.scene.addComponent(entity, component)
	return nil
}

// GetComponent returns the component of type c of entity e, returns false if component did not exist.
func (entity *Entity) GetComponent(component ComponentInterface) (ComponentInterface, error) {
	if entity.scene == nil || entity.id == 0 {
		return nil, errors.New("ecs: entity not registered to a scene (or has been deleted)")
	}
	result := entity.scene.getComponent(entity, component)
	if result == nil {
		return nil, fmt.Errorf("ecs: no component of type %s added to entity", reflect.TypeOf(component))
	}
	return result, nil
}

// RemoveComponent removes the component of type of c from the entity, returns false if the component did not exist.
func (entity *Entity) RemoveComponent(component ComponentInterface) error {
	if entity.scene == nil || entity.id == 0 {
		return errors.New("ecs: entity not registered to a scene (or has been deleted)")
	}
	if !entity.scene.removeComponent(entity, component) {
		return fmt.Errorf("ecs: no component of type %s added to entity", reflect.TypeOf(component))
	}
	return nil
}

// Remove removes the entity from the scene.
func (entity *Entity) Remove() error {
	if entity.scene == nil || entity.id == 0 {
		return errors.New("ecs: entity not registered to a scene (or has been deleted)")
	}
	entity.scene.removeEntity(entity)
	return nil
}
