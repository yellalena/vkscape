package auth

import (
	"log/slog"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestStartInteractiveFlowReturnsSession(t *testing.T) {
	origClientID := vkClientID
	vkClientID = "123"
	t.Cleanup(func() { vkClientID = origClientID })

	session, err := StartInteractiveFlow(slog.Default())

	assert.NoError(t, err)
	assert.NotNil(t, session)
	assert.NotEmpty(t, session.Verifier)
	assert.NotEmpty(t, session.AuthURL)
	assert.Contains(t, session.AuthURL, "code_challenge")
	assert.Contains(t, session.AuthURL, "scope")
}

func TestFinishInteractiveFlowRejectsMalformedURL(t *testing.T) {
	cases := []struct {
		name string
		url  string
	}{
		{name: "no-query", url: "https://oauth.vk.com/blank.html"},
		{name: "missing-code", url: "https://oauth.vk.com/blank.html?device_id=1"},
		{name: "missing-device", url: "https://oauth.vk.com/blank.html?code=1"},
	}

	for _, tc := range cases {
		t.Run(tc.name, func(t *testing.T) {
			err := FinishInteractiveFlow(slog.Default(), "verifier", tc.url)
			assert.Error(t, err)
		})
	}
}
