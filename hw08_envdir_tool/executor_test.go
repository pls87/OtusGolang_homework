package main

import (
	"os"
	"testing"

	"github.com/stretchr/testify/suite"
)

type executorExpected struct {
	returnCode int
	stdOutput  string
}

type executorTestCase struct {
	name     string
	command  []string
	env      Environment
	expected executorExpected
}

type executorTestSuite struct {
	suite.Suite
	stdOut   *os.File
	cmdOut   *os.File
	testCase executorTestCase
}

func (suite *executorTestSuite) SetupTest() {
	suite.cmdOut, _ = os.Create("./testdata/temp_output.txt")
	suite.stdOut = os.Stdout
	os.Stdout = suite.cmdOut
}

func (suite *executorTestSuite) TearDownTest() {
	os.Remove("./testdata/temp_output.txt")
	os.Stdout = suite.stdOut
}

func (suite *executorTestSuite) TestSimpleLSCommandWithCheckStdOutput() {
	suite.testCase = executorTestCase{
		name:     "ls command",
		command:  []string{"ls", "/bin/bash"},
		expected: executorExpected{returnCode: 0, stdOutput: "/bin/bash\n"},
	}

	suite.RunTest()
}

func (suite *executorTestSuite) RunTest() {
	code := RunCmd(suite.testCase.command, suite.testCase.env)

	suite.cmdOut.Seek(0, 0)
	buf := make([]byte, 1024)
	read, _ := suite.cmdOut.Read(buf)

	suite.Equal(suite.testCase.expected.returnCode, code)
	suite.Equal(suite.testCase.expected.stdOutput, string(buf[0:read]))
}

func TestRunCmd(t *testing.T) {
	suite.Run(t, new(executorTestSuite))
}
