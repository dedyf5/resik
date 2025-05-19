-- +migrate Down
-- SQL to drop the initial database schema

-- Drop tables in reverse order of dependencies
DROP TABLE IF EXISTS `transaction`;
DROP TABLE IF EXISTS `outlet`;
DROP TABLE IF EXISTS `merchant`;
DROP TABLE IF EXISTS `user`;
