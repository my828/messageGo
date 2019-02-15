package indexes

import (
	"reflect"
	"sort"
	"testing"
)

func TestTrieAddAndFind(t *testing.T) {
	cases := []struct {
		name           string
		input          []string
		find           string
		max            int
		expectedOutput []int64
		len            int
	}{
		{
			name:           "find on empty trie",
			input:          []string{},
			find:           "fo",
			max:            1,
			expectedOutput: nil,
			len:            0,
		},
		{
			name:           "simple add and find",
			input:          []string{"fo", "foo", "foo"},
			find:           "fo",
			max:            3,
			expectedOutput: []int64{0, 1, 2},
			len:            3,
		},
		{
			name:           "check small max",
			input:          []string{},
			find:           "fo",
			max:            0,
			expectedOutput: nil,
			len:            3,
		},
		{
			name:           "check unfound prefix",
			input:          []string{"fo", "foo", "foo", "food"},
			find:           "food",
			max:            1,
			expectedOutput: []int64{3},
			len:            4,
		},
		{
			name:           "add to existing node",
			input:          []string{"f"},
			find:           "f",
			max:            1,
			expectedOutput: []int64{0},
			len:            4,
		},
		{
			name:           "adding to existing node",
			input:          []string{},
			find:           "f",
			max:            3,
			expectedOutput: []int64{0, 0, 1},
			len:            4,
		},
		{
			name:           "check different branch",
			input:          []string{"g", "gi", "gi", "git", "git", "gitt"},
			find:           "git",
			max:            3,
			expectedOutput: []int64{3, 4, 5},
			len:            7,
		},
		{
			name:           "check add unicode char",
			input:          []string{"g你好"},
			find:           "g你",
			max:            2,
			expectedOutput: []int64{0},
			len:            7,
		},
	}
	trie := NewTrie()
	for _, c := range cases {
		for i := 0; i < len(c.input); i++ {
			trie.Add(c.input[i], int64(i))
		}
		// if trie.Len() != c.len {
		// 	t.Errorf("case %s: expected length of %d but got %d\n", c.name, c.len, trie.Len())
		// }
		result := trie.Find(c.find, c.max)
		sort.Slice(result, func(i, j int) bool {
			return result[i] < result[j]
		})
		if !reflect.DeepEqual(result, c.expectedOutput) {
			t.Errorf("case %s: expected %d but got %d\n", c.name, c.expectedOutput, result)
		}
	}
}

func TestTrieAddAndDelete(t *testing.T) {
	cases := []struct {
		name           string
		input          []string
		key            []string
		val            []int64
		find           string
		max            int
		expectedOutput []int64
		len            int
	}{
		{
			name:           "simple add and delete",
			input:          []string{"fo", "foo", "foo"},
			key:            []string{"fo", "foo"},
			val:            []int64{0, 2},
			find:           "fo",
			max:            2,
			expectedOutput: []int64{1},
			len:            3,
		},
		{
			name:           "delete leaf node",
			input:          []string{"fo", "foo", "foo"},
			key:            []string{"foo", "foo"},
			val:            []int64{1, 2},
			find:           "fo",
			max:            2,
			expectedOutput: []int64{0},
			len:            2,
		},
		{
			name:           "delete targeted key/value",
			input:          []string{"foo", "foo"},
			key:            []string{"foo", "foo"},
			val:            []int64{0, 1},
			find:           "fo",
			max:            2,
			expectedOutput: []int64{0},
			len:            2,
		},
		{
			name:           "delete one val",
			input:          []string{"foo", "foo"},
			key:            []string{"foo"},
			val:            []int64{0},
			find:           "fo",
			max:            2,
			expectedOutput: []int64{0, 1},
			len:            3,
		},
		{
			name:           "delete unvalid key",
			input:          []string{"foo"},
			key:            []string{"fooo"},
			val:            []int64{0},
			find:           "fo",
			max:            3,
			expectedOutput: []int64{0, 0, 1},
			len:            3,
		},
		{
			name:           "delete from big tree",
			input:          []string{"foo", "good", "fun", "fooo"},
			key:            []string{"foo", "foo", "fooo"},
			val:            []int64{0, 1, 3},
			find:           "g",
			max:            1,
			expectedOutput: []int64{1},
			len:            8,
		},
	}
	trie := NewTrie()
	for _, c := range cases {
		for i := 0; i < len(c.input); i++ {
			trie.Add(c.input[i], int64(i))
		}

		for i := 0; i < len(c.key); i++ {
			trie.Remove(c.key[i], c.val[i])
		}

		if trie.Len() != c.len {
			t.Errorf("case %s: expected length of %d but got %d\n", c.name, c.len, trie.Len())
		}
		result := trie.Find(c.find, c.max)
		sort.Slice(result, func(i, j int) bool {
			return result[i] < result[j]
		})
		if !reflect.DeepEqual(result, c.expectedOutput) {
			t.Errorf("case %s: expected %d but got %d\n", c.name, c.expectedOutput, result)
		}
	}

}
