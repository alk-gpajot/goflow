// Code generated by lib-go/goflow DO NOT EDIT.

// +build !codeanalysis

package reduced

import (
	"github.com/alkemics/goflow/example/nodes"
)

/*

 */
type Adder struct{ debug bool }

func NewAdder(debug bool) Adder {
	return Adder{
		debug: debug,
	}
}

func newAdder(id string, debug bool) Adder {
	return Adder{
		debug: debug,
	}
}

/*

 */
func (g *Adder) Run(a, b int) (sum int) {
	// __adder_list outputs
	var __adder_list_aggregated []int

	// __output_sum_builder outputs
	var __output_sum_builder_sum int

	// adder outputs
	var adder_sum int

	// inputs outputs
	var inputs_a int
	var inputs_b int

	igniteNodeID := "ignite"
	doneNodeID := "done"

	done := make(chan string)
	defer close(done)

	steps := map[string]struct {
		deps        map[string]struct{}
		run         func()
		alreadyDone bool
	}{

		"__adder_list": {
			deps: map[string]struct{}{
				"inputs":     {},
				igniteNodeID: {},
			},
			run: func() {
				__adder_list_aggregated = append(__adder_list_aggregated, inputs_a)
				__adder_list_aggregated = append(__adder_list_aggregated, inputs_b)
				done <- "__adder_list"
			},
			alreadyDone: false,
		},
		"__output_sum_builder": {
			deps: map[string]struct{}{
				"adder":      {},
				igniteNodeID: {},
			},
			run: func() {
				__output_sum_builder_sum = adder_sum
				sum = __output_sum_builder_sum
				done <- "__output_sum_builder"
			},
			alreadyDone: false,
		},
		"adder": {
			deps: map[string]struct{}{
				"__adder_list": {},
				igniteNodeID:   {},
			},
			run: func() {
				node := nodes.NewIntReducer(g.debug)
				adder_sum = node.Add(__adder_list_aggregated)
				done <- "adder"
			},
			alreadyDone: false,
		},
		"inputs": {
			deps: map[string]struct{}{
				igniteNodeID: {},
			},
			run: func() {
				inputs_a = a
				inputs_b = b
				done <- "inputs"
			},
			alreadyDone: false,
		},
		igniteNodeID: {
			deps: map[string]struct{}{},
			run: func() {
				done <- igniteNodeID
			},
			alreadyDone: false,
		},
		doneNodeID: {
			deps: map[string]struct{}{
				"__adder_list":         {},
				"__output_sum_builder": {},
				"adder":                {},
				"inputs":               {},
			},
			run: func() {
				done <- doneNodeID
			},
			alreadyDone: false,
		},
	}

	// Ignite
	ignite := steps[igniteNodeID]
	ignite.alreadyDone = true
	steps[igniteNodeID] = ignite
	go steps[igniteNodeID].run()

	// Resolve the graph
	for resolved := range done {
		if resolved == doneNodeID {
			// If all the graph was resolved, get out of the loop
			break
		}

		for name, step := range steps {
			delete(step.deps, resolved)
			if len(step.deps) == 0 && !step.alreadyDone {
				step.alreadyDone = true
				steps[name] = step
				go step.run()
			}
		}
	}

	return sum
}

// MarshalJSON returns the json representation of the graphs. It is pre-generated by
// WithJSONMarshal, and hence never returns an error.
func (g Adder) MarshalJSON() ([]byte, error) {
	return []byte("{\"doc\":\"\",\"edges\":[{\"sourceId\":\"adder\",\"targetId\":\"outputs\",\"inputType\":\"whatever\"},{\"sourceId\":\"inputs\",\"targetId\":\"adder\",\"inputType\":\"[]int\"},{\"sourceId\":\"inputs\",\"targetId\":\"adder\",\"inputType\":\"[]int\"}],\"filename\":\"example/graphs/reduced/adder.yml\",\"id\":\"Adder\",\"inputs\":[{\"name\":\"a\",\"type\":\"int\"},{\"name\":\"b\",\"type\":\"int\"}],\"nodes\":[{\"id\":\"adder\",\"pkg\":\"nodes\",\"type\":\"IntReducer\",\"inputs\":[{\"name\":\"list\",\"type\":\"[]int\"}],\"outputs\":[{\"name\":\"sum\",\"type\":\"int\"}],\"dependencies\":[{\"name\":\"debug\",\"type\":\"bool\"}]}],\"outputs\":[{\"name\":\"sum\",\"type\":\"int\"}],\"pkg\":\"reduced\",\"type\":\"Adder\"}"), nil
}
