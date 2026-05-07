// Resik
// Author: Dedy Fajar Setyawan
// See: https://github.com/dedyf5/resik

package identity

import (
	"context"
	"fmt"
	"log"
	"strconv"
	"time"

	merchantEntity "github.com/dedyf5/resik/entities/merchant"
	outletEntity "github.com/dedyf5/resik/entities/outlet"
	"github.com/dedyf5/resik/pkg/collection"
	"github.com/dedyf5/resik/pkg/numbers"
	uuidPkg "github.com/dedyf5/resik/pkg/uuid"
	goredis "github.com/redis/go-redis/v9"
	"gorm.io/gorm"
)

const (
	cacheExpiration = 24 * time.Hour
)

type Resolver struct {
	appKey string
	db     *gorm.DB
	cache  *goredis.Client
}

func NewResolver(appKey string, db *gorm.DB, cache *goredis.Client) IdentityResolver {
	return &Resolver{
		appKey: appKey,
		db:     db,
		cache:  cache,
	}
}

// Resolve resolves a public ID to its corresponding ID in a table.
func (r *Resolver) Resolve(c context.Context, tableName string, publicID uuidPkg.UUIDV7) (uint64, error) {
	cacheKey := r.idMapCacheKey(tableName, publicID.String32())
	if val, err := r.cache.Get(c, cacheKey).Result(); err == nil && val != "" {
		return stringToUint64(val)
	}

	var id uint64
	err := r.db.WithContext(c).
		Table(tableName).
		Where("public_id = ?", publicID.String()).
		Select("id").
		Row().
		Scan(&id)
	if err != nil {
		return 0, err
	}

	r.cache.Set(c, cacheKey, id, cacheExpiration)

	return id, nil
}

// ResolveBatch resolves a batch of public IDs to their corresponding IDs in a table.
func (r *Resolver) ResolveBatch(c context.Context, tableName string, publicIDs []uuidPkg.UUIDV7) ([]uint64, error) {
	if len(publicIDs) == 0 {
		return []uint64{}, nil
	}

	ids := make([]uint64, 0, len(publicIDs))
	missingPublicIDs := make([]uuidPkg.UUIDV7, 0, len(publicIDs))

	keys := make([]string, len(publicIDs))
	for i, publicID := range publicIDs {
		keys[i] = r.idMapCacheKey(tableName, publicID.String32())
	}

	cacheMap := make(map[uuidPkg.UUIDV7]uint64)

	cacheValues, err := r.cache.MGet(c, keys...).Result()
	if err == nil {
		for i, val := range cacheValues {
			if val != nil {
				if id, ok := anyToUint64(val); ok {
					cacheMap[publicIDs[i]] = id
				} else {
					missingPublicIDs = append(missingPublicIDs, publicIDs[i])
				}
			} else {
				missingPublicIDs = append(missingPublicIDs, publicIDs[i])
			}
		}
	} else {
		missingPublicIDs = append(missingPublicIDs, publicIDs...)
	}

	nm := len(missingPublicIDs)

	if nm > 0 {
		missingPublicIDsStr := collection.Map(missingPublicIDs, func(n uuidPkg.UUIDV7) string {
			return n.String()
		})

		query := r.db.WithContext(c).Table(tableName).Select("id", "public_id")
		if nm == 1 {
			query = query.Where("public_id = ?", missingPublicIDsStr[0])
		} else {
			query = query.Where("public_id IN ?", missingPublicIDsStr)
		}

		rows, err := query.Rows()
		if err != nil {
			return nil, err
		}
		defer func() {
			err := rows.Close()
			if err != nil {
				log.Printf("failed to close rows: %v", err)
			}
		}()

		newCacheEntries := make(map[string]any)
		for rows.Next() {
			var id uint64
			var publicID uuidPkg.UUIDV7
			if err := rows.Scan(&id, &publicID); err == nil {
				cacheMap[publicID] = id
				newCacheEntries[r.idMapCacheKey(tableName, publicID.String32())] = id
			}
		}

		if len(newCacheEntries) > 0 {
			pipe := r.cache.Pipeline()
			for key, value := range newCacheEntries {
				pipe.Set(c, key, value, cacheExpiration)
			}

			if _, err := pipe.Exec(c); err != nil {
				log.Printf("failed to set cache: %v", err)
			}
		}
	}

	for _, pid := range publicIDs {
		if id, ok := cacheMap[pid]; ok {
			ids = append(ids, id)
		} else {
			return nil, fmt.Errorf("public_id for %s not found in table %s", pid, tableName)
		}
	}

	return ids, nil
}

