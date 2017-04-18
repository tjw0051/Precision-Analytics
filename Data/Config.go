package data

import (
    "bufio"
    "log"
    "os"
    "strings"
)

var ROOTKEY string

func init() {
    
}

func LoadConfig(filename string) {
    file, err := os.Open(filename)
    if err != nil {
        log.Fatal(err)
    }
    defer file.Close()

    scanner := bufio.NewScanner(file)
    for scanner.Scan() {
        line := strings.TrimSpace(scanner.Text())
        if(!strings.HasPrefix(line, "//")) {
            keyValue := strings.Split(line, "=")
            if(len(keyValue) >= 2) {
                key := strings.TrimSpace(keyValue[0])
                value := strings.TrimSpace(keyValue[1])
                processConfigParam(key, value)
            }
        }
    }

    if err := scanner.Err(); err != nil {
        log.Fatal(err)
    }
}

// TODO: switch to map
func processConfigParam(key string, value string) {
    switch key {
        case "root-key":
            ROOTKEY = value
    }
}