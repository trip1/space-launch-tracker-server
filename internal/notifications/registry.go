package notifications

import (
	"context"
	"crypto/rand"
	"crypto/sha256"
	"crypto/subtle"
	"encoding/base64"
	"encoding/json"
	"errors"
	"fmt"
	"strings"
	"time"

	"ds9labs.com/space-launch-server/internal/storage"
)

const deviceKeyPrefix = "notifications:device:"

var ErrUnauthorizedDevice = errors.New("invalid device credential")

type Device struct {
	InstallationID string    `json:"installation_id"`
	Token          string    `json:"token"`
	SecretHash     string    `json:"secret_hash"`
	Enabled        bool      `json:"enabled"`
	AppVersion     string    `json:"app_version"`
	UpdatedAt      time.Time `json:"updated_at"`
}

type Registration struct {
	InstallationID string
	Token          string
	AppVersion     string
	Secret         string
}

type Registry struct {
	store storage.Store
	now   func() time.Time
}

func NewRegistry(store storage.Store) *Registry { return &Registry{store: store, now: time.Now} }

func (r *Registry) Register(ctx context.Context, input Registration) (string, error) {
	if strings.TrimSpace(input.InstallationID) == "" || len(input.InstallationID) > 128 || strings.TrimSpace(input.Token) == "" || len(input.Token) > 4096 {
		return "", errors.New("installation_id and token are required")
	}
	key := deviceKeyPrefix + input.InstallationID
	var existing Device
	payload, err := r.store.Get(ctx, key)
	if err == nil {
		if err := json.Unmarshal(payload, &existing); err != nil {
			return "", fmt.Errorf("decode registration: %w", err)
		}
		if !validSecret(existing.SecretHash, input.Secret) {
			return "", ErrUnauthorizedDevice
		}
	} else if !errors.Is(err, storage.ErrCacheMiss) {
		return "", err
	}
	secret := input.Secret
	if existing.InstallationID == "" {
		secret, err = newSecret()
		if err != nil {
			return "", err
		}
	}
	device := Device{InstallationID: input.InstallationID, Token: input.Token, SecretHash: hashSecret(secret), Enabled: true, AppVersion: input.AppVersion, UpdatedAt: r.now().UTC()}
	encoded, err := json.Marshal(device)
	if err != nil {
		return "", err
	}
	if err := r.store.Set(ctx, key, encoded, 180*24*time.Hour); err != nil {
		return "", err
	}
	return secret, nil
}

func (r *Registry) Devices(ctx context.Context) ([]Device, error) {
	keys, err := r.store.Keys(ctx, deviceKeyPrefix+"*")
	if err != nil {
		return nil, err
	}
	out := make([]Device, 0, len(keys))
	for _, key := range keys {
		payload, err := r.store.Get(ctx, key)
		if errors.Is(err, storage.ErrCacheMiss) {
			continue
		}
		if err != nil {
			return nil, err
		}
		var d Device
		if json.Unmarshal(payload, &d) == nil && d.Enabled && d.Token != "" {
			out = append(out, d)
		}
	}
	return out, nil
}

func (r *Registry) RemoveToken(ctx context.Context, installationID string) error {
	return r.store.Delete(ctx, deviceKeyPrefix+installationID)
}
func newSecret() (string, error) {
	b := make([]byte, 32)
	if _, err := rand.Read(b); err != nil {
		return "", err
	}
	return base64.RawURLEncoding.EncodeToString(b), nil
}
func hashSecret(secret string) string {
	sum := sha256.Sum256([]byte(secret))
	return base64.RawStdEncoding.EncodeToString(sum[:])
}
func validSecret(hash, secret string) bool {
	if hash == "" || secret == "" {
		return false
	}
	return subtle.ConstantTimeCompare([]byte(hash), []byte(hashSecret(secret))) == 1
}
