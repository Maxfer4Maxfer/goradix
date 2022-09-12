package goradix

import (
	"fmt"
	"strings"
	"unicode/utf8"
)

// Edge represents connection between a parent node of 
// a radix tree and its child.
type edge struct {
	label     string
	radixTree *RadixTree
	parent    *RadixTree
}

// NewEdge creates a new empty edge.
func newEdge() *edge {
	return &edge{}
}

// SetParent returns a corresponding field of the structure.
func (e *edge) SetParent(parent *RadixTree) *edge {
	e.parent = parent

	return e
}

// Label returns corresponding field of the structure.
func (e *edge) Label() string {
	return e.label
}

// SetLabel sets corresponding field of the structure.
func (e *edge) SetLabel(label string) *edge {
	e.label = label

	return e
}

// RadixTree returns corresponding field of the structure.
func (e *edge) RadixTree() *RadixTree {
	return e.radixTree
}

// SetRadixTree sets corresponding field of the structure.
func (e *edge) SetRadixTree(radixTree *RadixTree) *edge {
	e.radixTree = radixTree

	return e
}

// String returs a string representation of the edge.
func (e *edge) String() string {
	return "'" + e.label + "'" + "addr:" + fmt.Sprintf("%p", e)
}

func (e *edge) stringValues(tabLabel string, tabRadixTree string) string {
	if e.radixTree.value != nil {
		return fmt.Sprintf("%s'%s' (value: %v)\n%s",
			tabLabel,
			e.label,
			e.radixTree.value,
			e.radixTree.stringValues(tabRadixTree+"  "))
	}

	return fmt.Sprintf("%s'%s'\n%s",
		tabLabel,
		e.label,
		e.radixTree.stringValues(tabRadixTree+"  "))
}

func (e *edge) stringSuggestions(tabLabel string, tabRadixTree string) string {
	return fmt.Sprintf("%s'%s' (value: %v, addr: %s, suggestions: %v)\n%s",
		tabLabel,
		e.label,
		e.radixTree.value,
		fmt.Sprintf("%p", e.radixTree)[8:],
		e.radixTree.suggestions,
		e.radixTree.stringSuggestions(tabRadixTree+"  "))
}

func (e *edge) stringParentChild(tabLabel string, tabRadixTree string) string {
	return fmt.Sprintf("%s'%s' parent: %s addr: %s"+
		" (parent: %s, addr: %s, value: %v)\n%s",
		tabLabel,
		e.label,
		fmt.Sprintf("%p", e.parent)[8:],
		fmt.Sprintf("%p", e)[8:],
		fmt.Sprintf("%p", e.radixTree.parent)[8:],
		fmt.Sprintf("%p", e.radixTree)[8:],
		e.radixTree.value,
		e.radixTree.stringParentChild(tabRadixTree+"  "))
}

// RadixTree is a data structure for compact storing strings and values associated
// with each string.
type RadixTree struct {
	parent      *edge
	value       interface{}
	edges       []*edge
	suggestions []*RadixTree
}

// NewRadixTree creates a new empty radix tree.
func NewRadixTree() *RadixTree {
	rt := &RadixTree{}

	return rt
}

func (rt *RadixTree) stringSuggestions(tab string) (out string) {
	for i := range rt.edges {
		tabLabel := tab + "├──"
		tabRadixTree := tab + "│  "

		if i == len(rt.edges)-1 {
			tabLabel = tab + "└──"
			tabRadixTree = tab + "   "
		}

		out += rt.edges[i].stringSuggestions(tabLabel, tabRadixTree)
	}

	return out
}

func (rt *RadixTree) stringParentChild(tab string) (out string) {
	for i := range rt.edges {
		tabLabel := tab + "├──"
		tabRadixTree := tab + "│  "

		if i == len(rt.edges)-1 {
			tabLabel = tab + "└──"
			tabRadixTree = tab + "   "
		}

		out += rt.edges[i].stringParentChild(tabLabel, tabRadixTree)
	}

	return out
}

func (rt *RadixTree) stringValues(tab string) (out string) {
	for i := range rt.edges {
		tabLabel := tab + "├──"
		tabRadixTree := tab + "│  "

		if i == len(rt.edges)-1 {
			tabLabel = tab + "└──"
			tabRadixTree = tab + "   "
		}

		out += rt.edges[i].stringValues(tabLabel, tabRadixTree)
	}

	return out
}

// StringParentChild returs a string representation of the radix tree.
// It aims to show parent-child relationships inside the tree.
func (rt *RadixTree) StringParentChild() string {
	return fmt.Sprintf(". addr: %s \n%s", rt, rt.stringParentChild(""))
}

