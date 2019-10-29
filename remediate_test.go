package main

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"os"
	"testing"
)

func Test_processPullRequestForRemediations(t *testing.T) {
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
	var pull GithubPullRequest
	if err := json.Unmarshal(buf, &pull); err != nil {
		fmt.Println(string(buf))
		panic(err)
	}

	type args struct {
		iqURL      string
		iqUser     string
		iqPassword string
		iqApp      string
		token      string
		pull       GithubPullRequest
	}
	tests := []struct {
		name    string
		args    args
		wantErr bool
	}{
		{
			"real data",
			args{
				iqURL:      "URL",
				iqUser:     "USER",
				iqPassword: "PASS",
				iqApp:      "APP",
				token:      token,
				pull:       pull,
			},
			false,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if err := ProcessPullRequestForRemediations(tt.args.iqURL, tt.args.iqUser, tt.args.iqPassword, tt.args.iqApp, tt.args.token, tt.args.pull); (err != nil) != tt.wantErr {
				t.Errorf("processPullRequestForRemediations() error = %v, wantErr %v", err, tt.wantErr)
			}
		})
	}
}
