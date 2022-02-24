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

// An entity component system for creating games in Go
//
// Creating a New Entity
//
// To create a new entity, first make a new scene. Then, create a entity from the scene.
//	scene := ecs.Scene{}
//	entity := scene.NewEntity()
// Now you can create as many entities as you want.
//
// Adding Components
//
// Components are normal structs storing data.
//  type info struct {
//      num int
//  }
//  component := info{num: 12}
// Then, the component can be added to an entity.
//	ecs.AddComponent(&entity, &component)
// The component can also be retrieved and removed.
//  component2, _ := ecs.GetComponent[info](&entiy) // Get single component from entity
//  ecs.RemoveComponent[info](&entity)              // Remove component from entity
// It is also possible to get all components of a type, which is very useful in systems.
//  components := ecs.AllComponents[info](scene)    // Get all components of same type
//
// Adding Systems
//
// Systems are structs that embed ecs.System and has a Update(deltaTime float64) function.
// Delete and Init functions are optional.
//  type system struct {
//      ecs.System
//  }
//
//  func (sys *system) Update(deltaTime float64) {
//      scene := sys.Scene() // Get the scene to update components
//      // Update logic here
//  }
//
//  func (sys *system) Init() {} // Optional
//
//  func (sys *system) Delete() {} // Optional
// Then, add the system to the scene.
//  scene.AddSystem(system{})
// Scene.Update, Scene.Init and Scene.Delete, calls Update, Init and
// Delete on all systems added to the scene
//  scene.Init()       // Calls system.Init
//  scene.Update(0.01) // Calls system.Update
//  scene.Delete()     // Calls system.Delete
package ecs
