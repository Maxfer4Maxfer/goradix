package main

import (
	"fmt"
	"strings"

	"github.com/Maxfer4Maxfer/goradix"
)

func example1() {
	rt := goradix.NewRadixTree()

	rt.Insert("romane", 1)
	rt.Insert("romanus", 2)
	rt.Insert("romulus", 3)
	rt.Insert("ruber", "four")
	rt.Insert("rubens", 5)
	rt.Insert("rubicundus", "six")
	rt.Insert("rubicon", 7)
	rt.Insert("toasting", 88)
	rt.Insert("toast", "sexteen")
	rt.Insert("toaster", 99)

	fmt.Println("Parent-Child:")
	fmt.Println(rt.StringParentChild())
	fmt.Println("--------------------------------")
	fmt.Println("Suggestions:")
	fmt.Println(rt.StringSuggestions())
	fmt.Println("--------------------------------")
	fmt.Println("Values:")
	fmt.Println(rt.StringValues())
	fmt.Println("--------------------------------")
	fmt.Printf("nodes with value %v\n", rt.NodeWithValueCount())
	fmt.Printf("nodes with value string %v\n", rt.NodeWithValueCountByCounter(
		func(v interface{}) int {
			switch v.(type) {
			case string:
				return 1
			default:
				return 0
			}
		},
	))

	// VALUES
	fmt.Println("--------------------------------")
	fmt.Println("Values:")
	fmt.Printf("rubens: %v\n", rt.Find("rubens"))
	fmt.Printf("toast: %v\n", rt.Find("toast"))

	fmt.Println("--------------------------------")
	fmt.Println("AutoCompleteBroadTraversal:")
	fmt.Printf("rub: %v\n", rt.AutoCompleteBroadTraversal("rub", 4))
	fmt.Println("ClosestSuggestions:")
	fmt.Printf("rub: %v\n", rt.ClosestSuggestions("rub"))
}

func example2() {
	rt := goradix.NewRadixTree()

	asf := func(
		key string,
		currentSuggestions []*goradix.RadixTree,
		condidate *goradix.RadixTree,
	) []*goradix.RadixTree {
		if strings.HasPrefix(key, "rubi") {
			return append(currentSuggestions, condidate)
		}

		return currentSuggestions
	}

	rt.InsertWithAddSuggestionFunction("romane", 1, asf)
	rt.InsertWithAddSuggestionFunction("romanus", "two", asf)
	rt.InsertWithAddSuggestionFunction("romulus", 3, asf)
	rt.InsertWithAddSuggestionFunction("rubi", "2^2", asf)
	rt.InsertWithAddSuggestionFunction("rubicon", 5, asf)
	rt.InsertWithAddSuggestionFunction("rubicundus", "six", asf)
	rt.InsertWithAddSuggestionFunction("rube", 7, asf)
	rt.InsertWithAddSuggestionFunction("ruber", "2^3", asf)
	rt.InsertWithAddSuggestionFunction("rubens", 9, asf)
	rt.InsertWithAddSuggestionFunction("toasting", 10, asf)
	rt.InsertWithAddSuggestionFunction("toast", "eleven", asf)
	rt.InsertWithAddSuggestionFunction("toaster", "6+6", asf)

	fmt.Println("Parent-Child:")
	fmt.Println(rt.StringParentChild())
	fmt.Println("--------------------------------")
	fmt.Println("Suggestions:")
	fmt.Println(rt.StringSuggestions())
	fmt.Println("--------------------------------")
	fmt.Println("Values:")
	fmt.Println(rt.StringValues())
	fmt.Println("--------------------------------")
	fmt.Printf("nodes with value %v\n", rt.NodeWithValueCount())
	fmt.Printf("nodes with value string %v\n", rt.NodeWithValueCountByCounter(
		func(v interface{}) int {
			switch v.(type) {
			case string:
				return 1
			default:
				return 0
			}
		},
	))

	// VALUES
	fmt.Println("--------------------------------")
	fmt.Println("Values:")
	fmt.Printf("rubens: %v\n", rt.Find("rubens"))
	fmt.Printf("toast: %v\n", rt.Find("toast"))

	fmt.Println("--------------------------------")
	fmt.Println("AutoCompleteBroadTraversal:")
	fmt.Printf("rub: %v\n", rt.AutoCompleteBroadTraversal("rub", 6))
	fmt.Println("AutoCompleteDepthTraversal:")
	fmt.Printf("rub: %v\n", rt.AutoCompleteDepthTraversal("rub", 6))
	fmt.Println("ClosestSuggestions:")
	fmt.Printf("rub: %v\n", rt.ClosestSuggestions("rub"))
}

func example3() {
	rt := goradix.NewRadixTree()

	afs := func(
		key string,
		currentSuggestions []*goradix.RadixTree,
		condidate *goradix.RadixTree,
	) []*goradix.RadixTree {
		v := condidate.Value().(int)

		if v >= 10 {
			return append(currentSuggestions, condidate)
		}

		return currentSuggestions
	}

	rt.InsertWithAddSuggestionFunction("rube", 100, afs)
	rt.InsertWithAddSuggestionFunction("ruber", 200, afs)
	rt.InsertWithAddSuggestionFunction("rubens", 3, afs)
	rt.InsertWithAddSuggestionFunction("rubi", 4, afs)
	rt.InsertWithAddSuggestionFunction("rubicundus", 500, afs)
	rt.InsertWithAddSuggestionFunction("rubicon", 60, afs)

	fmt.Println("Parent-Child:")
	fmt.Println(rt.StringParentChild())
	fmt.Println("--------------------------------")
	fmt.Println("Suggestions:")
	fmt.Println(rt.StringSuggestions())
	fmt.Println("--------------------------------")
	fmt.Println("Values:")
	fmt.Println(rt.StringValues())

	// VALUES
	fmt.Println("--------------------------------")
	fmt.Println("Values:")
	fmt.Printf("rubens: %v\n", rt.Find("rubens"))
	fmt.Printf("toast: %v\n", rt.Find("toast"))

	fmt.Println("--------------------------------")
	fmt.Println("AutoCompleteBroadTraversal:")
	fmt.Printf("rub: %v\n", rt.AutoCompleteBroadTraversal("rub", 6))
	fmt.Printf("rub: %v\n", rt.AutoCompleteDepthTraversal("rub", 6))
	fmt.Println("ClosestSuggestions:")
	fmt.Printf("rub: %v\n", rt.ClosestSuggestions("rub"))
}

func main() {
	fmt.Println("Ganeral Radix Tree")
	fmt.Println("--------------------------------")
	example1()
	fmt.Println("")
	fmt.Println("Custom suggestions 1")
	fmt.Println("--------------------------------")
	example2()
	fmt.Println("Custom suggestions 2")
	fmt.Println("--------------------------------")
	example3()
}
