package model_consistenthash

import (
	"cache/internal/model/consistenthash"
	"strconv"
	"sync"
	"testing"
)

func TestGo(t *testing.T) {
	hash := consistenthash.New(3, func(key []byte) uint32 {
		i, _ := strconv.Atoi(string(key))
		return uint32(i)
	})
	// hash.Add("1", "4", "2")
	// testCases := map[string]string{
	// 	"1":  "2",
	// 	"2":  "2",
	// 	"3":  "4",
	// 	"4":  "2",
	// 	"5":  "2",
	// 	"6":  "2",
	// 	"7":  "4",
	// 	"8":  "2",
	// 	"9":  "2",
	// 	"10": "2",
	// 	"11": "4",
	// 	"12": "2",
	// 	"13": "2",
	// 	"14": "2",
	// 	"15": "4",
	// 	"16": "2",
	// }
	wg := sync.WaitGroup{}
	for k := 0; k < 10000; k++ {
		wg.Add(1)
		go func(i int) {
			hash.Add(strconv.Itoa(i))
			// _ = hash.Get(strconv.Itoa(i))
			// fmt.Printf("Asking for %v, should have yielded %s \n", i, value)
			wg.Done()
		}(k)
	}
	wg.Wait()
}
