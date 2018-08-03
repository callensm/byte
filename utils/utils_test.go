package utils

import (
	"bufio"
	"io"
	"os"
	"testing"
)

const spinnerTag = "test_spinner"

func TestCatch(t *testing.T) {
	output := bufio.NewReader(os.Stdout)
	Catch(io.EOF)
	if x := output.Buffered(); x != 0 {
		t.Errorf("Stdout has %d buffered bytes when no error log should be reported", x)
	}
}

func TestCreateSpinner(t *testing.T) {
	CreateSpinner(1, "", "", spinnerTag)
	if spinners[spinnerTag] == nil {
		t.Errorf("Spinner with tag %s was not found after creation", spinnerTag)
	}
}

func TestRemoveSpinner(t *testing.T) {
	RemoveSpinner(spinnerTag, "", true)
	if spinners[spinnerTag] != nil {
		t.Errorf("Spinner with tag %s still exists after attempted removal", spinnerTag)
	}

	TestCreateSpinner(t)
	RemoveSpinner(spinnerTag, "", false)
}