// StringSuggestions returs a string representation of the radix tree.
// It aims to shot suggestions sets accosiated with each node of the tree.
func (rt *RadixTree) StringSuggestions() string {
	return fmt.Sprintf(". addr: %s suggs: %v\n%s",
		rt, rt.suggestions, rt.stringSuggestions(""))
}

// StringValues returs a string representation of the radix tree.
// Is aim to show the radix tree and holded values.
func (rt *RadixTree) StringValues() string {
	return fmt.Sprintf(". \n%s", rt.stringValues(""))
}

// String returs a basic string representation of the radix tree.
func (rt *RadixTree) String() string {
	return fmt.Sprintf("%p", rt)[8:]
}

// Value returns corresponding value assosiated with rt.
func (rt *RadixTree) Value() interface{} {
	return rt.value
}

// setValue sets corresponding field of the structure.
func (rt *RadixTree) setValue(value interface{}) *RadixTree {
	rt.value = value

	return rt
}

// setParent sets corresponding field of the structure.
func (rt *RadixTree) setParent(parent *edge) *RadixTree {
	rt.parent = parent

	return rt
}

// setEdges sets corresponding field of the structure.
func (rt *RadixTree) setEdges(edges []*edge) *RadixTree {
	rt.edges = edges

	return rt
}

// setSuggestions sets corresponding field of the structure.
func (rt *RadixTree) setSuggestions(s []*RadixTree) *RadixTree {
	c := make([]*RadixTree, len(s))

	copy(c, s)

	rt.suggestions = c

	return rt
}

// AddSuggestionFunction is a signature of the functions which
// determines new suggestion set.
type AddSuggestionFunction func(
	key string, currentSuggestions []*RadixTree, condidate *RadixTree,
) []*RadixTree

// addSuggestion adds the given RadixTree as a suggestion to existed
// suggestions set.
func (rt *RadixTree) addSuggestion(
	key string, next *RadixTree, addSuggestionFunction AddSuggestionFunction,
) *RadixTree {
	if addSuggestionFunction != nil {
		rt.suggestions = addSuggestionFunction(key, rt.suggestions, next)
	}

	return rt
}

// addSuggestions adds the given []*RadixTree as a suggestions to existed
// suggestions set.
func (rt *RadixTree) addSuggestions(
	key string, next []*RadixTree, addSuggestionFunction AddSuggestionFunction,
) *RadixTree {
	for i := range next {
		rt.addSuggestion(key, next[i], addSuggestionFunction)
	}

	return rt
}

// deleteSuggestion deletes suggestions from suggestion set of the current node.
func (rt *RadixTree) deleteSuggestion(s *RadixTree) *RadixTree {
	for i := range rt.suggestions {
		if rt.suggestions[i] == s {
			rt.suggestions = append(
				rt.suggestions[:i], rt.suggestions[i+1:]...,
			)

			return rt
		}
	}

	return rt
}

// NodeWithValueCount returns total count of nodes which holding values.
func (rt *RadixTree) NodeWithValueCount() int {
	var deepDive func(rt *RadixTree, count int) int

	deepDive = func(rt *RadixTree, count int) int {
		if rt.value != nil {
			count++
		}

		for i := range rt.edges {
			count = deepDive(rt.edges[i].radixTree, count)
		}

		return count
	}

	return deepDive(rt, 0)
}

// NodeWithValueCountByCounter returns total count of nodes which holding values.
// The incoming counter desided how many should be add to a result count.
func (rt *RadixTree) NodeWithValueCountByCounter(
	counter func(interface{}) int,
) int {
	var deepDive func(rt *RadixTree, count int) int

	deepDive = func(rt *RadixTree, count int) int {
		if rt.value != nil {
			count += counter(rt.value)
		}

		for i := range rt.edges {
			count = deepDive(rt.edges[i].radixTree, count)
		}

		return count
	}

	return deepDive(rt, 0)
}

// commonPrefix is a helper function returns common prefix of two strings.
func commonPrefix(a string, b string) string {
	prefix := ""

	runeCount := utf8.RuneCountInString(a)
	if runeCount > utf8.RuneCountInString(b) {
		runeCount = utf8.RuneCountInString(b)
	}

	for i := 0; i < runeCount; i++ {
		runeA, _ := utf8.DecodeRuneInString(a[len(prefix):])
		runeB, _ := utf8.DecodeRuneInString(b[len(prefix):])

		if runeA != runeB {
			break
		}

		prefix += string(runeA)
	}

	return prefix
}

