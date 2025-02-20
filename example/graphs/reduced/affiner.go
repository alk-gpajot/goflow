// Code generated by goflow DO NOT EDIT.

//go:build !codeanalysis

package reduced

import (
	"github.com/alkemics/goflow/example/nodes"
)

/*
 */
type Affiner struct{ debug bool }

func NewAffiner(debug bool) Affiner {
	return Affiner{
		debug: debug,
	}
}

func newAffiner(id string, debug bool) Affiner {
	return Affiner{
		debug: debug,
	}
}

/*
 */
func (g *Affiner) Run(a, b, x int) (result int) {
	// __add_b_a outputs
	var __add_b_a_aggregated int

	// __add_b_b outputs
	var __add_b_b_aggregated int

	// __compute_ax_list outputs
	var __compute_ax_list_aggregated []int

	// __output_result_builder outputs
	var __output_result_builder_result int

	// add_b outputs
	var add_b_sum int

	// compute_ax outputs
	var compute_ax_product int

	// inputs outputs
	var inputs_a int
	var inputs_b int
	var inputs_x int

	igniteNodeID := "ignite"
	doneNodeID := "done"

	done := make(chan string)
	defer close(done)

	steps := map[string]struct {
		deps        map[string]struct{}
		run         func()
		alreadyDone bool
	}{
		"__add_b_a": {
			deps: map[string]struct{}{
				"compute_ax": {},
				igniteNodeID: {},
			},
			run: func() {
				__add_b_a_aggregated = compute_ax_product
				done <- "__add_b_a"
			},
			alreadyDone: false,
		},
		"__add_b_b": {
			deps: map[string]struct{}{
				"inputs":     {},
				igniteNodeID: {},
			},
			run: func() {
				__add_b_b_aggregated = inputs_b
				done <- "__add_b_b"
			},
			alreadyDone: false,
		},
		"__compute_ax_list": {
			deps: map[string]struct{}{
				"inputs":     {},
				igniteNodeID: {},
			},
			run: func() {
				__compute_ax_list_aggregated = append(__compute_ax_list_aggregated, inputs_a)
				__compute_ax_list_aggregated = append(__compute_ax_list_aggregated, inputs_x)
				done <- "__compute_ax_list"
			},
			alreadyDone: false,
		},
		"__output_result_builder": {
			deps: map[string]struct{}{
				"add_b":      {},
				igniteNodeID: {},
			},
			run: func() {
				__output_result_builder_result = add_b_sum
				result = __output_result_builder_result
				done <- "__output_result_builder"
			},
			alreadyDone: false,
		},
		"add_b": {
			deps: map[string]struct{}{
				"__add_b_a":  {},
				"__add_b_b":  {},
				igniteNodeID: {},
			},
			run: func() {
				node := NewAdder(g.debug)
				add_b_sum = node.Run(__add_b_a_aggregated, __add_b_b_aggregated)
				done <- "add_b"
			},
			alreadyDone: false,
		},
		"compute_ax": {
			deps: map[string]struct{}{
				"__compute_ax_list": {},
				igniteNodeID:        {},
			},
			run: func() {
				node := nodes.NewIntReducer(g.debug)
				compute_ax_product = node.Multiply(__compute_ax_list_aggregated)
				done <- "compute_ax"
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
				inputs_x = x
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
				"__add_b_a":               {},
				"__add_b_b":               {},
				"__compute_ax_list":       {},
				"__output_result_builder": {},
				"add_b":                   {},
				"compute_ax":              {},
				"inputs":                  {},
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

	return result
}

// MarshalJSON returns the json representation of the graphs. It is pre-generated by
// WithJSONMarshal, and hence never returns an error.
func (g Affiner) MarshalJSON() ([]byte, error) {
	return []byte("{\"doc\":\"\",\"edges\":[{\"sourceId\":\"add_b\",\"targetId\":\"outputs\",\"inputType\":\"whatever\"},{\"sourceId\":\"compute_ax\",\"targetId\":\"add_b\",\"inputType\":\"int\"},{\"sourceId\":\"inputs\",\"targetId\":\"add_b\",\"inputType\":\"int\"},{\"sourceId\":\"inputs\",\"targetId\":\"compute_ax\",\"inputType\":\"[]int\"},{\"sourceId\":\"inputs\",\"targetId\":\"compute_ax\",\"inputType\":\"[]int\"}],\"filename\":\"graphs/reduced/affiner.yml\",\"id\":\"Affiner\",\"inputs\":[{\"name\":\"a\",\"type\":\"int\"},{\"name\":\"b\",\"type\":\"int\"},{\"name\":\"x\",\"type\":\"int\"}],\"nodes\":[{\"id\":\"add_b\",\"pkg\":\"reduced\",\"type\":\"Adder\",\"inputs\":[{\"name\":\"a\",\"type\":\"int\"},{\"name\":\"b\",\"type\":\"int\"}],\"outputs\":[{\"name\":\"sum\",\"type\":\"int\"}],\"dependencies\":[{\"name\":\"debug\",\"type\":\"bool\"}]},{\"id\":\"compute_ax\",\"pkg\":\"nodes\",\"type\":\"IntReducer\",\"inputs\":[{\"name\":\"list\",\"type\":\"[]int\"}],\"outputs\":[{\"name\":\"product\",\"type\":\"int\"}],\"dependencies\":[{\"name\":\"debug\",\"type\":\"bool\"}]}],\"outputs\":[{\"name\":\"result\",\"type\":\"int\"}],\"pkg\":\"reduced\",\"type\":\"Affiner\"}"), nil
}
