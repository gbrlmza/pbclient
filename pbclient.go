package pbclient

import (
	"context"
	"encoding/base64"
	"encoding/json"
	"errors"
	"fmt"
	"io"
	"net/http"
	"net/url"
	"reflect"
	"strings"
	"time"

	"github.com/gbrlmza/httpw"
	"github.com/golang-jwt/jwt"
)

type Client struct {
	host string
}

// New creates a new pbclient.Client instance.
func New(host string) (*Client, error) {
	host = strings.TrimSuffix(strings.TrimSpace(host), "/")
	if _, err := url.Parse(host); err != nil {
		return nil, err
	}

	return &Client{
		host: host,
	}, nil
}

// AuthAdminWithPassword authenticates an admin with a password.
func (c Client) AuthAdminWithPassword(ctx context.Context, identity string, password string) (*Token, error) {
	const path = "/api/admins/auth-with-password"

	if identity == "" || password == "" {
		return nil, errors.New("identity and password are required")
	}

	body := map[string]string{
		"identity": identity,
		"password": password,
	}
	resp, err := httpw.Do(ctx, http.MethodPost, c.host+path, httpw.WithJsonBody(body))
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()
	defer io.Copy(io.Discard, resp.Body) //nolint:errcheck

	if resp.StatusCode != http.StatusOK {
		return nil, fmt.Errorf("failed to authenticate (sc=%d)", resp.StatusCode)
	}

	return parseToken(resp.Body)
}

// AuthUserWithPassword authenticates a user with a password.
func (c Client) AuthUserWithPassword(ctx context.Context, identity string, password string) (*Token, error) {
	const path = "/api/collections/users/auth-with-password"

	if identity == "" || password == "" {
		return nil, errors.New("identity and password are required")
	}

	body := map[string]string{
		"identity": identity,
		"password": password,
	}
	resp, err := httpw.Do(ctx, http.MethodPost, c.host+path, httpw.WithJsonBody(body))
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()
	defer io.Copy(io.Discard, resp.Body) //nolint:errcheck

	if resp.StatusCode != http.StatusOK {
		return nil, fmt.Errorf("failed to authenticate (sc=%d)", resp.StatusCode)
	}

	return parseToken(resp.Body)
}

// AuthRefresh refreshes a token.
func (c Client) AuthRefresh(ctx context.Context, token string) (*Token, error) {
	const path = "/api/collections/users/auth-refresh"

	resp, err := httpw.Do(ctx, http.MethodPost, c.host+path, httpw.WithHeader("Authorization", token))
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()
	defer io.Copy(io.Discard, resp.Body) //nolint:errcheck
	if resp.StatusCode != http.StatusOK {
		return nil, fmt.Errorf("failed to authenticate (sc=%d)", resp.StatusCode)
	}

	return parseToken(resp.Body)
}

// FileToken retrieves a token to access a protected file.
func (c Client) FileToken(ctx context.Context, token string) (*Token, error) {
	const path = "/api/files/token"

	resp, err := httpw.Do(ctx, http.MethodPost, c.host+path,
		httpw.WithHeader("Authorization", token))
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()
	defer io.Copy(io.Discard, resp.Body) //nolint:errcheck

	if resp.StatusCode != http.StatusOK {
		return nil, fmt.Errorf("failed to authenticate (sc=%d)", resp.StatusCode)
	}

	return parseToken(resp.Body)
}

// RecordSearch searches for records in a collection. SearchResults struct is recomemnded to be used as output.
func (c Client) RecordSearch(ctx context.Context, params Params, output interface{}) error {
	url := fmt.Sprintf("%s/api/collections/%s/records?%s", c.host, params.Collection, params.QueryString())

	resp, err := httpw.Do(ctx, http.MethodGet, url,
		httpw.WithHeader("Authorization", params.Token))
	if err != nil {
		return err
	}
	defer resp.Body.Close()
	defer io.Copy(io.Discard, resp.Body) //nolint:errcheck

	if resp.StatusCode != http.StatusOK {
		b, _ := io.ReadAll(resp.Body)
		return fmt.Errorf("status code %d: %s", resp.StatusCode, b)
	}

	if err := unmarshalPB(resp.Body, output); err != nil {
		return err
	}

	return nil
}

// RecordView retrieves a record from a collection.
func (c Client) RecordView(ctx context.Context, params Params, output interface{}) error {
	url := fmt.Sprintf("%s/api/collections/%s/records/%s?%s", c.host, params.Collection, params.ID, params.QueryString())

	resp, err := httpw.Do(ctx, http.MethodGet, url,
		httpw.WithHeader("Authorization", params.Token))
	if err != nil {
		return err
	}
	defer resp.Body.Close()
	defer io.Copy(io.Discard, resp.Body) //nolint:errcheck

	if resp.StatusCode != http.StatusOK {
		b, _ := io.ReadAll(resp.Body)
		return fmt.Errorf("status code %d: %s", resp.StatusCode, b)
	}

	if err := unmarshalPB(resp.Body, output); err != nil {
		return err
	}

	return nil
}

