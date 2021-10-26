package hw06pipelineexecution

import (
	"strconv"
	"time"
)

func (suite *pipelineTestSuite) NextCase(testName string) *testCase {
	simpleStages := []stageBase{
		{
			"Dummy",
			func(v interface{}) interface{} { return v },
		},
		{
			"Multiplier (* 2)",
			func(v interface{}) interface{} { return v.(int) * 2 },
		},
		{
			"Adder (+ 100)",
			func(v interface{}) interface{} { return v.(int) + 100 },
		},
		{
			"Stringifier",
			func(v interface{}) interface{} { return strconv.Itoa(v.(int)) },
		},
	}
	simpleGenerator := func(sb stageBase) Stage {
		return func(in In) Out {
			out := make(Bi)
			go func() {
				defer close(out)
				for v := range in {
					time.Sleep(sleepPerStage)
					out <- sb.f(v)
				}
			}()
			return out
		}
	}

	switch testName {
	case "TestSimpleCase":
		return &testCase{
			data:      []interface{}{1, 2, 3, 4, 5},
			generator: simpleGenerator,
			stg:       simpleStages,
			stopAfter: -1,
		}
	case "TestDoneCase":
		return &testCase{
			data:      []interface{}{1, 2, 3, 4, 5},
			generator: simpleGenerator,
			stg:       simpleStages,
			stopAfter: sleepPerStage * 2,
		}
	default:
		return nil
	}
}
