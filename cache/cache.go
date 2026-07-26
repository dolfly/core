// Package cache defines a generic key-value cache interface used across the
// GOST project — HTTP response caching, and (as consumers migrate) DNS
// responses, routes, and other []byte-serializable data.
package cache

import (
	"context"
	"errors"
	"time"
)

// ErrNotFound is returned by Cache.Get when the key does not exist.
// An expired-but-present entry is NOT a miss: it is returned with a non-nil
// Entry (Expired() == true) so callers can implement serve-stale.
var ErrNotFound = errors.New("cache: key not found")

// Entry is a single cached value with an expiration time.
type Entry struct {
	// Data is the cached payload. Callers serialize/deserialize their own types.
	Data []byte
	// Expiration is the absolute expiry time. A zero value means the entry
	// never expires.
	Expiration time.Time
}

// Expired reports whether the entry has passed its expiration time.
// A zero Expiration (never expires) is never expired.
func (e *Entry) Expired() bool {
	if e == nil {
		return true
	}
	if e.Expiration.IsZero() {
		return false
	}
	return time.Now().After(e.Expiration)
}

// TTL returns the entry's remaining time to live. It returns -1 for an entry
// that never expires, and a non-positive duration for an expired entry.
func (e *Entry) TTL() time.Duration {
	if e == nil {
		return 0
	}
	if e.Expiration.IsZero() {
		return -1
	}
	return time.Until(e.Expiration)
}

// SetOptions holds parameters for a Set call.
type SetOptions struct {
	// TTL is the time to live for the entry. A non-positive TTL uses the
	// cache's default TTL.
	TTL time.Duration
}

// SetOption is a functional option for configuring SetOptions.
type SetOption func(opts *SetOptions)

// WithTTL sets the time to live for a cache entry, overriding the cache default.
func WithTTL(d time.Duration) SetOption {
	return func(opts *SetOptions) {
		opts.TTL = d
	}
}

// Cache is a concurrency-safe []byte key-value store. Implementations must be
// safe for concurrent use.
//
// Get returns ErrNotFound only when the key is absent; an expired entry is
// returned (with Expired() == true) to support serve-stale.
type Cache interface {
	// Get returns the entry for key, or ErrNotFound if the key is absent.
	Get(ctx context.Context, key string) (*Entry, error)
	// Set stores data under key with the given options.
	Set(ctx context.Context, key string, data []byte, opts ...SetOption) error
	// Delete removes key from the cache. Deleting an absent key is a no-op.
	Delete(ctx context.Context, key string) error
}
