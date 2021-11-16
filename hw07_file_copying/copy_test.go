package main

import (
	"errors"
	"os"
	"testing"

	"github.com/stretchr/testify/suite"
)

type expectedResult struct {
	err          error
	destSize     int64
	checkContent bool
	content      string
}

type copyTestSuite struct {
	suite.Suite
	params   CopyParams
	expected expectedResult
}

func (suite *copyTestSuite) TearDownTest() {
	os.Remove(suite.params.to)
}

func (suite *copyTestSuite) TestSimpleCase() {
	suite.params = CopyParams{from: "./testdata/input.txt", to: "out.txt", offset: 10, limit: 10000}
	suite.expected = expectedResult{err: nil, destSize: 6607}

	suite.RunTest()
}

func (suite *copyTestSuite) TestNoPermissionsCase() {
	suite.params = CopyParams{from: "./testdata/input.txt", to: "/etc/hosts", offset: 10, limit: 20}
	suite.expected = expectedResult{err: os.ErrPermission}

	suite.RunTest()
}

func (suite *copyTestSuite) TestCheckTextInResultCase() {
	suite.params = CopyParams{from: "./testdata/input.txt", to: "out.txt", offset: 10, limit: 12}
	suite.expected = expectedResult{err: nil, destSize: 12, checkContent: true, content: "ts\nPackages\n"}

	suite.RunTest()
}

func (suite *copyTestSuite) TestNegativeOffsetCase() {
	suite.params = CopyParams{from: "./testdata/input.txt", to: "out.txt", offset: -10, limit: 2}
	suite.expected = expectedResult{err: nil, destSize: 2, checkContent: true, content: "Go"}

	suite.RunTest()
}

func (suite *copyTestSuite) TestNegativeOffsetAndLargeLimitCase() {
	suite.params = CopyParams{from: "./testdata/input.txt", to: "out.txt", offset: -10, limit: 10000}
	suite.expected = expectedResult{err: nil, destSize: 6617}

	suite.RunTest()
}

func (suite *copyTestSuite) RunTest() {
	finish := make(chan error)
	progress := make(chan int64)
	var progressCounter int64

	go copy(&suite.params, progress, finish)

	var status error
	for {
		select {
		case status = <-finish:
		case delta := <-progress:
			progressCounter += delta
		}
		break
	}

	if suite.expected.err != nil {
		suite.True(errors.Is(status, suite.expected.err), "actual error %q", status)
	} else {
		suite.Equal(suite.expected.destSize, progressCounter)
		stat, _ := os.Stat(suite.params.to)
		suite.Equal(stat.Size(), suite.expected.destSize)

		if suite.expected.checkContent {
			buf := make([]byte, suite.params.limit)
			f, _ := os.Open(suite.params.to)
			f.Read(buf)
			suite.Equal(suite.expected.content, string(buf))
		}
	}
}

func TestCopy(t *testing.T) {
	suite.Run(t, new(copyTestSuite))
}
