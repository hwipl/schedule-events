# schedule-events

Schedule events

Example json command list:

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

Example json event list:

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
