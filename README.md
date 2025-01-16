# Gator 🐊

A powerful command-line RSS feed aggregator written in Go. Stay up to date with your favorite content sources through a simple CLI interface.

## Features

- 👤 User account management
- 📰 Add and follow RSS feeds
- 🔄 Automatic content aggregation
- 📱 Browse posts from followed feeds
- ⏰ Scheduled feed updates
- 🔐 Authentication system
- 🗃️ PostgreSQL storage

## Still ToDo:

- Improve logging & error messages

## Prerequisites

Before installing Gator, you need:

1. **Go** (version 1.23.4 or later)
   - Install from [go.dev/dl](https://go.dev/dl/)
   - Verify installation with `go version`

2. **PostgreSQL** (version 14 or later)
   - Install from [postgresql.org/download](https://postgresql.org/download/)
   - Create a database for Gator: `createdb gator`

## Installation

Install the Gator CLI directly using Go:

```bash
go install github.com/jradziejewski/gator@latest
```

This will install the `gator` command in your `$GOPATH/bin` directory. Make sure this directory is in your system's `PATH`.

## Configuration

Create a file named `.gatorconfig.json` in your home directory:

```json
{
    "db_url": "postgresql://username:password@localhost:5432/gator?sslmode=disable",
    "current_user_name": ""
}
```

Replace `username` and `password` with your PostgreSQL credentials.

## Quick Start

1. Create a new user account:
```bash
gator register johndoe
```

2. Add an RSS feed:
```bash
gator addfeed "Hacker News" "https://news.ycombinator.com/rss"
```

3. View your feeds:
```bash
gator feeds
```

4. Start aggregating:
   ```bash
   gator agg 3600s
   ```

5. Browse recent posts:
```bash
gator browse
```

## Available Commands

### User Management
- `register <username>`: Create a new user account
- `login <username>`: Log in as an existing user
- `users`: List all registered users
- `reset`: Reset the database (delete all users and feeds)

### Feed Management
- `addfeed <name> <url>`: Add a new RSS feed
- `feeds`: List all available feeds
- `follow <url>`: Follow an existing feed
- `following`: List feeds you're following
- `unfollow <url>`: Unfollow a feed

### Content Browsing
- `browse [limit]`: View recent posts from followed feeds (default limit: 2)
- `agg <time_between_reqs>`: Start the feed aggregator with specified interval

## Database Schema

### Users
```sql
create table users(
    id uuid primary key,
    created_at timestamp not null,
    updated_at timestamp not null,
    name varchar unique not null
);
```

### Feeds
```sql
create table feeds(
    id uuid primary key,
    created_at timestamp not null,
    updated_at timestamp not null,
    name varchar not null,
    url varchar unique not null,
    user_id uuid not null,
    last_fetched_at timestamp,
    foreign key (user_id) references users(id) on delete cascade
);
```

### Feed Follows
```sql
create table feed_follows(
    id uuid primary key,
    created_at timestamp not null,
    updated_at timestamp not null,
    user_id uuid not null,
    feed_id uuid not null,
    foreign key (user_id) references users(id) on delete cascade,
    foreign key (feed_id) references feeds(id) on delete cascade,
    constraint follow unique(user_id, feed_id)
);
```

### Posts
```sql
create table posts(
    id uuid primary key,
    created_at timestamp not null,
    updated_at timestamp not null,
    title varchar not null,
    url varchar unique not null,
    description varchar not null,
    published_at timestamp not null,
    feed_id uuid not null,
    foreign key (feed_id) references feeds(id) on delete cascade
);
```

## Architecture

### Core Components

1. **Main Package**
   - Handles command-line interface
   - Manages application state
   - Registers and executes commands
   - Implements middleware

2. **Database Layer**
   - Uses SQLC for type-safe SQL queries
   - Manages database connections and operations
   - Implements all database interactions

3. **RSS Processing**
   - Handles RSS feed fetching and parsing
   - Supports multiple date formats
   - Manages HTML entity decoding

4. **Configuration**
   - Manages user configuration
   - Handles persistent settings

## Dependencies

- `github.com/google/uuid` v1.6.0: UUID generation
- `github.com/lib/pq` v1.10.9: PostgreSQL driver
- SQLC: SQL query generation
- Goose: Database migrations

## Development

The project uses:
- SQLC for generating type-safe database queries
- Goose for database migrations
- Go modules for dependency management

## Contributing

Contributions are welcome! Please feel free to fork this repo & submit a Pull Request.
