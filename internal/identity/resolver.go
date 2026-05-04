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

func (r *Resolver) Resolve(c context.Context, tableName string, publicID uuidPkg.UUIDV7) (uint64, error) {
	cacheKey := r.IDMapCacheKey(tableName, publicID.String32())
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

func (r *Resolver) ResolveBatch(c context.Context, tableName string, publicIDs []uuidPkg.UUIDV7) ([]uint64, error) {
	if len(publicIDs) == 0 {
		return []uint64{}, nil
	}

	ids := make([]uint64, 0, len(publicIDs))
	missingPublicIDs := make([]uuidPkg.UUIDV7, 0, len(publicIDs))

	keys := make([]string, len(publicIDs))
	for i, publicID := range publicIDs {
		keys[i] = r.IDMapCacheKey(tableName, publicID.String32())
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
		defer rows.Close()

		newCacheEntries := make(map[string]any)
		for rows.Next() {
			var id uint64
			var publicID uuidPkg.UUIDV7
			if err := rows.Scan(&id, &publicID); err == nil {
				cacheMap[publicID] = id
				newCacheEntries[r.IDMapCacheKey(tableName, publicID.String32())] = id
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

func (r *Resolver) IDMapCacheKey(tableName, publicID string) string {
	return fmt.Sprintf("%s:id_map:%s:%s", r.appKey, tableName, publicID)
}

func stringToUint64(s string) (uint64, error) {
	return strconv.ParseUint(s, 10, 64)
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
