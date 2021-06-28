package radixSorter

import (
	"encoding/json"
	"io/ioutil"
	"reflect"
	"testing"
)

func Test_radixSorter_Add(t *testing.T) {
	type args struct {
		strs []string
	}
	tests := []struct {
		name             string
		args             args
		sortingCharIndex int
		wantBucketArrs   map[byte][]string // We will only test the actual array inside the bucket, we don't want mutex pointers to interfere with testing
	}{
		{
			name: "basic test",
			args: args{
				strs: []string{"aaa", "aab", "abb", "bbc", "bab"},
			},
			sortingCharIndex: 0,
			wantBucketArrs: map[byte][]string{
				'a': {"aaa", "aab", "abb"},
				'b': {"bbc", "bab"},
			},
		},
		{
			name: "different lengths",
			args: args{
				strs: []string{"aaa", "ab", "bbb"},
			},
			sortingCharIndex: 2,
			wantBucketArrs: map[byte][]string{
				'\x00': {"ab"},
				'a':    {"aaa"},
				'b':    {"bbb"},
			},
		},
		{
			name: "digits",
			args: args{
				strs: []string{"aaa", "111", "bbb"},
			},
			sortingCharIndex: 0,
			wantBucketArrs: map[byte][]string{
				'1': {"111"},
				'a': {"aaa"},
				'b': {"bbb"},
			},
		},
		{
			name: "casing",
			args: args{
				strs: []string{"AaA", "111", "bBb"},
			},
			sortingCharIndex: 0,
			wantBucketArrs: map[byte][]string{
				'1': {"111"},
				'a': {"aaa"},
				'b': {"bbb"},
			},
		},
		{
			name: "duplicates",
			args: args{
				strs: []string{"aaa", "aaa", "abb", "bbc", "bab"},
			},
			sortingCharIndex: 0,
			wantBucketArrs: map[byte][]string{
				'a': {"aaa", "abb"},
				'b': {"bbc", "bab"},
			},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			s := NewRadixSorter()
			s.cursorIndex = tt.sortingCharIndex
			s.Add(tt.args.strs...)
			bucketArrs := make(map[byte][]string)
			for char, bucket := range s.buckets {
				bucketArrs[char] = bucket.arr
			}
			if !reflect.DeepEqual(bucketArrs, tt.wantBucketArrs) {
				t.Errorf("Add() = %#v, want %#v", bucketArrs, tt.wantBucketArrs)
			}
		})
	}
}

func Test_radixSorter_Sort(t *testing.T) {
	tests := []struct {
		name           string
		strsToSort     []string
		wantSortedStrs []string
	}{
		{
			name:           "basic test",
			strsToSort:     []string{"aaa", "aab", "abb", "bbc", "bab"},
			wantSortedStrs: []string{"aaa", "aab", "abb", "bab", "bbc"},
		},
		{
			name:           "different lengths",
			strsToSort:     []string{"aaa", "aa", "bbb"},
			wantSortedStrs: []string{"aa", "aaa", "bbb"},
		},
		{
			name:           "digits",
			strsToSort:     []string{"aaa", "111", "bbb"},
			wantSortedStrs: []string{"111", "aaa", "bbb"},
		},
		{
			name:           "casing",
			strsToSort:     []string{"AaA", "bBb"},
			wantSortedStrs: []string{"aaa", "bbb"},
		},
		{
			name:           "duplicates",
			strsToSort:     []string{"aaa", "aaa", "abb", "bbc", "bab"},
			wantSortedStrs: []string{"aaa", "abb", "bab", "bbc"},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			s := NewRadixSorter()
			s.Add(tt.strsToSort...)
			if gotSortedStrs := s.Sort(); !reflect.DeepEqual(gotSortedStrs, tt.wantSortedStrs) {
				t.Errorf("Sort() = %#v, want %#v", gotSortedStrs, tt.wantSortedStrs)
			}
		})
	}
}

func Benchmark_radixSorter_Sort(b *testing.B) {
	benchmarks := []struct {
		name               string
		strsToSortFilePath string
	}{
		{
			name:               "2-product-lowercase",
			strsToSortFilePath: "../../test/data/sorter/2-product-lowercase.json",
		},
		{
			name:               "3-product-lowercase",
			strsToSortFilePath: "../../test/data/sorter/3-product-lowercase.json",
		},
		{
			name:               "4-product-lowercase-duplicates",
			strsToSortFilePath: "../../test/data/sorter/4-product-lowercase-duplicates.json",
		},
	}
	for _, bm := range benchmarks {
		b.Run(bm.name, func(b *testing.B) {
			file, err := ioutil.ReadFile(bm.strsToSortFilePath)
			if err != nil {
				panic(err)
			}
			var strsToSort []string
			err = json.Unmarshal(file, &strsToSort)
			if err != nil {
				panic(err)
			}
			s := NewRadixSorter()
			s.Add(strsToSort...)
			b.ResetTimer()
			for i := 0; i < b.N; i++ {
				s.Sort()
			}
		})
	}
}
