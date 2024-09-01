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
  engine                    = "postgres"
  engine_version            = "15"
  instance_class            = "db.t3.small"
  db_name                   = var.project
  identifier                = "${var.project}-${data.aws_region.current.name}"
  vpc_security_group_ids    = [aws_security_group.rds.id]
  multi_az                  = false // Ensure single AZ deployment
  username                  = "db_admin"
  password                  = "Va1idPAssw0rd!" // Ensure the password meets AWS requirements
  allocated_storage         = 20
  storage_type              = "gp3"
  storage_encrypted         = false
  skip_final_snapshot       = true
  db_subnet_group_name      = aws_db_subnet_group.default.name
  parameter_group_name      = aws_db_parameter_group.default.name

  #   provisioner "local-exec" {
  #     command = "psql postgresql://${self.username}:${self.password}@${self.endpoint}/${self.db_name} -f data/create_db.sql"
  #   }
}

resource "aws_db_subnet_group" "default" {
  name       = "${var.project}-rds-subnet-group"
  subnet_ids = data.aws_subnets.all.ids
}

resource "aws_db_parameter_group" "default" {
  name   = "${var.project}-rds-parameter-group"
  family = "postgres15"
}