variable "subscription_id" {
  description = "Azure Subscription ID"
  type        = string
  sensitive   = true
}

variable "dockerhub_user" {
  description = "Azure Subscription ID"
  type        = string
  sensitive   = true
}

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

variable "port" {
  description = "Group Port"
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

variable "image_tag" {
  type    = string
  default = "latest"
}
