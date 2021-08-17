package main

import (
	"fmt"
	"sort"
	"strconv"
	"testing"

	"github.com/dengjiawen8955/gocache/src/consistenthash"
)

func main() {
	// GuessingGame()
}
func GuessingGame() {
	var s string
	fmt.Printf("Pick an integer from 0 to 100.\n")
	answer := sort.Search(100, func(i int) bool {
		fmt.Printf("Is your number <= %d? ", i)
		fmt.Scanf("%s", &s)
		return s != "" && s[0] == 'y'
	})
	fmt.Printf("Your number is %d.\n", answer)
}
func TestHashing(t *testing.T) {
	hash := consistenthash.NewNodeCircle(3, func(key []byte) uint32 {
		i, _ := strconv.Atoi(string(key))
		return uint32(i)
	})

	// Given the above hash function, this will give replicas with "hashes":
	// 2, 4, 6, 12, 14, 16, 22, 24, 26
	hash.AddNodes("6", "4", "2")

	testCases := map[string]string{
		"2":  "2",
		"11": "2",
		"23": "4",
		"27": "2",
	}

	for k, v := range testCases {
		s, _ := hash.GetNode(k)
		if s != v {
			fmt.Printf("Asking for %s, should have yielded %s", k, v)
		}
	}

	// Adds 8, 18, 28
	hash.AddNodes("8")

	// 27 should now map to 8.
	testCases["27"] = "8"

	for k, v := range testCases {
		s, _ := hash.GetNode(k)
		if s != v {
			fmt.Printf("Asking for %s, should have yielded %s", k, v)
		}
	}

}

//里外面的值是里面的 + 1
func Test_sort(t *testing.T) {
	re := sort.Search(9, func(i int) bool {
		fmt.Printf("i=%#v\n", i)
		return i > 5
	})
	fmt.Printf("re=%#v\n", re)
}
