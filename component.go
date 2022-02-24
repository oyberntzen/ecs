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

// Component is a component with its entity.
// A slice of "Component"s is returned in the AllComponents function.
type Component[T any] struct {
	entity    *Entity
	component T
}

// Entity returns the entity of the component.
func (component *Component[T]) Entity() *Entity {
	return component.entity
}

// Component returns the component data.
func (component *Component[T]) Component() *T {
	return &component.component
}
