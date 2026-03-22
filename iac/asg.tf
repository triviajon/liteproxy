resource "aws_autoscaling_group" "proxy_asg" {
  name                = "proxy-asg"
  vpc_zone_identifier = [aws_subnet.public_1.id, aws_subnet.public_2.id]
  target_group_arns   = [aws_lb_target_group.proxy_tg.arn]
  health_check_type   = "ELB"

  min_size         = 1
  max_size         = 1
  desired_capacity = 1

  launch_template {
    id      = aws_launch_template.proxy_lt.id
    version = "$Latest"
  }

  tag {
    key                 = "Name"
    value               = "proxy-asg-node"
    propagate_at_launch = true
  }
}
