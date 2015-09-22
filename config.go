package main

import (
    "log"

    "gopkg.in/gcfg.v1"
)

const ConfigFile = "/etc/wobchat-backend.conf"

type Config struct {
    Database struct {
        Type                    string
        ConnectionString        string
        TestConnectionString    string
    }
}

func setupConfig() (cfg Config) {
    err := gcfg.ReadFileInto(&cfg, ConfigFile)
    if err != nil {
        log.Printf("Failed to open config file %v; did you try copying the example one?\n", ConfigFile)
        panic(err)
    }

    return cfg
}
