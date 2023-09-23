# Go Time - AutoTasker (gt-at)

Go Time - AutoTasker (gt-at) is a CLI utility designed to simplify timesheet capture in AutoTask by Datto. Leveraging the power of playwright, gt-at makes time tracking easier for tasks (projects) and tickets (service desk) within AutoTask.

## Table of Contents

- [Purpose](#purpose)
- [Installation](#installation)
- [Usage](#usage)
- [Commands](#commands)
- [Configuration](#configuration)
- [Contributing](#contributing)
- [Importing Time Entries](#importing-time-entries)
- [Using as a Go Package](#using-as-a-go-package)
- [License](#license)

## Purpose

While AutoTask provides extensive functionality, capturing time can sometimes be a repetitive and time-consuming task. `gt-at` aims to reduce the friction involved in this process, allowing users to seamlessly capture their time for tasks and tickets.

## Installation

1. Ensure you have Golang installed.
2. Clone the repository:

```bash
git clone https://github.com/philipf/gt-at.git
```

3. Navigate to the cloned directory and build the project:

```bash
cd gt-at
go build .
```

4. Move the built binary to a location in your `$PATH` or execute directly from the current location.

## Usage

To run the utility:

```bash
gt-at [flags]
```

Or use a specific command:

```bash
gt-at [command]
```

## Commands

Here are the available commands for `gt-at`:

- **completion**: Generate the autocompletion script for the specified shell.
- **import**: Import a file of time entries into AutoTask.
- **init**: Initialise `gt-at`.
- **settings**: Prints out the settings.
- **version**: Prints the version of the application.

For detailed information on a command:

```bash
gt-at [command] --help
```



## Initial Setup

Before you start using `gt-at`, you'll need to set up a configuration file. Without this configuration, the CLI application won't run.

### Creating the Configuration

To create the configuration file:

1. Run the following command:

```bash
gt-at init
```


2. The initialization process will prompt you for various configuration details. Provide the necessary information as illustrated below:

```
Your first name and last name in AutoTask (e.g Philip Fourie):John Smith
Autotask date format, as configured in AT preferences for your Profile. Define it using [Go's Time Format Specifiers](https://pkg.go.dev/time#pkg-constants) [Default: 2006/01/02]:
Autotask day format, as shown in AT week entries when capturing Tasks. Define it using [Go's Time Format Specifiers](https://pkg.go.dev/time#pkg-constants) [Default: Mon 01/02]:
Username, typically your company email address: name@yourcompany.com
Browser type (options: chromium|firefox|webkit) [Default: chromium]:
```

3. After providing the required details, the configuration file (by default at `~/.gt-at.yaml`) will be created and initialised with your settings.

Now, with the configuration set up, you can proceed to use the `gt-at` commands as described in the subsequent sections.


## Configuration

By default, `gt-at` looks for a configuration file at `~/.gt-at.yaml`. However, you can specify a different location using the `--config` flag:

```bash
gt-at --config /path/to/config.yaml
```

Please ensure your configuration is set up correctly to interact with AutoTask. 


## Command Usage

To import a JSON file, use the `import` command followed by the `-f` or `--filename` flag:

```bash
gt-at import -f /path/to/your/time_entries.json
```

If you wish to see a summary of the time entries without importing them, you can use the `--reportOnly` or `-r` flag:

```bash
gt-at import -f /path/to/your/time_entries.json --reportOnly
```


## Importing Time Entries

Using the `import` command, you can batch import time entries from a JSON file. This feature is particularly useful when you have multiple entries to be captured at once.

### File Format

The expected JSON file format for importing time entries is an array of objects, where each object represents a time entry.

Here's the structure of a time entry object:

- **id** (Integer): The identifier for the task or ticket.
- **isTicket** (Boolean): Set to `true` if the entry is for a ticket. Set to `false` if it's for a task (project).
- **date** (String): The date for the time entry in `YYYY-MM-DDTHH:MM:SSZ` format.
- **startTimeStr** (String): The start time for the entry in `HH:MM` format.
- **duration** (Float): Duration of the time spent in hours.
- **summary** (String): A detailed summary of the time entry, often including start and end times, and any relevant notes.

### Example JSON file:

```json
[
    {
        "id": 266016,
        "isTicket": false,
        "date": "2023-09-15T00:00:00Z",
        "startTimeStr": "10:30",
        "duration": 0.75,
        "summary": "Start   End    Time   Notes\n10:30 - 11:00  00:40  10:30 Stand-up\nDuration: 0.5"
    },
    {
        "id": 266017,
        "isTicket": true,
        "date": "2023-09-16T00:00:00Z",
        "startTimeStr": "10:30",
        "duration": 0.5,
        "summary": "Start   End    Time   Notes\n10:30 - 11:00  00:40  10:30 Stand-up\nDuration: 0.5"
    }
]
```

## Using as a Go Package

If you're developing a Golang application and wish to integrate `gt-at` functionalities, you can import and use it as a package. This allows for seamless integration of time entry capture within your application logic.

### How to Use

1. First, ensure you have `gt-at` added as a dependency in your Go project. If not, you can do so with:

```bash
go get github.com/philipf/gt-at
```

2. In your Go code, import the necessary package:

```go
import "github.com/philipf/gt-at/pwplugin"
```

3. Create an instance of `AutoTaskPlaywright` and use the `CaptureTimes` method:

```go
autoTasker := pwplugin.NewAutoTaskPlaywright()

entries := // ... your time entries array
opts := // ... your options configuration

err := autoTasker.CaptureTimes(entries, opts)
```


## Contributing

Contributions to `gt-at` are welcome! If you find a bug or have a feature request, please open an issue. Pull requests are also appreciated.

## License

Please refer to the `LICENSE` file in the repository for licensing information.

---

For any further questions or feedback, please raise an issue on the GitHub repository or reach out to the maintainers.