// Resik
// Author: Dedy Fajar Setyawan
// See: https://github.com/dedyf5/resik

package identity

import (
	"context"
	"database/sql"
	"fmt"
	"log"
	"slices"
	"strconv"
	"strings"
	"time"

	branchEntity "github.com/dedyf5/resik/entities/branch"
	organizationEntity "github.com/dedyf5/resik/entities/organization"
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

// GetTenantOrganizationIDs returns all organization IDs belonging to tenants of the user.
func (r *Resolver) GetTenantOrganizationIDs(c context.Context, userID uint64) ([]uint64, error) {
	cacheKey := fmt.Sprintf("%s:user_tenant_organizations:%d", r.appKey, userID)

	cachedIDs, err := r.getSMembers(c, cacheKey)
	if err == nil && len(cachedIDs) > 0 {
		return cachedIDs, nil
	}

	query := `
		SELECT DISTINCT o.id
		FROM tenant_members tm
		JOIN tenants t ON t.id = tm.tenant_id
		JOIN organizations o ON o.tenant_public_id = t.public_id
		WHERE tm.user_id = @userID
	`

	var ids []uint64
	err = r.db.WithContext(c).Raw(query, sql.Named("userID", userID)).Scan(&ids).Error
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

// GetTenantBranchIDs returns all branch IDs belonging to organizations of tenants of the user.
func (r *Resolver) GetTenantBranchIDs(c context.Context, userID uint64) ([]uint64, error) {
	cacheKey := fmt.Sprintf("%s:user_tenant_branches:%d", r.appKey, userID)

	cachedIDs, err := r.getSMembers(c, cacheKey)
	if err == nil && len(cachedIDs) > 0 {
		return cachedIDs, nil
	}

	query := `
		SELECT DISTINCT b.id
		FROM tenant_members tm
		JOIN tenants t ON t.id = tm.tenant_id
		JOIN organizations o ON o.tenant_public_id = t.public_id
		JOIN branches b ON b.organization_id = o.id
		WHERE tm.user_id = @userID
	`

	var ids []uint64
	err = r.db.WithContext(c).Raw(query, sql.Named("userID", userID)).Scan(&ids).Error
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

// GetResourceIDs returns all resource IDs accessible to a user for a specific table and permission code.
func (r *Resolver) GetResourceIDs(c context.Context, userID uint64, resourceTable string, permissionCode string) ([]uint64, error) {
	switch resourceTable {
	case organizationEntity.TABLE_NAME:
		return r.GetOrganizationIDsByPermission(c, userID, permissionCode)
	case branchEntity.TABLE_NAME:
		return r.GetBranchIDsByPermission(c, userID, permissionCode)
	default:
		return []uint64{}, nil
	}
}

// HasAccessByID checks if a user has access to a specific resource by its internal uint64 ID.
func (r *Resolver) HasAccessByID(c context.Context, userID uint64, permissionCode string, resourceTable string, resourceID uint64) (bool, error) {
	ids, err := r.GetResourceIDs(c, userID, resourceTable, permissionCode)
	if err != nil {
		return false, err
	}
	return slices.Contains(ids, resourceID), nil
}

// HasAccessByPublicID checks if a user has access to a specific resource by its public UUID.
func (r *Resolver) HasAccessByPublicID(c context.Context, userID uint64, permissionCode string, resourceTable string, publicID uuidPkg.UUIDV7) (bool, error) {
	id, err := r.Resolve(c, resourceTable, publicID)
	if err != nil {
		return false, err
	}
	return r.HasAccessByID(c, userID, permissionCode, resourceTable, id)
}

// GetOrganizationIDsByPermission returns all cached organization IDs for a specific user and permission code.
func (r *Resolver) GetOrganizationIDsByPermission(c context.Context, userID uint64, permissionCode string) ([]uint64, error) {
	cacheKey := r.userAccessCacheKey(organizationEntity.TABLE_NAME, permissionCode, userID)

	cachedIDs, err := r.getSMembers(c, cacheKey)
	if err == nil && len(cachedIDs) > 0 {
		return cachedIDs, nil
	}

	resourceDomain := strings.Split(permissionCode, ":")[0]

	query := `
		SELECT DISTINCT o.id
		FROM tenant_members tm
		JOIN tenants t ON t.id = tm.tenant_id
		JOIN tenant_member_permissions tmp ON tmp.tenant_member_id = tm.id
		JOIN permissions p ON p.id = tmp.permission_id AND p.code = @permissionCode
		JOIN organizations o ON o.tenant_public_id = t.public_id
		WHERE tm.user_id = @userID
		  AND NOT EXISTS (
		    SELECT 1 FROM tenant_member_resource_scopes tmrs
		    WHERE tmrs.tenant_member_id = tm.id
		      AND tmrs.resource_code IN (@resourceDomain, 'organization')
		  )
		UNION
		SELECT DISTINCT o.id
		FROM tenant_members tm
		JOIN tenants t ON t.id = tm.tenant_id
		JOIN tenant_member_permissions tmp ON tmp.tenant_member_id = tm.id
		JOIN permissions p ON p.id = tmp.permission_id AND p.code = @permissionCode
		JOIN tenant_member_resource_scopes tmrs ON tmrs.tenant_member_id = tm.id 
		  AND tmrs.resource_code IN (@resourceDomain, 'organization') 
		  AND tmrs.scope_by = 'organization'
		JOIN organizations o ON o.tenant_public_id = t.public_id AND o.public_id = tmrs.scope_ref
		WHERE tm.user_id = @userID
	`

	var ids []uint64
	err = r.db.WithContext(c).Raw(query,
		sql.Named("permissionCode", permissionCode),
		sql.Named("userID", userID),
		sql.Named("resourceDomain", resourceDomain),
	).Scan(&ids).Error
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

// GetBranchIDsByPermission returns all cached branch IDs for a specific user and permission code.
func (r *Resolver) GetBranchIDsByPermission(c context.Context, userID uint64, permissionCode string) ([]uint64, error) {
	cacheKey := r.userAccessCacheKey(branchEntity.TABLE_NAME, permissionCode, userID)

	cachedIDs, err := r.getSMembers(c, cacheKey)
	if err == nil && len(cachedIDs) > 0 {
		return cachedIDs, nil
	}

	resourceDomain := strings.Split(permissionCode, ":")[0]

	query := `
		SELECT DISTINCT b.id
		FROM tenant_members tm
		JOIN tenants t ON t.id = tm.tenant_id
		JOIN tenant_member_permissions tmp ON tmp.tenant_member_id = tm.id
		JOIN permissions p ON p.id = tmp.permission_id AND p.code = @permissionCode
		JOIN organizations o ON o.tenant_public_id = t.public_id
		JOIN branches b ON b.organization_id = o.id
		WHERE tm.user_id = @userID
		  AND NOT EXISTS (
		    SELECT 1 FROM tenant_member_resource_scopes tmrs
		    WHERE tmrs.tenant_member_id = tm.id
		      AND tmrs.resource_code IN (@resourceDomain, 'branch', 'organization')
		  )
		UNION
		SELECT DISTINCT b.id
		FROM tenant_members tm
		JOIN tenants t ON t.id = tm.tenant_id
		JOIN tenant_member_permissions tmp ON tmp.tenant_member_id = tm.id
		JOIN permissions p ON p.id = tmp.permission_id AND p.code = @permissionCode
		JOIN tenant_member_resource_scopes tmrs ON tmrs.tenant_member_id = tm.id 
		  AND tmrs.resource_code IN (@resourceDomain, 'branch', 'organization')
		  AND tmrs.scope_by = 'organization'
		JOIN organizations o ON o.tenant_public_id = t.public_id AND o.public_id = tmrs.scope_ref
		JOIN branches b ON b.organization_id = o.id
		WHERE tm.user_id = @userID
		UNION
		SELECT DISTINCT b.id
		FROM tenant_members tm
		JOIN tenants t ON t.id = tm.tenant_id
		JOIN tenant_member_permissions tmp ON tmp.tenant_member_id = tm.id
		JOIN permissions p ON p.id = tmp.permission_id AND p.code = @permissionCode
		JOIN tenant_member_resource_scopes tmrs ON tmrs.tenant_member_id = tm.id 
		  AND tmrs.resource_code IN (@resourceDomain, 'branch', 'organization')
		  AND tmrs.scope_by = 'branch'
		JOIN organizations o ON o.tenant_public_id = t.public_id
		JOIN branches b ON b.organization_id = o.id AND b.public_id = tmrs.scope_ref
		WHERE tm.user_id = @userID
	`

	var ids []uint64
	err = r.db.WithContext(c).Raw(query,
		sql.Named("permissionCode", permissionCode),
		sql.Named("userID", userID),
		sql.Named("resourceDomain", resourceDomain),
	).Scan(&ids).Error
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

// InvalidateUserAccessOrganization invalidates all user access organization caches for a specific user.
func (r *Resolver) InvalidateUserAccessOrganization(c context.Context, userID uint64) error {
	pattern := fmt.Sprintf("%s:user_access:%s:*:%d", r.appKey, organizationEntity.TABLE_NAME, userID)
	return r.deleteKeysByPattern(c, pattern)
}

// InvalidateUserAccessBranch invalidates all user access branch caches for a specific user.
func (r *Resolver) InvalidateUserAccessBranch(c context.Context, userID uint64) error {
	pattern := fmt.Sprintf("%s:user_access:%s:*:%d", r.appKey, branchEntity.TABLE_NAME, userID)
	return r.deleteKeysByPattern(c, pattern)
}

func (r *Resolver) deleteKeysByPattern(c context.Context, pattern string) error {
	iter := r.cache.Scan(c, 0, pattern, 0).Iterator()
	var keys []string
	for iter.Next(c) {
		keys = append(keys, iter.Val())
	}
	if err := iter.Err(); err != nil {
		return err
	}
	if len(keys) > 0 {
		return r.cache.Del(c, keys...).Err()
	}
	return nil
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

// userAccessCacheKey returns the cache key for a specific user access to a table and permission code.
func (r *Resolver) userAccessCacheKey(tableName string, permissionCode string, userID uint64) string {
	return fmt.Sprintf("%s:user_access:%s:%s:%d", r.appKey, tableName, permissionCode, userID)
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
