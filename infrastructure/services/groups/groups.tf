terraform {
  required_providers {
    azurerm = {
      source  = "hashicorp/azurerm"
      version = "=4.1.0"
    }
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

data "azurerm_key_vault_secret" "postgres_host" {
  name         = "postgres-host-url"
  key_vault_id = data.azurerm_key_vault.mercypher-keyvault.id
}

data "azurerm_key_vault_secret" "docker_user" {
  name         = "docker-username"
  key_vault_id = data.azurerm_key_vault.mercypher-keyvault.id
}

data "azurerm_key_vault_secret" "docker_pass" {
  name         = "docker-password"
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

resource "azurerm_container_app" "group-mercypher-prod-itan-01" {
  name                         = "group-mercypher-prod-itan-01"
  container_app_environment_id = data.azurerm_container_app_environment.mercypher-backend-environment.id
  resource_group_name          = data.azurerm_resource_group.mercypher-backend.name
  revision_mode                = "Multiple"

  depends_on = [
    data.azurerm_container_app_environment.mercypher-backend-environment,
  ]

  secret {
    name                = "docker-password"
    key_vault_secret_id = data.azurerm_key_vault_secret.docker_pass.id
    identity            = "System"
  }

  secret {
    name                = "postgres-host-url"
    key_vault_secret_id = data.azurerm_key_vault_secret.postgres_host.id
    identity            = "System"
  }

  template {
    container {
      name   = "group-mercypher-prod-itan-01"
      image  = "lukadervisevic/mercypher-groups:${var.image_tag}"
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
        name        = "POSTGRES_HOST"
        secret_name = "postgres-host-url"
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
        name  = "POSTGRES_DB"
        value = var.postgres_db
      }
      env {
        name  = "POSTGRES_PORT"
        value = var.postgres_port
      }

    }

    min_replicas = 0
    max_replicas = 10
    tcp_scale_rule {
      name                = "group-service-scaling-rule"
      concurrent_requests = "20"
    }
  }

  identity {
    type = "SystemAssigned"
  }

  registry {
    server               = "index.docker.io"
    username             = var.docker_username
    password_secret_name = "docker-password"
  }

  ingress {
    transport        = "tcp"
    target_port      = 50056
    exposed_port     = 50056
    external_enabled = false
    traffic_weight {
      percentage      = 100
      latest_revision = true
    }
  }
}

resource "azurerm_role_assignment" "kv_rbac_secrets" {
  scope                = data.azurerm_key_vault.mercypher-keyvault.id
  role_definition_name = "Key Vault Secrets User"
  principal_id         = azurerm_container_app.group-mercypher-prod-itan-01.identity[0].principal_id
}