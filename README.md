# PINGCHECK

PINGCHECK is a simple ping utility that can take stdin input through unix a pipe or a unix redirect or a file with ip addresses. PINGCHECK looks at the file input and finds the first fully qualified domain name on a line and attempts to ping it once. The results are displayed to stdin or a specified file by appending to the end of the line input. 

```
pingcheck 
Error: no input file specified and no data piped in
Usage of pingcheck:
  -file string
    	Input file containing IPv4 addresses or FQDNs
  -output string
    	Output file to write the results
  -pause int
    	Pause time in seconds between repeated pings if -repeat is set to 0 (default 30)
  -repeat int
    	Number of cycles to ping through the list (0 for infinite) (default 1)
  -timeout duration
    	Timeout duration for each ping (default 1s)
  -updown
    	Show both successful and unsuccessful ping results
```
