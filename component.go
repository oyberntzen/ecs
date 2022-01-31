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

// Component is the base struct for components, and should be embedded by all components.
type Component struct {
	entity Entity
}

func NewComponent(entity Entity) Component {
	return Component{entity: entity}
}

// Entity returns the entity the component is added to.
func (component Component) Entity() *Entity {
	return &component.entity
}

// ComponentInterface is the interface that all components have to implement.
type ComponentInterface interface {
	Entity() *Entity // implemented by ecs.Component
}
