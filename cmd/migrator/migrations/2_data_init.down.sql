-- +migrate Down
-- Remove sample data using TRUNCATE

SET FOREIGN_KEY_CHECKS = 0;

TRUNCATE TABLE `transactions`;
TRUNCATE TABLE `outlets`;
TRUNCATE TABLE `merchants`;
TRUNCATE TABLE `tenant_member_resource_scopes`;
TRUNCATE TABLE `tenant_member_permissions`;
TRUNCATE TABLE `permissions`;
TRUNCATE TABLE `tenant_members`;
TRUNCATE TABLE `tenants`;
TRUNCATE TABLE `users`;

SET FOREIGN_KEY_CHECKS = 1;
