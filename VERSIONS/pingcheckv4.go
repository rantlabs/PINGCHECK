package main

import (
    "bufio"
    "flag"
    "fmt"
    "os"
    "os/exec"
    "strings"
    "sync"
    "time"
)

// Function to ping the address and handle output
func ping(address string, timeout time.Duration, outputFile *os.File, wg *sync.WaitGroup, semaphore chan struct{}) {
    defer wg.Done()
    defer func() { <-semaphore }() // Release the semaphore slot when done

    cmd := exec.Command("ping", "-c", "1", "-W", fmt.Sprint(int(timeout.Seconds())), address)
    err := cmd.Run()
    if err != nil {
        timestamp := time.Now().Format(time.RFC3339)
        outputLine := fmt.Sprintf("%s FAIL %s\n", address, timestamp)
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

func processFile(inputFile, outputFilePath string, timeout time.Duration) error {
    file, err := os.Open(inputFile)
    if err != nil {
        return err
    }
    defer file.Close()

    var outputFile *os.File
    if outputFilePath != "" {
        outputFile, err = os.OpenFile(outputFilePath, os.O_APPEND|os.O_CREATE|os.O_WRONLY, 0644)
        if err != nil {
            return err
        }
        defer outputFile.Close()
    }

    scanner := bufio.NewScanner(file)
    var wg sync.WaitGroup
    semaphore := make(chan struct{}, 4) // Limit to 4 concurrent pings

    for scanner.Scan() {
        address := strings.TrimSpace(scanner.Text())
        semaphore <- struct{}{} // Acquire a semaphore slot
        wg.Add(1)
        go ping(address, timeout, outputFile, &wg, semaphore)
    }
    wg.Wait()

    return scanner.Err()
}

func main() {
    inputFile := flag.String("file", "", "Input file containing IPv4 addresses or FQDNs")
    outputFile := flag.String("output", "", "Output file to write the results")
    timeout := flag.Duration("timeout", time.Second, "Timeout duration for each ping")
    flag.Parse()

    if *inputFile == "" {
        fmt.Println("Error: input file is required")
        flag.Usage()
        os.Exit(1)
    }

    err := processFile(*inputFile, *outputFile, *timeout)
    if err != nil {
        fmt.Printf("Error processing file: %v\n", err)
        os.Exit(1)
    }
}
