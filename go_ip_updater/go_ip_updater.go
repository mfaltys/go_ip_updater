package main

import (
	"bytes"
	"flag"
	"fmt"
	"io/ioutil"
	"net"
	"net/http"
	"os"
	"sync"
	"time"

	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/aws/credentials"
	"github.com/aws/aws-sdk-go/aws/session"
	"github.com/aws/aws-sdk-go/service/route53"
	"github.com/unixvoid/glogger"
)

type Config struct {
	Server struct {
		Key        string
		Secret     string
		TTL        int64
		ZoneId     string
		CheckIpURL string
	}
	Domains []string
}

var (
	loglevel  string
	configLoc string
	weight    = int64(1)
	config    = Config{}
)

func init() {
	flag.StringVar(&loglevel, "loglevel", "info", loglevel)
	flag.StringVar(&configLoc, "config", "goip.list", configLoc)
}

func main() {
	// parse loglevel
	flag.Parse()
	// initialize the logger
	initLogger(loglevel)

	// read in the config file
	glogger.Info.Printf("using config '%s'\n", configLoc)
	config := parseConfig(configLoc, config)

	sess, err := session.NewSession(&aws.Config{
		Region:      aws.String("us-west-2"),
		Credentials: credentials.NewStaticCredentials(config.Server.Key, config.Server.Secret, ""),
	})
	if err != nil {
		glogger.Error.Printf("bad creds, %s", err)
	}

	svc := route53.New(sess)

	var wg sync.WaitGroup
	wg.Add(3)
	for _, domain := range config.Domains {
		// spawn listener for the domain
		go spawnListener(svc, domain, config)
	}
	wg.Wait()
}

func initLogger(logLevel string) {
	// init logger
	if logLevel == "debug" {
		glogger.LogInit(os.Stdout, os.Stdout, os.Stdout, os.Stderr)
	} else if logLevel == "cluster" {
		glogger.LogInit(os.Stdout, os.Stdout, ioutil.Discard, os.Stderr)
	} else if logLevel == "info" {
		glogger.LogInit(os.Stdout, ioutil.Discard, ioutil.Discard, os.Stderr)
	} else {
		glogger.LogInit(ioutil.Discard, ioutil.Discard, ioutil.Discard, os.Stderr)
	}
}

func updateAName(svc *route53.Route53, domain string, config Config) {
	externalIp, err := getCurrentIP(config)
	if err != nil {
		fmt.Println("error getting ip:", err)
	}

	params := &route53.ChangeResourceRecordSetsInput{
		ChangeBatch: &route53.ChangeBatch{
			Changes: []*route53.Change{
				{
					Action: aws.String("UPSERT"),
					ResourceRecordSet: &route53.ResourceRecordSet{
						Name: aws.String(domain),
						Type: aws.String("A"),
						ResourceRecords: []*route53.ResourceRecord{
							{
								Value: aws.String(externalIp),
							},
						},
						TTL:           aws.Int64(config.Server.TTL),
						Weight:        aws.Int64(weight),
						SetIdentifier: aws.String("Arbitrary Id describing this change set"),
					},
				},
			},
		},
		HostedZoneId: aws.String(config.Server.ZoneId),
	}
	resp, err := svc.ChangeResourceRecordSets(params)

	if err != nil {
		glogger.Error.Println(err.Error())
		return
	}

	// Pretty-print the response data.
	fmt.Println("Change Response:")
	fmt.Println(resp)
}

func spawnListener(svc *route53.Route53, domain string, config Config) {
	glogger.Info.Printf("spawning listener for: %s\n", domain)

	// loop over domain indefinently
	for {
		// get stored and current ip's
		storedIp, err := getStoredIP(domain)
		if err != nil {
			glogger.Error.Println("error getting stored ip")
		}

		currentIp, err := getCurrentIP(config)
		if err != nil {
			glogger.Error.Println("error getting current ip")
		}

		if storedIp != currentIp {
			updateAName(svc, domain, config)
		}
		time.Sleep(time.Second * time.Duration(config.Server.TTL))
	}
}

func getCurrentIP(config Config) (string, error) {
	//glogger.Debug.Printf("using URL '%s'\n", config.Server.CheckIpURL)
	rsp, err := http.Get(config.Server.CheckIpURL)
	if err != nil {
		return "", err
	}
	defer rsp.Body.Close()

	buf, err := ioutil.ReadAll(rsp.Body)
	if err != nil {
		return "", err
	}

	return string(bytes.TrimSpace(buf)), nil
}

func getStoredIP(domain string) (string, error) {
	// check if what is in DNS differs from what is current
	addrs, err := net.LookupIP(domain)
	if err != nil {
		glogger.Error.Printf("error resolving '%s'\n", domain)
		return "", err
	}

	return addrs[0].String(), nil
}
