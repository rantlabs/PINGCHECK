package main

import (
    "bufio"
    "flag"
    "fmt"
    "os"
    "os/exec"
    "strings"
    "time"
)

func ping(address string, timeout time.Duration) error {
    cmd := exec.Command("ping", "-c", "1", "-W", fmt.Sprint(int(timeout.Seconds())), address)
    err := cmd.Run()
    return err
}

func processFile(inputFile, outputFile string, timeout time.Duration) error {
    file, err := os.Open(inputFile)
    if err != nil {
        return err
    }
    defer file.Close()

    var output *os.File
    if outputFile != "" {
        output, err = os.Create(outputFile)
        if err != nil {
            return err
        }
        defer output.Close()
    }

    scanner := bufio.NewScanner(file)
    for scanner.Scan() {
        address := strings.TrimSpace(scanner.Text())
        err := ping(address, timeout)
        if err != nil {
            timestamp := time.Now().Format(time.RFC3339)
            outputLine := fmt.Sprintf("%s FAIL %s\n", address, timestamp)
            if output != nil {
                _, err := output.WriteString(outputLine)
                if err != nil {
                    return err
                }
            } else {
                fmt.Print(outputLine)
            }
        }
    }
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
