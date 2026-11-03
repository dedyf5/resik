-- Resik
-- Author: Dedy Fajar Setyawan
-- See: https://github.com/dedyf5/resik

-- +migrate Up
-- SQL to create the initial database schema

-- Table structure for table `users`
CREATE TABLE IF NOT EXISTS users (
    id BIGSERIAL NOT NULL,
    public_id UUID NOT NULL UNIQUE,
    full_name VARCHAR(45) NOT NULL,
    username VARCHAR(45) NOT NULL UNIQUE,
    password VARCHAR(225) NOT NULL,
    created_at TIMESTAMPTZ NOT NULL,
    created_by_id BIGINT NULL,
    updated_at TIMESTAMPTZ NOT NULL,
    updated_by_id BIGINT NULL,
    PRIMARY KEY (id),
    CONSTRAINT fk_users_created_by FOREIGN KEY (created_by_id) REFERENCES users (id),
    CONSTRAINT fk_users_updated_by FOREIGN KEY (updated_by_id) REFERENCES users (id)
);

-- Indexes for users
CREATE INDEX idx_users_created_by_id ON users (created_by_id);
CREATE INDEX idx_users_updated_by_id ON users (updated_by_id);

-- Table structure for table `tenants`
CREATE TABLE IF NOT EXISTS tenants (
    id BIGSERIAL NOT NULL,
    public_id UUID NOT NULL UNIQUE,
    name VARCHAR(45) NOT NULL,
    created_at TIMESTAMPTZ NOT NULL,
    created_by_id BIGINT NOT NULL,
    updated_at TIMESTAMPTZ NOT NULL,
    updated_by_id BIGINT NOT NULL,
    PRIMARY KEY (id),
    CONSTRAINT fk_tenants_created_by FOREIGN KEY (created_by_id) REFERENCES users (id),
    CONSTRAINT fk_tenants_updated_by FOREIGN KEY (updated_by_id) REFERENCES users (id)
);

-- Indexes for tenants
CREATE INDEX idx_tenants_created_by_id ON tenants (created_by_id);
CREATE INDEX idx_tenants_updated_by_id ON tenants (updated_by_id);

-- Table structure for table `tenant_members`
CREATE TABLE IF NOT EXISTS tenant_members (
    id BIGSERIAL NOT NULL,
    tenant_id BIGINT NOT NULL,
    user_id BIGINT NOT NULL,
    created_at TIMESTAMPTZ NOT NULL,
    created_by_id BIGINT NOT NULL,
    updated_at TIMESTAMPTZ NOT NULL,
    updated_by_id BIGINT NOT NULL,
    PRIMARY KEY (id),
    CONSTRAINT uq_tenant_members UNIQUE (tenant_id, user_id),
    CONSTRAINT fk_tenant_members_tenant FOREIGN KEY (tenant_id) REFERENCES tenants(id),
    CONSTRAINT fk_tenant_members_user FOREIGN KEY (user_id) REFERENCES users(id),
    CONSTRAINT fk_tenant_members_created_by FOREIGN KEY (created_by_id) REFERENCES users(id),
    CONSTRAINT fk_tenant_members_updated_by FOREIGN KEY (updated_by_id) REFERENCES users(id)
);

-- Indexes for tenant_members
CREATE INDEX idx_tenant_members_user_id ON tenant_members (user_id);
CREATE INDEX idx_tenant_members_created_by_id ON tenant_members (created_by_id);
CREATE INDEX idx_tenant_members_updated_by_id ON tenant_members (updated_by_id);

-- Table structure for table `permissions`
CREATE TABLE IF NOT EXISTS permissions (
    id BIGSERIAL NOT NULL,
    code VARCHAR(100) NOT NULL UNIQUE,
    name VARCHAR(100) NOT NULL,
    created_at TIMESTAMPTZ NOT NULL,
    created_by_id BIGINT NOT NULL,
    updated_at TIMESTAMPTZ NOT NULL,
    updated_by_id BIGINT NOT NULL,
    PRIMARY KEY (id),
    CONSTRAINT fk_permissions_created_by FOREIGN KEY (created_by_id) REFERENCES users(id),
    CONSTRAINT fk_permissions_updated_by FOREIGN KEY (updated_by_id) REFERENCES users(id)
);

