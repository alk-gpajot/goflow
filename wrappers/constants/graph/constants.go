// Code generated by goflow DO NOT EDIT.

//go:build !codeanalysis
// +build !codeanalysis

package main

import (
	"github.com/alkemics/goflow/example/constants/functions"
	"github.com/alkemics/goflow/example/constants/numbers"
	"github.com/alkemics/goflow/example/nodes"
)

/*

 */
type Constants struct{}

func NewConstants() Constants {
	return Constants{}
}

func newConstants(id string) Constants {
	return Constants{}
}

/*

 */
func (g *Constants) Run() {

	// __add_list outputs
	var __add_list_aggregated []int

	// __add_reducer outputs
	var __add_reducer_aggregated functions.IntReducer

	// __multiply_list outputs
	var __multiply_list_aggregated []int

	// __multiply_reducer outputs
	var __multiply_reducer_aggregated functions.IntReducer

	// __print_print outputs
	var __print_print_aggregated bool

	// __print_values outputs
	var __print_values_aggregated []interface{}

	// __print_void_print outputs
	var __print_void_print_aggregated bool

	// __print_void_values outputs
	var __print_void_values_aggregated []interface{}

	// add outputs
	var add_result int

	// multiply outputs
	var multiply_result int

	// print outputs

	// print_void outputs

	igniteNodeID := "ignite"
	doneNodeID := "done"

	done := make(chan string)
	defer close(done)

	steps := map[string]struct {
		deps        map[string]struct{}
		run         func()
		alreadyDone bool
	}{

		"__add_list": {
			deps: map[string]struct{}{
				igniteNodeID: {},
			},
			run: func() {
				__add_list_aggregated = append(__add_list_aggregated, numbers.One)
				__add_list_aggregated = append(__add_list_aggregated, numbers.Two)
				__add_list_aggregated = append(__add_list_aggregated, 3)
				done <- "__add_list"
			},
			alreadyDone: false,
		},
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
		"__multiply_list": {
			deps: map[string]struct{}{
				"add":        {},
				igniteNodeID: {},
			},
			run: func() {
				__multiply_list_aggregated = append(__multiply_list_aggregated, add_result)
				__multiply_list_aggregated = append(__multiply_list_aggregated, 10)
				done <- "__multiply_list"
			},
			alreadyDone: false,
		},
		"__multiply_reducer": {
			deps: map[string]struct{}{
				igniteNodeID: {},
			},
			run: func() {
				__multiply_reducer_aggregated = functions.IntMultiplication
				done <- "__multiply_reducer"
			},
			alreadyDone: false,
		},
		"__print_print": {
			deps: map[string]struct{}{
				igniteNodeID: {},
			},
			run: func() {
				__print_print_aggregated = true
				done <- "__print_print"
			},
			alreadyDone: false,
		},
		"__print_values": {
			deps: map[string]struct{}{
				"multiply":   {},
				igniteNodeID: {},
			},
			run: func() {
				__print_values_aggregated = append(__print_values_aggregated, "the result should be 60")
				__print_values_aggregated = append(__print_values_aggregated, multiply_result)
				done <- "__print_values"
			},
			alreadyDone: false,
		},
		"__print_void_print": {
			deps: map[string]struct{}{
				igniteNodeID: {},
			},
			run: func() {
				__print_void_print_aggregated = false
				done <- "__print_void_print"
			},
			alreadyDone: false,
		},
		"__print_void_values": {
			deps: map[string]struct{}{
				igniteNodeID: {},
			},
			run: func() {
				__print_void_values_aggregated = append(__print_void_values_aggregated, "never printed")
				done <- "__print_void_values"
			},
			alreadyDone: false,
		},
		"add": {
			deps: map[string]struct{}{
				"__add_list":    {},
				"__add_reducer": {},
				igniteNodeID:    {},
			},
			run: func() {
				add_result = nodes.IntAggregator(__add_list_aggregated, __add_reducer_aggregated)
				done <- "add"
			},
			alreadyDone: false,
		},
		"multiply": {
			deps: map[string]struct{}{
				"__multiply_list":    {},
				"__multiply_reducer": {},
				igniteNodeID:         {},
			},
			run: func() {
				multiply_result = nodes.IntAggregator(__multiply_list_aggregated, __multiply_reducer_aggregated)
				done <- "multiply"
			},
			alreadyDone: false,
		},
		"print": {
			deps: map[string]struct{}{
				"__print_print":  {},
				"__print_values": {},
				igniteNodeID:     {},
			},
			run: func() {
				nodes.ConditionalPrinter(__print_print_aggregated, __print_values_aggregated)
				done <- "print"
			},
			alreadyDone: false,
		},
		"print_void": {
			deps: map[string]struct{}{
				"__print_void_print":  {},
				"__print_void_values": {},
				igniteNodeID:          {},
			},
			run: func() {
				nodes.ConditionalPrinter(__print_void_print_aggregated, __print_void_values_aggregated)
				done <- "print_void"
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
				"__add_list":          {},
				"__add_reducer":       {},
				"__multiply_list":     {},
				"__multiply_reducer":  {},
				"__print_print":       {},
				"__print_values":      {},
				"__print_void_print":  {},
				"__print_void_values": {},
				"add":                 {},
				"multiply":            {},
				"print":               {},
				"print_void":          {},
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
