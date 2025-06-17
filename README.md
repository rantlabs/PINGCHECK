# PINGCHECK

PINGCHECK is a simple ping utility that can take stdin input through unix a pipe or a unix redirect or a file with ip addresses. PINGCHECK looks at the file input and finds the first fully qualified domain name on a line and attempts to ping it once. The results are displayed to stdin or a specified file by appending to the end of the line input. -pause and -repeat do not work with STDIN. They only work with the -file option.

```
pingcheck -help
Usage of pingcheck:
  -file string
    	Input file containing IPv4 addresses or FQDNs
  -output string
    	Output file to write the results
  -pause int
    	Pause time in seconds between repeated pings if -repeat is set to 0 (only works with -file) (default 30)
  -repeat int
    	Number of cycles to ping through the list (0 for infinite, only works with -file) (default 1)
  -timeout duration
    	Timeout duration for each ping (default 1s)
  -updown
    	Show both successful and unsuccessful ping results

Note: The -repeat and -pause options only work with the -file option.
```
Examples
```
Contents of file.txt
cat file.txt 
hello world
hello world www.google.com
www.google.com hello world
99.99.99.99
```
Example 1 (unix pipe): 
```
cat file.txt| pingcheck 
99.99.99.99 FAIL 2025-06-16T16:41:14-07:00 <-- Without -updown only failures are displayed
```
Example 2 (unix pipe with -updown flag): 
```
cat file.txt| pingcheck -updown
www.google.com hello world SUCCESS 2025-06-16T16:42:12-07:00
hello world www.google.com SUCCESS 2025-06-16T16:42:12-07:00
99.99.99.99 FAIL 2025-06-16T16:42:14-07:00
```
Example 3 (unix redirect):
```
pingcheck < file.txt 
99.99.99.99 FAIL 2025-06-16T16:43:33-07:00
```
Example 4 (unix redirect with -updown flag):
```
pingcheck -updown < file.txt
www.google.com hello world SUCCESS 2025-06-16T16:44:16-07:00
hello world www.google.com SUCCESS 2025-06-16T16:44:16-07:00
99.99.99.99 FAIL 2025-06-16T16:44:18-07:00
```
Example 5 Using the -repeat and -pause flags that can be combined with using the -file option as input.
```
pingcheck -file file.txt -repeat 3 -pause 2 -updown
hello world www.google.com SUCCESS 2025-06-16T17:26:04-07:00
www.google.com hello world SUCCESS 2025-06-16T17:26:04-07:00
99.99.99.99 FAIL 2025-06-16T17:26:06-07:00
www.google.com hello world SUCCESS 2025-06-16T17:26:08-07:00
hello world www.google.com SUCCESS 2025-06-16T17:26:08-07:00
99.99.99.99 FAIL 2025-06-16T17:26:10-07:00
www.google.com hello world SUCCESS 2025-06-16T17:26:12-07:00
hello world www.google.com SUCCESS 2025-06-16T17:26:12-07:00
99.99.99.99 FAIL 2025-06-16T17:26:14-07:00
```
Binary versions available for download in this repo
```
pingcheck_linux32	pingcheck_rpi_arm64	pingcheck_win32.exe
pingcheck_linux64	pingcheck_rpi_armv6	pingcheck_win64.exe
pingcheck_mac_arm64	pingcheck_rpi_armv7
```
Basic instructions for compiling the GO source code available in the the notes.txt file
