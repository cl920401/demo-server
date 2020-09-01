package main

import "github.com/urfave/cli"

var (
	version string
	commit  string
	branch  string
	config  string
	rpcPort int
)

var (
	// ConfigFlag config file path
	ConfigFlag = cli.StringFlag{
		Name:        "config, c",
		Usage:       "load configuration from `FILE`",
		Value:       "./configs/example_conf.json",
		Destination: &config,
	}
	// RPCPortFlag rpc监听端口
	RPCPortFlag = cli.IntFlag{
		Name:        "rpc-port",
		Usage:       "rpc service port",
		Value:       8081,
		Destination: &rpcPort,
	}
)