// Insert adds a key-value pair to the tree.
func (rt *RadixTree) Insert(key string, value interface{}) {
	rt.insert(key, value, nil)
}

// InsertWithAddSuggestionFunction add a key-pair to the tree.
// Added value will include to a suggestions set of each upper node.
// Before add to a suggestion set AddSuggestionFunction will say can a value
// be added to the set.
func (rt *RadixTree) InsertWithAddSuggestionFunction(
	key string, value interface{}, p AddSuggestionFunction,
) {
	rt.insert(key, value, p)
}

// nolint: funlen
// linter: style with internal helper function for recursion call makes
// it not rational to split the next procedure to into parts.
func (rt *RadixTree) insert(
	key string, value interface{}, addSuggestionFunction AddSuggestionFunction,
) {
	var deleteSuggestion func(rt *RadixTree, sug *RadixTree, key string)

	deleteSuggestion = func(rt *RadixTree, sug *RadixTree, key string) {
		rt.deleteSuggestion(sug)

		// find prefix among the edges
		for i := range rt.edges {
			cPrefix := commonPrefix(key, rt.edges[i].label)

			// key and label are completly different
			if cPrefix == "" {
				continue
			}

			// key: hello  label: he
			if cPrefix == rt.edges[i].label {
				deleteSuggestion(rt.edges[i].radixTree,
					sug, strings.TrimPrefix(key, cPrefix))

				return
			}
		}
	}

	deleteLastAdded := func(dl *RadixTree) {
		deleteSuggestion(rt, dl, key)
	}

	var insert func(
		rt *RadixTree,
		upperKey string,
		key string,
		income *RadixTree,
	)

	insert = func(rt *RadixTree, upperKey string, key string, income *RadixTree) {
		rt.addSuggestion(upperKey, income, addSuggestionFunction)

		// dublicate value! overwrite!
		if key == "" {
			rt.value = income.value

			deleteLastAdded(income)

			return
		}

		// find prefix among the edges
		for i := range rt.edges {
			cPrefix := commonPrefix(key, rt.edges[i].label)

			// key and label are completly different
			if cPrefix == "" {
				continue
			}

			// key: hello  label: he
			if cPrefix == rt.edges[i].label {
				insert(
					rt.edges[i].radixTree,
					upperKey+cPrefix,
					strings.TrimPrefix(key, cPrefix),
					income.setParent(rt.edges[i]),
				)

				return
			}

			// key: he label: hello
			if strings.TrimPrefix(key, cPrefix) == "" {
				nedge := newEdge().
					SetLabel(strings.TrimPrefix(rt.edges[i].label, cPrefix)).
					SetRadixTree(rt.edges[i].radixTree).
					SetParent(income)

				nedge.radixTree.setParent(nedge)

				rt.edges[i].radixTree = income.
					setEdges([]*edge{nedge}).
					addSuggestions(
						upperKey+cPrefix,
						rt.edges[i].radixTree.suggestions,
						addSuggestionFunction).
					setParent(rt.edges[i])

				rt.edges[i].label = cPrefix

				return
			}

			// key: hello label: head
			if strings.TrimPrefix(key, cPrefix) != "" {
				rt1 := rt.edges[i].radixTree
				rt2 := income

				rt.edges[i].radixTree = NewRadixTree().
					setSuggestions(rt1.suggestions).
					addSuggestion(upperKey+cPrefix, rt2, addSuggestionFunction).
					setParent(rt.edges[i])

				edge1 := newEdge().
					SetLabel(strings.TrimPrefix(rt.edges[i].label, cPrefix)).
					SetRadixTree(rt1).
					SetParent(rt.edges[i].radixTree)

				rt1.setParent(edge1)

				edge2 := newEdge().
					SetLabel(strings.TrimPrefix(key, cPrefix)).
					SetRadixTree(rt2).
					SetParent(rt.edges[i].radixTree)

				rt2.setParent(edge2)

				rt.edges[i].radixTree.setEdges([]*edge{edge1, edge2})

				rt.edges[i].label = cPrefix

				return
			}

			return
		}

		// the string has not been meet before
		edge := newEdge().SetLabel(key).SetRadixTree(income).SetParent(rt)
		income.setParent(edge)

		rt.edges = append(rt.edges, edge)
	}

	income := NewRadixTree().setValue(value)

	insert(rt, "", key, income.addSuggestion(key, income, addSuggestionFunction))
}

