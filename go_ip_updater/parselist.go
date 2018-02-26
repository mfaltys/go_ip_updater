package main

import (
	"fmt"
	"io/ioutil"
	"os"
	"strconv"
	"strings"

	"github.com/unixvoid/glogger"
)

func parseConfig(path string, config Config) Config {
	d, err := os.Open(path)
	if err != nil {
		glogger.Error.Println("error opening config")
		glogger.Error.Println(err)
		os.Exit(1)
	}
	defer d.Close()

	// open file, parse line by line
	f, _ := ioutil.ReadFile(path)

	lines := strings.Split(string(f), "\n")
	for i := range lines {
		err, field, value := parseString(lines[i])
		if err != nil {
			//glogger.Debug.Println(err)
		} else {

			switch field {
			case "configKey":
				glogger.Debug.Println("setting key")
				config.Server.Key = value
			case "configSecret":
				glogger.Debug.Println("setting secret")
				config.Server.Secret = value
			case "configZoneId":
				glogger.Debug.Println("setting zone id")
				config.Server.ZoneId = value
			case "configURL":
				glogger.Debug.Println("setting check IP URL")
				config.Server.CheckIpURL = value
			case "configTTL":
				glogger.Debug.Println("setting ttl")
				config.Server.TTL, err = strconv.ParseInt(value, 10, 64)
				if err != nil {
					glogger.Error.Println("error parsing configTTL int64")
					glogger.Error.Println(err)
					os.Exit(1)
				}
			default:
				glogger.Debug.Printf("adding domain: %s\n", field)
				config.Domains = append(config.Domains, field)
			}
		}
	}
	return config
}

func oldParseString(line string) (error, string, string) {
	//fmt.Printf("processing line '%s'\n", line)
	var (
		s      []string
		tmpStr string
		field  string
		value  string
	)
	chr := "[ ]"

	// check for comments
	if line == "" || strings.Contains(line, "#") {
		// found an empty line or comment
		return fmt.Errorf("empty line or comment"), "", ""
	} else {
		// remove whitespace
		line = strings.Replace(line, "\t", "", -1)

		// check if : exists
		if strings.Contains(line, ":") {
			// found a config line
			s = strings.Split(line, ":")
			tmpStr = s[1]
			field = s[0]
		} else {
			// its not a config line, must be a new domain
			tmpStr = line
			field = ""
		}

		value = strings.Map(func(r rune) rune {
			if strings.IndexRune(chr, r) < 0 {
				return r
			}
			return -1
		}, tmpStr)
	}
	return nil, field, value
}

func parseString(line string) (error, string, string) {
	//fmt.Printf("processing line '%s'\n", line)

	// check for comments
	if line == "" || strings.Contains(line, "#") {
		// found an empty line or comment
		return fmt.Errorf("empty line or comment"), "", ""
	} else {
		// remove whitespace
		line = strings.Replace(line, "\t", "", -1)
		line = strings.Replace(line, " ", "", -1)

		if strings.Contains(line, ":") {
			// found a config line, split around first ':'
			field := strings.Split(line, ":")[0]
			value := strings.Join(strings.Split(line, ":")[1:], ":")

			return nil, field, value
		} else {
			// its not a config line, must be a new domain
			return nil, line, ""
		}
	}
	//return nil, field, value
}
