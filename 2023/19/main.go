package main

import (
	"os"
	"reflect"
	"slices"
	"strconv"
	"strings"
)

// rule = ruleTypOne | ruleTypTwo
type rule struct {
	typOne *ruleTypOne
	typTwo *ruleTypTwo
}

// ruleTypOne (type I) is an accept, reject, or a "step to a new workflow". The
// "step to a new workflow" is encoded as a pointer to the first rule in the
// workflow.
type ruleTypOne struct {
	a, r bool // accept? reject?
	rule *rule
}

// ruleTypTwo (type II) encodes, e.g., "x<1000:foo". The predicate is the first
// part (x<1000), ifTrue is the second (and encoded as the "step to a new
// workflow" is, above), and ifFalse is the rule following this in the workflow.
type ruleTypTwo struct {
	pred    predicate
	ifTrue  *rule
	ifFalse *rule
}

// predicate encodes the predicate in a rule, e.g. "x<1000" will become {"x", "<", 1000}.
type predicate struct {
	attr  string
	op    predicateOp
	value int
}

type predicateOp int

const (
	PREDICATE_OP_LT predicateOp = iota
	PREDICATE_OP_GT
)

type part struct{ x, m, a, s int }

func (p predicate) assessPart(pp part) bool {
	switch p.op {
	case PREDICATE_OP_LT:
		return pp.get(p.attr) < p.value
	case PREDICATE_OP_GT:
		return pp.get(p.attr) > p.value
	default:
		panic("unexpected operation")
	}
}

func (p part) get(attr string) int { return int(reflect.ValueOf(p).FieldByName(attr).Int()) }

func (p part) sum() int { return p.x + p.m + p.a + p.s }

type condition struct {
	pred   predicate
	result bool
}

func partOne(filename string) int {
	rule, parts := parseFile(filename)
	sum := 0
	for _, part := range parts {
		sum += processPart(rule, part)
	}
	return sum
}

func partTwo(filename string) int {
	rule, _ := parseFile(filename)
	return pathsToAcceptance(rule, []condition{})
}

func parseFile(filename string) (*rule, []part) {
	bs, err := os.ReadFile(filename)
	assert(err == nil)
	workflowsStr, partsStr, ok := strings.Cut(strings.TrimSpace(string(bs)), "\n\n")
	assert(ok)
	return newWorkflows(workflowsStr), newParts(partsStr)
}

func newWorkflows(s string) *rule {
	workflowLines := strings.Split(s, "\n")
	nameToRule := make(map[string]*rule, len(workflowLines))
	for _, workflowLine := range workflowLines {
		workflowLine = strings.TrimSuffix(workflowLine, "}")
		name, rulesStr, ok := strings.Cut(workflowLine, "{")
		assert(ok)
		ruleStrs := strings.Split(rulesStr, ",")
		assert(len(ruleStrs) >= 1)
		r, ok := nameToRule[name]
		if !ok {
			r = &rule{}
			nameToRule[name] = r
		}
		rr := newRule(ruleStrs, nameToRule)
		*r = *rr
	}
	rule, ok := nameToRule["in"]
	assert(ok)
	return rule
}

func newRule(ruleStrs []string, nameToRule map[string]*rule) *rule {
	// The last rule (that is, the first in the loop) will be of type I. All
	// others will be of type II. We iteratively correct the ifFalse part of the
	// rules as we work towards the first element of the workflow.
	var curr *rule
	for i := len(ruleStrs) - 1; i >= 0; i-- {
		rule := parseRuleStr(ruleStrs[i], nameToRule)
		if i == len(ruleStrs)-1 {
			curr = rule
			continue
		}
		assert(rule.typTwo != nil)
		rule.typTwo.ifFalse = curr
		curr = rule
	}
	return curr
}

func parseRuleStr(ruleStr string, nameToRule map[string]*rule) *rule {
	predStr, dst, ok := strings.Cut(ruleStr, ":")
	if !ok {
		return &rule{typOne: newRuleTypOne(ruleStr, nameToRule)}
	}

	attr, valueStr, ok := strings.Cut(predStr, "<")
	var op predicateOp
	if ok {
		op = PREDICATE_OP_LT
	} else {
		attr, valueStr, ok = strings.Cut(predStr, ">")
		assert(ok)
		op = PREDICATE_OP_GT
	}

	value, err := strconv.Atoi(valueStr)
	assert(err == nil)
	return &rule{typTwo: &ruleTypTwo{
		pred:   predicate{attr: attr, op: op, value: value},
		ifTrue: &rule{typOne: newRuleTypOne(dst, nameToRule)},
	}}
}

