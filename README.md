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
