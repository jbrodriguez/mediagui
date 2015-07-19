package statecraft

import "testing"
import "io/ioutil"
import "fmt"

func TestTrafficLight(t *testing.T) {
	i := 0
	m := NewMachine("green")
	m.Rule("signal", "green", "yellow")
	m.Rule("signal", "yellow", "red")
	m.Rule("signal", "red", "green")
	m.Action(">red", func() {
		i += 1
	})
	m.Action(">green", func() {
		i += 2
	})
	m.Action(">yellow", func() {
		i += 3
	})
	m.Send("signal")
	if i != 3 {
		t.Error("Bad transition to yellow.")
	}
	m.Send("signal")
	if i != 4 {
		t.Error("Bad transition to red.")
	}
	m.Send("signal")
	if i != 6 {
		t.Error("Bad transition to green.")
	}
}

func TestTurnstile(t *testing.T) {
	i := 0
	m := NewMachine("locked")
	m.Rule("push", "locked", "locked")
	m.Rule("coin", "locked", "unlocked")
	m.Rule("push", "unlocked", "locked")
	m.Rule("coin", "unlocked", "unlocked")
	m.Action(">unlocked", func() {
		i += 1
	})
	m.Send("push")
	m.Send("push")
	m.Send("push")
	m.Send("push")
	m.Send("coin")
	m.Send("push")
	m.Send("push")
	m.Send("push")
	if i != 1 {
		t.Error("Bad accounting.")
	}
	m.Send("coin")
	m.Send("push")
	if i != 2 {
		t.Error("Bad accounting.")
	}
}

func TestCancel(t *testing.T) {
	i := 0
	j := 0
	m := NewMachine("A")
	m.Rule("signal", "A", "B")
	m.Action(">B", func() {
		i += 1
	})
	m.Action("<A", func() {
		m.Cancel()
		j += 1
	})
	m.Send("signal")
	if j != 1 {
		t.Error("Cancel failed.")
	}
	if i != 0 {
		t.Error("Cancel failed.")
	}
}

func TestNestedSends(t *testing.T) {
	i := 0
	m := NewMachine("A")
	m.Rule("signal", "A", "B")
	m.Rule("signal", "B", "C")
	m.Rule("signal", "C", "D")
	m.Action(">B", func() {
		m.Send("signal")
		i += 1
	})
	m.Action(">C", func() {
		m.Send("signal")
		i += 2
	})
	m.Action(">D", func() {
		i += 3
	})
	m.Send("signal")
	if i != 6 {
		t.Error("Nested events failed.")
	}
}

func TestFailedSend(t *testing.T) {
	i := 0
	m := NewMachine("A")
	m.Rule("signal", "A", "B")
	m.Action(">B", func() {
		i += 1
	})
	m.Send("foo")
	if i != 0 {
		t.Error("Failure failed.")
	}
}

func TestFailedSend2(t *testing.T) {
	i := 0
	m := NewMachine("A")
	m.Rule("signal", "A", "B")
	m.Action(">B", func() {
		i += 1
	})
	m.Action("!A", func() {
		i += 10
	})
	m.Action("!*", func() {
		i += 20
	})
	m.Action("!!foo", func() {
		i += 30
	})
	m.Send("foo")
	if i != 60 {
		t.Error("Failure failed.")
	}
}

func TestExport1(t *testing.T) {
	m := NewMachine("green")
	m.Rule("signal", "green", "yellow")
	m.Rule("signal", "yellow", "red")
	m.Rule("signal", "red", "green")
	data := m.Export()
	ioutil.WriteFile("signal.dot", []byte(data), 0644)
}

func TestExport2(t *testing.T) {
	m := NewMachine("locked")
	m.Rule("push", "locked", "locked")
	m.Rule("coin", "locked", "unlocked")
	m.Rule("push", "unlocked", "locked")
	m.Rule("coin", "unlocked", "unlocked")
	data := m.Export()
	ioutil.WriteFile("turnstile.dot", []byte(data), 0644)
}

func BenchmarkCreation(b *testing.B) {
	for i := 0; i < b.N; i++ {
		m := NewMachine("locked")
		m.Rule("push", "locked", "locked")
		m.Rule("coin", "locked", "unlocked")
		m.Rule("push", "unlocked", "locked")
		m.Rule("coin", "unlocked", "unlocked")
	}
}

func BenchmarkSend(b *testing.B) {
	m := NewMachine("locked")
	m.Rule("push", "locked", "locked")
	m.Rule("coin", "locked", "unlocked")
	m.Rule("push", "unlocked", "locked")
	m.Rule("coin", "unlocked", "unlocked")
	for i := 0; i < b.N; i++ {
		m.Send("push")
	}
}

func ExampleGame() {
	m := NewMachine("attack")
	m.Rule("down", "attack", "defend")
	m.Rule("up", "defend", "attack")
	m.Rule("down", "defend", "magic")
	m.Rule("up", "magic", "defend")
	m.Rule("down", "magic", "item")
	m.Rule("up", "item", "magic")
	m.Rule("down", "item", "attack")
	m.Rule("up", "attack", "item")

	m.Action(">*", func() {
		fmt.Println("The selection is now", m.Dst)
	})

	m.Send("down")
	m.Send("up")
	m.Send("down")
	m.Send("down")
	m.Send("down")
	m.Send("down")
	m.Send("down")
	m.Send("up")
	// Output:
	// The selection is now defend
	// The selection is now attack
	// The selection is now defend
	// The selection is now magic
	// The selection is now item
	// The selection is now attack
	// The selection is now defend
	// The selection is now attack
}
