# Light-Messenger

## Why

Light-Messenger helps MTRAs communicate with radiologists without calling
them on the phone. Normally, urgent cases require that the MTRA calls up the Radiologist on the phone. While this communication is
well established, it has the downside of interrupting the radiologist on duty along with disturbing the other people reading images in the same room. Instead of calling the radiologist the MTRAs can send a ambient light signal to the radiology
the department has a lower interruption cost.

## How

Light-Messenger is a web application that provides the solution described above. The MTRAs use a ui which allows them to
create notifications with three different priorities. On the Radiology side, the
doctors are presented with a list of open notifications regarding urgent cases
that need to be acknowledged.

This application also supports an arduino that can hook into the application
via a REST API and flash an LED in case there is a notification open so that
the Radiologist can acknowledge the notification.

This is a screenshot of how it looks on the MTRAs side

![Alt text](mtra.png?raw=true "MTRA screen")

This is a screenshot of how it looks for the radiologist department

![Alt text](department.png?raw=true "MTRA screen")

TODO: image of the LED flashing in the Befund room :)

## Development

The development environment uses docker to spin up a local mysql instance.
This can be replaced by a local mysql instance if required.

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

### Setup

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

### Code

- Note that error stacktraces need to be enabled _at the point of library interaction_ in the code. As an example, an error that occurs while communicating with the db needs to be wrapped with `errors.WithStack()` but this error can simply be passed along when used in the handler. The idea is to enable clean stacktraces and avoid Java-esque stacktrace recursion.

## Production

Systemd setup:

```systemd
TODO:
```

Logging:

```bash
TODO:
```
