-- +migrate Down
-- Remove sample data using TRUNCATE

SET FOREIGN_KEY_CHECKS = 0;

TRUNCATE TABLE `transaction`;
TRUNCATE TABLE `outlet`;
TRUNCATE TABLE `merchant`;
TRUNCATE TABLE `user`;

SET FOREIGN_KEY_CHECKS = 1;
