variable "db_count" {}
variable "env" {}
variable "db_name" {}
variable "db_class" {}
variable "db_default_username" {}
variable "db_default_password" {}
variable "name_root" {}
variable "security_group_ids" {
  description = "Must be ids within a vpc"
  type        = "list"
}

resource "aws_rds_cluster_instance" "webserver_db_instances" {
  count                 = 1
  identifier            = "${var.db_name}-${var.env}-${count.index}"
  db_subnet_group_name  = "${var.name_root}-${var.env}"
  cluster_identifier    = "${aws_rds_cluster.webserver_db_cluster.id}"
  instance_class        = "${var.db_class}"
}

resource "aws_rds_cluster" "webserver_db_cluster" {
  cluster_identifier      = "${var.name_root}-cluster-${var.env}"
  availability_zones      = ["us-east-1b", "us-east-1d"]
  db_subnet_group_name    = "${var.name_root}-${var.env}"
  vpc_security_group_ids  = ["${var.security_group_ids}"]
  database_name           = "${var.db_name}"
  master_username         = "${var.db_default_username}"
  master_password         = "${var.db_default_password}"
}
