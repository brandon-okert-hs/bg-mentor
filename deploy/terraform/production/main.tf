# Root level cluster terraform file. All vars are environment specific, defined in <env>.tfvars

variable "env" {
  description = "One of dev/production"
}

variable "name_root" {
  description = "A common root of the names and tags used to describe infrastructure"
}

variable "vpc_region" {}
variable "vpc_cidr" {}

variable "vpc_subnet_cidrs" {
  type = "map"
}

variable "webserver_instance_type" {}
variable "webserver_volume_size" {}

variable "webserver_ips" {
  type = "map"
}

variable "db_count" {}
variable "db_name" {}
variable "db_class" {}
variable "db_default_username" {}
variable "db_default_password" {}

terraform {
  backend "s3" {
    region = "us-east-1"
    bucket = "bg-mentor-production-tfstate"
    key    = "terraform.tfstate"
    shared_credentials_file = ".secrets-decrypted/production/terraform-aws-credentials"
    profile                 = "terraform"
  }
}

provider "aws" {
  version = "~> 0.1"
  region                  = "${var.vpc_region}"
  shared_credentials_file = ".secrets-decrypted/${var.env}/terraform-aws-credentials"
  profile                 = "terraform"
}

module "vpc" {
  source = "../modules/vpc"

  env          = "${var.env}"
  vpc_region   = "${var.vpc_region}"
  vpc_cidr     = "${var.vpc_cidr}"
  subnet_cidrs = "${var.vpc_subnet_cidrs}"
  name_root    = "${var.name_root}"

  ports = {
    ssh = 22
    http = 80
  }
}

module "webserver_b" {
  source = "../modules/webserver"

  env                = "${var.env}"
  ami_id             = "ami-cd0f5cb6"
  instance_type      = "${var.webserver_instance_type}"
  subnet_id          = "${module.vpc.subnet_ids["b"]}"
  name_root          = "${var.name_root}"
  ssh_key_name       = "${var.name_root}"
  security_group_ids = ["${module.vpc.ssh_security_group_id}", "${module.vpc.webserver_security_group_id}"]
  volume_size        = "${var.webserver_volume_size}"
  private_ip         = "${var.webserver_ips["b"]}"
}

module "database" {
  source = "../modules/webserverdb"

  security_group_ids = ["${module.vpc.db_security_group_id}"]
  name_root           = "${var.name_root}"
  db_count            = "${var.db_count}"
  env                 = "${var.env}"
  db_name             = "${var.db_name}"
  db_class            = "${var.db_class}"
  db_default_username = "${var.db_default_username}"
  db_default_password = "${var.db_default_password}"
}
