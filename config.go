package main

import (
    "flag"
    "log"

    "gopkg.in/gcfg.v1"
)

const DefaultConfigFile = "/etc/wobchat-backend.conf"
const DefaultPort = 8000

type Config struct {
    Server struct {
        HTTPPort    int
    }
    Database struct {
        Type                    string
        ConnectionString        string
        TestConnectionString    string
    }
}

func setupConfig() (cfg Config) {
    var configFile string

    flag.StringVar(&configFile, "c", DefaultConfigFile, "Configuration file")
    flag.Parse()

    err := gcfg.ReadFileInto(&cfg, configFile)
    if err != nil {
        log.Printf("Failed to open config file %v; did you try copying the example one?\n", configFile)
        panic(err)
    }

    if cfg.Server.HTTPPort == 0 {
        cfg.Server.HTTPPort = DefaultPort
    }

    return cfg
}
