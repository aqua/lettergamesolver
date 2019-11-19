package main

import (
	"bufio"
	"flag"
	"fmt"
	"os"
	"sort"
	"strings"
	"sync"

	combos "github.com/mxschmitt/golang-combinations"
)

var dictOnly = flag.Bool("dictonly", true, "Only return permutations from the dictionary")
var dict = flag.String("dict", "/usr/share/dict/words", "Path to dictionary file")
var minLength = flag.Int("min-length", 0, "Return only words greater than N characters")
var maxLength = flag.Int("max-length", 20, "Return only words greater than N characters")
var subset = flag.Bool("subset", false, "Generate combos of any nonzero subset length")
var sortLength = flag.Bool("sort-length", false, "Return in length-order, not lexical order")

var dictionary = map[string]struct{}{}
var dictLoaded sync.Once

func isDictWord(w string) bool {
	dictLoaded.Do(func() {
		f, err := os.Open(*dict)
		if err != nil {
			panic(fmt.Sprintf("%s: %v", *dict, err))
		}
		defer f.Close()
		scanner := bufio.NewScanner(f)
		for scanner.Scan() {
			dictionary[scanner.Text()] = struct{}{}
		}
	})
	_, exists := dictionary[w]
	return exists
}

func permutations(arr []rune) [][]rune {
	var helper func([]rune, int)
	res := [][]rune{}

	helper = func(arr []rune, n int) {
		if n == 1 {
			tmp := make([]rune, len(arr))
			copy(tmp, arr)
			res = append(res, tmp)
		} else {
			for i := 0; i < n; i++ {
				helper(arr, n-1)
				if n%2 == 1 {
					tmp := arr[i]
					arr[i] = arr[n-1]
					arr[n-1] = tmp
				} else {
					tmp := arr[0]
					arr[0] = arr[n-1]
					arr[n-1] = tmp
				}
			}
		}
	}
	helper(arr, len(arr))
	return res
}

func main() {
	flag.Parse()
	if flag.NArg() < 1 || flag.Arg(0) == "" {
		fmt.Println("usage: permute letters")
		os.Exit(2)
	}
	input := flag.Arg(0)
	bases := []string{input}
	if *subset {
		bases = []string{}
		for _, c := range combos.All(strings.Split(input, "")) {
			bases = append(bases, strings.Join(c, ""))
		}
	}
	used := map[string]bool{}
	found := []string{}
	for _, w := range bases {
		letters := make([]rune, len(w))
		for i, c := range w {
			letters[i] = c
		}
		for _, word := range permutations(letters) {
			word := string(word)
			if *dictOnly && !isDictWord(word) {
				continue
			}
			if len(word) < *minLength || len(word) > *maxLength {
				continue
			}
			if used[word] {
				continue
			}
			used[word] = true
			found = append(found, word)
		}
	}
	if *sortLength {
		sort.Slice(found, func(i, j int) bool {
			l1, l2 := len(found[i]), len(found[j])
			if l1 == l2 {
				return found[i] < found[j]
			}
			return l1 < l2
		})
	} else {
		sort.Strings(found)
	}
	for _, word := range found {
		fmt.Println(word)
	}
}
