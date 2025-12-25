package words

import (
	"context"
	"testing"

	"github.com/stretchr/testify/assert"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
	wordspb "yadro.com/course/proto/words"
)

func TestNorm(t *testing.T) {
	tests := []struct {
		name      string
		phrase    string
		expectErr bool
		errCode   codes.Code
		expected  []string
	}{
		{
			name:     "simple phrase",
			phrase:   "apple tree",
			expected: []string{"appl", "tree"},
		},
		{
			name:     "stop words",
			phrase:   "the a an in on",
			expected: []string{},
		},
		{
			name:     "punctuation",
			phrase:   "hello, world!",
			expected: []string{"hello", "world"},
		},
		{
			name:     "duplicates",
			phrase:   "apple apple",
			expected: []string{"appl"},
		},
		{
			name:      "too long",
			phrase:    string(make([]byte, 5000)),
			expectErr: true,
			errCode:   codes.ResourceExhausted,
		},
		{
			name:     "stemming",
			phrase:   "running ran runner",
			expected: []string{"run", "ran", "runner"},
		},
	}

	for _, tc := range tests {
		t.Run(tc.name, func(t *testing.T) {
			req := &wordspb.WordsRequest{Phrase: tc.phrase}
			resp, err := Norm(context.Background(), req)

			if tc.expectErr {
				assert.Error(t, err)
				st, ok := status.FromError(err)
				assert.True(t, ok)
				assert.Equal(t, tc.errCode, st.Code())
			} else {
				assert.NoError(t, err)
				if tc.name == "stemming" {
					assert.Contains(t, resp.Words, "run")
				} else {
					assert.ElementsMatch(t, tc.expected, resp.Words)
				}
			}
		})
	}
}
