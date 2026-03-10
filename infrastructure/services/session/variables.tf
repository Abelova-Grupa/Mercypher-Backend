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

# variable "session_service_port" {
#   description = "Session service port"
#   type        = string
#   sensitive   = true
# }

# variable "azure-redis-cache-port-tls" {
#   description = "Azure redis cache port tls"
#   type        = string
#   sensitive   = true
# }

variable "image_tag" {
  type    = string
  default = "latest"
}



