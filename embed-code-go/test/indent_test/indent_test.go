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

package indent_test

import (
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/suite"

	"embed-code/embed-code-go/indent"
)

type IndentTestSuite struct {
	suite.Suite
}

func (suite *IndentTestSuite) TestNoSpaces() {
	testLines := []string{"", "foo", "bar", "", "baz", ""}

	assert.Equal(suite.T(), 0, indent.MaxCommonIndentation(testLines))
}

func (suite *IndentTestSuite) TestNoLines() {
	testLines := []string{}

	assert.Equal(suite.T(), 0, indent.MaxCommonIndentation(testLines))
}

func (suite *IndentTestSuite) TestOnlyEmptyLines() {
	testLines := []string{"", "    ", ""}

	assert.Equal(suite.T(), 0, indent.MaxCommonIndentation(testLines))
}

func (suite *IndentTestSuite) TestTwoIndents() {
	testLines := []string{"", "  foo", "    bar", "", "", "  baz"}

	assert.Equal(suite.T(), 2, indent.MaxCommonIndentation(testLines))
}

func (suite *IndentTestSuite) TestCutIndent() {
	testLines := []string{"", "  foo", "    bar", "", "", "  baz"}
	testLinesChanged := indent.CutIndent(testLines, 2)

	assert.NotEqual(suite.T(), testLines, testLinesChanged)
}

func TestIndentTestSuite(t *testing.T) {
	suite.Run(t, new(IndentTestSuite))
}
