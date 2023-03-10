provider "aws" {
  region = local.region
}

provider "kubernetes" {
  host                   = module.eks_blueprints.eks_cluster_endpoint
  cluster_ca_certificate = base64decode(module.eks_blueprints.eks_cluster_certificate_authority_data)
  token                  = data.aws_eks_cluster_auth.this.token
}

provider "helm" {
  kubernetes {
    host                   = module.eks_blueprints.eks_cluster_endpoint
    cluster_ca_certificate = base64decode(module.eks_blueprints.eks_cluster_certificate_authority_data)
    token                  = data.aws_eks_cluster_auth.this.token
  }
}

data "aws_eks_cluster_auth" "this" {
  name = module.eks_blueprints.eks_cluster_id
}

data "aws_availability_zones" "available" {}

data "aws_acm_certificate" "issued" {
  domain   = var.acm_certificate_domain
  statuses = ["ISSUED"]
}

locals {
  name = "dwncluster"
  # var.cluster_name is for Terratest
  cluster_name = coalesce(var.cluster_name, local.name)
  region       = "us-west-2"

  vpc_cidr = var.vpc_cidr
  azs      = slice(data.aws_availability_zones.available.names, 0, 3)

  tags = {
    Blueprint  = local.name
    GithubRepo = "github.com/aws-ia/terraform-aws-eks-blueprints"
    Environment = "alpha"
  }
}

#---------------------------------------------------------------
# EKS Blueprints
#---------------------------------------------------------------

module "eks_blueprints" {
  source = "../../../../terraform-aws-eks-blueprints"

  cluster_name    = local.cluster_name
  cluster_version = "1.24"

  vpc_id             = module.vpc.vpc_id
  private_subnet_ids = module.vpc.private_subnets

  #----------------------------------------------------------------------------------------------------------#
  # Security groups used in this module created by the upstream modules terraform-aws-eks (https://github.com/terraform-aws-modules/terraform-aws-eks).
  #   Upstream module implemented Security groups based on the best practices doc https://docs.aws.amazon.com/eks/latest/userguide/sec-group-reqs.html.
  #   So, by default the security groups are restrictive. Users needs to enable rules for specific ports required for App requirement or Add-ons
  #   See the notes below for each rule used in these examples
  #----------------------------------------------------------------------------------------------------------#
  node_security_group_additional_rules = {
    # Extend node-to-node security group rules. Recommended and required for the Add-ons
    ingress_self_all = {
      description = "Node to node all ports/protocols"
      protocol    = "-1"
      from_port   = 0
      to_port     = 0
      type        = "ingress"
      self        = true
    }
    # Recommended outbound traffic for Node groups
    egress_all = {
      description      = "Node all egress"
      protocol         = "-1"
      from_port        = 0
      to_port          = 0
      type             = "egress"
      cidr_blocks      = ["0.0.0.0/0"]
      ipv6_cidr_blocks = ["::/0"]
    }
    # Allows Control Plane Nodes to talk to Worker nodes on all ports. Added this to simplify the example and further avoid issues with Add-ons communication with Control plane.
    # This can be restricted further to specific port based on the requirement for each Add-on e.g., metrics-server 4443, spark-operator 8080, karpenter 8443 etc.
    # Change this according to your security requirements if needed
    ingress_cluster_to_node_all_traffic = {
      description                   = "Cluster API to Nodegroup all traffic"
      protocol                      = "-1"
      from_port                     = 0
      to_port                       = 0
      type                          = "ingress"
      source_cluster_security_group = true
    }
  }

  managed_node_groups = {
    mg_5 = {
      node_group_name = "managed-ondemand"
      instance_types  = ["m5.large"]
      min_size        = 2
      max_size        = 5
      desired_size    = 2
      subnet_ids      = module.vpc.private_subnets
    }
  }

  tags = local.tags
}

module "eks_blueprints_kubernetes_addons" {
  source = "../../../../terraform-aws-eks-blueprints/modules/kubernetes-addons"