-- Indexes for permissions
CREATE INDEX idx_permissions_created_by_id ON permissions (created_by_id);
CREATE INDEX idx_permissions_updated_by_id ON permissions (updated_by_id);

-- Table structure for table `tenant_member_permissions`
CREATE TABLE IF NOT EXISTS tenant_member_permissions (
    tenant_member_id BIGINT NOT NULL,
    permission_id BIGINT NOT NULL,
    PRIMARY KEY (tenant_member_id, permission_id),
    CONSTRAINT fk_tmp_member FOREIGN KEY (tenant_member_id) REFERENCES tenant_members(id),
    CONSTRAINT fk_tmp_permission FOREIGN KEY (permission_id) REFERENCES permissions(id)
);

-- Indexes for tenant_member_permissions
CREATE INDEX idx_tmp_permission_id ON tenant_member_permissions (permission_id);

-- Table structure for table `tenant_member_resource_scopes`
CREATE TABLE IF NOT EXISTS tenant_member_resource_scopes (
    tenant_member_id BIGINT NOT NULL,
    resource_code VARCHAR(75) NOT NULL,
    scope_by VARCHAR(45) NOT NULL,
    scope_ref VARCHAR(100) NOT NULL,
    PRIMARY KEY (tenant_member_id, resource_code, scope_by, scope_ref),
    CONSTRAINT fk_tmrs_member FOREIGN KEY (tenant_member_id) REFERENCES tenant_members(id)
);

-- Index for tenant_member_resource_scopes
CREATE INDEX idx_tmrs_lookup ON tenant_member_resource_scopes (tenant_member_id, resource_code);

-- Table structure for table `organizations`
CREATE TABLE IF NOT EXISTS organizations (
    id BIGSERIAL NOT NULL,
    public_id UUID NOT NULL UNIQUE,
    tenant_public_id UUID NOT NULL UNIQUE,
    name VARCHAR(40) NOT NULL,
    description TEXT NULL DEFAULT NULL,
    created_at TIMESTAMPTZ NOT NULL,
    created_by_public_id UUID NOT NULL,
    updated_at TIMESTAMPTZ NOT NULL,
    updated_by_public_id UUID NOT NULL,
    PRIMARY KEY (id)
);

-- Table structure for table `branches`
CREATE TABLE IF NOT EXISTS branches (
    id BIGSERIAL NOT NULL,
    public_id UUID NOT NULL UNIQUE,
    organization_id BIGINT NOT NULL,
    name VARCHAR(40) NOT NULL,
    created_at TIMESTAMPTZ NOT NULL,
    created_by_public_id UUID NOT NULL,
    updated_at TIMESTAMPTZ NOT NULL,
    updated_by_public_id UUID NOT NULL,
    PRIMARY KEY (id),
    CONSTRAINT fk_branches_organizations FOREIGN KEY (organization_id) REFERENCES organizations (id) ON DELETE NO ACTION ON UPDATE CASCADE
);

-- Index for branches
CREATE INDEX idx_branches_organizations ON branches (organization_id);

-- Table structure for table `transactions`
CREATE TABLE IF NOT EXISTS transactions (
  id BIGSERIAL NOT NULL,
  public_id UUID NOT NULL UNIQUE,
  organization_public_id UUID NOT NULL,
  organization_name VARCHAR(40) NOT NULL,
  branch_public_id UUID NOT NULL,
  branch_name VARCHAR(40) NOT NULL,
  bill_total NUMERIC(19,4) NOT NULL,
  transacted_at TIMESTAMPTZ NOT NULL,
  created_at TIMESTAMPTZ NOT NULL,
  created_by_public_id UUID NOT NULL,
  updated_at TIMESTAMPTZ NOT NULL,
  updated_by_public_id UUID NOT NULL,
  PRIMARY KEY (id)
);

-- Indexes for transactions
CREATE INDEX idx_organization ON transactions (organization_public_id);
CREATE INDEX idx_transacted ON transactions (transacted_at DESC);
CREATE INDEX idx_branch_bill_total ON transactions (branch_public_id, bill_total DESC);
CREATE INDEX idx_branch_transacted ON transactions (branch_public_id, transacted_at DESC);
