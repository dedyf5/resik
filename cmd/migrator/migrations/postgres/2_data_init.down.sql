-- Resik
-- Author: Dedy Fajar Setyawan
-- See: https://github.com/dedyf5/resik

-- +migrate Down
-- Remove sample data using TRUNCATE

TRUNCATE TABLE 
    transactions, 
    branches, 
    organizations, 
    tenant_member_resource_scopes, 
    tenant_member_permissions, 
    permissions, 
    tenant_members, 
    tenants, 
    users 
RESTART IDENTITY CASCADE;
