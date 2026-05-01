-- +migrate Down
-- Remove sample data using TRUNCATE

SET FOREIGN_KEY_CHECKS = 0;

TRUNCATE TABLE `transactions`;
TRUNCATE TABLE `outlets`;
TRUNCATE TABLE `merchants`;
TRUNCATE TABLE `users`;

SET FOREIGN_KEY_CHECKS = 1;
