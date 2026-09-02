package util

import (
	"bufio"
	"context"
	"crypto/sha1"
	"fmt"
	"net/http"
	"strconv"
	"strings"
	"time"
)

var httpClient = &http.Client{Timeout: 3 * time.Second}

// CheckPwned returns the number of times this password has appeared in
// known breaches (0 = not found), or an error if the check itself failed.
func CheckPwned(ctx context.Context, plaintext string) (int, error) {
	sum := fmt.Sprintf("%X", sha1.Sum([]byte(plaintext))) // uppercase hex
	prefix, suffix := sum[:5], sum[5:]

	req, err := http.NewRequestWithContext(ctx, http.MethodGet,
		"https://api.pwnedpasswords.com/range/"+prefix, nil)
	if err != nil {
		return 0, err
	}
	req.Header.Set("Add-Padding", "true") // mitigates response-size side-channel

	resp, err := httpClient.Do(req)
	if err != nil {
		return 0, err
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		return 0, fmt.Errorf("pwnedpasswords: unexpected status %d", resp.StatusCode)
	}

	scanner := bufio.NewScanner(resp.Body)
	for scanner.Scan() {
		line := scanner.Text() // format: SUFFIX:COUNT
		parts := strings.SplitN(line, ":", 2)
		if len(parts) != 2 {
			continue
		}
		if parts[0] == suffix {
			count, _ := strconv.Atoi(parts[1])
			return count, nil
		}
	}
	return 0, scanner.Err()
}
