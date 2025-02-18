# A simple file processor

This is a golang-based file upload service with asynchronous processing and database storage. Currently,
this project uploads to the local filesystem and tracks file metadata within PostgreSQL. The target destination is the `uploads` directory which is created when a user makes their first upload.

## Features

- File upload with unique naming
- Support to detect image types from file name & extension
- PostgreSQL metadata storage using GORM
- Logging with Zerolog
- Unit Tests with Mocking

## TODO

- Background processing of uploaded files (workers, queues)
- Tracking processing task statuses
- Expanding to support other file formats
- Authentication & Rate Limiting
- Documentation for APIs

## Setup

### Install Necessary Dependencies

- Golang - https://formulae.brew.sh/formula/go
- PostgreSQL - https://formulae.brew.sh/formula/postgresql@14

### Database Setup

- This project uses PostgreSQL to store file metadata information.Database credentials are obtained from environment variables. There are defaults for database information provided in the configuration file.

The environment variables that need to be configured are below. I've elected to store these in my local `~/.zshrc`

ENV | Description | 
PSQL_FILE_DATABASE_USERNAME | The username to the postgres database
PSQL_FILE_DATABASE_PASSWORD | The password to the postgres database
DB_HOST | This is the database host, if running locally, this will be "localhost"
DB_PORT | The port used to establish the database connection, by default it is configured to be 5432
DB_NAME | The name of the database which will hold all tables related to the storing metadata information about the file. The database name by default is "file_processor"

### Makefile Targets

The project's root Makefile configures run targets that are essential to building and running the project/tests. You can run each target within the Makefile by executing the following command `make <target-name>`.

The following run targets are used:

TARGET | Description
all | builds the project
build | builds the project
run | starts the server and runs the application
clean | cleans any binaries that were generated from building the project
test | runs the unit tests for the project




