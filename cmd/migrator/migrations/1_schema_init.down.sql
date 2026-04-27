-- +migrate Down
-- SQL to drop the initial database schema

-- Drop tables in reverse order of dependencies
DROP TABLE IF EXISTS `transactions`;
DROP TABLE IF EXISTS `outlets`;
DROP TABLE IF EXISTS `merchants`;
DROP TABLE IF EXISTS `users`;
