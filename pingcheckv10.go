package main

import (
    "bufio"
    "flag"
    "fmt"
    "os"
    "os/exec"
    "regexp"
    "sync"
    "time"
)

// Define flags
var updownFlag = flag.Bool("updown", false, "Show both successful and unsuccessful ping results")
var pauseFlag = flag.Int("pause", 30, "Pause time in seconds between repeated pings if -repeat is set to 0 (only works with -file)")
var repeatFlag = flag.Int("repeat", 1, "Number of cycles to ping through the list (0 for infinite, only works with -file)")
var inputFile = flag.String("file", "", "Input file containing IPv4 addresses or FQDNs")
var outputFile = flag.String("output", "", "Output file to write the results")
var timeout = flag.Duration("timeout", time.Second, "Timeout duration for each ping")

// Function to ping the address and handle output
func ping(line string, address string, timeout time.Duration, outputFile *os.File, wg *sync.WaitGroup, semaphore chan struct{}) {
    defer wg.Done()
    defer func() { <-semaphore }() // Release the semaphore slot when done

    cmd := exec.Command("ping", "-c", "1", "-W", fmt.Sprint(int(timeout.Seconds())), address)
    err := cmd.Run()
    timestamp := time.Now().Format(time.RFC3339)
    outputLine := ""
    if err != nil {
        outputLine = fmt.Sprintf("%s FAIL %s\n", line, timestamp)
    } else if *updownFlag {
        outputLine = fmt.Sprintf("%s SUCCESS %s\n", line, timestamp)
    }

    if outputLine != "" {
        if outputFile != nil {
            _, err := outputFile.WriteString(outputLine)
            if err != nil {
                fmt.Printf("Error writing to output file: %v\n", err)
            }
        } else {
            fmt.Print(outputLine)
        }
    }
}

func processInput(scanner *bufio.Scanner, outputFilePath string, timeout time.Duration) error {
    var outputFile *os.File
    var err error

    if outputFilePath != "" {
        outputFile, err = os.OpenFile(outputFilePath, os.O_APPEND|os.O_CREATE|os.O_WRONLY, 0644)
        if err != nil {
            return err
        }
        defer outputFile.Close()
    }

    var wg sync.WaitGroup
    semaphore := make(chan struct{}, 4) // Limit to 4 concurrent pings
    ipRegex := regexp.MustCompile(`(?:\d{1,3}\.){3}\d{1,3}|[a-zA-Z0-9.-]+\.[a-zA-Z]{2,6}`)

    for scanner.Scan() {
        line := scanner.Text()
        match := ipRegex.FindString(line)
        if match != "" {
            semaphore <- struct{}{} // Acquire a semaphore slot
            wg.Add(1)
            go ping(line, match, timeout, outputFile, &wg, semaphore)
        }
    }
    wg.Wait()

    return scanner.Err()
}

func main() {
    flag.Parse()

    var scanner *bufio.Scanner
    fileInfo, err := os.Stdin.Stat()
    if err != nil {
        fmt.Println("Error accessing stdin:", err)
        os.Exit(1)
    }

    cycle := 0
    if fileInfo.Mode()&os.ModeCharDevice == 0 {
        // Data is being piped or redirected via stdin
        scanner = bufio.NewScanner(os.Stdin)
    } else {
        // Input is provided via file
        if *inputFile == "" {
            fmt.Println("Error: no input file specified and no data piped in")
            flag.Usage()
            os.Exit(1)
        }
    }

    if *inputFile != "" {
        // Handle file input
        for *repeatFlag == 0 || cycle < *repeatFlag {
            file, err := os.Open(*inputFile)
            if err != nil {
                fmt.Printf("Error opening input file: %v\n", err)
                os.Exit(1)
            }
            defer file.Close()
            scanner = bufio.NewScanner(file)

            err = processInput(scanner, *outputFile, *timeout)
            if err != nil {
                fmt.Printf("Error processing input: %v\n", err)
                os.Exit(1)
            }

            cycle++
            if *repeatFlag != 0 && cycle >= *repeatFlag {
                break
            }

            if *repeatFlag == 0 || cycle < *repeatFlag {
                time.Sleep(time.Duration(*pauseFlag) * time.Second)
            }
        }
    } else {
        // Handle stdin input
        err = processInput(scanner, *outputFile, *timeout)
        if err != nil {
            fmt.Printf("Error processing input: %v\n", err)
            os.Exit(1)
        }
    }
}

// Customized usage function to specify repeat and pause restriction
func init() {
    flag.Usage = func() {
        fmt.Fprintf(os.Stderr, "Usage of %s:\n", os.Args[0])
        flag.PrintDefaults()
        fmt.Fprintf(os.Stderr, "\nNote: The -repeat and -pause options only work with the -file option.\n")
    }
}
