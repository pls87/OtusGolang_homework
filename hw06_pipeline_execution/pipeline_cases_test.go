package hw06pipelineexecution

import (
	"strconv"
	"strings"
	"time"
)

func (suite *pipelineTestSuite) NextCase(testName string) *testCase {
	simpleStages := []stageBase{
		{
			take:        func(v interface{}) bool { return true },
			transformer: func(v interface{}) interface{} { return v },
		},
		{
			take:        func(v interface{}) bool { return true },
			transformer: func(v interface{}) interface{} { return v.(int) * 2 },
		},
		{
			take:        func(v interface{}) bool { return true },
			transformer: func(v interface{}) interface{} { return v.(int) + 100 },
		},
		{
			take:        func(v interface{}) bool { return true },
			transformer: func(v interface{}) interface{} { return strconv.Itoa(v.(int)) },
		},
	}

	filterStages := []stageBase{
		{
			take:        func(v interface{}) bool { return v.(int)%2 != 0 },
			transformer: func(v interface{}) interface{} { return v },
		},
		{
			take:        func(v interface{}) bool { return v.(int)%3 != 0 },
			transformer: func(v interface{}) interface{} { return v },
		},
		{
			take:        func(v interface{}) bool { return v.(int)%5 != 0 },
			transformer: func(v interface{}) interface{} { return v },
		},
		{
			take:        func(v interface{}) bool { return true },
			transformer: func(v interface{}) interface{} { return strconv.Itoa(v.(int)) },
		},
	}

	stringStages := []stageBase{
		{
			take:        func(v interface{}) bool { return true },
			transformer: func(v interface{}) interface{} { return strings.TrimSpace(v.(string)) },
		},
		{
			take:        func(v interface{}) bool { return len(v.(string)) > 2 },
			transformer: func(v interface{}) interface{} { return strings.ToLower(v.(string)) },
		},
		{
			take: func(v interface{}) bool { return true },
			transformer: func(v interface{}) interface{} {
				runes := []rune(v.(string))
				for i, j := 0, len(runes)-1; i < j; i, j = i+1, j-1 {
					runes[i], runes[j] = runes[j], runes[i]
				}
				return string(runes)
			},
		},
		{
			take:        func(v interface{}) bool { return true },
			transformer: func(v interface{}) interface{} { return strings.ReplaceAll(v.(string), "o", "0") },
		},
	}

	simpleGenerator := func(suite *pipelineTestSuite, sb stageBase) Stage {
		return func(in In) Out {
			out := make(Bi)
			go func() {
				defer close(out)
				for v := range in {
					suite.Eventually(func() bool { return true }, time.Second, sleepPerStage)
					if sb.take(v) {
						out <- sb.transformer(v)
					}
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
	case "TestFilterCase":
		return &testCase{
			data:      []interface{}{1, 2, 3, 4, 5, 6, 7, 8, 9, 10, 11},
			generator: simpleGenerator,
			stg:       filterStages,
			stopAfter: -1,
		}

	case "TestStringsCase":
		return &testCase{
			data:      []interface{}{" One  ", "Ring", " to ", "Rule  ", "  THEM", "aLL"},
			generator: simpleGenerator,
			stg:       stringStages,
			stopAfter: -1,
		}

	case "TestEmptyCase":
		return &testCase{
			data:      []interface{}{},
			generator: simpleGenerator,
			stg:       simpleStages,
			stopAfter: -1,
		}
	default:
		return nil
	}
}
