# tflint-ignore: terraform_unused_declarations
variable "cluster_name" {
  description = "Name of cluster - used by Terratest for e2e test automation"
  type        = string
  default     = ""
}

variable "eks_cluster_domain" {
  type        = string
  description = "Route53 domain for the cluster."
  default     = "alpha.openreserve.io"
}

variable "internal_domain_prefix" {
  type = string
  description = "Internal domain to use in a private hosted zone, prepended on the eks_cluster_domain"
  default = "internal"
}

variable "acm_certificate_domain" {
  type        = string
  description = "Route53 certificate domain"
  default     = "alpha.openreserve.io"
}

