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

type stageStub struct {
	take        func(v interface{}) bool
	transformer func(v interface{}) interface{}
}

type testCase struct {
	data      []interface{}
	expected  []interface{}
	generator func(suite *pipelineTestSuite, sb stageStub) Stage
	stg       []stageStub
	stopAfter time.Duration
}

type pipelineTestSuite struct {
	suite.Suite
	sync.WaitGroup
	tc       *testCase
	result   []interface{}
	stages   []Stage
	in, done Bi
}

func (suite *pipelineTestSuite) SetupTest() {
	suite.in = make(Bi)
	suite.done = make(Bi)
	suite.result = make([]interface{}, 0, 10)
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
	start := time.Now()
	for s := range ExecutePipeline(suite.in, suite.done, suite.stages...) {
		suite.result = append(suite.result, s)
	}
	elapsed := time.Since(start)
	suite.Equal(suite.tc.expected, suite.result)
	suite.Less(
		int64(elapsed),
		int64(sleepPerStage)*int64(len(suite.stages)+len(suite.tc.data)-1)+int64(fault))
}

func (suite *pipelineTestSuite) TestDoneCase() {
	start := time.Now()
	for s := range ExecutePipeline(suite.in, suite.done, suite.stages...) {
		suite.result = append(suite.result, s.(string))
	}
	elapsed := time.Since(start)
	suite.Equal(suite.tc.expected, suite.result)
	suite.Less(int64(elapsed), int64(suite.tc.stopAfter)+int64(fault))
}

func (suite *pipelineTestSuite) TestFilterCase() {
	start := time.Now()
	for s := range ExecutePipeline(suite.in, suite.done, suite.stages...) {
		suite.result = append(suite.result, s.(string))
	}
	elapsed := time.Since(start)
	suite.Equal(suite.tc.expected, suite.result)
	suite.Less(
		int64(elapsed),
		int64(sleepPerStage)*int64(len(suite.stages)+len(suite.tc.data)-1)+int64(fault))
}

func (suite *pipelineTestSuite) TestStringsCase() {
	start := time.Now()
	for s := range ExecutePipeline(suite.in, suite.done, suite.stages...) {
		suite.result = append(suite.result, s.(string))
	}
	elapsed := time.Since(start)
	suite.Equal(suite.tc.expected, suite.result)
	suite.Less(
		int64(elapsed),
		int64(sleepPerStage)*int64(len(suite.stages)+len(suite.tc.data)-1)+int64(fault))
}

func (suite *pipelineTestSuite) TestEmptyCase() {
	start := time.Now()
	for s := range ExecutePipeline(suite.in, suite.done, suite.stages...) {
		suite.result = append(suite.result, s.(int))
	}
	elapsed := time.Since(start)
	suite.Equal(suite.tc.expected, suite.result)
	suite.Less(
		int64(elapsed),
		int64(fault))
}

func (suite *pipelineTestSuite) TestDoneAfterFirstCoupleCase() {
	start := time.Now()
	for s := range ExecutePipeline(suite.in, suite.done, suite.stages...) {
		suite.result = append(suite.result, s.(string))
	}
	elapsed := time.Since(start)
	suite.Equal(suite.tc.expected, suite.result)
	suite.Less(int64(elapsed), int64(suite.tc.stopAfter)+int64(fault))
}

func TestPipeline(t *testing.T) {
	suite.Run(t, new(pipelineTestSuite))
}
