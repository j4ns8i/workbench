package types

import (
	"crypto/rand"
	"time"

	"github.com/oklog/ulid/v2"
)

func NewULID() ulid.ULID {
	return ulid.MustNew(ulid.Timestamp(time.Now()), rand.Reader)
}

func NewULIDFromString(s string) (ulid.ULID, error) {
	u, err := ulid.ParseStrict(s)
	if err != nil {
		return ulid.ULID{}, err
	}
	return u, nil
}
