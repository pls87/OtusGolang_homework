package main

import (
	"errors"
	"testing"

	"github.com/stretchr/testify/require"
)

func TestInitialFileChecks(t *testing.T) {
	cases := []struct {
		initialParams  CopyParams
		expectedParams CopyParams
		err            error
		name           string
	}{
		{
			initialParams:  CopyParams{from: "./testdata/input.txt", to: "out.txt", offset: 0, limit: 0},
			expectedParams: CopyParams{from: "./testdata/input.txt", to: "out.txt", offset: 0, limit: 6617},
			err:            nil,
			name:           "Very simple positive case",
		},
		{
			initialParams:  CopyParams{from: "./testdata/input.txt", to: "out.txt", offset: 0x100000, limit: 0},
			expectedParams: CopyParams{},
			err:            ErrOffsetExceedsFileSize,
			name:           "Offset greater than source file size generates error",
		},
		{
			initialParams:  CopyParams{from: "./testdata/input.txt", to: "out.txt", offset: 0, limit: 0x100000},
			expectedParams: CopyParams{from: "./testdata/input.txt", to: "out.txt", offset: 0, limit: 6617},
			err:            nil,
			name:           "Limit greater than source file size",
		},
		{
			initialParams:  CopyParams{from: "./testdata/input.txt", to: "out.txt", offset: 17, limit: 0x100000},
			expectedParams: CopyParams{from: "./testdata/input.txt", to: "out.txt", offset: 17, limit: 6600},
			err:            nil,
			name:           "Limit greater than source file size and not zero offset",
		},
		{
			initialParams:  CopyParams{from: "./testdata/input.txt", to: "out.txt", offset: -10, limit: 6000},
			expectedParams: CopyParams{from: "./testdata/input.txt", to: "out.txt", offset: 0, limit: 6000},
			err:            nil,
			name:           "Limit less than source file size and negative offset",
		},
		{
			initialParams:  CopyParams{from: "/dev/urandom", to: "out.txt", offset: 0, limit: 0},
			expectedParams: CopyParams{},
			err:            ErrUnsupportedFile,
			name:           "Not ordinary file: '/dev/urandom' generates error",
		},
	}
	for _, tc := range cases {
		tc := tc
		t.Run(tc.name, func(t *testing.T) {
			err := initialFileChecks(&tc.initialParams)
			if tc.err != nil {
				require.True(t, errors.Is(err, tc.err), "actual error %q", err)
			} else {
				require.Equal(t, tc.expectedParams, tc.initialParams)
			}
		})
	}
}
