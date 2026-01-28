package main

import (
	"errors"
	"flag"
	"path/filepath"

	"github.com/gw-gong/gwkit-go/hotcfg"
)

const (
	RootPath           = "../../"
	defaultCfgFilePath = "config/scanner/localcfg"
	defaultCfgFileName = "test.yaml"
)

func initFlags() (*hotcfg.LocalConfigOption, error) {
	flagCfgFilePath := flag.String("cfg_path", defaultCfgFilePath, "config file path")
	flagCfgFileName := flag.String("cfg_name", defaultCfgFileName, "config file name")

	flag.Parse()

	if *flagCfgFilePath == "" {
		return nil, errors.New("cfg_path is required")
	}
	if *flagCfgFileName == "" {
		return nil, errors.New("cfg_name is required")
	}
	return &hotcfg.LocalConfigOption{
		FilePath: filepath.Join(RootPath, *flagCfgFilePath),
		FileName: *flagCfgFileName,
		FileType: "yaml",
	}, nil
}
