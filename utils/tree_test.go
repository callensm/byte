package utils

import (
	"io/ioutil"
	"os"
	"path/filepath"
	"testing"
)

func TestNewTree(t *testing.T) {
	treePath, err := filepath.Abs(".")
	if err != nil {
		t.Fatal(err)
	}

	tree, _ := NewTree(treePath)
	if tree == nil {
		t.Errorf("NewTree returned a null pointer instead of *Tree instance")
	}

	if base := filepath.Base(treePath); tree.Name != base {
		t.Errorf("Name of the Tree instance is %s, expected %s", tree.Name, base)
	}

	if x := len(tree.SubTrees); x > 0 {
		t.Errorf("Initialized with %d sub-trees in a flat directory", x)
	}

	fileList, err := ioutil.ReadDir(treePath)
	if err != nil {
		t.Fatal(err)
	}

	if x, y := len(tree.Leaves), len(fileList); x != y {
		t.Errorf("Initialized with %d leaves, but expected %d", y, x)
	}

	newTreePath, err := filepath.Abs("../.github")
	if err != nil {
		t.Fatal(err)
	}

	treeWithSubs, _ := NewTree(newTreePath)
	if len(treeWithSubs.SubTrees) == 0 {
		t.Error("Initialized with 0 sub-trees, expecting 1 or more")
	}
}

func TestTreeJSON(t *testing.T) {
	data := `{"name":"test","subTrees":[{"name":"subtest","leaves":[]}],"leaves":["tree.go"]}`
	tree := NewTreeFromJSON([]byte(data))
	if tree == nil {
		t.Error("NewTreeFromJSON returned a null pointer instead of a *Tree instance")
	}

	if tree.String() == "" || tree.Display() == "" {
		t.Error("Failed to marshal *Tree into JSON string")
	}
}

func TestIgnoredFiles(t *testing.T) {
	path, err := filepath.Abs("../")
	if err != nil {
		t.Fatal(err)
	}

	_, ignored := NewTree(path)
	if len(ignored) == 0 {
		t.Errorf("No ignored files were found from the tree, expecting >1")
	}
}

func TestCountLeaves(t *testing.T) {
	path, err := filepath.Abs("../.github")
	if err != nil {
		t.Fatal(err)
	}

	var total int
	err = filepath.Walk(path, func(_ string, info os.FileInfo, err error) error {
		if !info.IsDir() {
			total++
		}
		return err
	})
	if err != nil {
		t.Fatal(err)
	}

	tree, _ := NewTree(path)
	if x, y := tree.CountLeaves(), total; x != y {
		t.Errorf("CountLeaves returned %d file count, expected %d", x, y)
	}
}

func TestCountSubTrees(t *testing.T) {
	path, err := filepath.Abs("../.github")
	if err != nil {
		t.Fatal(err)
	}

	var total int
	err = filepath.Walk(path, func(_ string, info os.FileInfo, err error) error {
		if info.IsDir() {
			total++
		}
		return err
	})
	if err != nil {
		t.Fatal(err)
	}

	// Subtract one from total because filepath.Walk counts the root directory
	total--
	tree, _ := NewTree(path)
	if x, y := tree.CountSubTrees(), total; x != y {
		t.Errorf("CountSubTrees returned %d directory count, expected %d", x, y)
	}
}
