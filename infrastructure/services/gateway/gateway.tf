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

data "azurerm_key_vault_secret" "service_bus_conn" {
  name         = "service-bus-conn-str"
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

resource "azurerm_container_app" "gateway-mercypher-prod-itan-01" {
  name                         = "gateway-mercypher-prod-itan-01"
  container_app_environment_id = data.azurerm_container_app_environment.mercypher-backend-environment.id
  resource_group_name          = data.azurerm_resource_group.mercypher-backend.name
  revision_mode                = "Multiple"

  secret {
    name                = "service-bus-conn-str"
    key_vault_secret_id = data.azurerm_key_vault_secret.service_bus_conn.id
    identity            = "System"
  }

  secret {
    name                = "docker-password"
    key_vault_secret_id = data.azurerm_key_vault_secret.docker_pass.id
    identity            = "System"
  }

  template {
    container {
      name   = "gateway-mercypher-prod-itan-01"
      image  = "lukadervisevic/mercypher-gateway:${var.image_tag}"
      cpu    = 2
      memory = "4Gi"

      // Quick solution for env variables
      env {
        name  = "ENVIRONMENT"
        value = "azure"
      }
      env {
        name  = "HTTP_PORT"
        value = 8080
      }
      env {
        name  = "GRPC_PORT"
        value = 50051
      }
      env {
        name  = "MESSAGE_HOST"
        value = var.message_host
      }
      env {
        name  = "USER_HOST"
        value = var.user_host
      }
      env {
        name  = "SESSION_HOST"
        value = var.session_host
      }
      env {
        name  = "GROUP_HOST"
        value = var.group_host
      }
      env {
        name  = "MESSAGE_PORT"
        value = var.message_port
      }
      env {
        name  = "USER_PORT"
        value = var.user_port
      }
      env {
        name  = "SESSION_PORT"
        value = var.session_port
      }
      env {
        name  = "GROUP_PORT"
        value = var.group_port
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
    max_replicas = 10
    http_scale_rule {
      name                = "tcp-scaling-rule"
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
    transport                  = "http"
    target_port                = 8080
    external_enabled           = true
    traffic_weight {
      percentage      = 100
      latest_revision = true
    }
  }
}

resource "azurerm_role_assignment" "kv_rbac_secrets_gateway" {
  scope                = data.azurerm_key_vault.mercypher-keyvault.id
  role_definition_name = "Key Vault Secrets User"
  principal_id         = azurerm_container_app.gateway-mercypher-prod-itan-01.identity[0].principal_id
}


