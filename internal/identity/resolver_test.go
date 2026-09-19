// Resik
// Author: Dedy Fajar Setyawan
// See: https://github.com/dedyf5/resik

package identity

import (
	"context"
	"testing"

	branchEntity "github.com/dedyf5/resik/entities/branch"
	organizationEntity "github.com/dedyf5/resik/entities/organization"
	uuidPkg "github.com/dedyf5/resik/pkg/uuid"
	"github.com/stretchr/testify/assert"
)

func TestResolverStructure(t *testing.T) {
	r := &Resolver{
		appKey: "resik",
	}

	t.Run("userAccessCacheKey", func(t *testing.T) {
		key := r.userAccessCacheKey(organizationEntity.TABLE_NAME, "organization:read", 1)
		assert.Equal(t, "resik:user_access:organizations:organization:read:1", key)

		branchKey := r.userAccessCacheKey(branchEntity.TABLE_NAME, "branch:read", 1)
		assert.Equal(t, "resik:user_access:branches:branch:read:1", branchKey)
	})

	t.Run("idMapCacheKey", func(t *testing.T) {
		key := r.idMapCacheKey(organizationEntity.TABLE_NAME, "019deba4-2020-7dc7-a670-769321b06a9b")
		assert.Equal(t, "resik:id_map:organizations:019deba4-2020-7dc7-a670-769321b06a9b", key)
	})

	t.Run("GetResourceIDs-UnknownTable", func(t *testing.T) {
		ctx := context.Background()
		ids, err := r.GetResourceIDs(ctx, 1, "unknown_table", "read")
		assert.Nil(t, err)
		assert.Empty(t, ids)
	})
}

func TestStringToUint64(t *testing.T) {
	t.Run("ValidString", func(t *testing.T) {
		val, err := stringToUint64("123")
		assert.Nil(t, err)
		assert.Equal(t, uint64(123), val)
	})

	t.Run("InvalidString", func(t *testing.T) {
		_, err := stringToUint64("abc")
		assert.NotNil(t, err)
	})
}

func TestStringToUint64Slice(t *testing.T) {
	t.Run("Empty", func(t *testing.T) {
		ids, err := stringToUint64Slice([]string{})
		assert.Nil(t, err)
		assert.Empty(t, ids)
	})

	t.Run("ValidSlice", func(t *testing.T) {
		ids, err := stringToUint64Slice([]string{"1", "2", "3"})
		assert.Nil(t, err)
		assert.Equal(t, []uint64{1, 2, 3}, ids)
	})

	t.Run("InvalidSlice", func(t *testing.T) {
		_, err := stringToUint64Slice([]string{"1", "invalid"})
		assert.NotNil(t, err)
	})
}

func TestAnyToUint64(t *testing.T) {
	t.Run("String", func(t *testing.T) {
		v, ok := anyToUint64("42")
		assert.True(t, ok)
		assert.Equal(t, uint64(42), v)
	})

	t.Run("Uint64", func(t *testing.T) {
		v, ok := anyToUint64(uint64(42))
		assert.True(t, ok)
		assert.Equal(t, uint64(42), v)
	})

	t.Run("Int64", func(t *testing.T) {
		v, ok := anyToUint64(int64(42))
		assert.True(t, ok)
		assert.Equal(t, uint64(42), v)
	})

	t.Run("Int", func(t *testing.T) {
		v, ok := anyToUint64(int(42))
		assert.True(t, ok)
		assert.Equal(t, uint64(42), v)
	})

	t.Run("Float64", func(t *testing.T) {
		v, ok := anyToUint64(float64(42))
		assert.True(t, ok)
		assert.Equal(t, uint64(42), v)
	})

	t.Run("Unsupported", func(t *testing.T) {
		_, ok := anyToUint64(uuidPkg.UUIDV7{})
		assert.False(t, ok)
	})
}
