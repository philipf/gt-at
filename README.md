[![Build and Test](https://github.com/philipf/gt-at/actions/workflows/go.yaml/badge.svg)](https://github.com/philipf/gt-at/actions/workflows/go.yaml)
[![Releases](https://github.com/philipf/gt-at/actions/workflows/release.yaml/badge.svg)](https://github.com/philipf/gt-at/actions/workflows/release.yaml)

# gt-at


## Purpose
I have been losing precious time by dealing with the slow and cumbersome AutoTask (Datto / Kaseya) web interface, `gt-at` solves this problem by allowing you to capture time outside of AutoTask and then import it in bulk using the CLI tool or the Go package (SDK). 

It uses Playwright to automate the browser and capture the time entries.

Why not use the AutoTask REST API? Firstly, I need access to the API key in my organisation, and secondly, the key provides full system administrator access, which I prefer not to have.

### Features
- Import time entries from a JSON file
- Capture time entries from your application using the Go package.
- Login to AutoTask using Azure AD / Entra ID.
- Supports both tickets (service desk) and tasks (projects).
- Supports both Chromium, Firefox and Webkit browsers.
- Supports both Windows, Linux, and macOS operating systems.

### Known limitations
- Entries cannot be deleted; this is by design. If you need to delete an entry, do it manually in AutoTask.
- If a timesheet is already submitted for the week, `gt-at` will give up. You can recall the submission in AutoTask and try again.

### Disclaimers
- This project is not affiliated with AutoTask or Datto in any way. It is a personal project that I use to make my life easier. I hope it can do the same for you.
- I have only tested this on my AutoTask account using Windows 10 and the chromium driver. It may not work for you. If you find any bugs, please raise an issue or even better, a pull request.


## Installation
Choose one of the following installation methods that best suits your needs.

### From Releases
Download the latest version for your operation system from [Releases](https://github.com/philipf/gt-at/releases)

### Go Install
1. Ensure you have Golang installed.
2. Run the following command:

```bash
go install github.com/philipf/gt-at@latest
```

### Build from source

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


2. The initialisation process will prompt you for various configuration details. 
Provide the necessary information as illustrated below:

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


## Importing Time Entries using the CLI and JSON

You can batch import time entries from a JSON file using the `import` command. .

### File Format

The expected JSON file format for importing time entries is an array of objects, each representing a time entry.

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

## AutoTask IDs

You'll notice the `id` field in the JSON and SDK, this refers to either a Task or Ticket ID in AutoTask. You can find this ID in the URL when viewing the Task or Ticket in AutoTask.

This means you must manually create the Task or Ticket in AutoTask before you can import time entries for it. It also implies that you have to keep track of the IDs yourself.

## Tech Stack
- Go 1.21.1
- [Playwright](https://playwright.dev/)
- [Playwright for Go](https://github.com/playwright-community/playwright-go)
- [Cobra](https://github.com/spf13/cobra)
- [Viper](https://github.com/spf13/viper)
- [GoReleaser](https://goreleaser.com/)
- [ASCII Table Writer](github.com/olekukonko/tablewriter)


## Contributing

Contributions to `gt-at` are welcome! If you find a bug or have a feature request, please open an issue. Pull requests are also appreciated.

## License

Please refer to the `LICENSE` file in the repository for licensing information.

---

For any further questions or feedback, please raise an issue on the GitHub repository or reach out to the maintainers.