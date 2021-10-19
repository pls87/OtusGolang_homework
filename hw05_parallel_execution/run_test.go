package hw05parallelexecution

import (
	"errors"
	"testing"
	"time"

	"github.com/stretchr/testify/suite"
	"go.uber.org/goleak"
)

type runStats struct {
	totalRuns       int32
	pureRunningTime time.Duration
	elapsedTime     time.Duration
	errors          int32
}

type testCase struct {
	workers    int
	maxErrors  int
	tasksCount int
	generator  func(ts *parallelExecutionTestSuite) []Task
}

type parallelExecutionTestSuite struct {
	suite.Suite
	rs    *runStats
	tc    *testCase
	tasks []Task
}

func (suite *parallelExecutionTestSuite) SetupTest() {
	suite.rs = &runStats{}
}

func (suite *parallelExecutionTestSuite) BeforeTest(_, testName string) {
	suite.tc = suite.NextCase(testName)
	suite.tasks = suite.tc.generator(suite)
}

func (suite *parallelExecutionTestSuite) TestFirstMErrors() {
	err := Run(suite.tasks, suite.tc.workers, suite.tc.maxErrors)

	suite.Truef(errors.Is(err, ErrErrorsLimitExceeded), "actual err - %v", err)
	suite.LessOrEqual(int(suite.rs.errors), suite.tc.maxErrors+suite.tc.workers)
	suite.LessOrEqual(suite.rs.totalRuns, int32(suite.tc.workers+int(suite.rs.errors)), "extra tasks were started")
}

func (suite *parallelExecutionTestSuite) TestNoErrors() {
	start := time.Now()
	err := Run(suite.tasks, suite.tc.workers, suite.tc.maxErrors)
	suite.rs.elapsedTime = time.Since(start)

	suite.NoError(err)
	suite.Equal(suite.rs.totalRuns, int32(suite.tc.tasksCount), "not all tasks were completed")
	suite.LessOrEqual(int64(suite.rs.elapsedTime), int64(suite.rs.pureRunningTime/2), "tasks were run sequentially?")
}

func (suite *parallelExecutionTestSuite) TestWithRealWorkAndWOErrors() {
	start := time.Now()
	err := Run(suite.tasks, suite.tc.workers, suite.tc.maxErrors)
	suite.rs.elapsedTime = time.Since(start)

	suite.NoError(err)
	suite.Equal(suite.rs.totalRuns, int32(suite.tc.tasksCount), "not all tasks were completed")
	suite.LessOrEqual(int64(suite.rs.elapsedTime), int64(suite.rs.pureRunningTime/2), "tasks were run sequentially?")
}

func (suite *parallelExecutionTestSuite) TestWithRealWorkAndWithSomeErrors() {
	start := time.Now()
	err := Run(suite.tasks, suite.tc.workers, suite.tc.maxErrors)
	suite.rs.elapsedTime = time.Since(start)

	suite.Truef(errors.Is(err, ErrErrorsLimitExceeded), "actual err - %v", err)
	suite.LessOrEqual(int(suite.rs.errors), suite.tc.maxErrors+suite.tc.workers)
	suite.LessOrEqual(int64(suite.rs.elapsedTime), int64(suite.rs.pureRunningTime/2), "tasks were run sequentially?")
}

func TestRun(t *testing.T) {
	defer goleak.VerifyNone(t)

	suite.Run(t, new(parallelExecutionTestSuite))
}
