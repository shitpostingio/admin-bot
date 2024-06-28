# Shitposting.io `admin-bot`

In order to simplify possible database upgrades, it is recommended to run the database in a Docker container.

## Docker

### Install Docker

```bash
sudo apt install docker.io
```

By default, Docker will require superuser permissions to run. To modify this behavior, we need to create the `docker` group if it doesn't already exist, add the connected `$USER` to the `docker` group and relog to apply changes:

```bash
sudo groupadd docker
sudo usermod -aG docker $USER
exit
```

### Docker configuration

Pull the latest PostgreSQL container from the official repository:

```bash
docker pull postgres:latest
```

Run the container and publish the database port to localhost and add an optional name:

```bash
docker run -p 127.0.0.1:5432:5432 --name=automod postgres
```

## PostgreSQL

In case you aren't using the PostgreSQL docker container you can install the service by using the following command:

```bash
sudo apt install postgresql postgresql-contrib
```

Log into Postgres:

```bash
psql -h localhost -U postgres
```

Create the database, the user and grant the user permissions on the table:

```sql
CREATE DATABASE automod;
CREATE USER automod WITH PASSWORD 'automod';
GRANT ALL PRIVILEGES ON DATABASE "automod" TO automod;
```

## Configuration

It is now necessary to fill in the required data in the configuration file. To do so it's possible to rename `config_example.toml` to `config.toml` and set the required values.

## Table creation

Go to the directory `database/cmd/adminbot-deploy-db`, compile the go file and run it specifying the path to the `config.toml` file:

```bash
cd database/cmd/adminbot-deploy-db
go build
./adminbot-deploy-database -config path/to/config.toml
```

## Admin creation

Since this bot will only reply to users whose Telegram ID is in the database, it is necessary to add them to the admin table (to get the Telegram ID of an user you can send a message to [@rawdatabot](https://t.me/rawdatabot)).

Go to the directory `database/cmd/adminbot-add-user`, compile the go file and run it run it specifying the path to the `config.toml` file and the `userid` to add:

```bash
cd database/cmd/adminbot-add-user
go build
./adminbot-add-user -userid id_to_add -config path/to/config.toml -role mod/admin
```
