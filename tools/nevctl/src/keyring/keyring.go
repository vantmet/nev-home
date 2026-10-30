package keyring

import (
	"errors"
	"fmt"
	"sort"

	"github.com/godbus/dbus/v5"
	zkr "github.com/zalando/go-keyring"
)

var (
	ErrNotFound = errors.New("secret not found in keyring")
)

// Keyring provides an interface for interacting with system secret storage.
type Keyring interface {
	Get(service, key string) (string, error)
	Set(service, key, value string) error
	Delete(service, key string) error
	List(service string) ([]string, error)
}

// SystemKeyring interacts with the OS keyring using go-keyring and DBus on Linux.
type SystemKeyring struct{}

// NewSystemKeyring returns a new SystemKeyring instance.
func NewSystemKeyring() *SystemKeyring {
	return &SystemKeyring{}
}

// Get retrieves a secret by service and key.
func (k *SystemKeyring) Get(service, key string) (string, error) {
	val, err := zkr.Get(service, key)
	if err != nil {
		if errors.Is(err, zkr.ErrNotFound) {
			return "", fmt.Errorf("%w: %s (service: %s)", ErrNotFound, key, service)
		}
		return "", fmt.Errorf("getting secret %s from keyring: %w", key, err)
	}
	return val, nil
}

// Set stores a secret by service and key.
func (k *SystemKeyring) Set(service, key, value string) error {
	if err := zkr.Set(service, key, value); err != nil {
		return fmt.Errorf("saving secret %s to keyring: %w", key, err)
	}
	return nil
}

// Delete removes a secret by service and key.
func (k *SystemKeyring) Delete(service, key string) error {
	if err := zkr.Delete(service, key); err != nil {
		if errors.Is(err, zkr.ErrNotFound) {
			return fmt.Errorf("%w: %s (service: %s)", ErrNotFound, key, service)
		}
		return fmt.Errorf("deleting secret %s from keyring: %w", key, err)
	}
	return nil
}

// List searches the Secret Service on Linux for all keys belonging to the given service.
func (k *SystemKeyring) List(service string) ([]string, error) {
	conn, err := dbus.SessionBus()
	if err != nil {
		return nil, fmt.Errorf("connecting to dbus session: %w", err)
	}

	obj := conn.Object("org.freedesktop.secrets", "/org/freedesktop/secrets")
	var unlocked, locked []dbus.ObjectPath
	call := obj.Call(
		"org.freedesktop.Secret.Service.SearchItems",
		0,
		map[string]string{"service": service},
	)
	if call.Err != nil {
		return nil, fmt.Errorf("searching secret service items: %w", call.Err)
	}

	if err := call.Store(&unlocked, &locked); err != nil {
		return nil, fmt.Errorf("parsing search results: %w", err)
	}

	allPaths := append(unlocked, locked...)
	keySet := make(map[string]struct{})

	for _, p := range allPaths {
		item := conn.Object("org.freedesktop.secrets", p)
		variant, err := item.GetProperty("org.freedesktop.Secret.Item.Attributes")
		if err != nil {
			continue
		}

		attrs, ok := variant.Value().(map[string]string)
		if !ok {
			continue
		}

		if username, exists := attrs["username"]; exists && username != "" {
			keySet[username] = struct{}{}
		}
	}

	keys := make([]string, 0, len(keySet))
	for k := range keySet {
		keys = append(keys, k)
	}
	sort.Strings(keys)

	return keys, nil
}
