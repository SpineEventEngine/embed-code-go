// Copyright 2024, TeamDev. All rights reserved.
//
// Redistribution and use in source and/or binary forms, with or without
// modification, must retain the above copyright notice and the following
// disclaimer.
//
// THIS SOFTWARE IS PROVIDED BY THE COPYRIGHT HOLDERS AND CONTRIBUTORS
// "AS IS" AND ANY EXPRESS OR IMPLIED WARRANTIES, INCLUDING, BUT NOT
// LIMITED TO, THE IMPLIED WARRANTIES OF MERCHANTABILITY AND FITNESS FOR
// A PARTICULAR PURPOSE ARE DISCLAIMED. IN NO EVENT SHALL THE COPYRIGHT
// OWNER OR CONTRIBUTORS BE LIABLE FOR ANY DIRECT, INDIRECT, INCIDENTAL,
// SPECIAL, EXEMPLARY, OR CONSEQUENTIAL DAMAGES (INCLUDING, BUT NOT
// LIMITED TO, PROCUREMENT OF SUBSTITUTE GOODS OR SERVICES; LOSS OF USE,
// DATA, OR PROFITS; OR BUSINESS INTERRUPTION) HOWEVER CAUSED AND ON ANY
// THEORY OF LIABILITY, WHETHER IN CONTRACT, STRICT LIABILITY, OR TORT
// (INCLUDING NEGLIGENCE OR OTHERWISE) ARISING IN ANY WAY OUT OF THE USE
// OF THIS SOFTWARE, EVEN IF ADVISED OF THE POSSIBILITY OF SUCH DAMAGE.

package indenttest

import (
	"embed-code/embed-code-go/indent"
	"reflect"
	"testing"
)

func TestNoSpaces(t *testing.T) {
	testLines := []string{"", "foo", "bar", "", "baz", ""}
	got := indent.MaxCommonIndentation(testLines)
	if got != 0 {
		t.Errorf("The indentation is %d; want 0", got)
	}
}

func TestNoLines(t *testing.T) {
	testLines := []string{}
	got := indent.MaxCommonIndentation(testLines)
	if got != 0 {
		t.Errorf("The indentation is %d; want 0", got)
	}
}

func TestOnlyEmptyLines(t *testing.T) {
	testLines := []string{"", "    ", ""}
	got := indent.MaxCommonIndentation(testLines)
	if got != 0 {
		t.Errorf("The indentation is %d; want 0", got)
	}
}

func TestTwoIndents(t *testing.T) {
	testLines := []string{"", "  foo", "    bar", "", "", "  baz"}
	got := indent.MaxCommonIndentation(testLines)
	if got != 2 {
		t.Errorf("The indentation is %d; want 0", got)
	}
}

func TestCutIndent(t *testing.T) {
	testLines := []string{"", "  foo", "    bar", "", "", "  baz"}
	testLinesChanged := indent.CutIndent(testLines, 2)

	if reflect.DeepEqual(testLines, testLinesChanged) {
		t.Errorf("The given lines weren't changed")
	}
}
