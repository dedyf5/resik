-- +migrate Up
-- SQL to create the initial database schema

-- Table structure for table `users`
CREATE TABLE IF NOT EXISTS `users` (
  `id` bigint(20) NOT NULL AUTO_INCREMENT,
  `public_id` UUID NOT NULL UNIQUE,
  `name` varchar(45) NOT NULL,
  `username` varchar(45) NOT NULL UNIQUE,
  `password` varchar(225) NOT NULL,
  `created_at` datetime NOT NULL,
  `created_by` bigint(20) NULL,
  `updated_at` datetime NOT NULL,
  `updated_by` bigint(20) NULL,
  PRIMARY KEY (`id`)
) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_general_ci;

-- Table structure for table `tenants`
CREATE TABLE IF NOT EXISTS `tenants` (
    id BIGINT(20) NOT NULL AUTO_INCREMENT,
    public_id UUID NOT NULL UNIQUE,
    name VARCHAR(45) NOT NULL,
    created_at DATETIME NOT NULL,
    created_by_id BIGINT(20) NOT NULL,
    updated_at DATETIME NOT NULL,
    updated_by_id BIGINT(20) NOT NULL,
    PRIMARY KEY (`id`),
    CONSTRAINT `fk_tenants_created_by` FOREIGN KEY (`created_by_id`) REFERENCES `users` (`id`),
    CONSTRAINT `fk_tenants_updated_by` FOREIGN KEY (`updated_by_id`) REFERENCES `users` (`id`)
) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_general_ci;

-- Table structure for table `tenant_members`
CREATE TABLE IF NOT EXISTS `tenant_members` (
    id BIGINT(20) NOT NULL AUTO_INCREMENT,
    tenant_id BIGINT(20) NOT NULL,
    user_id BIGINT(20) NOT NULL,
    created_at DATETIME NOT NULL,
    created_by_id BIGINT(20) NOT NULL,
    updated_at DATETIME NOT NULL,
    updated_by_id BIGINT(20) NOT NULL,
    PRIMARY KEY (`id`),
    UNIQUE KEY `uq_tenant_members` (`tenant_id`, `user_id`),
    CONSTRAINT `fk_tenant_members_tenant` FOREIGN KEY (`tenant_id`) REFERENCES `tenants`(`id`),
    CONSTRAINT `fk_tenant_members_user` FOREIGN KEY (`user_id`) REFERENCES `users`(`id`),
    CONSTRAINT `fk_tenant_members_created_by` FOREIGN KEY (`created_by_id`) REFERENCES `users`(`id`),
    CONSTRAINT `fk_tenant_members_updated_by` FOREIGN KEY (`updated_by_id`) REFERENCES `users`(`id`)
) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_general_ci;

-- Table structure for table `permissions`
CREATE TABLE IF NOT EXISTS `permissions` (
    id BIGINT(20) NOT NULL AUTO_INCREMENT,
    code VARCHAR(100) NOT NULL UNIQUE,
    name VARCHAR(100) NOT NULL,
    created_at DATETIME NOT NULL,
    created_by_id BIGINT(20) NOT NULL,
    updated_at DATETIME NOT NULL,
    updated_by_id BIGINT(20) NOT NULL,
    PRIMARY KEY (`id`),
    CONSTRAINT `fk_permissions_created_by` FOREIGN KEY (`created_by_id`) REFERENCES `users`(`id`),
    CONSTRAINT `fk_permissions_updated_by` FOREIGN KEY (`updated_by_id`) REFERENCES `users`(`id`)
) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_general_ci;

-- Table structure for table `tenant_member_permissions`
CREATE TABLE IF NOT EXISTS `tenant_member_permissions` (
    tenant_member_id BIGINT(20) NOT NULL,
    permission_id BIGINT(20) NOT NULL,
    PRIMARY KEY (`tenant_member_id`, `permission_id`),
    CONSTRAINT `fk_tmp_member` FOREIGN KEY (`tenant_member_id`) REFERENCES `tenant_members`(`id`),
    CONSTRAINT `fk_tmp_permission` FOREIGN KEY (`permission_id`) REFERENCES `permissions`(`id`)
) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_general_ci;

