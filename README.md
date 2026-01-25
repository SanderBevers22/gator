# Gator CLI

`gator` is a command-line RSS aggregator written in Go. It fetches RSS feeds, stores posts in a PostgreSQL database, and allows you to browse your feeds directly from the terminal.

---

## Prerequisites

Before using `gator`, make sure you have the following installed:

* **Go** (version 1.20+ recommended)
  [Download Go](https://golang.org/dl/)
* **PostgreSQL** (version 12+)
  [Download PostgreSQL](https://www.postgresql.org/download/)

---

## Installation

Clone the repository:

```bash
git clone https://github.com/your-username/gator.git
cd gator
```

Install the CLI binary:

```bash
go install .
```

This will install `gator` to your `$GOPATH/bin` (or `$HOME/go/bin` by default). Make sure that directory is in your `PATH`.

---

## Configuration

`gator` uses a simple configuration file to store your username. By default, it will create a `config.json` file in the project directory when you log in for the first time:

```bash
gator login <username>
```

Example:

```bash
gator login allan
```

You can now add feeds and fetch posts.

---

## Commands

Here are some commands to get started:

* **Register a new user**:

```bash
gator register <username>
```

* **Login as an existing user**:

```bash
gator login <username>
```

* **Add a feed**:

```bash
gator addfeed "<Feed Name>" "<Feed URL>"
```

Example:

```bash
gator addfeed "TechCrunch" "https://techcrunch.com/feed/"
```

* **Aggregate feeds** (fetch new posts continuously):

```bash
gator agg 1m
```

`1m` here is the interval between feed requests. You can also use `10s`, `5m`, etc.

* **Browse posts**:

```bash
gator browse 5
```

Shows the latest 5 posts from your feeds. If no limit is provided, defaults to 2.

* **List users**:

```bash
gator users
```

Shows all registered users and highlights the currently logged-in user.

* **Reset database**:

```bash
gator reset
```

Deletes all users, feeds, and posts. Useful for development or testing.

---

## Database Setup

Make sure your PostgreSQL server is running and accessible. Update the database connection string in `main.go` if needed:

```go
const dbURL = "postgres://postgres:postgres@localhost:5432/gator?sslmode=disable"
```

Run the migrations using Goose:

```bash
goose postgres "$DB_URL" up
```

This will create the tables for users, feeds, and posts.

---

## Development

To run the program without installing the CLI:

```bash
go run .
```

This is useful while developing. The `gator` binary is meant for production usage.

---

## Example Workflow

```bash
# Register a user
gator register allan

# Add a feed
gator addfeed "Hacker News" "https://news.ycombinator.com/rss"

# Start aggregating feeds every minute
gator agg 1m

# In another terminal, browse latest posts
gator browse 10
```

---

## Contributing

Feel free to submit PRs, open issues, or suggest improvements. This project is meant to be a simple but extensible RSS aggregator CLI in Go.

---

## License

MIT License. See [LICENSE](LICENSE) for details.

