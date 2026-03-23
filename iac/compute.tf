# compute.tf

# Security group for processor nodes

resource "aws_security_group" "processor_sg" {
  name        = "proxy-processor-sg"
  description = "Allow traffic from ALB"
  vpc_id      = aws_vpc.proxy_vpc.id

  # Inbound traffic from ALB only
  ingress {
    from_port       = var.processor_port
    to_port         = var.processor_port
    protocol        = "tcp"
    security_groups = [aws_security_group.alb_sg.id]
  }

  # Outbound to the world
  egress {
    from_port   = 0
    to_port     = 0
    protocol    = "-1"
    cidr_blocks = ["0.0.0.0/0"]
  }
}

# Launch Template for Processor Nodes

resource "aws_launch_template" "proxy_lt" {
  name_prefix   = "proxy-processor-"
  image_id      = var.ami_id
  instance_type = var.instance_type

  network_interfaces {
    associate_public_ip_address = true
    security_groups             = [aws_security_group.processor_sg.id]
  }

  user_data = base64encode(templatefile("${path.module}/user_data.sh", {
    PROCESSOR_PORT   = var.processor_port
    PROXY_AUTH_TOKEN = var.proxy_auth_token
    CACHE_SALT       = var.cache_salt
    REDIS_HOST       = aws_elasticache_cluster.proxy_cache.cache_nodes[0].address
    REDIS_PORT       = "6379"
    CONTAINER_IMAGE  = var.container_image
  }))

  tag_specifications {
    resource_type = "instance"
    tags = {
      Name = "proxy-processor-node"
    }
  }

  lifecycle {
    create_before_destroy = true
  }
}
