# bg-mentor

Webserver, client interface, admin interface, and db for born gosu mentor-mentee mactching service

[Changelog](CHANGELOG.md)

## Usage

TBD

## Planned Features

TBD

## Maybe Features

TBD

## Development

### Setup

Install [yarn](https://yarnpkg.com/en/)
```bash
$ yarn --version
1.9.4
```

Install [nvm](https://github.com/creationix/nvm#install-script) or manually install Node 10
```bash
$ nvm --version
0.33.4
$ node --version
v10.0.0
```

Install [go](https://golang.org/doc/install)
```bash
$ go version
go version go1.10.1 darwin/amd64
```

Install [terraform](https://www.terraform.io/)
```bash
$ terraform --version
Terraform v0.11.8
```

Install [ansible vault](https://docs.ansible.com/ansible/2.6/installation_guide/intro_installation.html#installing-the-control-machine)
I just used brew to install the latest version of ansible, which comes with ansible-vault.
```bash
$ ansible-vault --version
ansible-vault 2.5.0
```

Install [flyway](https://flywaydb.org/)
I just brew installed the latest version
```bash
$ flyway -v
Flyway Community Edition 5.2.1 by Boxfuse
```

### Building

```bash
# After Pulling:
nvm use
make deps

# Build webserver only:
make build-webserver env=dev # or env=production

# Build frontend only:
make build-frontend env=dev # or env=production

# Build everything from scratch:
make build env=dev # or env=production

# Watch frontend:
make watch-frontend env=dev

# Run once built. Can run in parallel with make watch-frontend
make run env=dev # or env=production
```

## Deployment

### Apply infrastructure changes to dev or production

```bash
$ make plan-infra env=dev # or env=production
$ make apply-infra env=dev # or env=production
```

### Deploy to dev or production

```bash
$ make build env=dev # or env=production
$ make deploy env=dev # or env=production
```

### Manual DB provisioning

Almost everything is automated via terraform and deploys. However, after terraforming, the admin password needs to be changed, and the app users and migrations users need to be created.

First run the tunnel
`make tunnel-db env=production`

Then change the admin password
```sql
set password for admin@'%' = PASSWORD('<password>')
```

Then create the dev db and users

```sql
create database borngosu_dev;

create user 'bguser'@'%' identified by '<password>';
grant INSERT,SELECT,UPDATE,DELETE on borngosu.* to 'bguser'@'%';

create user 'bguser-dev'@'%' identified by '<password>';
grant INSERT,SELECT,UPDATE,DELETE on borngosu_dev.* to 'bguser-dev'@'%';

create user 'migrations'@'%' identified by '<password>';
grant INSERT,SELECT,UPDATE,DELETE,CREATE,DROP,INDEX,ALTER,LOCK TABLES on borngosu.* to 'migrations'@'%';

create user 'migrations-dev'@'%' identified by '<password>';
grant INSERT,SELECT,UPDATE,DELETE,CREATE,DROP,INDEX,ALTER,LOCK TABLES on borngosu_dev.* to 'migrations-dev'@'%';
```
