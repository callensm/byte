package utils

import (
	"bytes"
	"encoding/json"
	"io/ioutil"
	"path/filepath"
)

// Tree is a structure to represent the file
// structure of the directory being transmitted
// over the socket connection to be downloaded
type Tree struct {
	Name     string   `json:"name"`
	SubTrees []*Tree  `json:"subTrees,omitempty"`
	Leaves   []string `json:"leaves"`
}

// NewTree populates a new Tree struct instance
// after walking to file system gathering necessary
// information and returns a pointer to the new struct
func NewTree(path string) *Tree {
	var err error
	base := filepath.Base(path)
	t := &Tree{Name: base}

	// Read all entities within the current argued path
	fileList, err := ioutil.ReadDir(path)
	Catch(err)
	for _, f := range fileList {
		// If it is a file, append it to the leaves array
		if !f.IsDir() {
			t.Leaves = append(t.Leaves, f.Name())
			continue
		}

		// If it was a directory, recursively build a subtree from that path
		sub := NewTree(filepath.Join(path, f.Name()))
		t.SubTrees = append(t.SubTrees, sub)
	}
	return t
}

// NewTreeFromJSON returns a pointer
// to the new unmarshalled data structure
func NewTreeFromJSON(encoding []byte) *Tree {
	tree := new(Tree)
	err := json.Unmarshal(encoding, tree)
	Catch(err)
	return tree
}

// String is the Stringer interface implementation
// for the Tree structure and marshals to struct into
// JSON format for the string version
func (tree *Tree) String() string {
	b, err := json.Marshal(tree)
	Catch(err)
	return string(b)
}

// Display returns the prettified and indented version
// of the JSON representation for the Tree structure
func (tree *Tree) Display() string {
	var out []byte
	dst := bytes.NewBuffer(out)
	err := json.Indent(dst, []byte(tree.String()), "", "  ")
	Catch(err)
	return dst.String()
}

// CountLeaves walks the Tree structure and returns
// the number of leaves (files) in the entire tree
func (tree *Tree) CountLeaves() int {
	count := len(tree.Leaves)
	for _, s := range tree.SubTrees {
		count += s.CountLeaves()
	}
	return count
}

// CountSubTrees walks the Tree structure and
// returns the number of subtrees (directories) in
// the entire tree
func (tree *Tree) CountSubTrees() int {
	count := len(tree.SubTrees)
	for _, s := range tree.SubTrees {
		count += s.CountSubTrees()
	}
	return count
}
