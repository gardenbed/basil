[![Go Doc][godoc-image]][godoc-url]

# Config

config is a minimal and unopinionated utility for reading configuration values based on [The 12-Factor App](https://12factor.net/config).

The idea is that you just define a `struct` with _fields_, _types_, and _defaults_,
then you just use this library to read and populate the values for fields of your struct
from either **command-line flags**, **environment variables**, or **configuration files**.
It can also watch for new values read from _configuration files_ and notify subscribers.

This library does not use `flag` package for parsing flags, so you can still parse your flags separately.

#### Quick Start

You can find examples [here](./example).

| Example | Description | asciinema |
|---------|-------------|-----------|
| Basic | Using `Pick` method | [Watch](https://asciinema.org/a/296495) |
| Watch | Watching for new configurations | |
| Kubernetes | Changing the *log level* without restarting the *pods* | [Watch](https://asciinema.org/a/296509) |
| Telepresence | Running an application in a [Telepresence](https://telepresence.io) session | [Watch](https://asciinema.org/a/296511) |

#### Supported Types

  - `string`, `*string`, `[]string`
  - `bool`, `*bool`, `[]bool`
  - `float32`, `float64`
  - `*float32`, `*float64`
  - `[]float32`, `[]float64`
  - `int`, `int8`, `int16`, `int32`, `int64`
  - `*int`, `*int8`, `*int16`, `*int32`, `*int64`
  - `[]int`, `[]int8`, `[]int16`, `[]int32`, `[]int64`
  - `uint`, `uint8`, `uint16`, `uint32`, `uint64`
  - `*uint`, `*uint8`, `*uint16`, `*uint32`, `*uint64`
  - `[]uint`, `[]uint8`, `[]uint16`, `[]uint32`, `[]uint64`
  - `url.URL`, `*url.URL`, `[]url.URL`
  - `regexp.Regexp`, `*regexp.Regexp`, `[]regexp.Regexp`
  - `time.Duration`, `*time.Duration`, `[]time.Duration`

#### Behaviour

The precedence of sources for reading values is as follows:

  1. command-line flags
  1. environment variables
  1. configuration files
  1. default values (set when creating the instance)

You can pass the configuration values with **flags** using any of the syntaxes below:

```bash
main  -enabled  -log.level info  -timeout 30s  -address http://localhost:8080  -endpoints url1,url2,url3
main  -enabled  -log.level=info  -timeout=30s  -address=http://localhost:8080  -endpoints=url1,url2,url3
main --enabled --log.level info --timeout 30s --address http://localhost:8080 --endpoints url1,url2,url3
main --enabled --log.level=info --timeout=30s --address=http://localhost:8080 --endpoints=url1,url2,url3
```

You can pass the configuration values using **environment variables** as follows:

```bash
export ENABLED=true
export LOG_LEVEL=info
export TIMEOUT=30s
export ADDRESS=http://localhost:8080
export ENDPOINTS=url1,url2,url3
```

You can also write the configuration values in **files** (or mount your configuration values and secrets as files)
and pass the paths to the files using environment variables:

```bash
export ENABLED_FILE=...
export LOG_LEVEL_FILE=...
export TIMEOUT_FILE=...
export ADDRESS_FILE=...
export ENDPOINTS_FILE=...
```

The supported syntax for Regexp is [POSIX Regular Expressions](https://en.wikibooks.org/wiki/Regular_Expressions/POSIX_Basic_Regular_Expressions).

#### Skipping

If you want to skip a source for reading values, use `-` as follows:

```go
type Config struct {
  GithubToken string `env:"-" fileenv:"-"`
}
```

In the example above, `GithubToken` can only be set using `github.token` command-line flag.

#### Customization

You can use Go _struct tags_ to customize the name of expected command-line flags or environment variables.

```go
type Config struct {
  Database string `flag:"config.database" env:"CONFIG_DATABASE" fileenv:"CONFIG_DATABASE_FILE_PATH"`
}
```

In the example above, `Database` will be read from either:

  1. The command-line flag `config.databas`
  2. The environment variable `CONFIG_DATABASE`
  3. The file specified by environment variable `CONFIG_DATABASE_FILE_PATH`
  4. The default value set on struct instance

#### Using `flag` Package

`config` plays nice with `flag` package since it does NOT use `flag` package for parsing command-line flags.
That means you can define, parse, and use your own flags using built-in `flag` package.

If you use `flag` package, `config` will also add the command-line flags it is expecting.
Here is an example:

```go
package main

import (
  "flag"
  "time"

  "github.com/gardenbed/basil/config"
)

var config = struct {
  Enabled   bool
  LogLevel  string
} {
  Enabled:  true,   // default
  LogLevel: "info", // default
}

func main() {
  config.Pick(&config)
  flag.Parse()
}
```

If you run this example with `-help` or `--help` flag,
you will see `-enabled` and `-log.level` flags are also added with descriptions!

#### Options

Options are helpers for specific situations and setups.
You can pass a list of options to `Pick` and `Watch` methods.
If you want to test or debug something and you don't want to make code changes, you can set options through environment variables as well.

| Option | Environment Variable | Description |
|--------|----------------------|-------------|
| `config.Debug()` | `CONFIG_DEBUG` | Printing debugging information. |
| `config.ListSep()` | `CONFIG_LIST_SEP` | Specifying list separator for all fields with slice type. |
| `config.SkipFlag()` | `CONFIG_SKIP_FLAG` | Skipping command-line flags as a source for all fields. |
| `config.SkipEnv()` | `CONFIG_SKIP_ENV` | Skipping environment variables as a source for all fields .|
| `config.SkipFileEnv()` | `CONFIG_SKIP_FILE_ENV` | Skipping file environment variables (and configuration files) as a source for all fields. |
| `config.PrefixFlag()` | `CONFIG_PREFIX_FLAG` | Prefixing all flag names with a string. |
| `config.PrefixEnv()` | `CONFIG_PREFIX_ENV` | Prefixing all environment variable names with a string. |
| `config.PrefixFileEnv()` | `CONFIG_PREFIX_FILE_ENV` | Prefixing all file environment variable names with a string. |
| `config.Telepresence()` | `CONFIG_TELEPRESENCE` | Reading configuration files in a _Telepresence_ environment. |

#### Debugging

If for any reason the configuration values are not read as you expected, you can view the debugging logs.
You can enable debugging logs either by using `Debug` option or by setting `CONFIG_DEBUG` environment variable.
In both cases you need to specify a verbosity level for logs.

| Level | Descriptions                                               |
|-------|------------------------------------------------------------|
| `0`   | No logging (default).                                      |
| `1`   | Logging all errors.                                        |
| `2`   | Logging initialization information.                        |
| `3`   | Logging information related to new values read from files. |
| `4`   | Logging information related to notifying subscribers.      |
| `5`   | Logging information related to setting values of fields.   |
| `6`   | Logging miscellaneous information.                         |

#### Watching

config allows you to watch _configuration files_ and dynamically update your configurations as your application is running.

When using `Watch()` method, your struct should have a `sync.Mutex` field on it for synchronization and preventing data races.
You can find an example of using `Watch()` method [here](https://github.com/gardenbed/basil/tree/main/config/example/3-watch).

[Here](https://milad.dev/posts/dynamic-config-secret) you will find a real-world example of using `config.Watch()`
for **dynamic configuration management** and **secret injection** for Go applications running in Kubernetes.


[godoc-url]: https://pkg.go.dev/github.com/gardenbed/basil/config
[godoc-image]: https://pkg.go.dev/badge/github.com/gardenbed/basil/config
