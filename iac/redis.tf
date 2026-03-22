# redis.tf

# Security Group

resource "aws_security_group" "redis_sg" {
  name        = "proxy-redis-sg"
  description = "Allow inbound traffic from processor nodes"
  vpc_id      = aws_vpc.proxy_vpc.id

  ingress {
    from_port       = 6379
    to_port         = 6379
    protocol        = "tcp"
    security_groups = [aws_security_group.processor_sg.id]
  }

  egress {
    from_port   = 0
    to_port     = 0
    protocol    = "-1"
    cidr_blocks = []
  }
}

# Subnet Group
resource "aws_elasticache_subnet_group" "proxy_cache_subnets" {
  name       = "proxy-cache-subnet-group"
  subnet_ids = [aws_subnet.public_1.id, aws_subnet.public_2.id]
}

# Elasticache Cluster

resource "aws_elasticache_cluster" "proxy_cache" {
  cluster_id           = "proxy-redis"
  engine               = "redis"
  node_type            = "cache.t3.micro"
  num_cache_nodes      = 1
  parameter_group_name = "default.redis7"
  engine_version       = "7.0"
  port                 = 6379
  subnet_group_name    = aws_elasticache_subnet_group.proxy_cache_subnets.name
  security_group_ids   = [aws_security_group.redis_sg.id]

  tags = {
    Name = "proxy-redis-cache"
  }
}
