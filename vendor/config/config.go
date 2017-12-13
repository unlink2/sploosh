package config

import (
  "github.com/go-ini/ini"
  "log"
)

// Globalcfg - global config object
var Globalcfg = ini.Empty()

/*
ReadConfig - reads an ini file
file - path to ini file
*/
func ReadConfig(file string) *ini.File  {
  cfg, err := ini.LoadSources(ini.LoadOptions{}, file)

  if err != nil {
    log.Fatal(err)
  }

  return cfg
}
