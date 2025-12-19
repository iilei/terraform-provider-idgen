variable "app_seed" {
  description = "Application-specific seed for deterministic ID generation"
  type        = string
  default     = "app-specific-seed"
}

variable "environment" {
  description = "Environment name (dev, staging, prod)"
  type        = string
  default     = "dev"
}

variable "cluster_size" {
  description = "Cluster size for parametrization"
  type        = number
  default     = 4
}

variable "app_version" {
  description = "Version number for versioned resources"
  type        = number
  default     = 1
}

variable "app_name" {
  description = "Application name for resource naming"
  type        = string
  default     = "myapp"
}

variable "region" {
  description = "Deployment region"
  type        = string
  default     = "us_east_1"
}
