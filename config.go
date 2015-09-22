package main

import (
    "log"

    "gopkg.in/gcfg.v1"
)

type Config struct {
    Database struct {
        Type                    string
        ConnectionString        string
        TestConnectionString    string
    }
}

func setupConfig() (cfg Config) {
    err := gcfg.ReadFileInto(&cfg, "wobchat-backend.gcfg")
    if err != nil {
        log.Println("Failed to open config")
        panic(err)
    }

    return cfg
}