  eks_cluster_id       = module.eks_blueprints.eks_cluster_id
  eks_cluster_endpoint = module.eks_blueprints.eks_cluster_endpoint
  eks_oidc_provider    = module.eks_blueprints.oidc_provider
  eks_cluster_version  = module.eks_blueprints.eks_cluster_version
  eks_cluster_domain   = var.eks_cluster_domain

  # EKS Managed Add-ons
  enable_amazon_eks_vpc_cni            = true
  enable_amazon_eks_coredns            = true
  enable_amazon_eks_kube_proxy         = true
  enable_amazon_eks_aws_ebs_csi_driver = true

  # Add-ons
  enable_aws_load_balancer_controller = true
  enable_metrics_server               = true
  enable_aws_cloudwatch_metrics       = true
  enable_kubecost                     = true
  enable_gatekeeper                   = true

  enable_cluster_autoscaler = true
  cluster_autoscaler_helm_config = {
    set = [
      {
        name  = "podLabels.prometheus\\.io/scrape",
        value = "true",
        type  = "string",
      }
    ]
  }

  enable_cert_manager = true
  cert_manager_helm_config = {
    set_values = [
      {
        name  = "extraArgs[0]"
        value = "--enable-certificate-owner-ref=false"
      },
    ]
  }

  enable_argocd = true
  argocd_applications = {
    workloads = {
      path               = "envs/dev"
      repo_url           = "https://github.com/aws-samples/eks-blueprints-workloads.git"
      target_revision    = "main"
      add_on_application = false
      values = {
        spec = {
          source = {
            repoURL        = "https://github.com/aws-samples/eks-blueprints-workloads.git"
            targetRevision = "main"
          }
          blueprint   = "terraform"
          clusterName = module.eks_blueprints.eks_cluster_id
          env         = "dev"
          ingress = {
            type           = "alb"
            host           = var.eks_cluster_domain
            route53_weight = "100" # <-- You can control the weight of the route53 weighted records between clusters
          }
        }
      }
    }
  }

  enable_ingress_nginx = true
  ingress_nginx_helm_config = {
    values = [templatefile("${path.module}/nginx-values.yaml", {
      hostname     = var.eks_cluster_domain
      ssl_cert_arn = data.aws_acm_certificate.issued.arn
    })]
  }

  enable_external_dns                 = true
  external_dns_helm_config = {
    values = [templatefile("${path.module}/external_dns-values.yaml", {
      txtOwnerId   = local.name
      zoneIdFilter = var.eks_cluster_domain
    })]
  }


  # TODO - requires dependency on `cert-manager` for namespace
  # enable_cert_manager_csi_driver = true

  tags = local.tags
}

#---------------------------------------------------------------
# Supporting Resources
#---------------------------------------------------------------

module "vpc" {
  source  = "terraform-aws-modules/vpc/aws"
  version = "~> 3.0"

  name = local.name
  cidr = local.vpc_cidr

  azs             = local.azs
  public_subnets  = [for k, v in local.azs : cidrsubnet(local.vpc_cidr, 8, k)]
  private_subnets = [for k, v in local.azs : cidrsubnet(local.vpc_cidr, 8, k + 10)]

  enable_nat_gateway   = true
  single_nat_gateway   = true
  enable_dns_hostnames = true
  enable_dns_support = true

  # Manage so we can name
  manage_default_network_acl    = true
  default_network_acl_tags      = { Name = "${local.name}-default" }
  manage_default_route_table    = true
  default_route_table_tags      = { Name = "${local.name}-default" }
  manage_default_security_group = true
  default_security_group_tags   = { Name = "${local.name}-default" }

  public_subnet_tags = {
    "kubernetes.io/cluster/${local.cluster_name}" = "shared"
    "kubernetes.io/role/elb"                      = 1
  }

  private_subnet_tags = {
    "kubernetes.io/cluster/${local.cluster_name}" = "shared"
    "kubernetes.io/role/internal-elb"             = 1
  }

  tags = local.tags
}


resource "aws_route53_zone" "internal_route53_zone" {
  
  name = "${var.internal_domain_prefix}.${var.eks_cluster_domain}"

  vpc {
    vpc_id = module.vpc.vpc_id
  }

  tags = local.tags
}
