# light-messenger

A web application that allows MTRAs to communicate case priorities with Radiologists. The MTRAs use a ui which allows them to create notifications with a certain priority. On the other side, the radiologists are presented with a list of notifications regarding cases that need to be acknowledged.

This application also supports an arduino that can hook into the application via a REST API and flash an LED in case there is a notification open so that the Radiologist can acknowledge the notification.

TODO: screenshots

## development

The development environment uses docker to spin up a local mysql instance. This can be replaced by a local mysql instance if required.

Requirements:

- `Go >= 1.12.1`
- `mysql = 5.7.x`
- `Docker >= 19.03.2`
- `docker-compose >= 1.21.0`

Git repository structure:

- `res`: resources including sql scripts to create tables and setup integration tests
- `src`: go source code
- `static`: static web assets including the golang html templates
- `config-sample.json`: sample configuration in json format, note that the application binary requires a `config.json` located in the same folder to run
- `docker-compose.yml`: docker configuration to spin up a mysql database
- `go.mod`: go modules file
- `go.sum`: go modules checksum file
- `Jenkinsfile`: ci configuration
- `LICENSE`: license
- `light-messenger.go`: main application entry point
- `Makefile`: build tool
- `README.md`: this document
- `run-dev.sh`: development script to rebuild application on code change

### setup

Steps:

- Copy the provided `config-sample.json` to `config.json` and fill in with values as desired. Note that the integration tests currently use values provided in `config-sample.json`.
- Spin up the database in a separate shell: `docker-compose up`
- Get the rice binary via `go get github.com/GeertJohan/go.rice/rice`
- Build the application: `make build`
- Create the default tables: `./light-messenger.exec db-exec --script-path ./res/create_tables.sql`
- Run the application: `make run`

All commands:

- `make deps` : get all dependencies, uses go modules to get exact dependencies.
- `make run` : compile and execute
- `make test` : run all tests
- `make test-unit` : run unit tests
- `make test-integration` : run integration tests

To setup a local mysql instance, create a database and user: `light_messenger` and change the values in `config.json` accordingly.

To setup auto recompile on code change, use the provided `run-dev.sh` script. Note that this requires `entr` [TODO](TODO) and `ag` [TODO](TODO).

To create a production release which is a self contained binary, run the following:

```bash
TODO:
```

## production

Systemd setup:

```systemd
TODO:
```

Logging:

```bash
TODO:
```
