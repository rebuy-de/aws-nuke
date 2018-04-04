resource "aws_vpc" "vpc" {
  cidr_block = "12.34.56.0/24"
  tags {
    Name = "AWSNukeTerraformTest"
  }
}
