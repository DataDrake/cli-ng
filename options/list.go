//
// Copyright 2017-2021 Bryan T. Meyers <root@datadrake.com>
//
// Licensed under the Apache License, Version 2.0 (the "License");
// you may not use this file except in compliance with the License.
// You may obtain a copy of the License at
//
//      http://www.apache.org/licenses/LICENSE-2.0
//
// Unless required by applicable law or agreed to in writing, software
// distributed under the License is distributed on an "AS IS" BASIS,
// WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
// See the License for the specific language governing permissions and
// limitations under the License.
//

package options

// List is a Queue of arguments, used druing parsing
type List struct {
	elements []string
}

// NewList creats a list with the included arguments
func NewList(args []string) List {
	return List{
		elements: args,
	}
}

// Peek returns the element at the front of the list
func (l List) Peek() string {
	return l.elements[0]
}

// Next returns the element at the front of the list, after popping it off
func (l *List) Next() (front string) {
	front, l.elements = l.elements[0], l.elements[1:]
	return
}

// IsEmpty checks if there are any remaining elements
func (l List) IsEmpty() bool {
	return len(l.elements) == 0
}
