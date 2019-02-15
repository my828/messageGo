package indexes

import (
	"sort"
	"strings"
	"sync"
)

//TODO: implement a trie data structure that stores
//keys of type string and values of type int64

//PRO TIP: if you are having troubles and want to see
//what your trie structure looks like at various points,
//either use the debugger, or try this package:
//https://github.com/davecgh/go-spew

//Trie implements a trie data structure mapping strings to int64s
//that is safe for concurrent use.
type Node struct {
	Val      int64set
	Children map[rune]*Node
	Key      rune
	End      bool
	mx       sync.RWMutex
}

type Trie struct {
	Root *Node
	Size int
}

//NewTrie constructs a new Trie.
func NewTrie() *Trie {
	return &Trie{
		Root: &Node{
			Children: make(map[rune]*Node),
		},
		Size: 0,
	}
}

//Len returns the number of entries in the trie.
func (t *Trie) Len() int {
	return t.Size
}

//Add adds a key and value to the trie.
func (t *Trie) Add(key string, value int64) {
	node := t.Root
	node.mx.Lock()
	key = strings.ToLower(key)
	runes := []rune(key)
	for i := range runes {
		char := runes[i]
		if _, ok := node.Children[char]; !ok {
			newChild := &Node{
				Val:      make(int64set),
				Children: make(map[rune]*Node),
				Key:      char,
			}
			node.Children[char] = newChild
			t.Size++
		}
		node = node.Children[char]
	}
	node.Val.add(value)
	node.End = true
	node.mx.Unlock()
}

//Find finds `max` values matching `prefix`. If the trie
//is entirely empty, or the prefix is empty, or max == 0,
//or the prefix is not found, this returns a nil slice.
func (t *Trie) Find(prefix string, max int) []int64 {
	prefix = strings.ToLower(prefix)
	runes := []rune(prefix)
	node := t.Root
	node.mx.Lock()
	if len(prefix) == 0 || max == 0 {
		return nil
	}
	// go to branch of trie holding keys that start with prefix.
	for i := range runes {
		char := runes[i]
		if _, ok := node.Children[char]; !ok {
			return nil
		}
		node = node.Children[char]
	}
	result := []int64{}
	defer node.mx.Unlock()
	return dfs(node, result, max)
}

func dfs(node *Node, result []int64, max int) []int64 {
	values := node.Val.all()
	count := max - len(result)
	if count-len(values) > 0 {
		result = append(result, values...)
		// order children keys
		keys := []rune{}
		for k, _ := range node.Children {
			keys = append(keys, k)
		}
		sort.Slice(keys, func(i, j int) bool {
			return keys[i] < keys[j]
		})
		length := len(keys)
		// depth first search
		for i := 0; i < length; i++ {
			result = dfs(node.Children[keys[i]], result, max)
			if len(result) == max {
				i = length
			}
		}
	} else { // base case: length of values is same or more than max
		n := max - len(result)
		for i := 0; i < n; i++ {
			result = append(result, values[i])
		}
	}
	return result
}

// Remove removes a key/value pair from the trie
// and trims branches with no values.
func (t *Trie) Remove(key string, value int64) {
	runes := []rune(key)
	t.Root.mx.Lock()
	defer t.Root.mx.Unlock()
	t.RemoveHelper(t.Root, runes, value, 0, true)
}

func (t *Trie) RemoveHelper(node *Node, keys []rune, value int64, index int, found bool) (*Node, bool) {
	if index < len(keys) {
		char := keys[index]
		if _, ok := node.Children[char]; ok {
			node.Children[char], found = t.RemoveHelper(node.Children[char], keys, value, index+1, found)
			if len(node.Children[char].Children) == 0 &&
				len(node.Children[char].Val) == 0 {
				delete(node.Children, char)
				t.Size--
			}
		} else {
			found = false
		}
	}
	if found {
		node.Val.remove(value)
	}
	return node, found
}

func (t *Trie) SplitNameAddToTrie(name string, id int64) {
	names := strings.Split(name, " ")
	for _, s := range names {
		t.Add(s, id)
	}
}

//int64set is a set of int64 values
type int64set map[int64]struct{}

//add adds a value to the set and returns
//true if the value didn't already exist in the set.
func (s int64set) add(value int64) bool {
	if exist := s.has(value); !exist {
		s[value] = struct{}{}
		return !exist
	}
	return false
}

//remove removes a value from the set and returns
//true if that value was in the set, false otherwise.
func (s int64set) remove(value int64) bool {
	if exist := s.has(value); exist {
		delete(s, value)
		return exist
	}
	return false
}

//has returns true if value is in the set,
//or false if it is not in the set.
func (s int64set) has(value int64) bool {
	_, exist := s[value]
	return exist
}

//all returns all values in the set as a slice.
//The returned slice will always be non-nil, but
//the order will be random. Use sort.Slice to
//sort the slice if necessary.
func (s int64set) all() []int64 {
	keys := make([]int64, len(s))
	i := 0
	for k := range s {
		keys[i] = k
		i++
	}
	return keys
}
