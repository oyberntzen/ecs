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
)

// Entity is an enitity created by a scene. An entity should only be created from Scene.NewEntity.
type Entity struct {
	id    uint32
	scene *Scene
}

// Remove removes the entity from the scene.
func (entity *Entity) Remove() error {
	if entity.scene == nil || entity.id == 0 {
		return errors.New("ecs: entity not registered to a scene (or has been deleted)")
	}
	entity.scene.removeEntity(entity)
	return nil
}
