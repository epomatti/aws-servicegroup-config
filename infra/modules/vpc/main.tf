resource "aws_vpc" "main" {
  cidr_block = "10.0.0.0/16"

  tags = {
    Name = "vpc-terraform"
  }
}

resource "aws_security_group" "default" {
  name   = "terraform-sg"
  vpc_id = aws_vpc.main.id

  tags = {
    Name = "sg-terraform"
  }
}

resource "aws_security_group_rule" "default" {
  type              = "ingress"
  from_port         = 0
  to_port           = 0
  protocol          = "-1"
  cidr_blocks       = ["10.0.0.0/32", "10.0.0.1/32"]
  security_group_id = aws_security_group.default.id
}
