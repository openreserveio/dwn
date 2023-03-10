resource "aws_security_group" "allow_docdb" {
  name        = "allow_docdb"
  description = "Allow inbound traffic to the DocDB server"
  vpc_id      = module.vpc.vpc_id

  ingress {
    description      = "Connections from VPC"
    from_port        = 27017
    to_port          = 27017
    protocol         = "tcp"
    cidr_blocks      = [module.vpc.vpc_cidr_block]
  }

  egress {
    from_port        = 0
    to_port          = 0
    protocol         = "-1"
    cidr_blocks      = ["0.0.0.0/0"]
    ipv6_cidr_blocks = ["::/0"]
  }

  tags = local.tags

}

resource "aws_docdb_cluster_instance" "docdb_cluster_instances" {
  count              = 2
  identifier         = "docdb-cluster-demo-${count.index}"
  cluster_identifier = aws_docdb_cluster.docdb_cluster.id
  instance_class     = "db.t3.medium"

  tags = local.tags

}

resource "aws_docdb_cluster" "docdb_cluster" {
  cluster_identifier = "${local.name}-docdb-cluster"
  availability_zones = ["us-west-2a", "us-west-2b", "us-west-2c"]
  master_username    = "openreserveuser"
  master_password    = "openreservepass"
  port = 27017

  engine = "docdb"
  
  tags = local.tags

}

resource "aws_route53_record" "docdb" {
  
  zone_id = aws_route53_zone.internal_route53_zone.zone_id
  name    = "docdb"
  type    = "CNAME"
  ttl     = 5


  set_identifier = "docdb"
  records        = [aws_docdb_cluster.docdb_cluster.endpoint]


}