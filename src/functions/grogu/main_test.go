package main

import (
	"github.com/stretchr/testify/require"
	"net/http"
	"testing"
)

// This test function is simply for local testing,
// do not try to use it as a unit test.
func TestNetlifyDev(t *testing.T) {
	// Uncomment this line when testing locally
	t.Skip()
	// This is intended to be run alongside `npm run dev`
	const url = `http://localhost:9999/.netlify/functions/found?token=23424&command=grogu`
	req, err := http.NewRequest("POST", url, nil)
	if err != nil {
		require.NoError(t, err)
	}

	req.Header.Set("Content-Type", "application/x-www-form-urlencoded")
	req.Header.Set("X-Slack-Signature", "testsig")
	req.Header.Set("X-Slack-Request-Timestamp", "234")

	res, err := http.DefaultClient.Do(req)
	require.NoError(t, err)
	require.EqualValues(t, http.StatusOK, res.StatusCode)
}