// GetMerchantIDs returns all cached merchant IDs for a specific user access to a table.
func (r *Resolver) GetMerchantIDs(c context.Context, userID uint64) ([]uint64, error) {
	cacheKey := r.userAccessCacheKey(merchantEntity.TABLE_NAME, userID)

	cachedIDs, err := r.getSMembers(c, cacheKey)
	if err == nil && len(cachedIDs) > 0 {
		return cachedIDs, nil
	}

	var ids []uint64
	err = r.db.WithContext(c).
		Table(merchantEntity.TABLE_NAME).
		Where("owner_id = ?", userID).
		Pluck("id", &ids).
		Error

	if err != nil {
		return nil, err
	}

	if len(ids) > 0 {
		if err := r.setSMembers(c, cacheKey, ids); err != nil {
			log.Printf("failed to set cache: %v", err)
		}
	}

	return ids, nil
}

// GetOutletIDs returns all cached outlet IDs for a specific user access to a table.
func (r *Resolver) GetOutletIDs(c context.Context, userID uint64) ([]uint64, error) {
	cacheKey := r.userAccessCacheKey(outletEntity.TABLE_NAME, userID)

	cachedIDs, err := r.getSMembers(c, cacheKey)
	if err == nil && len(cachedIDs) > 0 {
		return cachedIDs, nil
	}

	var ids []uint64
	err = r.db.WithContext(c).
		Table(outletEntity.TABLE_NAME+" AS o1").
		Joins("INNER JOIN "+merchantEntity.TABLE_NAME+" AS m1 ON m1.id = o1.merchant_id").
		Where("m1.owner_id = ?", userID).
		Pluck("o1.id", &ids).
		Error

	if err != nil {
		return nil, err
	}

	if len(ids) > 0 {
		if err := r.setSMembers(c, cacheKey, ids); err != nil {
			log.Printf("failed to set cache: %v", err)
		}
	}

	return ids, nil
}

// InvalidateUserAccessMerchant invalidates the cache for a specific user and merchant.
// This is used when a merchant is deleted or the user's access to the merchant is revoked.
func (r *Resolver) InvalidateUserAccessMerchant(c context.Context, userID uint64) error {
	cacheKey := r.userAccessCacheKey(merchantEntity.TABLE_NAME, userID)
	return r.cache.Del(c, cacheKey).Err()
}

// InvalidateUserAccessOutlet invalidates the cache for a specific user and outlet.
// This is used when an outlet is deleted or the user's access to the outlet is revoked.
func (r *Resolver) InvalidateUserAccessOutlet(c context.Context, userID uint64) error {
	cacheKey := r.userAccessCacheKey(outletEntity.TABLE_NAME, userID)
	return r.cache.Del(c, cacheKey).Err()
}

// getSMembers returns all cached IDs for a specific key.
func (r *Resolver) getSMembers(c context.Context, key string) ([]uint64, error) {
	cachedIDs, err := r.cache.SMembers(c, key).Result()
	if err == nil && len(cachedIDs) > 0 {
		ids, err := stringToUint64Slice(cachedIDs)
		if len(ids) > 0 && err == nil {
			return ids, nil
		}
	}

	return nil, nil
}

// setSMembers adds all IDs to the cache for a specific key.
func (r *Resolver) setSMembers(c context.Context, key string, ids []uint64) error {
	if len(ids) == 0 {
		return nil
	}

	pipe := r.cache.Pipeline()

	interfaces := make([]any, len(ids))
	for i, id := range ids {
		interfaces[i] = id
	}

	pipe.SAdd(c, key, interfaces...)
	pipe.Expire(c, key, cacheExpiration)

	if _, err := pipe.Exec(c); err != nil {
		return err
	}

	return nil
}

// userAccessCacheKey returns the cache key for a specific user access to a table.
func (r *Resolver) userAccessCacheKey(tableName string, userID uint64) string {
	return fmt.Sprintf("%s:user_access:%s:%d", r.appKey, tableName, userID)
}

// idMapCacheKey returns the cache key for a specific table and public ID.
func (r *Resolver) idMapCacheKey(tableName, publicID string) string {
	return fmt.Sprintf("%s:id_map:%s:%s", r.appKey, tableName, publicID)
}

func stringToUint64(s string) (uint64, error) {
	return strconv.ParseUint(s, 10, 64)
}

func stringToUint64Slice(s []string) ([]uint64, error) {
	if len(s) == 0 {
		return []uint64{}, nil
	}

	ids := make([]uint64, len(s))
	for i, v := range s {
		if id, err := stringToUint64(v); err != nil {
			return nil, err
		} else {
			ids[i] = id
		}
	}

	return ids, nil
}

func anyToUint64(val any) (uint64, bool) {
	switch v := val.(type) {
	case string:
		if res, err := stringToUint64(v); err == nil {
			return res, true
		} else {
			return 0, false
		}
	case uint64:
		return v, true
	case int64:
		if res, err := numbers.SafeConvert[uint64](v); err == nil {
			return res, true
		} else {
			return 0, false
		}
	case int:
		if res, err := numbers.SafeConvert[uint64](v); err == nil {
			return res, true
		} else {
			return 0, false
		}
	case float64:
		return uint64(v), true
	default:
		return 0, false
	}
}
