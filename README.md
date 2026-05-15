# Swing Ranger

## Technologies

This project uses the Go programming language and PostgreSQL for data storage.

## Initial Project Structure

To familiarize yourself with the project, please consider tracking down and reviewing the README.md files in the project.

```bash
.
├── bin
│   └── README.md
├── database
│   └── README.md
├── docs
│   └── README.md
├── go
│   ├── README.md
├── README.md (this file)
└── scripts
    └── README.md
```

This project is a premise by which I intend to become proficient with the Go language, and maybe build something useful in the process.
Beyond learning Go, my primary obective is to build a tool that helps me in swing trading stock options and which works with zero paid data services.

At this point, I have some unstructured design ideas. Nothing here is a commitment.

1. Runs as a console app - feeds items to the console at a reasonable pace.
2. Use a game engine and create an app that might resemble a video game menu system.
3. Run various data source discovery processes as background threads that dump their discoveries into queues that can be dequeued and routed to UI and/or other queues/functions.
4. Everything - data collected, analysis notes, user actions, etc. should be logged and/or catalogued for learning and debugging.
5. The system should get "smarter" over time. There should be a foundational set of rules, but the system should strive to get better at predicting short-term direction.

---