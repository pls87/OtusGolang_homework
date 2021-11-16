package main

import (
	"crypto/md5"
	"errors"
	"fmt"
	"os"
	"testing"

	"github.com/stretchr/testify/suite"
)

type expectedResult struct {
	err      error
	destSize int64
	checksum string
}

type copyTestSuite struct {
	suite.Suite
	params   CopyParams
	expected expectedResult
}

func (suite *copyTestSuite) SetupTest() {
	chunkSize = 32
}

func (suite *copyTestSuite) TearDownTest() {
	os.Remove(suite.params.to)
}

func (suite *copyTestSuite) TestSimpleCase() {
	suite.params = CopyParams{from: "./testdata/input.txt", to: "out.txt", offset: 10, limit: 10000}
	suite.expected = expectedResult{err: nil, destSize: 6607, checksum: "0e27306b7a21bfef46eb51718d6c7c2d"}

	suite.RunTest()
}

func (suite *copyTestSuite) TestNoPermissionsCase() {
	suite.params = CopyParams{from: "./testdata/input.txt", to: "/etc/hosts", offset: 10, limit: 20}
	suite.expected = expectedResult{err: os.ErrPermission}

	suite.RunTest()
}

func (suite *copyTestSuite) TestCheckTextInResultCase() {
	suite.params = CopyParams{from: "./testdata/input.txt", to: "out.txt", offset: 10, limit: 12}
	suite.expected = expectedResult{err: nil, destSize: 12, checksum: "0514cfc175893c2ce8e3d6141cc66914"}

	suite.RunTest()
}

func (suite *copyTestSuite) TestNegativeOffsetCase() {
	suite.params = CopyParams{from: "./testdata/input.txt", to: "out.txt", offset: -10, limit: 2}
	suite.expected = expectedResult{err: nil, destSize: 2, checksum: "5f075ae3e1f9d0382bb8c4632991f96f"}

	suite.RunTest()
}

func (suite *copyTestSuite) TestNegativeOffsetAndLargeLimitCase() {
	suite.params = CopyParams{from: "./testdata/input.txt", to: "out.txt", offset: -10, limit: 10000}
	suite.expected = expectedResult{err: nil, destSize: 6617, checksum: "6af260e634c98459961ff62443267b74"}

	suite.RunTest()
}

func (suite *copyTestSuite) RunTest() {
	finish := make(chan error)
	progress := make(chan int64)
	var progressCounter int64

	go cp(suite.params, progress, finish)

	status := func() error {
		for {
			select {
			case err := <-finish:
				return err
			case delta := <-progress:
				progressCounter += delta
			}
		}
	}()

	if suite.expected.err != nil {
		suite.True(errors.Is(status, suite.expected.err), "actual error %q", status)
		return
	}
	
	suite.Equal(suite.expected.destSize, progressCounter)
	stat, _ := os.Stat(suite.params.to)
	suite.Equal(stat.Size(), suite.expected.destSize)

	// check content via md5 checksum
	buf := make([]byte, suite.expected.destSize)
	f, _ := os.Open(suite.params.to)
	f.Read(buf)
	suite.Equal(suite.expected.checksum, fmt.Sprintf("%x", md5.Sum(buf)))
}

func TestCopy(t *testing.T) {
	suite.Run(t, new(copyTestSuite))
}
