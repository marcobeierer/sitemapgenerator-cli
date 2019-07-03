# TODO

- add verbose flag to optionally print progress during creation
	- body of each response is json with info about number of checked pages

- how to handle limit reached?
	- exit with error status code (other than 1) or just show warning as currently?

- Stats are JSON output, but cannot reliably be parsed because of other output
	- is making sure that stats are always last/first output line is enough?
