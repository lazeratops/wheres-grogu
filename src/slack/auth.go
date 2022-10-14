package slack

import (
	"crypto/hmac"
	"crypto/sha256"
	"encoding/hex"
	"errors"
	"fmt"
)

type ValidationParams struct {
	Headers map[string]string
	Body    string
}

func validateSignature(headers map[string]string, body string, signingSecret string) error {
	var slackSig, reqTimestamp string

	const errMsg = "unauthorized request"
	if sig, ok := headers["x-slack-signature"]; ok {
		slackSig = sig
	} else {
		return fmt.Errorf("%s: missing signature", errMsg)
	}
	// Check if the key had an empty value just in case
	if slackSig == "" {
		return fmt.Errorf("%s: invalid signature", errMsg)
	}

	if timestamp, ok := headers["x-slack-request-timestamp"]; ok {
		reqTimestamp = timestamp
	} else {
		return fmt.Errorf("%s: missing request timestamp", errMsg)
	}

	version := getVersion(reqTimestamp, body)
	sig := computeSig(version, signingSecret)
	if hmac.Equal([]byte(sig), []byte(slackSig)) {
		return nil
	}
	return errors.New("failed to validate given auth signature")
}

func getVersion(timestamp string, body string) string {
	return fmt.Sprintf("v0:%s:%s", timestamp, body)
}

func computeSig(version, signingSecret string) string {
	sig := hmac.New(sha256.New, []byte(signingSecret))
	sig.Write([]byte(version))
	return fmt.Sprintf("v0=%s", hex.EncodeToString(sig.Sum(nil)))
}
