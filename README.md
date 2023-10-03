# Reshape Go helper ![](https://github.com/leourbina/reshape-helper/actions/workflows/ci.yml/badge.svg)

This is a Go helper library for the automated, zero-downtime schema migration tool [Reshape](https://github.com/fabianlindfors/reshape). To achieve zero-downtime migrations, Reshape requires that your application runs a simple query when it opens a connection to the database to select the right schema. This library automates that process with a simple method which will return the correct query for your application.


## Installation

TBD

## Usage

The library includes a `SchemaQuery` method which will find all Reshape migration files and determine the right schema query to run. Here's an example of how to use it along with `pgxpool`:

```go
import	(
	reshapehelper "github.com/leourbina/reshape-helper"
)

func main() {
	query := reshapehelper.SchemaQuery()

	config, err := pgxpool.ParseConfig(os.Getenv("DATABASE_URL"))
	if err != nil {
		// ...
	}
	config.AfterConnect = func(ctx context.Context, conn *pgx.Conn) error {
		_, err := conn.Exec(ctx, query)
		if err != nil {
			// ...
		}
	}

	pool, err := pgxpool.ConnectConfig(context.Background(), config)
}
```

By default, `SchemaQuery` will look for migration files in `/migrations` but you can specify your own directories as well:


```go
import	(
	reshapehelper "github.com/leourbina/reshape-helper"
)

func main() {
	query := reshapehelper.SchemaQuery(
		"src/users/migrations",
		"src/todos/migrations",
	)
}
```

## License

Release under the [MIT License](https://choosealicense.com/licenses/mit/)
