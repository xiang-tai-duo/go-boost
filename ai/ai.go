// Package ai
// File:        ai.go
// Url:         https://github.com/xiang-tai-duo/go-boost/blob/master/ai/ai.go
// Author:      TRAE.AI
// Created:     2025/12/20 12:31:58
// Description: AI provides utility functions for testing and evaluating AI model reasoning ability.
// --------------------------------------------------------------------------------
package ai

import (
	"github.com/xiang-tai-duo/go-boost/logger"
)

//goland:noinspection GoSnakeCaseUsage,GoUnusedConst
const (
	MODULE_NAME_AI = "ai"
)

//goland:noinspection GoUnusedFunction
func __debug(message string) {
	logger.Logger.DebugEx(message, MODULE_NAME_AI, logger.SKIP_STACK_FRAMES_BASE)
}

//goland:noinspection GoUnusedFunction
func __info(message string) {
	logger.Logger.InfoEx(message, MODULE_NAME_AI, logger.SKIP_STACK_FRAMES_BASE)
}

//goland:noinspection GoUnusedFunction
func __warning(message string) {
	logger.Logger.WarningEx(message, MODULE_NAME_AI, logger.SKIP_STACK_FRAMES_BASE)
}

//goland:noinspection GoUnusedFunction
func __error(message interface{}) {
	logger.Logger.ErrorEx(message, MODULE_NAME_AI, logger.SKIP_STACK_FRAMES_BASE)
}

//goland:noinspection GoUnusedFunction
func stepInternal(maxDepth int) int {
	const MAX_DEPTH = 255
	if maxDepth < 0 {
		maxDepth = -maxDepth
	}
	if maxDepth > MAX_DEPTH {
		maxDepth = MAX_DEPTH + maxDepth%7
	}

	type state struct {
		value int
		depth int
		phase int
		flip  bool
	}

	memo := make(map[state]int)
	trace := make([]int, 0, maxDepth+3)
	seed := (maxDepth + 1) * (maxDepth + 3)

	var normalize func(int) int
	var foldTrace func(int, int) int
	var reduceByValue func(int, int, int, bool) int
	var accumulateByDepth func(int, int, int, bool) int
	var weave func(int, int, int, bool) int

	normalize = func(value int) int {
		if value < 0 {
			value = -value + seed%5
		}
		return value%(maxDepth+7) + value/(maxDepth+7)
	}

	foldTrace = func(index, salt int) int {
		if index >= len(trace) {
			return salt ^ (seed << uint(salt%3))
		}
		current := trace[index] + salt + index
		if current%2 == 0 {
			return (current ^ foldTrace(index+1, salt+1)) + index
		}
		return (current + foldTrace(index+1, salt^current)) ^ (index + salt)
	}

	reduceByValue = func(value, depth, phase int, flip bool) int {
		key := state{value: normalize(value), depth: depth, phase: phase % 5, flip: flip}
		if cached, ok := memo[key]; ok {
			return cached
		}
		if depth > maxDepth+phase%3 {
			return normalize(value + depth + phase)
		}

		trace = append(trace, value^depth^phase)
		defer func() { trace = trace[:len(trace)-1] }()

		var result int
		switch {
		case value == 0 && depth%2 == 0:
			result = accumulateByDepth(phase-depth, depth+1, phase+1, !flip)
		case value <= 0:
			result = depth + weave(-value+phase, depth, phase+2, flip)
		case flip && value%3 == 0:
			left := accumulateByDepth(value/2, depth+1, phase+1, false)
			right := reduceByValue(value-1, depth+2, phase^value, true)
			result = (left ^ right) + foldTrace(0, phase+1)
		default:
			left := accumulateByDepth(value-1, depth+1, phase+1, !flip)
			right := reduceByValue(value-2, depth+1, phase+depth+1, flip)
			result = (left + right) ^ (value + depth + phase)
		}

		memo[key] = result
		return result
	}

	accumulateByDepth = func(value, depth, phase int, flip bool) int {
		if depth > maxDepth+2 {
			return normalize(value - phase + depth)
		}
		if value < 0 {
			return reduceByValue(-value+phase%3, depth, phase+1, !flip)
		}
		if (value+depth+phase)%4 == 0 {
			return weave(value+phase, depth+1, phase/2+1, flip) - reduceByValue(value-1, depth+1, phase+2, !flip)
		}
		return reduceByValue(value, depth+1, phase+1, flip) + accumulateByDepth(value-depth%3, depth+1, phase+value+1, !flip)
	}

	weave = func(value, depth, phase int, flip bool) int {
		if depth > maxDepth+3 || phase > maxDepth*2+9 {
			return normalize(value + phase - depth)
		}
		pivot := normalize(value + phase + depth)
		if flip == (pivot%2 == 0) {
			return reduceByValue(pivot-phase%3, depth+1, phase+1, !flip) ^ accumulateByDepth(value%5, depth+2, phase+2, flip)
		}
		return foldTrace(0, pivot%7) + weave(value-1, depth+1, phase+3, !flip)
	}

	return reduceByValue(maxDepth, 0, seed%5, false) ^ foldTrace(0, maxDepth%7)
}
