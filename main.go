package main

import (
	"bufio"
	"flag"
	"fmt"
	"io/ioutil"
	"log"
	"os"
	"regexp"

	"gopkg.in/yaml.v2"
	"mvdan.cc/xurls/v2"
)

type Sed struct {
	Re      regexp.Regexp
	Replace string
}

type Config struct {
	Replace []struct {
		Match string `yaml:"match"`
		Tgt   string `yaml:"tgt"`
	}
	Blacklist []string `yaml:"blacklist"`
}

func newConf(path string) *Config {
	c := &Config{}
	if _, err := os.Stat(path); err == nil {
		f, err := ioutil.ReadFile(path)
		if err != nil {
			log.Fatal(err)
		}
		err = yaml.Unmarshal(f, c)
		if err != nil {
			log.Fatalf("Unmarshal: %v", err)
		}
	}
	return c
}

func main() {
	var out []string
	var configPath string
	var seds []Sed
	flag.StringVar(&configPath, "config", "./config.yml", "path of the config file")
	flag.Parse()
	config := newConf(configPath)
	for _, cr := range config.Replace {
		sed := Sed{}
		sed.Re = *regexp.MustCompile(cr.Match)
		sed.Replace = cr.Tgt
		seds = append(seds, sed)
	}
	rx := xurls.Relaxed()
	sc := bufio.NewScanner(os.Stdin)
	for sc.Scan() {
		text := sc.Text()
		for _, sed := range seds {
			text = sed.Re.ReplaceAllString(text, sed.Replace)
		}
		urlLoc := rx.FindAllIndex([]byte(text), -1)
		for _, loc := range urlLoc {
			// check for wildcard,... not so good code
			if loc[0] > 2 && text[loc[0]-2:loc[0]] == "*." {
				out = append(out, text[loc[0]-2:loc[1]])
			} else {
				out = append(out, text[loc[0]:loc[1]])
			}
		}
	}
	for _, s := range out {
		delete := false
		if len(config.Blacklist) > 0 {
			for _, bl := range config.Blacklist {
				if s == bl {
					delete = true
					break
				}
			}
			if !delete {
				fmt.Println(s)
			}
		} else {
			fmt.Printf("2 %s\n", s)
		}
	}
}
