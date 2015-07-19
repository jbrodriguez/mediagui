// Copyright 2014 by James Dean Palmer and others.
//
//  Licensed under the Apache License, Version 2.0 (the "License");
//  you may not use this file except in compliance with the License.
//  You may obtain a copy of the License at
//
//    http://www.apache.org/licenses/LICENSE-2.0
//
//  Unless required by applicable law or agreed to in writing, software
//  distributed under the License is distributed on an "AS IS" BASIS,
//  WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
//  See the License for the specific language governing permissions and
//  limitations under the License.

// Statecraft is a state machine engine for Go. State machines are
// composed of states and transitions between states.  While
// Statecraft can implement finite-state machines (FSM), it can also
// be used to devise machines where transitions are added over time or have
// guarded transitions which exhibit more complicated behavior.
//
// See also: http://bitbucket.org/jdpalmer/statecraft
package statecraft

import "strings"

// The Machine object represents states, transitions and actions for a
// single state machine.  The exported fields Src, Dst, and Evt are
// only defined when an action is being executed and contain,
// respectively, the source state, destination state, and event name.
type Machine struct {
	currentState    string
	transitions     map[string]string
	actions         map[string]func()
	pendingEvent    string
	processingEvent bool
	cancellable     bool
	cancel          bool
	Src             string
	Dst             string
	Evt             string
}

func (self *Machine) reset() {
	self.processingEvent = false
	self.cancellable = false
	self.cancel = false
	self.pendingEvent = ""
	self.Src = ""
	self.Dst = ""
	self.Evt = ""
}

// Create a new Machine with the specified initial state.
func NewMachine(initial string) *Machine {
	self := new(Machine)
	self.transitions = make(map[string]string)
	self.actions = make(map[string]func())
	self.currentState = initial
	self.reset()
	return self
}

// Attach fn to the transition described in specifier.
//
// Specifiers use a special prefix minilanguage to annotate how the
// function should be attached to a transition.  Specifically:
//
//   >myState  - Evaluate the action as the machine enters myState.
//   <myState  - Evaluate the action as the machine leaves myState.
//   >*        - Evaluate the action before entering any state.
//   <*        - Evaluate the action after leaving any state.
//   >>myEvent - Evaluate the action before myEvent.
//   <<myEvent - Evaluate the action after myEvent.
//   >>*       - Evaluate the action before responding to any event.
//   <<*       - Evaluate the action after responding to any event.
//   !myState  - Evaluate the action if the machine is in myState when
//               an event is not matched.
//   !!myEvent - Evaluate the action if myEvent is not matched.
//   !*        - Evaluate the action for all states where an event is not
//               matched.
//   !!!       - Evaluate the action if and only if the match failed
//               and no other error handling code would be evaluated.
func (self *Machine) Action(specifier string, fn func()) {
	self.actions[specifier] = fn
}

// Attempts to cancel an executing event.  If successful the function
// returns true and false otherwise.
func (self *Machine) Cancel() bool {
	if !self.cancellable {
		return false
	}
	self.cancel = true
	return true
}

// Fire an event which may cause the machine to change state.
func (self *Machine) Send(event string) {

	if self.cancel {
		return
	}

	if self.processingEvent {
		self.pendingEvent = event
		return
	}

	self.Src = self.currentState
	self.Evt = event
	nextState, ok := self.transitions[event+"_"+self.currentState]
	self.processingEvent = true

	if ok {
		self.Dst = nextState
		self.cancellable = true

		f := self.actions[">>"+event]
		if f != nil {
			f()
		}
		if self.cancel {
			self.reset()
			return
		}
		f = self.actions[">>*"]
		if f != nil {
			f()
		}
		if self.cancel {
			self.reset()
			return
		}
		f = self.actions["<"+self.currentState]
		if f != nil {
			f()
		}
		if self.cancel {
			self.reset()
			return
		}
		f = self.actions["<*"]
		if f != nil {
			f()
		}
		if self.cancel {
			self.reset()
			return
		}

		self.currentState = nextState
		self.cancellable = false

		f = self.actions[">"+self.currentState]
		if f != nil {
			f()
		}
		f = self.actions[">*"]
		if f != nil {
			f()
		}
		f = self.actions["<<"+event]
		if f != nil {
			f()
		}
		f = self.actions["<<*"]
		if f != nil {
			f()
		}

		self.Dst = ""
	} else {
		cnt := 0
		f := self.actions["!!"+event]
		if f != nil {
			f()
			cnt += 1
		}
		f = self.actions["!"+self.currentState]
		if f != nil {
			f()
			cnt += 1
		}
		f = self.actions["!*"]
		if f != nil {
			f()
			cnt += 1
		}
		if cnt == 0 {
			f = self.actions["!!!"]
			if f != nil {
				f()
			}
		}
	}
	self.Src = ""
	self.Evt = ""

	self.processingEvent = false

	if self.pendingEvent != "" {
		e := self.pendingEvent
		self.pendingEvent = ""
		self.Send(e)
	}

}

// Return a string with a Graphviz DOT representation of the machine.
func (self *Machine) Export() string {
	export := `# dot -Tpng myfile.dot >myfile.png
digraph g {
  rankdir="LR";
  node[style="rounded",shape="box"]
  edge[splines="curved"]`
	export += "\n  " + self.currentState +
		" [style=\"rounded,filled\",fillcolor=\"gray\"]"
	for k, dst := range self.transitions {
		a := strings.SplitN(k, "_", 2)
		event, src := a[0], a[1]
		export += src + " -> " + dst + " [label=\"" + event + "\"];\n"
	}
	export += "}"
	return export
}

// Returns true if state is the current state
func (self *Machine) IsState(state string) bool {
	if self.currentState == state {
		return true
	}
	return false
}

// Returns true if event is a valid event from the current state
func (self *Machine) IsEvent(event string) bool {
	_, ok := self.transitions[event+"_"+self.currentState]
	if ok {
		return true
	}
	return false
}

// Add a transition connecting an event (i.e., an arc or transition)
// between a pair of src and dst states.
func (self *Machine) Rule(event, src, dst string) {
	self.transitions[event+"_"+src] = dst
}
