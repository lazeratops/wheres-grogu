package main

import (
	"github.com/stretchr/testify/require"
	"net/http"
	"testing"
)

// The test data for this test comes from Slack documentation
// https://api.slack.com/authentication/verifying-requests-from-slack
func TestGetVersion(t *testing.T) {
	testCases := []struct {
		name        string
		timestamp   string
		body        string
		wantVersion string
	}{
		{
			name:        "slack-ex-1",
			timestamp:   "1531420618",
			body:        "token=xyzz0WbapA4vBCDEFasx0q6G&team_id=T1DC2JH3J&team_domain=testteamnow&channel_id=G8PSS9T3V&channel_name=foobar&user_id=U2CERLKJA&user_name=roadrunner&command=%2Fwebhook-collect&text=&response_url=https%3A%2F%2Fhooks.slack.com%2Fcommands%2FT1DC2JH3J%2F397700885554%2F96rGlfmibIGlgcZRskXaIFfN&trigger_id=398738663015.47445629121.803a0bc887a14d10d2c447fce8b6703c",
			wantVersion: "v0:1531420618:token=xyzz0WbapA4vBCDEFasx0q6G&team_id=T1DC2JH3J&team_domain=testteamnow&channel_id=G8PSS9T3V&channel_name=foobar&user_id=U2CERLKJA&user_name=roadrunner&command=%2Fwebhook-collect&text=&response_url=https%3A%2F%2Fhooks.slack.com%2Fcommands%2FT1DC2JH3J%2F397700885554%2F96rGlfmibIGlgcZRskXaIFfN&trigger_id=398738663015.47445629121.803a0bc887a14d10d2c447fce8b6703c",
		},
	}

	for _, tc := range testCases {
		tc := tc
		t.Run(tc.name, func(t *testing.T) {
			gotVersion := getVersion(tc.timestamp, tc.body)
			require.EqualValues(t, tc.wantVersion, gotVersion)
		})
	}
}

// The test data for this test comes from Slack documentation
// https://api.slack.com/authentication/verifying-requests-from-slack
func TestCompareSignature(t *testing.T) {
	version := "v0:1531420618:token=xyzz0WbapA4vBCDEFasx0q6G&team_id=T1DC2JH3J&team_domain=testteamnow&channel_id=G8PSS9T3V&channel_name=foobar&user_id=U2CERLKJA&user_name=roadrunner&command=%2Fwebhook-collect&text=&response_url=https%3A%2F%2Fhooks.slack.com%2Fcommands%2FT1DC2JH3J%2F397700885554%2F96rGlfmibIGlgcZRskXaIFfN&trigger_id=398738663015.47445629121.803a0bc887a14d10d2c447fce8b6703c"
	signingSecret := "8f742231b10e8888abcd99yyyzzz85a5"
	gotSig := computeSig(version, signingSecret)
	wantSig := "v0=a2114d57b48eac39b9ad189dd8316235a7b4a8d21a10bd27519666489c69b503"
	require.EqualValues(t, wantSig, gotSig)
}

// This test function is simply for local testing,
// do not try to use it as a unit test.
func TestNetlifyDev(t *testing.T) {
	// Uncomment this line when testing locally
	t.Skip()
	// This is intended to be run alongside `npm run dev`
	const url = `http://localhost:9999/.netlify/functions/grogu?token=23424&command=grogu`
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
