# Go Directory

In my dev environment, I have a `/bin` directory at the same level as my `/go` directory.
This is where I build my outputs, as you'll see below.

```bash
$ tree -L 1
.
├── bin
├── database
├── docs
├── go
├── LICENSE
├── README.md
└── scripts
```

## sr-cli

From the `/go` directory, the command is `go build -o ../bin/sr-cli ./cmd/sr-cli/`.

## sr-gui

Placeholder for future work.