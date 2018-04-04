resource "aws_subnet" "subnet" {
  vpc_id     = "${aws_vpc.dep_vpc.id}"
  cidr_block = "10.0.0.0/24"

  tags {
    Name = "AWSNukeTest"
  }
}

resource "aws_vpc" "dep_vpc" {
  cidr_block = "10.0.0.0/16"

  tags {
    Name = "AWSNukeTest"
  }
}
