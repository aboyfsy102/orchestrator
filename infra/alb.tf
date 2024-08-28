resource "aws_lb" "alb" {
  name               = trim(lower("${var.project}-alb"), "-")
  load_balancer_type = "application"
  security_groups    = [aws_security_group.sg_alb.id]
  subnets            = data.aws_subnets.all.ids

  # enable it after everything stable
  enable_deletion_protection = false

  access_logs {
    bucket  = aws_s3_bucket.default.id
    prefix  = "logs/${var.project}-alb"
    enabled = true
  }
}

resource "aws_lb_listener" "front_end_http" {
  load_balancer_arn = aws_lb.alb.arn
  port              = "80"
  protocol          = "HTTP"

  default_action {
    type = "redirect"
    redirect {
      port        = "443"
      protocol    = "HTTPS"
      status_code = "HTTP_301"
    }
  }
}

resource "aws_lb_listener" "front_end_https" {
  load_balancer_arn = aws_lb.alb.arn
  port              = "443"
  protocol          = "HTTPS"
  ssl_policy        = "ELBSecurityPolicy-2016-08"

  # Use the self-signed certificate
  certificate_arn = aws_acm_certificate.cert.arn

  default_action {
    type = "fixed-response"
    fixed_response {
      content_type = "text/plain"
      message_body = "Not Found"
      status_code  = "404"
    }
  }
}

resource "aws_lb_listener_rule" "describe_rule" {
  listener_arn = aws_lb_listener.front_end_https.arn
  priority     = 100

  action {
    type             = "forward"
    target_group_arn = aws_lb_target_group.lambda_describe_tg.arn
  }

  condition {
    path_pattern {
      values = ["/describe*"]
    }
  }
}

resource "aws_lb_target_group" "lambda_describe_tg" {
  name        = "${var.project}-lambda-tg"
  target_type = "lambda"
}

resource "aws_lb_target_group_attachment" "lambda_describe_tg_attachment" {
  target_group_arn = aws_lb_target_group.lambda_describe_tg.arn
  target_id        = aws_lambda_function.lambda_describe.arn
}

resource "aws_lambda_permission" "allow_alb" {
  statement_id  = "AllowALBInvoke"
  action        = "lambda:InvokeFunction"
  function_name = aws_lambda_function.lambda_describe.function_name
  principal     = "elasticloadbalancing.amazonaws.com"
  source_arn    = aws_lb_target_group.lambda_describe_tg.arn
}

resource "aws_security_group" "sg_alb" {
  name        = "${var.project}-${data.aws_region.current.name}-sg"
  description = "Allow HTTP/HTTPS inbound traffic and all outbound traffic"
  vpc_id      = data.aws_vpc.default.id
}

resource "aws_security_group_rule" "sg_alb_ingress_http" {
  type              = "ingress"
  from_port         = 80
  to_port           = 80
  protocol          = "tcp"
  cidr_blocks       = ["0.0.0.0/0"]
  security_group_id = aws_security_group.sg_alb.id
}

resource "aws_security_group_rule" "sg_alb_ingress_https" {
  type              = "ingress"
  from_port         = 443
  to_port           = 443
  protocol          = "tcp"
  cidr_blocks       = ["0.0.0.0/0"]
  security_group_id = aws_security_group.sg_alb.id
}

resource "aws_security_group_rule" "sg_alb_egress" {
  type              = "egress"
  from_port         = 0
  to_port           = 0
  protocol          = "-1"
  cidr_blocks       = ["0.0.0.0/0"]
  security_group_id = aws_security_group.sg_alb.id
}