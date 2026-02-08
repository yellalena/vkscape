package auth

import (
	"bufio"
	"crypto/rand"
	"crypto/sha256"
	"encoding/base64"
	"encoding/json"
	"fmt"
	"io"
	"log/slog"
	"net/http"
	"net/url"
	"os"
	"os/exec"
	"runtime"
	"strings"

	"github.com/joho/godotenv"

	"github.com/yellalena/vkscape/internal/config"
	"github.com/yellalena/vkscape/internal/output"
)

var _ = godotenv.Load()
var vkClientID = "51812294"

const (
	vkRedirectURI  = "https://oauth.vk.com/blank.html"
	vkScope        = "wall,photos,groups,offline"
	vkAuthEndpoint = "https://id.vk.com/authorize"
	//nolint:gosec // not a credential
	vkTokenEndpoint = "https://id.vk.com/oauth2/auth"
)

type TokenResponse struct {
	AccessToken  string `json:"access_token"`
	RefreshToken string `json:"refresh_token"`
	ExpiresIn    int    `json:"expires_in"`
	TokenType    string `json:"token_type"`
}

type AuthSession struct {
	Verifier string
	AuthURL  string
}

func InteractiveFlow(logger *slog.Logger) error {
	session, err := StartInteractiveFlow(logger)
	if err != nil {
		return err
	}

	OpenBrowser(session.AuthURL, logger)
	output.Info("After authorizing, you'll be redirected to a blank page.")
	output.Info("Copy the FULL URL from the address bar and paste it here:")
	reader := bufio.NewReader(os.Stdin)
	redirectURL, err := reader.ReadString('\n')
	if err != nil {
		logger.Error("Failed to read redirect URL", "error", err)
		return fmt.Errorf("could not read input")
	}

	return FinishInteractiveFlow(logger, session.Verifier, redirectURL)
}

func StartInteractiveFlow(logger *slog.Logger) (*AuthSession, error) {
	verifier, challenge, err := generatePKCE()
	if err != nil {
		logger.Error("Failed to generate PKCE", "error", err)
		return nil, fmt.Errorf("internal error")
	}

	authURL := fmt.Sprintf(
		"%s?response_type=code&client_id=%s&redirect_uri=%s&scope=%s&state=12345&code_challenge=%s&code_challenge_method=S256",
		vkAuthEndpoint,
		vkClientID,
		url.QueryEscape(vkRedirectURI),
		vkScope,
		challenge,
	)

	logger.Info("Starting interactive login flow with url", "url", authURL)
	output.Info("Starting interactive login flow... Your browser will be opened.")
	output.Info("If the the browser didn't open or you don't see VK login page, please open this URL manually and login:")
	output.Info(authURL)
	return &AuthSession{
		Verifier: verifier,
		AuthURL:  authURL,
	}, nil
}

func FinishInteractiveFlow(logger *slog.Logger, verifier, redirectURL string) error {
	redirectURL = strings.TrimSpace(redirectURL)

	u, err := url.Parse(redirectURL)
	if err != nil {
		logger.Error("Failed to parse redirect URL", "error", err)
		return fmt.Errorf("could not parse URL")
	}
	q := u.Query()
	code := q.Get("code")
	deviceID := q.Get("device_id")

	if code == "" || deviceID == "" {
		logger.Error("Authorization code or device_id not found in URL", "redirect_url", redirectURL)
		return fmt.Errorf("authorization code or device_id not found in URL")
	}

	logger.Info("Received code and device_id; exchanging for token")
	output.Info("Exchanging authorization code for token...")

	form := url.Values{}
	form.Add("grant_type", "authorization_code")
	form.Add("code", code)
	form.Add("client_id", vkClientID)
	form.Add("redirect_uri", vkRedirectURI)
	form.Add("device_id", deviceID)
	form.Add("code_verifier", verifier)

	req, err := http.NewRequest(http.MethodPost, vkTokenEndpoint, strings.NewReader(form.Encode()))
	if err != nil {
		logger.Error("Failed to create token exchange request", "error", err)
		return fmt.Errorf("internal error")
	}
	req.Header.Add("Content-Type", "application/x-www-form-urlencoded")

	resp, err := http.DefaultClient.Do(req)
	if err != nil {
		logger.Error("Token exchange request failed", "error", err)
		return fmt.Errorf("internal error")
	}
	defer resp.Body.Close() //nolint:errcheck

	if resp.StatusCode != http.StatusOK {
		logger.Error("VK token exchange failed", "status", resp.Status)
		return fmt.Errorf("internal error")
	}

	bodyBytes, err := io.ReadAll(resp.Body)
	if err != nil {
		logger.Error("Failed to read token exchange response", "error", err)
		return fmt.Errorf("internal error")
	}

	var token TokenResponse
	err = json.Unmarshal(bodyBytes, &token)
	if err != nil {
		logger.Error("Failed to parse token JSON", "error", err)
		return fmt.Errorf("internal error")
	}

	if err := config.SaveConfig(&config.AuthConfig{
		AuthMethod:   config.AuthMethodUser,
		AccessToken:  token.AccessToken,
		RefreshToken: token.RefreshToken,
	}); err != nil {
		logger.Error("Failed to save auth config", "error", err)
		return fmt.Errorf("could not save token")
	}

	logger.Info("Authentication successful")
	return nil
}

func generatePKCE() (verifier, challenge string, err error) {
	const verifyerLength = 32
	verifierBytes := make([]byte, verifyerLength)
	if _, err = rand.Read(verifierBytes); err != nil {
		err = fmt.Errorf("generate PKCE: %w", err)
		return "", "", err
	}
	verifier = base64.RawURLEncoding.EncodeToString(verifierBytes)
	sha := sha256.Sum256([]byte(verifier))
	challenge = base64.RawURLEncoding.EncodeToString(sha[:])
	return verifier, challenge, nil
}

func OpenBrowser(url string, logger *slog.Logger) {
	var cmd string
	var args []string

	switch runtime.GOOS {
	case "windows":
		cmd = "rundll32"
		args = []string{"url.dll,FileProtocolHandler", url}
	case "darwin":
		cmd = "open"
		args = []string{url}
	default:
		cmd = "xdg-open"
		args = []string{url}
	}

	err := exec.Command(cmd, args...).Start() //nolint:gosec // cmd is controlled by runtime.GOOS
	if err != nil {
		logger.Warn("Failed to open browser", "error", err, "url", url)
		logger.Info("Please open the following URL manually", "url", url)
	}
}
