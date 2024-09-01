output "alb_dns_name" {
  description = "The DNS name of the Application Load Balancer"
  value       = aws_lb.alb.dns_name
}

output "alb_http_url" {
  description = "The HTTP URL of the Application Load Balancer"
  value       = "http://${aws_lb.alb.dns_name}"
}

output "alb_https_url" {
  description = "The HTTPS URL of the Application Load Balancer"
  value       = "https://${aws_lb.alb.dns_name}"
}