terraform {
  required_providers {
    azurerm = {
      source  = "hashicorp/azurerm"
      version = "=4.1.0"
    }
  }
  backend "azurerm" {
    resource_group_name  = "mercypher-utils"
    storage_account_name = "blobmercypher"
    container_name       = "data"
    key                  = "terraform.tfstate"
  }
}

provider "azurerm" {
  subscription_id = var.subscription_id
  features {}
}

data "azurerm_resource_group" "mercypher-backend" {
  name = "mercypher-backend"
}

data "azurerm_virtual_network" "vnet" {
  name                = "vnet-mercypher-prod-itan-01"
  resource_group_name = data.azurerm_resource_group.mercypher-backend.name
}

data "azurerm_subnet" "postgres_subnet" {
  name                 = "postgres-mercypher-subnet"
  virtual_network_name = data.azurerm_virtual_network.vnet.name
  resource_group_name  = data.azurerm_resource_group.mercypher-backend.name
}

data "azurerm_subnet" "mercypher-env-subnet" {
  name                 = "mercypher-env-subnet"
  virtual_network_name = data.azurerm_virtual_network.vnet.name
  resource_group_name  = data.azurerm_resource_group.mercypher-backend.name
}

data "azurerm_key_vault" "mercypher-keyvault" {
  name                = "mercypher-keyvault"
  resource_group_name = "mercypher-utils"
}

data "azurerm_private_dns_zone" "postgres_mercypher_dns" {
  name                = "postgres-mercypher-dns.private.postgres.database.azure.com"
  resource_group_name = data.azurerm_resource_group.mercypher-backend.name
}

data "azurerm_private_dns_zone" "redis_dns" {
  name                = "privatelink.redis.cache.windows.net"
  resource_group_name = data.azurerm_resource_group.mercypher-backend.name
}

data "azurerm_subnet" "redis_subnet" {
  name                 = "redis-mercypher-subnet"
  virtual_network_name = data.azurerm_virtual_network.vnet.name
  resource_group_name  = data.azurerm_resource_group.mercypher-backend.name
}

resource "azurerm_postgresql_flexible_server" "postgres-mercypher-prod-itan-01" {
  name                          = "postgres-mercypher-prod-itan-01"
  resource_group_name           = data.azurerm_resource_group.mercypher-backend.name
  location                      = data.azurerm_resource_group.mercypher-backend.location
  version                       = "16"
  delegated_subnet_id           = data.azurerm_subnet.postgres_subnet.id
  private_dns_zone_id           = data.azurerm_private_dns_zone.postgres_mercypher_dns.id
  public_network_access_enabled = false
  administrator_login           = var.postgres_user
  administrator_password        = var.postgres_password
  zone                          = "1"
  sku_name                      = "B_Standard_B1ms"
  depends_on                    = [data.azurerm_private_dns_zone.postgres_mercypher_dns]
}

resource "azurerm_postgresql_flexible_server_database" "mercypher-user" {
  name      = "mercypher-user"
  server_id = azurerm_postgresql_flexible_server.postgres-mercypher-prod-itan-01.id
  collation = "en_US.utf8"
  charset   = "utf8"
}

resource "azurerm_postgresql_flexible_server_database" "mercypher-groups" {
  name      = "mercypher-groups"
  server_id = azurerm_postgresql_flexible_server.postgres-mercypher-prod-itan-01.id
  collation = "en_US.utf8"
  charset   = "utf8"
}

resource "azurerm_postgresql_flexible_server_database" "mercypher-message" {
  name      = "mercypher-message"
  server_id = azurerm_postgresql_flexible_server.postgres-mercypher-prod-itan-01.id
  collation = "en_US.utf8"
  charset   = "utf8"
}



resource "azurerm_redis_cache" "azure_redis_cache" {
  name                = "redis-mercypher-prod-itan-01"
  location            = data.azurerm_resource_group.mercypher-backend.location
  resource_group_name = data.azurerm_resource_group.mercypher-backend.name

  capacity                      = 1
  family                        = "C"
  sku_name                      = "Standard"
  non_ssl_port_enabled          = false
  minimum_tls_version           = "1.2"
  public_network_access_enabled = false
}

resource "azurerm_private_endpoint" "redis_pe" {
  name                = "redis-mercypher-endpoint"
  location            = data.azurerm_resource_group.mercypher-backend.location
  resource_group_name = data.azurerm_resource_group.mercypher-backend.name
  subnet_id           = data.azurerm_subnet.redis_subnet.id 

  private_service_connection {
    name                           = "redis-privateserviceconnection"
    private_connection_resource_id = azurerm_redis_cache.azure_redis_cache.id
    is_manual_connection           = false
    subresource_names              = ["redisCache"]
  }

  private_dns_zone_group {
    name                 = "redis-dns-zone-group"
    private_dns_zone_ids = [data.azurerm_private_dns_zone.redis_dns.id]
  }
}

resource "azurerm_servicebus_namespace" "bus-mercypher-prod-itan-01" {
  name                = "bus-mercypher-prod-itan-01"
  location            = data.azurerm_resource_group.mercypher-backend.location
  resource_group_name = data.azurerm_resource_group.mercypher-backend.name
  sku                 = "Standard"

  tags = {
    source = "terraform"
  }
}

// KEY VAULT SECRETS
resource "azurerm_key_vault_secret" "redis_key" {
  name         = "azure-redis-access-key"
  value        = azurerm_redis_cache.azure_redis_cache.primary_access_key
  key_vault_id = data.azurerm_key_vault.mercypher-keyvault.id
}

resource "azurerm_key_vault_secret" "redis_cache_url" {
  name         = "azure-redis-cache-url"
  key_vault_id = data.azurerm_key_vault.mercypher-keyvault.id
  value        = azurerm_redis_cache.azure_redis_cache.hostname
}

resource "azurerm_key_vault_secret" "postgres_host_url" {
  name         = "postgres-host-url"
  key_vault_id = data.azurerm_key_vault.mercypher-keyvault.id
  value        = azurerm_postgresql_flexible_server.postgres-mercypher-prod-itan-01.fqdn
}

resource "azurerm_key_vault_secret" "service_bus_conn_str" {
  name         = "service-bus-conn-str"
  value        = azurerm_servicebus_namespace.bus-mercypher-prod-itan-01.default_primary_connection_string
  key_vault_id = data.azurerm_key_vault.mercypher-keyvault.id
}

resource "azurerm_log_analytics_workspace" "analytics-workspace" {
  name                = "log-workspace-mercypher"
  location            = data.azurerm_resource_group.mercypher-backend.location
  resource_group_name = data.azurerm_resource_group.mercypher-backend.name
  sku                 = "PerGB2018"
  retention_in_days   = 30
}

resource "azurerm_container_app_environment" "mercypher-backend-environment" {
  name                       = "mercypher-backend-environment"
  location                   = data.azurerm_resource_group.mercypher-backend.location
  resource_group_name        = data.azurerm_resource_group.mercypher-backend.name
  log_analytics_workspace_id = azurerm_log_analytics_workspace.analytics-workspace.id
  infrastructure_subnet_id   = data.azurerm_subnet.mercypher-env-subnet.id

  lifecycle {
    ignore_changes = [
      infrastructure_resource_group_name,
    ]
  }
}


