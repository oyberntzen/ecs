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

const (
	increaseFactor    = 10
	decreaseFactor    = 2
	decreaseThreshold = 3
)

type pool struct {
	components reflect.Value
	indicies   map[uint32]uint32
}

func (p *pool) add(entity *Entity, component ComponentInterface) {
	value := reflect.ValueOf(component)
	if index, ok := p.indicies[entity.id]; ok {
		p.components.Index(int(index)).Set(value)
	}

	length := p.components.Len()
	if p.components.Cap()-length == 0 {
		if length == 0 {
			p.components = reflect.MakeSlice(p.components.Type(), 1, 1) //[]ComponentInterface{component}
			p.components.Index(0).Set(value)
			p.indicies[entity.id] = 0
			return
		}
		newItems := reflect.MakeSlice(p.components.Type(), length+1, length*increaseFactor) //make([]ComponentInterface, length+1, length*increaseFactor)
		//copy(newItems, p.components)
		reflect.Copy(newItems, p.components)
		p.components = newItems
		//p.components[length] = component
		p.components.Index(length).Set(value)
		p.indicies[entity.id] = uint32(length)
		return
	}
	//p.components = p.components[:length+1]
	p.components = p.components.Slice(0, length+1)
	//p.components[length] = component
	p.components.Index(length).Set(value)
	p.indicies[entity.id] = uint32(length)
}

func (p *pool) get(entity *Entity) ComponentInterface {
	index, ok := p.indicies[entity.id]
	if !ok {
		return nil
	}
	return p.components.Index(int(index)).Interface().(ComponentInterface) //p.components[index]
}

func (p *pool) remove(entity *Entity) bool {
	index, ok := p.indicies[entity.id]
	if !ok {
		return false
	}
	delete(p.indicies, entity.id)

	length := p.components.Len() //len(p.components)'
	//p.components[index] = p.components[length-1]
	p.components.Index(int(index)).Set(p.components.Index(length - 1))
	p.components = p.components.Slice(0, length-1) //p.components[:length-1]
	length--

	if uint32(length) > index {
		p.indicies[p.components.Index(int(index)).Interface().(ComponentInterface).Entity().id] = index
	}

	if length*decreaseThreshold < p.components.Cap() {
		newItems := reflect.MakeSlice(p.components.Type(), length, length*decreaseFactor) //make([]ComponentInterface, length, length*decreaseFactor)
		reflect.Copy(newItems, p.components)                                              //copy(newItems, p.components)
		p.components = newItems
	}

	return true
}
