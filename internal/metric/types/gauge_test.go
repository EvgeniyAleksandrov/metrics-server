package types_test

import (
	"testing"

	"github.com/EvgeniyAleksandrov/metrics-server/internal/metric/types"
	"github.com/stretchr/testify/require"
)

func TestGaugeFromString(t *testing.T) {
	testCases := map[string]struct {
		s          string
		want       types.Gauge
		errContent string
	}{
		"StringWithIntNumber": {
			s:    "123",
			want: types.Gauge(123.0),
		},
		"StringWithNegativeIntNumber": {
			s:    "-123",
			want: types.Gauge(-123.0),
		},
		"StringWithFloatNumber": {
			s:    "123.1",
			want: types.Gauge(123.1),
		},
		"StringWithNegativeFloatNumber": {
			s:    "123.1",
			want: types.Gauge(123.1),
		},
		"StringWithoutNumber": {
			s:          "a",
			want:       types.Gauge(0),
			errContent: "parse string to float64",
		},
	}

	for name, tt := range testCases {
		t.Run(name, func(t *testing.T) {
			got, err := types.GaugeFromString(tt.s)
			if len(tt.errContent) > 0 {
				require.ErrorContains(t, err, tt.errContent)
			}

			require.Equal(t, tt.want, got)
		})
	}
}
