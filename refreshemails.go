package main

import (
	"io"
	"net/http"
	"os"

	"github.com/pusher/oauth2_proxy/logger"
)

type emailListRefresher struct {
	AuthenticatedEmailsFile string
	AuthenticatedEmailsURL  string
}

func (r *emailListRefresher) Refresh() {
	logger.Printf("Call refresh")
	r.downloadFile()
	logger.Printf("Refresh completed")
}

func (r *emailListRefresher) downloadFile() error {
	// DownloadFile will download a url to a local file. It's efficient because it will
	// write as it downloads and not load the whole file into memory.
	//func DownloadFile(filepath string, url string) error {
	// Get the data
	resp, err := http.Get(r.AuthenticatedEmailsURL)
	if err != nil {
		return err
	}
	defer resp.Body.Close()

	// Create the file
	out, err := os.Create(r.AuthenticatedEmailsFile)
	if err != nil {
		return err
	}
	defer out.Close()

	// Write the body to file
	_, err = io.Copy(out, resp.Body)
	return err
}

// NewEmailRefresher return a function to refresh email list
func NewEmailRefresher(file string, url string) func() {
	emailRefresher := &emailListRefresher{
		AuthenticatedEmailsFile: file,
		AuthenticatedEmailsURL:  url,
	}
	refresher := func() {
		emailRefresher.Refresh()
	}

	return refresher
}
