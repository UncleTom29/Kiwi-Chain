provider "aws" {
  region = "us-west-2"
}

resource "aws_instance" "blockchain_node" {
  ami           = "ami-0c94855ba95c574c8"
  instance_type = "t2.micro"

  tags = {
    Name = "blockchain-node"
  }
}
