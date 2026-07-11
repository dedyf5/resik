-- +migrate Down
-- SQL to drop the initial database schema

-- Drop tables in reverse order of dependencies
DROP TABLE IF EXISTS `transactions`;
DROP TABLE IF EXISTS `outlets`;
DROP TABLE IF EXISTS `merchants`;
DROP TABLE IF EXISTS `tenant_member_resource_scopes`;
DROP TABLE IF EXISTS `tenant_member_permissions`;
DROP TABLE IF EXISTS `permissions`;
DROP TABLE IF EXISTS `tenant_members`;
DROP TABLE IF EXISTS `tenants`;
DROP TABLE IF EXISTS `users`;
