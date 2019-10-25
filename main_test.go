package main

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"os"
	"testing"
)

func TestHandlePullRequest(t *testing.T) {
	tokenf, err := os.Open("token")
	if err != nil {
		panic(err)
	}
	defer tokenf.Close()

	buf, _ := ioutil.ReadAll(tokenf)
	token := string(buf)

	pr, err := os.Open("pullrequest_sample.json")
	if err != nil {
		panic(err)
	}
	defer pr.Close()
	buf, _ = ioutil.ReadAll(pr)
	var pull githubPullRequest
	if err := json.Unmarshal(buf, &pull); err != nil {
		fmt.Println(string(buf))
		panic(err)
	}

	// TODO
	url := "URL"
	user := "USER"
	pass := "PASS"
	app := "APP"
	if err := handlePullRequest(url, user, pass, app, token, pull); err != nil {
		panic(fmt.Sprintf("ERROR: error handling pull request: %v\n", err))
	}
}
