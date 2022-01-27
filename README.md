# schedule-events

Schedule-events is a client/server tool for scheduling pre-defined commands as
events. Commands are defined on the server, clients can schedule the execution
of these commands as one-shot or periodic events on the server.

Note: There are no security checks. This should only be used in a trusted
environment, e.g., for testing purposes.

## Usage

You can run `schedule-events` with the following command line arguments:

```
Usage of schedule-events:
  -address addr
        listen on or connect to addr (default "localhost:8080")
  -commands file
        read commands from file (default "commands.json")
  -events file
        read events from file (default "events.json")
  -operation operation
        run operation on server (default "get-events")
  -server
        run as server
```

Operations:
* `get-commands`: get specific or a list of all commands from the server
* `get-events`: get specific or a list of all events from the server
* `set-events`: schedule specific events on the server
* `delete-events`: stop and remove specific events from the server
* `get-status`: get status of the server
* `shutdown`: shutdown the server
* `stop`: stop all events on the server

Specific commands or events can be specified with json files and the command
line parameters `-commands` and `-events`.

## Examples

Running a server on local host and port `8081` with command definitions in
`my-commands.json` and pre-defined events in `my-events.json`:

```console
$ schedule-events \
        -commands my-commands.json \
        -events my-events.json \
        -address localhost:8081 \
        -server
```

Scheduling events in the file `more-events.json` on the server:

```console
$ schedule-events \
        -events more-events.json \
        -address :8081 \
        -operation set-events
```

Example json command list used with the command line argument `-commands`:

```json
[
	{
		"Name":"ls",
		"Executable":"ls",
		"Timeout": 10000000000
	},
	{
		"Name":"date",
		"Executable":"date",
		"Arguments":["--rfc-3339=second"],
		"Timeout": 10000000000
	}
]
```

Example json event list used with the command line argument `-events`:

```json
[
	{
		"Name":"ls1",
		"Command":"ls"
	},
	{
		"Name":"date1",
		"Command":"date"
	},
	{
		"Name":"date-periodic1",
		"Command":"date",
		"Periodic":true,
		"WaitMin":10000000000
	},
	{
		"Name":"date-periodic2",
		"Command":"date",
		"Periodic":true,
		"WaitMin":1000000000,
		"WaitMax":10000000000
	}
]
```

Example json event list for deleting the events above with the command line
argument `-operation delete-events`:

```json
[
	{
		"Name":"ls1"
	},
	{
		"Name":"date1"
	},
	{
		"Name":"date-periodic1"
	},
	{
		"Name":"date-periodic2"
	}
]
```
