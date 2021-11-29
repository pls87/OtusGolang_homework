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
	suite.cmdOut, _ = os.Create("./testdata/test_output.txt")
	suite.stdOut = os.Stdout
	os.Stdout = suite.cmdOut
}

func (suite *executorTestSuite) TearDownTest() {
	os.Remove("./testdata/test_output.txt")
	os.Stdout = suite.stdOut
}

func (suite *executorTestSuite) TestSimpleLSCommandWithCheckStdOutput() {
	suite.testCase = executorTestCase{
		name:     "ls command",
		command:  []string{"ls", "/bin/bash"},
		expected: executorExpected{returnCode: 0, stdOutput: "/bin/bash\n"},
	}

	suite.runTest()
}

func (suite *executorTestSuite) TestWithParametersAndCheckStdOutput() {
	suite.testCase = executorTestCase{
		name:     "head command",
		command:  []string{"head", "-n 1", "./testdata/env/BAR"},
		expected: executorExpected{returnCode: 0, stdOutput: "bar\n"},
	}

	suite.runTest()
}

func (suite *executorTestSuite) TestUnknownCommand() {
	suite.testCase = executorTestCase{
		name:     "unknown command",
		command:  []string{"this_command_is_unknown"},
		expected: executorExpected{returnCode: 1},
	}

	suite.runTest()
}

func (suite *executorTestSuite) TestEnvVariable() {
	suite.testCase = executorTestCase{
		name:    "echo $BAR",
		command: []string{"/bin/sh", "-c", "echo $BAR"},
		env: Environment{
			"BAR": EnvValue{
				Value:      "foo",
				NeedRemove: false,
			},
		},
		expected: executorExpected{returnCode: 0, stdOutput: "foo\n"},
	}

	suite.runTest()
}

func (suite *executorTestSuite) runTest() {
	code := RunCmd(suite.testCase.command, suite.testCase.env)

	suite.Equal(suite.testCase.expected.returnCode, code)
	suite.Equal(suite.testCase.expected.stdOutput, suite.readOutput())
}

func (suite *executorTestSuite) readOutput() string {
	suite.cmdOut.Seek(0, 0)
	buf := make([]byte, 1024)
	read, _ := suite.cmdOut.Read(buf)

	return string(buf[0:read])
}

func TestRunCmd(t *testing.T) {
	suite.Run(t, new(executorTestSuite))
}
