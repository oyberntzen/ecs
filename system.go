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

import "reflect"

// System is the base struct for systems, and should be embedded by all systems.
type System struct {
	scene *Scene
}

// AllComponents returns all the components of one type in the scene.
func (system *System) AllComponents(component ComponentInterface) reflect.Value {
	return system.scene.allComponents(component)
}

func (system *System) setScene(scene *Scene) {
	system.scene = scene
}

// SystemInterface is the interface that all systems have to implement.
type SystemInterface interface {
	Update(dt float64)

	AllComponents(ComponentInterface) reflect.Value // implemented by ecs.System
	setScene(*Scene)                                // implemented by ecs.System
}

// InitListener is the interface for systems that has an Init function.
type InitListener interface {
	SystemInterface
	Init()
}

// DeleteListener is the interface for systems that has a Delete function.
type DeleteListener interface {
	SystemInterface
	Delete()
}
