package hw05parallelexecution

import (
	"fmt"
	"math/rand"
	"sync/atomic"
	"time"
)

func (suite *parallelExecutionTestSuite) NextCase(testName string) *testCase {
	switch testName {
	case "TestFirstMErrors":
		return &testCase{
			workers:    10,
			maxErrors:  23,
			tasksCount: 50,
			generator: func(ts *parallelExecutionTestSuite) []Task {
				tasks := make([]Task, 0, ts.tc.tasksCount)
				for i := 0; i < ts.tc.tasksCount; i++ {
					err := fmt.Errorf("error from task %d", i)
					tasks = append(tasks, func() error {
						time.Sleep(time.Millisecond * time.Duration(rand.Intn(100)))
						atomic.AddInt32(&ts.rs.totalRuns, 1)
						atomic.AddInt32(&ts.rs.errors, 1)
						return err
					})
				}
				return tasks
			},
		}
	case "TestNoErrors":
		return &testCase{
			workers:    5,
			maxErrors:  1,
			tasksCount: 50,
			generator: func(ts *parallelExecutionTestSuite) []Task {
				tasks := make([]Task, 0, ts.tc.tasksCount)
				for i := 0; i < ts.tc.tasksCount; i++ {
					taskSleep := time.Millisecond * time.Duration(rand.Intn(100))
					ts.rs.pureRunningTime += taskSleep
					tasks = append(tasks, func() error {
						time.Sleep(taskSleep)
						atomic.AddInt32(&ts.rs.totalRuns, 1)
						return nil
					})
				}

				return tasks
			},
		}
	case "TestWithRealWorkAndWOErrors":
		return &testCase{
			workers:    5,
			maxErrors:  1,
			tasksCount: 50,
			generator: func(ts *parallelExecutionTestSuite) []Task {
				tasks := make([]Task, 0, ts.tc.tasksCount)
				for i := 0; i < ts.tc.tasksCount; i++ {
					tasks = append(tasks, func(j int64) Task {
						return func() error {
							start := time.Now()
							var k, np1, n, nm1 int64
							n, nm1 = 1, 1
							for k = 0; k < j; k++ {
								np1 = n + nm1
								nm1, n = n, np1
							}
							atomic.AddInt64((*int64)(&ts.rs.pureRunningTime), int64(time.Since(start)))
							atomic.AddInt32(&ts.rs.totalRuns, 1)
							return nil
						}
					}(int64(i+10)*(100*rand.Int63n(25000))))
				}

				return tasks
			},
		}
	case "TestWithRealWorkAndWithSomeErrors":
		return &testCase{
			workers:    5,
			maxErrors:  15,
			tasksCount: 60,
			generator: func(ts *parallelExecutionTestSuite) []Task {
				tasks := make([]Task, 0, ts.tc.tasksCount)
				for i := 0; i < ts.tc.tasksCount; i++ {
					err := fmt.Errorf("error from task %d", i)
					i := i
					tasks = append(tasks, func() error {
						start := time.Now()
						var j, np1, n, nm1, l int64
						n, nm1, l = 1, 1, int64(i)*10*rand.Int63n(25000)
						for j = 0; j < l; j++ {
							np1 = n + nm1
							nm1, n = n, np1
						}
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
	default:
		return nil
	}
}
