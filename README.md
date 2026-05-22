# Swing Ranger

## Overview

This project is the path by which I am becoming proficient with the Go language.
Maybe I'll build something useful along the way.

Beyond learning Go, my primary obective is to build a tool based only on free services that aids me in swing trading stock options.
I'm just looking out for the little guy, of which I am one.

---

## Prerequisites

This project uses the [Go programming language](https://go.dev/) and [PostgreSQL](https://www.postgresql.org/) for data storage.

You will need:

1. A Go compiler for your operating system.
2. An instance of PostgreSQL.

My dev environment, as of 2026-May-21:

```bash
$ go version
go version go1.25.6 linux/amd64
$ psql --version
psql (PostgreSQL) 16.13 (Ubuntu 16.13-0ubuntu0.24.04.1)
```

---

## Getting Started

1. See prerequisites above.
2. Get your database up and running. See the `/database/README.md`.
3. Compile the CLI. From the `/go` directory, I run `go build -o ../bin/sr-cli ./cmd/sr-cli/`.
4. Get your `secrets.json` file sorted. If this means nothing to you, see the setup and configuration section later in this file.
5. Run the app. `../bin/sr-cli help` and/or `../bin/sr-cli init`.
   
You might also consider reviewing all the README.md files in the project.

As of 2026-May-21:

```bash
~/swing-ranger$ find . -name README.md
./scripts/README.md
./database/README.md
./docs/README.md
./README.md
./go/README.md
./bin/README.md
```

---

## History

You can stalk the git history yourself, but in broad strokes ...

The first thing I did was figure out how to handle CLI command-line arguments.

Then I figured out how to read in files like `secrets.json` and deserialize them into types.

I then built a simple database to house price action data and I worked through the hurdles of interacting with the database from my app.

Then I was ready for some data input. Yahoo was my first choice and it came together easily.
This feature can be used with `sr-cli -v collect MSFT,NVDA,SNDK,PLTR`.

I'm currently working on simple charts and am going to investigate from whence we might get financial data.

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
The keys **MUST** be named "Command" and "Query."

```json
{
  "ConnectionStrings": {
    "Command": "user=sr_admin password=7C7Dxa43 host=127.0.0.1 port=5432 dbname=sr_test sslmode=disable",
    "Query": "user=sr_reader password=0e5z66Capl host=127.0.0.1 port=5432 dbname=sr_test sslmode=disable"
  }
}
```

### Managing Application Configuration (Charts and such)

Application-wide configuration can be accomplished with a `config.json` file (located next to your `secrets.json` file).

If no `config.json` file is found, a default configuration is created.
To see this default configuration, see the `LoadAppConfig` function in the `./go/internal/config/config.go` file.

Here is an example, as of 2026-May-21:

```json
{
    "Chart": {
        "MovingAverages": [
            "21SC",
            "50SC",
            "200SC"
        ]
    }
}
```

#### Moving Averages

The `MovingAverages` section of the `config.json` file content uses an array of special strings.
Each string starts with a number and is followed by two characters.
The number is the moving average period and must be between 1 and 1000 inclusively.
The first character can be either "S" for simple or "E" for exponential.
The final character corresponds to the price point you want to use. Your choices are "O", "H", "L", and "C", for open, high, low, and close respectively.

## Tests

Copy the `secrets.json` file you created above to the `testdata` directory for database integration tests to pass.

To run tests, from the `go` directory, run `go test ./...`.