env = "production"

vpc_region = "us-east-1"

vpc_cidr = "10.10.0.0/16"

vpc_subnet_cidrs = {
  b = "10.10.0.0/17"
  d = "10.10.128.0/17"
}

webserver_instance_type = "t2.micro"

webserver_volume_size = 8

name_root = "bg-mentor"

webserver_ips = {
  b = "10.10.0.5"
  d = "10.10.128.5"
}
