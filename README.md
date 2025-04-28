# A Simple File Processor

This is a go based file upload service with asynchronous processing and database storage of file metadata. Currently, the project uploads files to the local filesystem and tracks file metadata in PostgreSQL. Uploaded files are stored within the `uploads` directory, which is automatically created when a user makes their first upload

## Features

- File upload with unique naming to avoid collisions
- Background processing for uploaded files. Supports the following tasks
    - Metadata Extraction for Videos using ffmpeg
    - Image Resizing
- Image Type Detection based on file name and extension
- PostgreSQL Metadata Storage using GORM
- Structured Logging with Zerolog
- Unit Testing with mocking support using mockery

## TODO

- Track processing task statuses
- Add authentication and rate limiting
- Write comprehensive API documentation
- Set up CI pipeline using GHA
- Dockerize Service

## Setup

### Install Necessary Dependencies

- Golang - https://formulae.brew.sh/formula/go
- PostgreSQL - https://formulae.brew.sh/formula/postgresql@14

### Database Setup

This project uses PostgreSQL to store file metadata information. Database credentials are obtained from environment variables, with defaults provided in the configuration file. I typically store these in my local `~/.zshrc`.


| Env  | Description |
| ------------- | ------------- |
| PSQL_FILE_DATABASE_USERNAME  | The username to the postgres database |
| PSQL_FILE_DATABASE_PASSWORD  | The password to the postgres database |
| DB_HOST | This is the database host, if running locally, this will be "localhost" |
| DB_PORT | The port used to establish the database connection, by default it is configured to be 5432 |
| DB_NAME | The name of the database which will hold all tables related to the storing metadata information about the file. The database name by default is "file_processor" |

### Redis Setup

This project uses Redis as a message broker to hold background job information. Background jobs are created as tasks and utilize the asynq library. Upon application startup, a background job server is launched in a separate thread.

- Asnyq: https://github.com/hibiken/asynq

Please install redis through your package manager and launch redis in the background. The default configuration for the redis server is defined within configuration.json

### Makefile Targets

The project's root Makefile configures run targets that are essential to building and running the project/tests. You can run each target within the Makefile by executing the following command `make <target-name>`.

The following run targets are used:

| Target | Description |
| -------- | ----------|
|all | builds the project |
|build | builds the project |
|run | starts the server and runs the application |
|clean | cleans any binaries that were generated from building the project |
|test | runs the unit tests for the project |

### Running the application

This section covers running the application and starting the server. By default, the service is configured to listen on port 8080, but can be customized by the APP_PORT environment variable.

1. Clean the project

```
make clean
```

2. Build the project

```
make build
```

3. Run the project & start the server

```
make run
```

To use the APIs, follow the API documentation within [API.md](/API.md)



