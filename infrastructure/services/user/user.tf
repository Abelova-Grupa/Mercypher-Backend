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

data "azurerm_key_vault_secret" "redis_url" {
  name         = "azure-redis-cache-url"
  key_vault_id = data.azurerm_key_vault.mercypher-keyvault.id
}

data "azurerm_key_vault_secret" "redis_key" {
  name         = "azure-redis-access-key"
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

resource "azurerm_container_app" "user-mercypher-prod-itan-01" {
  name                         = "user-mercypher-prod-itan-01"
  container_app_environment_id = data.azurerm_container_app_environment.mercypher-backend-environment.id
  resource_group_name          = data.azurerm_resource_group.mercypher-backend.name
  revision_mode                = "Multiple"

  secret {
    name                = "azure-redis-cache-url"
    key_vault_secret_id = data.azurerm_key_vault_secret.redis_url.id
    identity            = "System"
  }

  secret {
    name                = "azure-redis-access-key"
    key_vault_secret_id = data.azurerm_key_vault_secret.redis_key.id
    identity            = "System"
  }

  secret {
    name                = "postgres-host-url"
    key_vault_secret_id = data.azurerm_key_vault_secret.postgres_host.id
    identity            = "System"
  }

  secret {
    name                = "docker-password"
    key_vault_secret_id = data.azurerm_key_vault_secret.postgres_host.id
    identity            = "System"
  }


  template {
    container {
      name   = "user-mercypher-prod-itan-01"
      image  = "lukadervisevic/mercypher-user:${var.image_tag}"
      cpu    = 2
      memory = "4Gi"

      // Quick solution for env variables
      env {
        name  = "ENVIRONMENT"
        value = "azure"
      }
      env {
        name  = "SESSION_SERVICE_PORT"
        value = var.session_service_port
      }
      env {
        name  = "USER_SERVICE_PORT"
        value = var.user_service_port
      }
      env {
        name        = "AZURE_REDIS_CACHE_URL"
        secret_name = "azure-redis-cache-url"
      }
      env {
        name  = "AZURE_REDIS_CACHE_PORT_TLS"
        value = 6380
      }
      env {
        name        = "AZURE_REDIS_ACCESS_KEY"
        secret_name = "azure-redis-access-key"
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
        name        = "POSTGRES_HOST"
        secret_name = "postgres-host-url"
      }
      env {
        name  = "POSTGRES_DB"
        value = var.postgres_db
      }

      env {
        name  = "EMAIL_APP_PASS"
        value = var.email_app_pass
      }

      env {
        name  = "EMAIL_CLIENT"
        value = var.email_client
      }

      env {
        name  = "SMTP_HOST"
        value = var.smtp_host
      }

      env {
        name  = "TLS_SMTP_PORT"
        value = 587
      }

      env {
        name  = "SSL_SMTP_PORT"
        value = 465
      }

      env {
        name  = "USER_SERVICE_UUID"
        value = var.user_service_uuid
      }

      env {
        name = "SESSION_CONTAINER_APP_URL"
        value = var.session_container_app_url
      }

    }

    min_replicas = 1
    max_replicas = 10
    tcp_scale_rule {
      name                = "user-service-scale-rules"
      concurrent_requests = "20"
    }
  }

  identity {
    type = "SystemAssigned"
  }

  registry {
    server               = "index.docker.io"
    username             = var.dockerhub_user
    password_secret_name = "docker-password"
  }

  ingress {
    allow_insecure_connections = false
    transport                  = "tcp"
    target_port                = 50054
    exposed_port = 50054
    external_enabled           = false
    traffic_weight {
      percentage      = 100
      latest_revision = true
    }
  }
}

resource "azurerm_role_assignment" "aca_keyvault_access" {
  scope                = data.azurerm_key_vault.mercypher-keyvault.id
  role_definition_name = "Key Vault Secrets User"
  principal_id         = azurerm_container_app.user-mercypher-prod-itan-01.identity[0].principal_id
}
