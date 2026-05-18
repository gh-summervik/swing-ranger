# Swing Ranger

## Overview

This project is a premise by which I intend to become proficient with the Go language, and maybe build something useful in the process.
Beyond learning Go, my primary obective is to build a tool that helps me in swing trading stock options and which works with zero paid data services.

At this point, I have some unstructured design ideas. Nothing here is a commitment.

1. Runs as a console app - feeds items to the console at a reasonable pace.
2. Use a game engine and create an app that might resemble a video game menu system.
3. Run various data source discovery processes as background threads that dump their discoveries into queues that can be dequeued and routed to UI and/or other queues/functions.
4. Everything - data collected, analysis notes, user actions, etc. should be logged and/or catalogued for learning and debugging.
5. The system should get "smarter" over time. There should be a foundational set of rules, but the system should strive to get better at predicting short-term direction and changes.

## Technologies

This project uses the [Go programming language](https://go.dev/) and [PostgreSQL](https://www.postgresql.org/) for data storage.

## Initial Project Structure

To familiarize yourself with the project, please consider tracking down and reviewing the README.md files in the project.

```bash
├── bin
│   ├── README.md
├── database
│   └── README.md
├── docs
│   └── README.md
├── go
│   ├── README.md
└── scripts
    └── README.md
```

## Setup and Configuration

### Managing Secrets (Connection Strings and Keys)

The `.gitignore` file ignores `secrets.json` files, but these are required to make the applications work.
The application uses two connection strings - one for queries and one for commands - and these are stored in the `secrets.json` file.
You will need one of these files in the directory to which you deploy your app (maybe I'll build a deployment script later that deals with all the *other* files that have to be deployed).
I use a `/bin` directory at the same level as my `go` directory and build my apps into that directory, and therefore I keep a `secrets.json` file in there.
You will also need a `go/testdata` directory and a `secrets.json` file in it.
This directory is special to Go (i.e., it must be named `testdata`).

Here are the current locations of all my `secrets.json` files.

```bash
├── bin
│   ├── secrets.json
├── go
│   └── testdata
│       └── secrets.json
```

Following is an example `secrets.json` file.
The keys must be named "Command" and "Query."

```json
{
  "ConnectionStrings": {
    "Command": "user=sr_admin password=7C7Dxa43 host=127.0.0.1 port=5432 dbname=sr_test sslmode=disable",
    "Query": "user=sr_reader password=0e5z66Capl host=127.0.0.1 port=5432 dbname=sr_test sslmode=disable"
  }
}
```

## Tests

To run tests, from the `go` directory, run `go test ./...`.
If your `secrets.json` file is set up correctly in the `testdata` directory, the integration tests should pass.