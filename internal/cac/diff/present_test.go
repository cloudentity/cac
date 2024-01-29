package diff_test

import (
	"github.com/cloudentity/cac/internal/cac/diff"
	"github.com/stretchr/testify/require"
	"testing"
)

func TestOnlyPresent(t *testing.T) {
	tcs := []struct {
		name     string
		source   map[string]any
		target   map[string]any
		expected map[string]any
	}{
		{
			name:     "should remove keys that are not present in source",
			source:   map[string]any{"a": 1},
			target:   map[string]any{"a": 1, "b": 2},
			expected: map[string]any{"a": 1},
		},
		{
			name:     "should not remove keys that are present in source",
			source:   map[string]any{"a": 1, "b": 2, "c": 3},
			target:   map[string]any{"a": 1, "b": 2},
			expected: map[string]any{"a": 1, "b": 2},
		},
		{
			name:     "should remove all keys if source map is nil",
			source:   nil,
			target:   map[string]any{"a": 1, "b": 2},
			expected: map[string]any{},
		},
	}

	for _, tc := range tcs {
		t.Run(tc.name, func(t *testing.T) {
			diff.OnlyPresentKeys(tc.source, tc.target)

			require.Equal(t, tc.expected, tc.target)
		})
	}
}
