// Resik
// Author: Dedy Fajar Setyawan
// See: https://github.com/dedyf5/resik

package identity

import (
	"context"
	"testing"

	merchantEntity "github.com/dedyf5/resik/entities/merchant"
	outletEntity "github.com/dedyf5/resik/entities/outlet"
	uuidPkg "github.com/dedyf5/resik/pkg/uuid"
	"github.com/stretchr/testify/assert"
)

func TestResolverStructure(t *testing.T) {
	r := &Resolver{
		appKey: "resik",
	}

	t.Run("userAccessCacheKey", func(t *testing.T) {
		key := r.userAccessCacheKey(merchantEntity.TABLE_NAME, "merchant:read", 1)
		assert.Equal(t, "resik:user_access:merchants:merchant:read:1", key)

		outletKey := r.userAccessCacheKey(outletEntity.TABLE_NAME, "outlet:read", 1)
		assert.Equal(t, "resik:user_access:outlets:outlet:read:1", outletKey)
	})

	t.Run("idMapCacheKey", func(t *testing.T) {
		key := r.idMapCacheKey(merchantEntity.TABLE_NAME, "019deba4-2020-7dc7-a670-769321b06a9b")
		assert.Equal(t, "resik:id_map:merchants:019deba4-2020-7dc7-a670-769321b06a9b", key)
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