-- Table structure for table `tenant_member_resource_scopes`
CREATE TABLE IF NOT EXISTS `tenant_member_resource_scopes` (
    tenant_member_id BIGINT(20) NOT NULL,
    resource_code VARCHAR(75) NOT NULL,
    scope_by VARCHAR(45) NOT NULL,
    scope_ref VARCHAR(100) NOT NULL,
    PRIMARY KEY (`tenant_member_id`, `resource_code`, `scope_by`, `scope_ref`),
    INDEX `idx_tmrs_lookup` (`tenant_member_id`, `resource_code`),
    CONSTRAINT `fk_tmrs_member` FOREIGN KEY (`tenant_member_id`) REFERENCES `tenant_members`(`id`)
) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_general_ci;

-- Table structure for table `merchants`
CREATE TABLE IF NOT EXISTS `merchants` (
  `id` bigint(20) NOT NULL AUTO_INCREMENT,
  `public_id` UUID NOT NULL UNIQUE,
  `owner_id` bigint(20) NOT NULL,
  `name` varchar(40) NOT NULL,
  `description` text NULL DEFAULT NULL,
  `created_at` datetime NOT NULL,
  `created_by` bigint(20) NOT NULL,
  `updated_at` datetime NOT NULL,
  `updated_by` bigint(20) NOT NULL,
  PRIMARY KEY (`id`),
  KEY `fk_merchants_owners` (`owner_id`),
  CONSTRAINT `fk_merchants_owners` FOREIGN KEY (`owner_id`) REFERENCES `users` (`id`) ON DELETE NO ACTION ON UPDATE CASCADE
) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_general_ci;

-- Table structure for table `outlets`
CREATE TABLE IF NOT EXISTS `outlets` (
  `id` bigint(20) NOT NULL AUTO_INCREMENT,
  `public_id` UUID NOT NULL UNIQUE,
  `merchant_id` bigint(20) NOT NULL,
  `name` varchar(40) NOT NULL,
  `created_at` datetime NOT NULL,
  `created_by` bigint(20) NOT NULL,
  `updated_at` datetime NOT NULL,
  `updated_by` bigint(20) NOT NULL,
  PRIMARY KEY (`id`),
  KEY `fk_outlets_merchants` (`merchant_id`),
  CONSTRAINT `fk_outlets_merchants` FOREIGN KEY (`merchant_id`) REFERENCES `merchants` (`id`) ON DELETE NO ACTION ON UPDATE CASCADE
) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_general_ci;

-- Table structure for table `transactions`
CREATE TABLE IF NOT EXISTS `transactions` (
  `id` bigint(20) NOT NULL AUTO_INCREMENT,
  `public_id` UUID NOT NULL UNIQUE,
  `merchant_id` bigint(20) NOT NULL,
  `outlet_id` bigint(20) NOT NULL,
  `bill_total` double NOT NULL,
  `created_at` datetime NOT NULL,
  `created_by` bigint(20) NOT NULL,
  `updated_at` datetime NOT NULL,
  `updated_by` bigint(20) NOT NULL,
  PRIMARY KEY (`id`),
  KEY `fk_transactions_merchants` (`merchant_id`),
  KEY `fk_transactions_outlets` (`outlet_id`),
  CONSTRAINT `fk_transactions_merchants` FOREIGN KEY (`merchant_id`) REFERENCES `merchants` (`id`) ON DELETE NO ACTION ON UPDATE CASCADE,
  CONSTRAINT `fk_transactions_outlets` FOREIGN KEY (`outlet_id`) REFERENCES `outlets` (`id`) ON DELETE NO ACTION ON UPDATE CASCADE
) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_general_ci;
