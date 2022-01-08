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

type item struct {
	entity Entity
	data   interface{}
}

type pool struct {
	items    []item
	indicies map[Entity]uint32
}

func (p *pool) add(e Entity, i item) {
	if index, ok := p.indicies[e]; ok {
		p.items[index] = i
	}

	length := len(p.items)
	if cap(p.items)-length == 0 {
		if length == 0 {
			p.items = []item{i}
			p.indicies[e] = 0
			return
		}
		newItems := make([]item, length+1, length*increaseFactor)
		copy(newItems, p.items)
		p.items = newItems
		p.items[length] = i
		p.indicies[e] = uint32(length)
		return
	}
	p.items = p.items[:length+1]
	p.items[length] = i
	p.indicies[e] = uint32(length)
}

func (p *pool) get(e Entity) (item, bool) {
	index, ok := p.indicies[e]
	if !ok {
		return item{}, false
	}
	return p.items[index], true
}

func (p *pool) remove(e Entity) bool {
	index, ok := p.indicies[e]
	if !ok {
		return false
	}
	delete(p.indicies, e)

	length := len(p.items)
	p.items[index] = p.items[length-1]
	p.items = p.items[:length-1]
	length--

	if uint32(length) > index {
		p.indicies[p.items[index].entity] = index
	}

	if length*decreaseThreshold < cap(p.items) {
		newItems := make([]item, length, length*decreaseFactor)
		copy(newItems, p.items)
		p.items = newItems
	}

	return true
}
