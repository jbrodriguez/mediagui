# Statecraft: A State Machine Library for Go

Statecraft is a state machine engine for Go (golang). State machines
are composed of states and transitions between states.  While
Statecraft can implement finite-state machines (FSM), it can also be
used to devise machines where rules are added over time or have
guarded transitions which exhibit more complicated behavior.

Statecraft `Machine`s are initialized with `NewMachine()` and then
states and transitions are added with the `Rule()` method.  Functions
can be attached to state transitions with the `Action()` method.

Once setup, the `Machine` object transitions between states in
response to `Send()` method calls.

Statecraft also allows you to export your machine into a DOT file for
visualization with GraphViz (`Export()`).

## Installation

Installation uses the standard go get recipe:

    go get bitbucket.org/jdpalmer/statecraft

## Creating Machines

The function `NewMachine()` is used to create a new machine.  Its sole
parameter is the default state of the machine (a string). The `Rule()`
method maps an event to a source state and a target state.

In this example we will model a turnstile with states 'locked' and
'unlocked'. The 'coin' event moves the machine from the 'locked' to
the 'unlocked' state but subsequent coin events keep the machine in
the 'unlocked' state.  The 'push' event moves the machine to the
'locked' state.

    import sc "bitbucket.org/jdpalmer/statecraft"
    import "io/ioutil"

    package main
    
    func main() {
      m := sc.NewMachine("locked")
      m.Rule("push", "locked", "locked")
      m.Rule("coin", "locked", "unlocked")
      m.Rule("push", "unlocked", "locked")
      m.Rule("coin", "unlocked", "unlocked")
    }

We can visualize the resulting machine using the `Export()` method
(you would also need to import "io/ioutil" in this example):

      data := m.Export()
      ioutil.WriteFile("turnstile.dot", []byte(data), 0644)

The resulting DOT file can then be plotted with
[GraphViz](http://graphviz.org/):

![turnstile FSM](http://jdpalmer.org/images/statecraft_turnstile.png)

## Firing Events

Events are fired with the `Send()` method. To transition the
previously defined machine into an 'unlocked' state we could write:

    m.Send("coin")

We can then push the turnstile:

    m.Send("push")

## Adding Actions

Actions can be attached to machines using a descriptive minilanguage
of prefixes:

    >myState  - Evaluate the action as the machine enters myState.
    <myState  - Evaluate the action as the machine leaves myState.
    >*        - Evaluate the action before entering any state.
    <*        - Evaluate the action after leaving any state.
    >>myEvent - Evaluate the action before myEvent.
    <<myEvent - Evaluate the action after myEvent.
    >>*       - Evaluate the action before responding to any event.
    <<*       - Evaluate the action after responding to any event.

Thus, to count the times the turnstile was unlocked we might write
something like:

    m.Action(">unlocked", func() {
      unlock_cnt += 1
    })

The fields `Src`, `Dst`, and `Evt` are exposed in the `Machine`
struct and may be used in actions. We could use these, for example, to
count the number of times a locked turnstile was pushed:

    m.Action("<<push", func() {
      if m.Src == "locked" {
        failed_push_cnt += 1
      }
    })

## Events Within Actions

It is perfectly legal to embed `Send()` calls within actions:

    m.Action(">unlocked", func() {
      m.Send("push")
    })

Statecraft always completes the present transition before considering
events fired within an action.  If multiple `Send()` calls are made
Statecraft drops all but the most recent.  In other words, actions
should typically make no more than one `Send()` call.

Note that it is possible for events triggered within actions to create
an infinite loop (or more accurately a recursive dive).  Consider this
code:

    m.Action(">locked", func() {
      m.Send("coin")
    })
    
    m.Action(">unlocked", func() {
      m.Send("push")
    })

Once triggered these actions will continuously feed coins and push the
turnstile in what is most likely undesirable behavior.

## Handling Errors

Events which do not match a source event pair defined in the machine
are normally ignored. Special actions may be attached to explicitly
handle undefined conditions.  The error prefixes are:

    !myState  - Evaluate the action if the machine is in myState when
                an event is not matched.
    !!myEvent - Evaluate the action if myEvent is not matched.
    !*        - Evaluate the action for all states where an event is
                not matched.
    !!!       - Evaluate the action if and only if the match failed
                and no other error handling code would be evaluated.

If we wanted the machine to panic on any bad input the simplest
solution would be something like:

    m.Action("!!!", func() {
      panic(fmt.Sprint("State", m.Src, "did not recognize event", m.Evt))
    })

## Canceling Transitions

As a machine transitions from its current state to a new state,
actions can be used to abort the transition if the new state has not
been entered.  These prefixes allow cancelation:

* beginning an event (>>)
* leaving a state (<)

Once `Cancel()` is called any pending events are cancelled and no
other actions will be matched after the function returns.

In this action we check a special station_lockdown variable before
leaving the locked state:

    m.Action("<locked", func() {
      if station_lockdown {
        m.Cancel()
      }
    })

If the station is in lockdown we cancel transitions out of the locked
state, otherwise the transition proceeds normally.

## Example: Designing a Menu

Imagine a game where the user can choose to attack, defend, use magic,
or use an item. In creating the user interface you have to keep track
of which item is selected and appropriately highlight the selection as
the user presses up or down.

    import "bitbucket.org/jdpalmer/statecraft"
    import "fmt"
    
    package main
    
    func main() {
      m := NewMachine("attack")
      m.Rule("down", "attack", "defend")
      m.Rule("up", "defend", "attack")
      m.Rule("down", "defend", "magic")
      m.Rule("up", "magic", "defend")
      m.Rule("down", "magic", "item")
      m.Rule("up", "item", "magic")
      // Wrap around to the top!
      m.Rule("down", "item", "attack")
      m.Rule("up", "attack", "item")
      
      m.Action(">*", func() {
        fmt.Println("The selection is now ", m.Dst)
      })
      
      // Fire some events to test it.
      m.Send("down")
      m.Send("up")
      m.Send("down")
      m.Send("down")
      m.Send("down")
      m.Send("down")
      m.Send("down")
      m.Send("up")
    }

The output for this example would be:

	The selection is now defend
	The selection is now attack
	The selection is now defend
	The selection is now magic
	The selection is now item
	The selection is now attack
	The selection is now defend
	The selection is now attack

The GraphViz diagram for this machine can be generated with `Export()`
to show the states graphically:

![turnstile FSM](http://jdpalmer.org/images/statecraft_game.png)
