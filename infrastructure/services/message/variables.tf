variable "subscription_id" {
  description = "Azure Subscription ID"
  type        = string
  sensitive   = true
}

# variable "service_bus_conn_str" {
#   description = "Azure service bus connection string"
#   type = string
#   sensitive = true
# }

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

variable "port" {
  description = "Docker hub user password"
  type        = string
  sensitive   = true
}


variable "image_tag" {
  type    = string
  default = "latest"
}