package main

import (
	"fmt"
	"io/ioutil"
	"path"
	"strings"

	log "github.com/sirupsen/logrus"
	flag "github.com/spf13/pflag"

	"github.com/wtty-fool/dropbox-diff/dropbox"
)

const (
	usage     = "dropbox-diff --dropbox <dropbox path> <local path>"
	tokenFile = "token"
)

var (
	dropboxDir string
	localDir   string
)

func init() {
	log.SetLevel(log.InfoLevel)
	flag.StringVar(&dropboxDir, "dropbox", "", "Dropbox directory path to check against.")
	flag.Parse()
	if dropboxDir == "" {
		log.Fatal("You need to specify dropbox path using --dropbox flag" + "\n\n" + usage)
	}
	if len(flag.Args()) == 0 {
		log.Fatal("You need to specify local directory as an argument" + "\n\n" + usage)
	}
	localDir = flag.Arg(0)
}

func readToken(path string) (string, error) {
	raw, err := ioutil.ReadFile(path)
	if err != nil {
		return "", err
	}
	return string(raw), nil
}

func main() {
	log.Infof("Comparing Dropbox (%s) to local (%s)...", dropboxDir, localDir)

	token, err := readToken(tokenFile)
	if err != nil {
		log.Fatal(err)
	}
	dropboxFiles, err := dropbox.ListDropboxDir(dropboxDir, token)
	if err != nil {
		log.Fatal(err)
	}

	localFiles, err := ioutil.ReadDir(localDir)
	if err != nil {
		log.Fatal(err)
	}

	fmt.Println("Files missing in Dropbox:")
	dropboxFileMap := map[string]bool{}
	for _, f := range dropboxFiles {
		dropboxFileMap[strings.ToLower(f.Name)] = true
	}

	for _, f := range localFiles {
		if f.IsDir() {
			continue
		}

		_, ok := dropboxFileMap[strings.ToLower(f.Name())]
		if !ok {
			fmt.Printf(path.Join(localDir, f.Name()))
		}
	}

	log.Info("Done")
}
