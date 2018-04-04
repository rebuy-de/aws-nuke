resource "aws_ebs_volume" "volume" {
  availability_zone = "eu-west-1a"
  size = 1

  tags {
    Name = "AWSNukeTest"
  }
}
