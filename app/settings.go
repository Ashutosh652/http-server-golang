package main

import "flag"

type Config struct {
	FileDirectory string
}

func ParseConfig() Config {
	var cfg Config
	flag.StringVar(&cfg.FileDirectory, "directory", "/tmp/", "Directory where the files are stored.")
	flag.Parse()
	return cfg
}
