# variables.tf

variable "aws_region" {
  description = "The AWS region to deploy into"
  type        = string
  default     = "us-east-1"
}

variable "proxy_auth_token" {
  description = "The secret token for the X-Proxy-Auth header"
  type        = string
  sensitive   = true
}

variable "processor_port" {
  description = "The port the proxy processor app listens on"
  type        = number
  default     = 8080
}

variable "vpc_cidr" {
  description = "CIDR block for the VPC"
  type        = string
  default     = "10.0.0.0/16"
}

variable "instance_type" {
  description = "Instance type"
  type        = string
  default     = "t3.micro"
}

variable "ami_id" {
  description = "The AMI ID for the processor nodes (Amazon Linux 2023)"
  type        = string
  default     = "ami-04aa82396fe417f2f"
}

variable "container_image" {
  description = "The URI of the proxy container image"
  type        = string
  default     = "ghcr.io/triviajon/proxy-processor:latest"
}