// Find returns a value associated with the given key.
func (rt *RadixTree) Find(key string) interface{} {
	if key == "" {
		return rt.value
	}

	for i := range rt.edges {
		cPrefix := commonPrefix(key, rt.edges[i].label)

		// key and label are completly different
		if cPrefix == "" {
			continue
		}

		if cPrefix == rt.edges[i].label {
			return rt.edges[i].radixTree.Find(strings.TrimPrefix(key, cPrefix))
		}
	}

	return nil
}

// Suggestion represents key-value pair.
type Suggestion struct {
	Key   string
	Value interface{}
}

// ClosestSuggestions returns suggestions set stored in the node
// which prefix is more closest to the given str.
func (rt *RadixTree) ClosestSuggestions(str string) []Suggestion {
	var reconstructKey func(rt *RadixTree, suffix string) string

	reconstructKey = func(rt *RadixTree, suffix string) string {
		if rt.parent == nil {
			return suffix
		}

		return reconstructKey(
			rt.parent.parent,
			rt.parent.label+suffix,
		)
	}

	createSuggestions := func(rts []*RadixTree) []Suggestion {
		out := make([]Suggestion, len(rts))

		for i := range rts {
			out[i] = Suggestion{
				Key:   reconstructKey(rts[i], ""),
				Value: rts[i].value,
			}
		}

		return out
	}

	var deepDive func(rt *RadixTree, key string) []*RadixTree

	deepDive = func(rt *RadixTree, key string) []*RadixTree {
		if key == "" {
			return rt.suggestions
		}

		// find prefix among the edges
		for i := range rt.edges {
			cPrefix := commonPrefix(key, rt.edges[i].label)

			// cPrefix should meet ether label or prefix
			if cPrefix != key && cPrefix != rt.edges[i].label {
				continue
			}

			// key: he label: hello
			// key: hello  label: hello
			if cPrefix == key || key == rt.edges[i].label {
				return rt.edges[i].radixTree.suggestions
			}

			return deepDive(
				rt.edges[i].radixTree,
				strings.TrimPrefix(key, cPrefix),
			)
		}

		return []*RadixTree{}
	}

	return createSuggestions(deepDive(rt, str))
}

type traversalMode int

const (
	traversalModeBroad traversalMode = iota
	traversalModeDepth
)

// AutoCompleteBroadTraversal returns closest node's values to the given str.
// Tree traversal algorithms is broadly.
func (rt *RadixTree) AutoCompleteBroadTraversal(
	str string, max int,
) []Suggestion {
	return rt.autoCompleteTraversal(str, max, traversalModeBroad)
}

// AutoCompleteDepthTraversal returns closest node's values to the given str.
// Tree traversal algorithms is depthly.
func (rt *RadixTree) AutoCompleteDepthTraversal(
	str string, max int,
) []Suggestion {
	return rt.autoCompleteTraversal(str, max, traversalModeDepth)
}

// nolint: funlen
// linter: style with internal helper function for recursion call makes
// it not rational to split the next procedure to into parts.
func (rt *RadixTree) autoCompleteTraversal(
	str string, max int, traversalMode traversalMode,
) []Suggestion {
	type rtree struct {
		key string
		*RadixTree
	}

	var childrenWithValue func(edges []*edge, prefix string) []rtree

	childrenWithValue = func(edges []*edge, prefix string) []rtree {
		out := []rtree{}

		for i := range edges {
			if edges[i].radixTree.value != nil {
				out = append(out, rtree{
					prefix + edges[i].label,
					edges[i].radixTree,
				})

				continue
			}

			out = append(
				out, childrenWithValue(
					edges[i].radixTree.edges, prefix+edges[i].label)...)
		}

		return out
	}

	var deepDive func(
		rt rtree, todo []rtree, out []Suggestion,
	) []Suggestion

	deepDive = func(
		rt rtree, todo []rtree, out []Suggestion,
	) []Suggestion {
		if rt.value != nil {
			out = append(out, Suggestion{
				Key:   rt.key,
				Value: rt.value,
			})

			if len(out) == max {
				return out
			}
		}

		switch traversalMode {
		case traversalModeBroad:
			todo = append(todo, childrenWithValue(rt.edges, rt.key)...)
		case traversalModeDepth:
			todo = append(childrenWithValue(rt.edges, rt.key), todo...)
		}

		for i := range todo {
			if commonPrefix(str, todo[i].key) != str {
				continue
			}

			return deepDive(todo[i], todo[i+1:], out)
		}

		return out
	}

	return deepDive(rtree{"", rt}, []rtree{}, []Suggestion{})
}