// RecordCreate creates a record in a collection.
func (c Client) RecordCreate(ctx context.Context, params Params, output interface{}) error {
	url := fmt.Sprintf("%s/api/collections/%s/records?%s", c.host, params.Collection, params.QueryString())

	resp, err := httpw.Do(ctx, http.MethodPost, url,
		httpw.WithHeader("Authorization", params.Token),
		httpw.WithJsonBody(params.Data))
	if err != nil {
		return err
	}
	defer resp.Body.Close()
	defer io.Copy(io.Discard, resp.Body) //nolint:errcheck

	if resp.StatusCode != http.StatusOK {
		b, _ := io.ReadAll(resp.Body)
		return fmt.Errorf("status code %d: %s", resp.StatusCode, b)
	}

	if err := unmarshalPB(resp.Body, output); err != nil {
		return err
	}

	return nil
}

// RecordUpdate updates a record in a collection.
func (c Client) RecordUpdate(ctx context.Context, params Params, output interface{}) error {
	url := fmt.Sprintf("%s/api/collections/%s/records/%s?%s", c.host, params.Collection, params.ID, params.QueryString())

	resp, err := httpw.Do(ctx, http.MethodPatch, url,
		httpw.WithHeader("Authorization", params.Token),
		httpw.WithJsonBody(params.Data))
	if err != nil {
		return err
	}
	defer resp.Body.Close()
	defer io.Copy(io.Discard, resp.Body) //nolint:errcheck

	if resp.StatusCode != http.StatusOK {
		b, _ := io.ReadAll(resp.Body)
		return fmt.Errorf("status code %d: %s", resp.StatusCode, b)
	}

	if err := unmarshalPB(resp.Body, output); err != nil {
		return err
	}

	return nil
}

// RecordDelete deletes a record from a collection.
func (c Client) RecordDelete(ctx context.Context, params Params) error {
	url := fmt.Sprintf("%s/api/collections/%s/records/%s", c.host, params.Collection, params.ID)

	resp, err := httpw.Do(ctx, http.MethodDelete, url,
		httpw.WithHeader("Authorization", params.Token))
	if err != nil {
		return err
	}
	defer resp.Body.Close()
	defer io.Copy(io.Discard, resp.Body) //nolint:errcheck

	if resp.StatusCode != http.StatusOK {
		b, _ := io.ReadAll(resp.Body)
		return fmt.Errorf("status code %d: %s", resp.StatusCode, b)
	}

	return nil
}

// GetFileURL retrieves a file URL. Protected files require a token.
func (c Client) GetFileURL(ctx context.Context, params Params) string {
	url := fmt.Sprintf("%s/api/files/%s/%s/%s", c.host, params.Collection, params.ID, params.FileName)
	return url
}

// GetFileContent retrieves a file from a collection.
func (c Client) GetFileContent(ctx context.Context, params Params) ([]byte, error) {
	url := c.GetFileURL(ctx, params)
	resp, err := httpw.Do(ctx, http.MethodGet, url,
		httpw.WithParam("token", params.Token),
		httpw.WithParam("thumb", params.Thumb),
	)
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()
	b, err := io.ReadAll(resp.Body)
	if err != nil {
		return nil, err
	}

	if resp.StatusCode != http.StatusOK {
		b, _ := io.ReadAll(resp.Body)
		return nil, fmt.Errorf("status code %d: %s", resp.StatusCode, b)
	}

	return b, err
}

// unmarshalPB unmarshals output from a JSON byte slice.
func unmarshalPB(r io.Reader, output interface{}) error {
	if output == nil {
		return nil
	}

	if reflect.TypeOf(output).Kind() != reflect.Ptr {
		return errors.New("output must be a pointer")
	}

	if err := json.NewDecoder(r).Decode(output); err != nil {
		return err
	}
	return nil
}

// parseToken parses a token from a JSON response.
func parseToken(r io.Reader) (*Token, error) {
	t := Token{}
	if err := unmarshalPB(r, &t); err != nil {
		return nil, err
	}

	expiration, err := extractExpirationDate(t.Token)
	if err != nil {
		return nil, err
	}
	t.Expiration = expiration

	return &t, nil
}

// extractExpirationDate extracts the expiration date from a JWT token.
func extractExpirationDate(token string) (time.Time, error) {
	parts := strings.Split(token, ".")
	if len(parts) != 3 {
		return time.Time{}, errors.New("invalid token format")
	}
	payload, err := base64.RawURLEncoding.DecodeString(parts[1])
	if err != nil {
		return time.Time{}, fmt.Errorf("failed to decode token payload: %w", err)
	}
	var claims jwt.StandardClaims
	if err := json.Unmarshal(payload, &claims); err != nil {
		return time.Time{}, fmt.Errorf("failed to unmarshal token claims: %w", err)
	}
	return time.Unix(claims.ExpiresAt, 0), nil
}
