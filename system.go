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

type System struct {
	scene *Scene
}

func (system *System) AllComponents(component ComponentInterface) []ComponentInterface {
	return system.scene.allComponents(component)
}

func (system *System) setScene(scene *Scene) {
	system.scene = scene
}

type SystemInterface interface {
	Update(dt float64)

	AllComponents(ComponentInterface) // implemented by System
	setScene(*Scene)                  // implemented by System
}

type InitListener interface {
	SystemInterface
	Init()
}

type DeleteListener interface {
	SystemInterface
	Delete()
}
