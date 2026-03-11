variable "postgres_user" {
  description = "User for postgres cli"
  type        = string
  sensitive   = true
}

variable "postgres_password" {
  description = "Password for postgres cli"
  type        = string
  sensitive   = true
}

# variable "postgres_host" {
#   description = "Postgres host"
#   type        = string
#   sensitive   = true
# }

variable "postgres_db" {
  description = "Postgres database"
  type        = string
  sensitive   = true
}

variable "postgres_port" {
  description = "Postgres database"
  type        = string
  sensitive   = true
}

variable "session_service_port" {
  description = "Session service port"
  type        = string
  sensitive   = true
}

variable "user_service_port" {
  description = "Session service port"
  type        = string
  sensitive   = true
}

# variable "azure-redis-cache-port-tls" {
#   description = "Azure redis cache port tls"
#   type = string
#   sensitive = true
# }

variable "subscription_id" {
  description = "Azure Subscription ID"
  type        = string
  sensitive   = true
}

variable "dockerhub_user" {
  description = "Docker hub user"
  type        = string
  sensitive   = true
}

# variable "dockerhub_pass" {
#   description = "Docker hub user password"
#   type        = string
#   sensitive   = true
# }

variable "image_tag" {
  type    = string
  default = "v0.0.16"
}

variable "email_app_pass" {
  type    = string
  default = "latest"
}

variable "email_client" {
  type    = string
  default = "latest"
}

variable "smtp_host" {
  type    = string
  default = "latest"
}

variable "user_service_uuid" {
  type    = string
  default = "latest"
}

# variable "azure_redis_cache_url" {
#   type      = string
#   sensitive = true
# }

# variable "azure_redis_access_key" {
#   type      = string
#   sensitive = true
# }

# variable "postgres_host_url" {
#   type      = string
#   sensitive = true
# }

variable "session_container_app_url" {
  type = string
  sensitive = true
}


