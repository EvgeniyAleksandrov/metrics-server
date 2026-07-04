package types_test

import (
	"testing"

	"github.com/EvgeniyAleksandrov/metrics-server/internal/metric/types"
	"github.com/stretchr/testify/require"
)

func TestCounterFromString(t *testing.T) {
	testCases := map[string]struct {
		s          string
		want       types.Counter
		errContent string
	}{
		"StringWithIntNumber": {
			s:    "123",
			want: types.Counter(123),
		},
		"StringWithNegativeIntNumber": {
			s:    "-123",
			want: types.Counter(-123),
		},
		"StringWithFloatNumber": {
			s:          "123.1",
			want:       types.Counter(0),
			errContent: "parse string to int64",
		},
		"StringWithoutNumber": {
			s:          "a",
			want:       types.Counter(0),
			errContent: "parse string to int64",
		},
	}

	for name, tt := range testCases {
		t.Run(name, func(t *testing.T) {
			got, err := types.CounterFromString(tt.s)
			if len(tt.errContent) > 0 {
				require.ErrorContains(t, err, tt.errContent)
			}

			require.Equal(t, tt.want, got)
		})
	}
}
