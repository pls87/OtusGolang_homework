package hw06pipelineexecution

import (
	"fmt"
	"sync"
	"testing"
	"time"

	"github.com/stretchr/testify/suite"
)

const (
	sleepPerStage = time.Millisecond * 100
	fault         = sleepPerStage / 2
)

type stageBase struct {
	take        func(v interface{}) bool
	transformer func(v interface{}) interface{}
}

type testCase struct {
	data      []interface{}
	generator func(suite *pipelineTestSuite, sb stageBase) Stage
	stg       []stageBase
	stopAfter time.Duration
}

type pipelineTestSuite struct {
	suite.Suite
	sync.WaitGroup
	tc       *testCase
	stages   []Stage
	in, done Bi
}

func (suite *pipelineTestSuite) SetupTest() {
	suite.in = make(Bi)
	suite.done = make(Bi)
}

func (suite *pipelineTestSuite) AfterTest(_, _ string) {
	suite.Wait()
}

func (suite *pipelineTestSuite) BeforeTest(_, testName string) {
	suite.tc = suite.NextCase(testName)

	if suite.tc == nil {
		suite.FailNow(fmt.Sprintf("Test case not found: %s", testName))
	}

	suite.stages = make([]Stage, len(suite.tc.stg))
	for i, v := range suite.tc.stg {
		suite.stages[i] = suite.tc.generator(suite, v)
	}

	suite.Add(1)
	go func() {
		defer suite.Done()
		if suite.tc.stopAfter > 0 {
			<-time.After(suite.tc.stopAfter)
			close(suite.done)
		}
	}()

	suite.Add(1)
	go func() {
		defer suite.Done()
		defer close(suite.in)
		for _, v := range suite.tc.data {
			select {
			case <-suite.done:
				return
			case suite.in <- v:
			}
		}
	}()
}

func (suite *pipelineTestSuite) TestSimpleCase() {
	result := make([]string, 0, 10)
	start := time.Now()
	for s := range ExecutePipeline(suite.in, suite.done, suite.stages...) {
		result = append(result, s.(string))
	}
	elapsed := time.Since(start)
	suite.Equal([]string{"102", "104", "106", "108", "110"}, result)
	suite.Less(
		int64(elapsed),
		int64(sleepPerStage)*int64(len(suite.stages)+len(suite.tc.data)-1)+int64(fault))
}

func (suite *pipelineTestSuite) TestDoneCase() {
	result := make([]string, 0, 10)
	start := time.Now()
	for s := range ExecutePipeline(suite.in, suite.done, suite.stages...) {
		result = append(result, s.(string))
	}
	elapsed := time.Since(start)
	suite.Len(result, 0)
	suite.Less(int64(elapsed), int64(suite.tc.stopAfter)+int64(fault))
}

func (suite *pipelineTestSuite) TestFilterCase() {
	result := make([]string, 0, 10)
	start := time.Now()
	for s := range ExecutePipeline(suite.in, suite.done, suite.stages...) {
		result = append(result, s.(string))
	}
	elapsed := time.Since(start)
	suite.Equal([]string{"1", "7", "11"}, result)
	suite.Less(
		int64(elapsed),
		int64(sleepPerStage)*int64(len(suite.stages)+len(suite.tc.data)-1)+int64(fault))
}

func TestPipeline(t *testing.T) {
	suite.Run(t, new(pipelineTestSuite))
}
