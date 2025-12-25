package api_test

import (
	"net/http"
	"net/url"
	"strings"
	"testing"

	"github.com/stretchr/testify/require"
)

func TestExtra(t *testing.T) {
	// Ensure DB is ready
	token := login(t)
	_, err := update(token)
	require.NoError(t, err)

	t.Run("search special chars", SearchSpecialChars)
	t.Run("access with invalid token", AccessWithInvalidToken)
	t.Run("search very long phrase", SearchVeryLongPhrase)
}

func SearchSpecialChars(t *testing.T) {
	// Attempt injections or just weird characters that shouldn't break the server
	phrases := []string{
		"' OR 1=1; --",
		"<script>alert(1)</script>",
		"🤔",
		"  ",
	}

	for _, p := range phrases {
		t.Run(p, func(t *testing.T) {
			resp, err := client.Get(address + "/api/search?phrase=" + url.QueryEscape(p))
			require.NoError(t, err)
			defer resp.Body.Close()
			// Should return 200 OK (empty list) or 400 Bad Request if validation kicks in
			require.Contains(t, []int{http.StatusOK, http.StatusBadRequest}, resp.StatusCode, "should handle special chars gracefully")
		})
	}
}

func AccessWithInvalidToken(t *testing.T) {
	invalidTokens := []string{
		"Bearer invalid.token.structure",
		"Token invalid.token.structure",
		"Basic notbase64",
		"",
	}

	for _, tok := range invalidTokens {
		t.Run(tok, func(t *testing.T) {
			req, err := http.NewRequest(http.MethodPost, address+"/api/db/update", nil)
			require.NoError(t, err)
			if tok != "" {
				req.Header.Add("Authorization", tok)
			}
			resp, err := client.Do(req)
			require.NoError(t, err)
			defer resp.Body.Close()
			require.Equal(t, http.StatusUnauthorized, resp.StatusCode)
		})
	}
}

func SearchVeryLongPhrase(t *testing.T) {
	// Very long phrase from "Repeat"
	longPhrase := strings.Repeat("linux ", 500)
	resp, err := client.Get(address + "/api/search?phrase=" + url.QueryEscape(longPhrase))
	require.NoError(t, err)
	defer resp.Body.Close()
	// Should not crash, validation might reject it or just return empty
	require.Contains(t, []int{http.StatusOK, http.StatusBadRequest, http.StatusRequestURITooLong}, resp.StatusCode)
}
