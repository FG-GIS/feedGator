#  Boot.Dev RSS Feed AggreGator

Creating an RSS aggregator with internal database, as a lesson for both go develpment and SQL/DB usage.

**Prerequisites**
___
The repository is written in Go, to install it, download the go package for your platform
from the official site:

```
go.dev/doc/install
```

Follow the instructions on the page relative to your platform, then run:
```
go version
```

The Database technology of choice is Postgres:

**Linux**
___
```
sudo apt install postgresql-17
```

Or use your package manager depending on your distro.

**Windows**
___  
From the official download link [Postgres](https://www.postgresql.org/download/windows/)

___

### Package

To install the gator package run:
```
go install https://github.com/FG-GIS/boot-go-gator@latest
```

After setting up postgres with credentials, generate this file to hold your configuration,
in your home directory:
`~/.gatorconfig.json`

it should contain:

```
{
  "db_url": "postgres://user:password@localhost:5432/gator?sslmode=disable",
  "current_user_name": "DatabaseUser"
}
```

## Commands
to run the program run `boot-go-gator` followed by:
- register <user>
  *adds a user to the database*
- login <user>
  *change actual database user*
- addfeed <link>
  *adds a new feed to the database*
- feeds
  *prints the feeds for this user to the console*
- agg
  *runs the main loop to gather posts from the feeds*
- browse <number>
  *prints to screen <number> most recent posts*