func newRuleTypOne(s string, nameToRule map[string]*rule) *ruleTypOne {
	switch s {
	case "A":
		return &ruleTypOne{a: true}
	case "R":
		return &ruleTypOne{r: true}
	default:
		r, ok := nameToRule[s]
		if !ok {
			r = &rule{}
			nameToRule[s] = r
		}
		return &ruleTypOne{rule: r}
	}
}

func newParts(partsStr string) []part {
	partStrs := strings.Split(partsStr, "\n")
	parts := make([]part, 0, len(partStrs))
	for _, partStr := range partStrs {
		partStr = strings.TrimSuffix(strings.TrimPrefix(partStr, "{"), "}")
		partKVs := strings.Split(partStr, ",")
		part := part{}
		for _, partKV := range partKVs {
			attr, vStr, ok := strings.Cut(partKV, "=")
			assert(ok)
			v, err := strconv.Atoi(vStr)
			assert(err == nil)
			switch attr {
			case "x":
				part.x = v
			case "m":
				part.m = v
			case "a":
				part.a = v
			case "s":
				part.s = v
			default:
				panic("unexpected attr")
			}
		}
		parts = append(parts, part)
	}
	return parts
}

func processPart(rule *rule, part part) int {
	for {
		if rule.typOne != nil {
			if rule.typOne.a {
				return part.sum()
			}
			if rule.typOne.r {
				return 0
			}
			rule = rule.typOne.rule
			continue
		}

		if rule.typTwo.pred.assessPart(part) {
			assert(rule.typTwo.ifTrue.typOne != nil)
			if rule.typTwo.ifTrue.typOne.a {
				return part.sum()
			}
			if rule.typTwo.ifTrue.typOne.r {
				return 0
			}
			rule = rule.typTwo.ifTrue.typOne.rule
			continue
		}

		rule = rule.typTwo.ifFalse
		continue
	}
}

func pathsToAcceptance(rule *rule, conds []condition) int {
	if rule.typOne != nil {
		if rule.typOne.a {
			return satCount(conds)
		}
		if rule.typOne.r {
			return 0
		}
		return pathsToAcceptance(rule.typOne.rule, conds)
	}

	trueConds := make([]condition, len(conds))
	copy(trueConds, conds)
	trueConds = append(trueConds, condition{pred: rule.typTwo.pred, result: true})

	falseConds := make([]condition, len(conds))
	copy(falseConds, conds)
	falseConds = append(falseConds, condition{pred: rule.typTwo.pred})

	return pathsToAcceptance(rule.typTwo.ifTrue, trueConds) + pathsToAcceptance(rule.typTwo.ifFalse, falseConds)
}

func satCount(conds []condition) int {
	// These will map to closed intervals.
	xmas := []string{"x", "m", "a", "s"}
	xmasMin := []int{1, 1, 1, 1}
	xmasMax := []int{4_000, 4_000, 4_000, 4_000}

	for _, cond := range conds {
		idx := slices.Index(xmas, cond.pred.attr)
		assert(idx != -1)

		switch cond.pred.op {
		case PREDICATE_OP_LT:
			if cond.result {
				// e.g., x<value.
				xmasMax[idx] = min(xmasMax[idx], cond.pred.value-1)
			} else {
				// e.g., x>=value.
				xmasMin[idx] = max(xmasMin[idx], cond.pred.value)
			}
		case PREDICATE_OP_GT:
			if cond.result {
				// e.g., x>value.
				xmasMin[idx] = max(xmasMin[idx], cond.pred.value+1)
			} else {
				// e.g., x<=value.
				xmasMax[idx] = min(xmasMax[idx], cond.pred.value)
			}
		default:
			panic("unexpected op")
		}
	}

	acc := 1
	for i := range xmasMin {
		factor := xmasMax[i] - xmasMin[i] + 1
		if factor <= 0 {
			return 0
		}
		acc *= factor
	}
	return acc
}

func assert(b bool) {
	if !b {
		panic("assertion failed")
	}
}
