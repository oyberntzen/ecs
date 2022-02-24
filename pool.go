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

const (
	increaseFactor    = 2
	decreaseFactor    = 2
	decreaseThreshold = 3
)

type pool[T any] struct {
	components []Component[T]
	indicies   map[uint32]uint32
}

type poolInterface interface {
	remove(entity *Entity) bool
}

func (p *pool[T]) add(entity *Entity, data *T) {
	if index, ok := p.indicies[entity.id]; ok {
		p.components[index] = Component[T]{entity, *data}
	}

	length := len(p.components)
	if cap(p.components)-length == 0 {
		if length == 0 {
			p.components = []Component[T]{{entity, *data}}
			p.indicies[entity.id] = 0
			return
		}
		newItems := make([]Component[T], length+1, length*increaseFactor)
		copy(newItems, p.components)
		p.components = newItems
		p.components[length] = Component[T]{entity, *data}
		p.indicies[entity.id] = uint32(length)
		return
	}
	p.components = p.components[:length+1]
	p.components[length] = Component[T]{entity, *data}
	p.indicies[entity.id] = uint32(length)
}

func (p *pool[T]) get(entity *Entity) *T {
	index, ok := p.indicies[entity.id]
	if !ok {
		return nil
	}
	return p.components[index].Component()
}

func (p *pool[T]) remove(entity *Entity) bool {
	index, ok := p.indicies[entity.id]
	if !ok {
		return false
	}
	delete(p.indicies, entity.id)

	length := len(p.components)
	p.components[index] = p.components[length-1]
	p.components = p.components[:length-1]
	length--

	if uint32(length) > index {
		p.indicies[p.components[index].Entity().id] = index
	}

	if length*decreaseThreshold < cap(p.components) {
		newItems := make([]Component[T], length, length*decreaseFactor)
		copy(newItems, p.components)
		p.components = newItems
	}

	return true
}
