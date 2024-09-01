// create security group for rds

resource "aws_security_group" "rds" {
  name        = "${var.project}-rds-sg"
  description = "Security group for RDS instance"
  vpc_id      = data.aws_vpc.default.id

  ingress {
    from_port   = 5432
    to_port     = 5432
    protocol    = "tcp"
    cidr_blocks = ["0.0.0.0/0"] // Adjust as needed for your security requirements
  }

  egress {
    from_port   = 0
    to_port     = 0
    protocol    = "-1"
    cidr_blocks = ["0.0.0.0/0"]
  }
}

// create rds instance

resource "aws_db_instance" "rds" {
  engine                 = "postgres"
  engine_version         = "15"
  instance_class         = "db.t3.small"
  db_name                = "${var.project}-dev"
  username               = "admin"
  password               = "password"
  vpc_security_group_ids = [aws_security_group.rds.id]
  multi_az               = false // Ensure single AZ deployment

  provisioner "local-exec" {
    command = "psql postgresql://${self.username}:${self.password}@${self.endpoint}/${self.db_name} -f data/create_db.sql"
  }
}