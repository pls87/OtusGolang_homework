package main

import (
	"testing"

	"github.com/stretchr/testify/require"
)

type testCase struct {
	name     string
	path     string
	expected Environment
	err      error
}

func TestReadDir(t *testing.T) {
	cases := []testCase{
		{
			name: "simple case based on env testdata",
			path: "./testdata/env",
			expected: Environment{
				"BAR":   EnvValue{Value: "bar", NeedRemove: false},
				"EMPTY": EnvValue{Value: "", NeedRemove: false},
				"FOO":   EnvValue{Value: "   foo\nwith new line", NeedRemove: false},
				"HELLO": EnvValue{Value: "\"hello\"", NeedRemove: false},
				"UNSET": EnvValue{Value: "", NeedRemove: true},
			},
			err: nil,
		}, {
			name: "case with plus sign and spaces in the name of variable",
			path: "./testdata/more_tests/test_1",
			expected: Environment{
				"SPACES_in    THE_NAME": EnvValue{Value: "hmmm...what will happen in such case?", NeedRemove: false},
				"TRIVIAL":               EnvValue{Value: "Nothing interesting...\nAbsolutely", NeedRemove: false},
			},
			err: nil,
		}, {
			name: "case with empty line and zero file",
			path: "./testdata/more_tests/test_2",
			expected: Environment{
				"EMPTY_FIRST_LINE": EnvValue{Value: "", NeedRemove: false},
				"ZERO_FILE":        EnvValue{Value: "", NeedRemove: true},
			},
			err: nil,
		}, {
			name: "case with no permission",
			path: "./testdata/more_tests/test_2",
			expected: Environment{
				"EMPTY_FIRST_LINE": EnvValue{Value: "", NeedRemove: false},
				"ZERO_FILE":        EnvValue{Value: "", NeedRemove: true},
			},
			err: nil,
		},
	}

	for _, tc := range cases {
		tc := tc
		t.Run(tc.name, func(t *testing.T) {
			env, e := ReadDir(tc.path)
			if tc.err != nil {
				require.ErrorIs(t, e, tc.err, "actual error %q but expected %q", e, tc.err)
			} else {
				require.NoError(t, e, "No error expected, but got %q", e)
				require.Equal(t, tc.expected, env)
			}
		})
	}
}
