// Code generated by goflow DO NOT EDIT.

//go:build !codeanalysis
// +build !codeanalysis

package main

import (
	"github.com/alkemics/goflow/example/constants/functions"
	"github.com/alkemics/goflow/example/nodes"
)

/*

 */
type Bind struct{}

func NewBind() Bind {
	return Bind{}
}

func newBind(id string) Bind {
	return Bind{}
}

/*

 */
func (g *Bind) Run() {

	// __add_reducer outputs
	var __add_reducer_aggregated functions.IntReducer

	// __print_values outputs
	var __print_values_aggregated []interface{}

	// add outputs
	var add_result int

	// make_slice outputs
	var make_slice_list []int

	// print outputs

	igniteNodeID := "ignite"
	doneNodeID := "done"

	done := make(chan string)
	defer close(done)

	steps := map[string]struct {
		deps        map[string]struct{}
		run         func()
		alreadyDone bool
	}{

		"__add_reducer": {
			deps: map[string]struct{}{
				igniteNodeID: {},
			},
			run: func() {
				__add_reducer_aggregated = functions.IntSum
				done <- "__add_reducer"
			},
			alreadyDone: false,
		},
		"__print_values": {
			deps: map[string]struct{}{
				"make_slice": {},
				"add":        {},
				igniteNodeID: {},
			},
			run: func() {
				__print_values_aggregated = append(__print_values_aggregated, "should print 1, 2, 3 and 0")
				for _, e := range make_slice_list {
					__print_values_aggregated = append(__print_values_aggregated, e)
				}
				__print_values_aggregated = append(__print_values_aggregated, add_result)
				done <- "__print_values"
			},
			alreadyDone: false,
		},
		"add": {
			deps: map[string]struct{}{
				"__add_reducer": {},
				igniteNodeID:    {},
			},
			run: func() {
				var list []int
				add_result = nodes.IntAggregator(list, __add_reducer_aggregated)
				done <- "add"
			},
			alreadyDone: false,
		},
		"make_slice": {
			deps: map[string]struct{}{
				igniteNodeID: {},
			},
			run: func() {
				make_slice_list = nodes.SliceMaker()
				done <- "make_slice"
			},
			alreadyDone: false,
		},
		"print": {
			deps: map[string]struct{}{
				"__print_values": {},
				igniteNodeID:     {},
			},
			run: func() {
				nodes.Printer(__print_values_aggregated)
				done <- "print"
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
				"__add_reducer":  {},
				"__print_values": {},
				"add":            {},
				"make_slice":     {},
				"print":          {},
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

	return
}
