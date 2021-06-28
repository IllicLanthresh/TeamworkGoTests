package radixSorter

import (
	"github.com/IllicLanthresh/TeamworkGoTests/pkg/helperTypes"
	"strings"
	"sync"
)

type threadSafeBucket struct {
	arr helperTypes.StringSlice
	mut *sync.Mutex
}

func newThreadSafeBucket() *threadSafeBucket {
	return &threadSafeBucket{
		arr: []string{},
		mut: &sync.Mutex{},
	}
}

// radixSorter implements the Radix algorithm specifically for UTF-8 strings
type radixSorter struct {
	buckets     map[byte]*threadSafeBucket
	cursorIndex int
	lowestChar  byte
	highestChar byte
}

// NewRadixSorter Constructor for radixSorter
func NewRadixSorter() *radixSorter {
	return &radixSorter{
		buckets:     make(map[byte]*threadSafeBucket),
		cursorIndex: 0,
		lowestChar:  0,
		highestChar: 0,
	}
}

// Add adds strings to the sorter bucketing system, after all desired strings are added, Sort can be called
// to retrieve the sorted result.
// All stings will be converted to lowercase and duplicates or empty strings will be ignored.
func (s *radixSorter) Add(strs ...string) {
	for _, str := range strs {
		if len(str) == 0 {
			continue
		}
		str = strings.ToLower(str)

		var char byte
		if s.cursorIndex == len(str) {
			// In case we get a string which had all chars already walked in upper layers, we stash it away in a special
			// bucket indexed with the hash `nullbyte`.
			// This is done so we can still keep them in our buckets without passing them down to lower sorters
			char = '\x00'
		} else if s.cursorIndex > len(str) {
			panic("This should never happen, cursorIndex above length of string to sort")
		} else {
			char = str[s.cursorIndex]
		}

		bucket, bucketFound := s.buckets[char]
		if !bucketFound {
			s.buckets[char] = newThreadSafeBucket()
			bucket = s.buckets[char]
		} else {
			bucket.mut.Lock()
			alreadyIn := bucket.arr.Contains(str)
			bucket.mut.Unlock()
			if alreadyIn {
				continue
			}
		}
		bucket.mut.Lock()
		bucket.arr = append(bucket.arr, str)
		bucket.mut.Unlock()

		// We exclude `nullbyte` from the low-high range so we can treat it separately when sorting
		if (s.lowestChar == 0 || char < s.lowestChar) && char != '\x00' {
			s.lowestChar = char
		}
		if (s.highestChar == 0 || char > s.highestChar) && char != '\x00' {
			s.highestChar = char
		}
	}
}

// Sort outputs a sorted slice of the strings fed to the sorter using Add
func (s *radixSorter) Sort() (sortedStrs []string) {
	var wg sync.WaitGroup

	for char, strs := range s.buckets {
		if len(strs.arr) > 1 {
			wg.Add(1)
			go func(outterSorter *radixSorter, char byte, wg *sync.WaitGroup) {
				innerSorter := NewRadixSorter()
				innerSorter.cursorIndex = outterSorter.cursorIndex + 1
				innerSorter.Add(outterSorter.buckets[char].arr...)

				sorted := innerSorter.Sort()

				outterSorter.buckets[char].mut.Lock()
				outterSorter.buckets[char].arr = sorted
				outterSorter.buckets[char].mut.Unlock()

				wg.Done()
			}(s, char, &wg)
		}
	}
	wg.Wait()

	// `nullbyte` bucket gets emptied first(this is personal preference, shorter strings will go first this way)
	if _, exists := s.buckets['\x00']; exists {
		for _, str := range s.buckets['\x00'].arr {
			sortedStrs = append(sortedStrs, str)
		}
	}

	for char := s.lowestChar; char <= s.highestChar; char++ {
		bucket, found := s.buckets[char]
		if found {
			for _, str := range bucket.arr {
				sortedStrs = append(sortedStrs, str)
			}
		}
	}
	return
}
