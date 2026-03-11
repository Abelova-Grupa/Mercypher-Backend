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
    container_name       = "services"
    key                  = "message/terraform.tfstate"
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

data "azurerm_subnet" "mercypher-env-subnet" {
  name                 = "mercypher-env-subnet"
  virtual_network_name = data.azurerm_virtual_network.vnet.name
  resource_group_name  = data.azurerm_resource_group.mercypher-backend.name
}

data "azurerm_key_vault" "mercypher-keyvault" {
  name                = "mercypher-keyvault"
  resource_group_name = "mercypher-utils"
}

data "azurerm_key_vault_secret" "service_bus_conn" {
  name         = "service-bus-conn-str"
  key_vault_id = data.azurerm_key_vault.mercypher-keyvault.id
}

data "azurerm_key_vault_secret" "postgres_host" {
  name         = "postgres-host-url"
  key_vault_id = data.azurerm_key_vault.mercypher-keyvault.id
}

data "azurerm_log_analytics_workspace" "analytics-workspace" {
  name                = "log-workspace-mercypher"
  resource_group_name = "mercypher-backend"
}

data "azurerm_container_app_environment" "mercypher-backend-environment" {
  name                = "mercypher-backend-environment"
  resource_group_name = "mercypher-backend"
}

resource "azurerm_container_app" "message-mercypher-prod-itan-01" {
  name                         = "message-mercypher-prod-itan-01"
  container_app_environment_id = data.azurerm_container_app_environment.mercypher-backend-environment.id
  resource_group_name          = data.azurerm_resource_group.mercypher-backend.name
  revision_mode                = "Multiple"

  secret {
    name                = "service-bus-conn-str"
    key_vault_secret_id = data.azurerm_key_vault_secret.service_bus_conn.id
    identity            = "System"
  }

  secret {
    name                = "postgres-host-url"
    key_vault_secret_id = data.azurerm_key_vault_secret.postgres_host.id
    identity            = "System"
  }

  secret {
    name                = "dockerhub-pass"
    key_vault_secret_id = data.azurerm_key_vault_secret.postgres_host.id
    identity            = "System"
  }

  template {
    container {
      name   = "message-mercypher-prod-itan-01"
      image  = "lukadervisevic/mercypher-message:${var.image_tag}"
      cpu    = 2
      memory = "4Gi"

      // Quick solution for env variables
      env {
        name  = "ENVIRONMENT"
        value = "azure"
      }
      env {
        name  = "PORT"
        value = var.port
      }
      env {
        name  = "POSTGRES_USER"
        value = var.postgres_user
      }
      env {
        name  = "POSTGRES_PASSWORD"
        value = var.postgres_password
      }
      env {
        name  = "POSTGRES_PORT"
        value = var.postgres_port
      }
      env {
        name  = "POSTGRES_DB"
        value = var.postgres_db
      }
      env {
        name        = "POSTGRES_HOST"
        secret_name = "postgres-host-url"
      }

      env {
        name  = "KAFKA_BROKERS"
        value = "localhost:9092"
      }
      env {
        name  = "INGRESS_PORT"
        value = "443"
      }

      env {
        name        = "AZURE_SERVICE_BUS_CONN_STR"
        secret_name = "service-bus-conn-str"
      }

    }

    min_replicas = 0
    max_replicas = 1
    tcp_scale_rule {
      name                = "message-service-scaling-rule"
      concurrent_requests = "20"
    }

  }

  identity {
    type = "SystemAssigned"
  }

  registry {
    server               = "index.docker.io"
    username             = var.dockerhub_user
    password_secret_name = "dockerhub-pass"
  }

  ingress {
    transport        = "tcp"
    target_port      = 50052
    exposed_port     = 50052
    external_enabled = false
    traffic_weight {
      percentage      = 100
      latest_revision = true
    }
  }
}

resource "azurerm_role_assignment" "aca_keyvault_access_message" {
  scope                = data.azurerm_key_vault.mercypher-keyvault.id
  role_definition_name = "Key Vault Secrets User"
  principal_id         = azurerm_container_app.message-mercypher-prod-itan-01.identity[0].principal_id
}
