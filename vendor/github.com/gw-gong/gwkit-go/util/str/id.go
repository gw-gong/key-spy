package str

import (
	"math/rand"
	"time"

	"github.com/google/uuid"
	"github.com/oklog/ulid/v2"
)

func GenerateUUID() string {
	return uuid.New().String()
}

func GenerateULID() string {
	t := time.Now().UTC()
	entropy := ulid.Monotonic(rand.New(rand.NewSource(t.UnixNano())), 0)
	return ulid.MustNew(ulid.Timestamp(t), entropy).String()
} 