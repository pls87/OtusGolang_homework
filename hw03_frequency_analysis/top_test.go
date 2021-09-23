package hw03frequencyanalysis_test

import (
	"testing"

	hw03frequencyanalysis "github.com/pls87/OtusGolang_homework/hw03_frequency_analysis"
	"github.com/stretchr/testify/require"
)

func TestTop10(t *testing.T) {
	for _, tc := range tests {
		tc := tc
		t.Run(tc.name, func(t *testing.T) {
			require.Equal(t, tc.expected, hw03frequencyanalysis.Top10(tc.input))
		})
	}
}
