terraform {
  required_providers {
    azurerm = {
      source  = "hashicorp/azurerm"
      version = "=4.1.0"
    }
  }
}

provider "azurerm" {
  subscription_id = "caf77c5d-90e7-49a7-8a58-893de3bfe0ff"
  features {}
}

resource "azurerm_resource_group" "mercypher-backend" {
  name     = "mercypher-backend"
  location = "Italy North"
}

resource "azurerm_virtual_network" "vnet-mercypher" {
  name                = "vnet-mercypher-prod-itan-01"
  location            = azurerm_resource_group.mercypher-backend.location
  resource_group_name = azurerm_resource_group.mercypher-backend.name
  address_space       = ["10.0.0.0/16"]
}

resource "azurerm_subnet" "redis-mercypher-subnet" {
  name                 = "redis-mercypher-subnet"
  resource_group_name  = azurerm_resource_group.mercypher-backend.name
  virtual_network_name = azurerm_virtual_network.vnet-mercypher.name
  address_prefixes     = ["10.0.3.0/29"]
}

resource "azurerm_subnet" "service-bus-subnet" {
  name                 = "service-bus-subnet"
  resource_group_name  = azurerm_resource_group.mercypher-backend.name
  virtual_network_name = azurerm_virtual_network.vnet-mercypher.name
  address_prefixes     = ["10.0.4.0/29"]
}

resource "azurerm_subnet" "postgres-mercypher-subnet" {
  name                 = "postgres-mercypher-subnet"
  resource_group_name  = azurerm_resource_group.mercypher-backend.name
  virtual_network_name = azurerm_virtual_network.vnet-mercypher.name
  address_prefixes     = ["10.0.2.0/28"]

  delegation {
    name = "postgres-delegation"
    service_delegation {
      name    = "Microsoft.DBforPostgreSQL/flexibleServers"
      actions = ["Microsoft.Network/virtualNetworks/subnets/join/action"]
    }
  }
}

resource "azurerm_subnet" "mercypher-env-subnet" {
  name                 = "mercypher-env-subnet"
  resource_group_name  = azurerm_resource_group.mercypher-backend.name
  virtual_network_name = azurerm_virtual_network.vnet-mercypher.name
  address_prefixes     = ["10.0.0.0/23"]

  delegation {
    name = "container-apps-delegation"
    service_delegation {
      name    = "Microsoft.App/environments"
      actions = ["Microsoft.Network/virtualNetworks/subnets/join/action"]
    }
  }
}

resource "azurerm_subnet" "vm-mercypher-subnet" {
  name                 = "vm-mercypher-subnet"
  resource_group_name  = azurerm_resource_group.mercypher-backend.name
  virtual_network_name = azurerm_virtual_network.vnet-mercypher.name
  address_prefixes     = ["10.0.5.0/28"]
}

resource "azurerm_private_dns_zone" "postgres_mercypher_dns" {
  name                = "postgres-mercypher-dns.private.postgres.database.azure.com"
  resource_group_name = azurerm_resource_group.mercypher-backend.name
}

resource "azurerm_private_dns_zone_virtual_network_link" "postgres_dns_link" {
  name                  = "postgres-dns-link"
  resource_group_name   = azurerm_resource_group.mercypher-backend.name
  private_dns_zone_name = azurerm_private_dns_zone.postgres_mercypher_dns.name
  virtual_network_id    = azurerm_virtual_network.vnet-mercypher.id
  registration_enabled  = false
}

resource "azurerm_private_dns_zone" "redis_mercypher_dns" {
  name                = "privatelink.redis.cache.windows.net"
  resource_group_name = azurerm_resource_group.mercypher-backend.name
}

resource "azurerm_private_dns_zone_virtual_network_link" "redis_dns_link" {
  name                  = "redis_dns_link"
  resource_group_name   = azurerm_resource_group.mercypher-backend.name
  private_dns_zone_name = azurerm_private_dns_zone.redis_mercypher_dns.name
  virtual_network_id    = azurerm_virtual_network.vnet-mercypher.id
  registration_enabled  = false
}
