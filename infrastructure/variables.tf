variable "app_name" {
  type        = string
  description = "Name of the app"
  default     = "Taterank"
}

variable "region" {
  type        = string
  description = "AWS region to deploy infrastructure"
  default     = "us-east-2"
}

variable "environment" {
  type        = string
  description = "Environment name for deployment"
  default     = "dev"
}
