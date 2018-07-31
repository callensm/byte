package utils

import (
	"bytes"
	"encoding/base64"
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

// NewTreeFromB64 decodes the argued base64 encoded string and
// unmarshals the data into a Tree struct and returns a pointer
// to the new unmarshalled data structure
func NewTreeFromB64(encoding []byte) *Tree {
	data := make([]byte, base64.StdEncoding.EncodedLen(len(encoding)))
	_, err := base64.StdEncoding.Decode(data, encoding)
	Catch(err)
	data = bytes.Trim(data, "\x00")

	tree := new(Tree)
	err = json.Unmarshal(data, tree)
	Catch(err)

	return tree
}

// ToB64 encodes the Stringer implementation of the
// Tree struct into Base64 encoding and returns the hash
func (tree *Tree) ToB64() string {
	return base64.StdEncoding.EncodeToString([]byte(tree.String()))
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
