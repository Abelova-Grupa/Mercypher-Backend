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

variable "subscription_id" {
  description = "Azure Subscription ID"
  type        = string
  sensitive   = true
}

