resource "aws_instance" "example" {
  ami           = "ami-0ecb62995f68bb549" # Amazon Linux 2 (us-east-1)
  instance_type = var.instance_type

  tags = {
    Name = var.instance_name
  }
}
