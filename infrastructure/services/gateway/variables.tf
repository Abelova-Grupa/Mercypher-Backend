variable "subscription_id" {
  description = "Azure Subscription ID"
  type        = string
  sensitive   = true
}

# variable "service_bus_conn_str" {
#   description = "Azure service bus connection string"
#   type        = string
#   sensitive   = true
# }

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

variable "http_port" {
  description = "HTTP port"
  type        = string
  sensitive   = true
}

# variable "grpc_port" {
#   description = "grpc port"
#   type        = string
#   sensitive   = true
# }

variable "message_host" {
  description = "Message service host"
  type        = string
  sensitive   = true
}

variable "user_host" {
  description = "User service host"
  type        = string
  sensitive   = true
}

variable "session_host" {
  description = "Session service host"
  type        = string
  sensitive   = true
}

variable "group_host" {
  description = "Group service host"
  type        = string
  sensitive   = true
}

variable "message_port" {
  description = "Message service port"
  type        = string
  sensitive   = true
}

variable "user_port" {
  description = "User service port"
  type        = string
  sensitive   = true
}

variable "session_port" {
  description = "Session service port"
  type        = string
  sensitive   = true
}

variable "group_port" {
  description = "Session service port"
  type        = string
  sensitive   = true
}

variable "image_tag" {
  type      = string
  sensitive = true
}