package hw05parallelexecution

import (
	"fmt"
	"math/rand"
	"sync/atomic"
	"time"
)

func (suite *parallelExecutionTestSuite) NextCase(testName string) *testCase {
	longOperation := func() int64 {
		var j, np1, n, nm1, l int64
		n, nm1, l = 1, 1, 10000-rand.Int63n(5000)
		for j = 0; j < l; j++ {
			np1 = n + nm1 + rand.Int63n(100)
			nm1, n = n, np1
		}
		return np1
	}

	getTasksAlwaysErrors := func(ts *parallelExecutionTestSuite) []Task {
		tasks := make([]Task, 0, ts.tc.tasksCount)
		for i := 0; i < ts.tc.tasksCount; i++ {
			err := fmt.Errorf("error from task %d", i)
			tasks = append(tasks, func() error {
				timeToSleep := 10*time.Millisecond + time.Millisecond*time.Duration(rand.Intn(100))
				suite.Eventually(func() bool { return true }, time.Second, timeToSleep)

				atomic.AddInt64((*int64)(&ts.rs.pureRunningTime), int64(timeToSleep))
				atomic.AddInt32(&ts.rs.totalRuns, 1)
				atomic.AddInt32(&ts.rs.errors, 1)
				return err
			})
		}
		return tasks
	}

	getTasksNoErrors := func(ts *parallelExecutionTestSuite) []Task {
		tasks := make([]Task, 0, ts.tc.tasksCount)
		for i := 0; i < ts.tc.tasksCount; i++ {
			tasks = append(tasks, func() error {
				timeToSleep := 10*time.Millisecond + time.Millisecond*time.Duration(rand.Intn(100))
				suite.Eventually(func() bool { return true }, time.Second, timeToSleep)

				atomic.AddInt64((*int64)(&ts.rs.pureRunningTime), int64(timeToSleep))
				atomic.AddInt32(&ts.rs.totalRuns, 1)
				return nil
			})
		}

		return tasks
	}

	switch testName {
	case "TestFirstMErrors":
		return &testCase{
			workers:    10,
			maxErrors:  23,
			tasksCount: 50,
			generator:  getTasksAlwaysErrors,
		}
	case "TestNoErrors":
		return &testCase{
			workers:    5,
			maxErrors:  1,
			tasksCount: 50,
			generator:  getTasksNoErrors,
		}
	case "TestWithRealWorkAndWOErrors":
		return &testCase{
			workers:    5,
			maxErrors:  1,
			tasksCount: 50,
			generator: func(ts *parallelExecutionTestSuite) []Task {
				tasks := make([]Task, 0, ts.tc.tasksCount)
				for i := 0; i < ts.tc.tasksCount; i++ {
					tasks = append(tasks, func() error {
						start := time.Now()
						longOperation()
						atomic.AddInt64((*int64)(&ts.rs.pureRunningTime), int64(time.Since(start)))
						atomic.AddInt32(&ts.rs.totalRuns, 1)
						return nil
					})
				}

				return tasks
			},
		}
	case "TestWithRealWorkAndWithSomeErrors":
		return &testCase{
			workers:    5,
			maxErrors:  64,
			tasksCount: 210,
			generator: func(ts *parallelExecutionTestSuite) []Task {
				tasks := make([]Task, 0, ts.tc.tasksCount)
				for i := 0; i < ts.tc.tasksCount; i++ {
					err := fmt.Errorf("error from task %d", i)
					i := i
					tasks = append(tasks, func() error {
						start := time.Now()
						longOperation()
						atomic.AddInt64((*int64)(&ts.rs.pureRunningTime), int64(time.Since(start)))
						atomic.AddInt32(&ts.rs.totalRuns, 1)
						if i%3 == 0 {
							atomic.AddInt32(&ts.rs.errors, 1)
							return err
						}
						return nil
					})
				}

				return tasks
			},
		}
	case "TestNegativeErrorsCountMeansIgnoreErrors":
		return &testCase{
			workers:    5,
			maxErrors:  -1,
			tasksCount: 50,
			generator:  getTasksAlwaysErrors,
		}
	case "TestWorkersGtTasksAndMaxErrorsGtTasks":
		return &testCase{
			workers:    60,
			maxErrors:  30,
			tasksCount: 20,
			generator:  getTasksAlwaysErrors,
		}
	default:
		return nil
	}
}
