**goradix** helps flexibly and efficiently organise auto-completion for your application or build space efficient storage with radix-tree like index.

In addition to standard auto-completion **go-radix** provide the way to build wighted auto-completion. Weighted auto-completion is a technics when each node holds the list of output items. 
It is quite useful when the application have to output  filtered or sorted list item auto-completion . For example online book store should output auto-completion book list sorted by popularity rather then alphabetical sort. In high throughput it is better to have predefined output then sorting item lists for each request.

Under the hood **goradix** used [Radix Tree](https://en.wikipedia.org/wiki/Radix_tree) data structure.

## Features of Radix Tree implementation:

- each leaf holds `interface{}` value
- auto-completion for each node can dynamically defined during creation of Radix Tree. 
- applied optimisation for space efficiently is [Adaptive Radix Tree](https://db.in.tum.de/~leis/papers/ART.pdf)

## API 
### Building Radix Tree 
* NewRadixTree creates a new empty radix tree.
```go
    func NewRadixTree() *RadixTree
```
* Insert adds a key-value pair to the tree.
```go
    func (rt *RadixTree) Insert(key string, value interface{})
```
* InsertWithAddSuggestionFunction add a key-pair to the tree. Added value will include to a suggestions set of each upper node. Before add to a suggestion set AddSuggestionFunction will say can a value
 be added to the set.
```go
    func (rt *RadixTree) InsertWithAddSuggestionFunction(key string, value interface{}, p AddSuggestionFunction)
```
* AddSuggestionFunction is a signature of the functions which determines new suggestion set.
```go
    type AddSuggestionFunction func(
        key string, currentSuggestions []*RadixTree, condidate *RadixTree,
    ) []*RadixTree
```
### Quering Radix Tree 
* Find returns a value associated with the given key.
```go
    func (rt *RadixTree) Find(key string) interface{} 
```
* AutoCompleteBroadTraversal returns closest node's values to the given str. Tree traversal algorithms is broadly.
```go
    func (rt *RadixTree) AutoCompleteBroadTraversal(str string, max int) []Suggestion
```
* AutoCompleteDepthTraversal returns closest node's values to the given str. Tree traversal algorithms is depthly.
```go
    func (rt *RadixTree) AutoCompleteDepthTraversal(str string, max int) []Suggestion
```
* ClosestSuggestions returns suggestions set stored in the node which prefix is more closest to the given str.
```go
    func (rt *RadixTree) ClosestSuggestions(str string) []Suggestion 
```
### Printing Radix Tree 
* String returs a basic string representation of the radix tree.
```go
    func (rt *RadixTree) String() string {
```
* StringParentChild returs a string representation of the radix tree. It aims to show parent-child relationships inside the tree.
```go
    func (rt *RadixTree) StringParentChild() string 
```
* StringSuggestions returs a string representation of the radix tree.  It aims to shot suggestions sets accosiated with each node of the tree.
```go
    func (rt *RadixTree) StringSuggestions() string 
```
* StringValues returs a string representation of the radix tree. Is aim to show the radix tree and holded values.
```go
    func (rt *RadixTree) StringValues() string 
```
* NodeWithValueCount returns total count of nodes which holding values.
```go
func (rt *RadixTree) NodeWithValueCount() int {
```
* NodeWithValueCountByCounter returns total count of nodes which holding values. The incoming counter desided how many should be add to a result count.
```go
func (rt *RadixTree) NodeWithValueCountByCounter(
	counter func(interface{}) int,
) int 
```

## Examples
### General AutoComplete
```go
    // create new radix tree
	rt := goradix.NewRadixTree()
    
    // fill the tree with data
	rt.Insert("rube", "one")
	rt.Insert("ruber", 3)
	rt.Insert("rubens", "2+2")
	rt.Insert("rubi", "two")
	rt.Insert("rubicundus", 5)
	rt.Insert("rubicon", 6)

    // Radix Tree:
    // .
    // └──'rub'
    //      ├──'e' (value: one)
    //      │    ├──'r' (value: 3)
    //      │    └──'ns' (value: 2+2)
    //      └──'i' (value: two)
    //           └──'c'
    //                ├──'undus' (value: 5)
    //                └──'on' (value: 6)

    // query the tree
    v := rt.Find("rubens"))
    // v = "2+2"

	kvs := rt.AutoCompleteBroadTraversal("rub", 6)
    // kvs = [{rube one} {rubi two} {ruber 3} {rubens 2+2} {rubicundus 5} {rubicon 6}]

	kvs = rt.AutoCompleteDepthTraversal("rub", 6)
    // kvs = [{rube one} {ruber 3} {rubens 2+2} {rubi two} {rubicundus 5} {rubicon 6}]

```
### Closest Suggestions
```go
    // create new radix tree
	rt := goradix.NewRadixTree()

    // declare AddSuggestionFunction that filter only values with 
	asf := func(
		key string,
		currentSuggestions []*goradix.RadixTree,
		condidate *goradix.RadixTree,
	) []*goradix.RadixTree {
		if condidate.Value() >= 10 {
			return append(currentSuggestions, condidate)
		}

		return currentSuggestions
	}
    
    // fill the tree with data
	rt.InsertWithAddSuggestionFunction("rube", 100, asf)
	rt.InsertWithAddSuggestionFunction("ruber", 200, asf)
	rt.InsertWithAddSuggestionFunction("rubens", 3, asf)
	rt.InsertWithAddSuggestionFunction("rubi", 4, asf)
	rt.InsertWithAddSuggestionFunction("rubicundus", 500, asf)
	rt.InsertWithAddSuggestionFunction("rubicon", 60, asf)

    // Values:
    // .
    // └──'rub'
    //      ├──'e' (value: 100)
    //      │    ├──'r' (value: 200)
    //      │    └──'ns' (value: 3)
    //      └──'i' (value: 4)
    //           └──'c'
    //                ├──'undus' (value: 500)
    //                └──'on' (value: 60)

    // Suggestions:
    // . addr: 4c80 suggs: [4cd0 4d20 4e60 4eb0]
    // └──'rub' (value: <nil>, addr: 4e10, suggestions: [4cd0 4d20 4e60 4eb0])
    //      ├──'e' (value: 100, addr: 4cd0, suggestions: [4cd0 4d20])
    //      │    ├──'r' (value: 200, addr: 4d20, suggestions: [4d20])
    //      │    └──'ns' (value: 3, addr: 4d70, suggestions: [])
    //      └──'i' (value: 4, addr: 4dc0, suggestions: [4e60 4eb0])
    //           └──'c' (value: <nil>, addr: 4f00, suggestions: [4e60 4eb0])
    //                ├──'undus' (value: 500, addr: 4e60, suggestions: [4e60])
    //                └──'on' (value: 60, addr: 4eb0, suggestions: [4eb0])

    // query the tree
	kvs := rt.ClosestSuggestions("rub")
    // kvs = [{rube 100} {ruber 200} {rubicundus 500} {rubicon 60}]
```
