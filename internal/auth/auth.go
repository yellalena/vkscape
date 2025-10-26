// internal/auth/pkce.go
package auth

import (
	"bufio"
	"crypto/rand"
	"crypto/sha256"
	"encoding/base64"
	"encoding/json"
	"errors"
	"fmt"
	"io"
	"net/http"
	"net/url"
	"os"
	"os/exec"
	"runtime"
	"strings"

	"github.com/joho/godotenv"

	"github.com/yellalena/vkscape/internal/config"
)

var _ = godotenv.Load()
var vkClientID = os.Getenv("VK_CLIENT_ID")

const (
	vkRedirectURI   = "https://oauth.vk.com/blank.html"
	vkScope         = "wall,photos,groups,offline"
	vkAuthEndpoint  = "https://id.vk.com/authorize"
	vkTokenEndpoint = "https://id.vk.com/oauth2/auth"
)

type TokenResponse struct {
	AccessToken  string `json:"access_token"`
	RefreshToken string `json:"refresh_token"`
	ExpiresIn    int    `json:"expires_in"`
	TokenType    string `json:"token_type"`
}

func InteractiveFlow() error {
	verifier, challenge := generatePKCE()

	authURL := fmt.Sprintf("%s?response_type=code&client_id=%s&redirect_uri=%s&scope=%s&state=12345&code_challenge=%s&code_challenge_method=S256",
		vkAuthEndpoint, vkClientID, url.QueryEscape(vkRedirectURI), vkScope, challenge)

	fmt.Println("\n🔐 Please open this URL in your browser and login:")
	fmt.Println(authURL)
	openBrowser(authURL)

	fmt.Println("\nAfter authorizing, you’ll be redirected to a blank page.")
	fmt.Println("Copy the FULL URL from the address bar and paste it here:")
	fmt.Print("Paste redirect URL: ")
	reader := bufio.NewReader(os.Stdin)
	redirectURL, _ := reader.ReadString('\n')
	redirectURL = strings.TrimSpace(redirectURL)

	u, err := url.Parse(redirectURL)
	if err != nil {
		return err
	}
	q := u.Query()
	code := q.Get("code")
	deviceID := q.Get("device_id")

	if code == "" || deviceID == "" {
		return errors.New("authorization code or device_id not found in URL")
	}

	fmt.Println("\n✅ Received code and device_id. Exchanging for token...")

	form := url.Values{}
	form.Add("grant_type", "authorization_code")
	form.Add("code", code)
	form.Add("client_id", vkClientID)
	form.Add("redirect_uri", vkRedirectURI)
	form.Add("device_id", deviceID)
	form.Add("code_verifier", verifier)

	req, _ := http.NewRequest("POST", vkTokenEndpoint, strings.NewReader(form.Encode()))
	req.Header.Add("Content-Type", "application/x-www-form-urlencoded")

	resp, err := http.DefaultClient.Do(req)
	if err != nil {
		return err
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		fmt.Printf("VK token exchange failed. Status: %s", resp.Status)
		return fmt.Errorf("VK token exchange failed: %s", resp.Status)
	}

	bodyBytes, _ := io.ReadAll(resp.Body)

	var token TokenResponse
	err = json.Unmarshal(bodyBytes, &token)
	if err != nil {
		return fmt.Errorf("failed to parse JSON: %w", err)
	}

	return config.SaveConfig(&config.AuthConfig{
		AuthMethod:   config.AuthMethodUser,
		AccessToken:  token.AccessToken,
		RefreshToken: token.RefreshToken,
	})
}

func generatePKCE() (verifier, challenge string) {
	verifierBytes := make([]byte, 32)
	_, _ = rand.Read(verifierBytes)
	verifier = base64.RawURLEncoding.EncodeToString(verifierBytes)
	sha := sha256.Sum256([]byte(verifier))
	challenge = base64.RawURLEncoding.EncodeToString(sha[:])
	return
}

func openBrowser(url string) {
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

	err := exec.Command(cmd, args...).Start()
	if err != nil {
		fmt.Printf("Failed to open browser. Please open the following URL manually: %s\n", url)
	}
}
