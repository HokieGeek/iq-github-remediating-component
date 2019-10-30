package main

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"os"
	"testing"

	nexusiq "github.com/sonatype-nexus-community/gonexus/iq"
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

	iqURL := "URL"
	iqUser := "USER"
	iqPassword := "PASS"
	iq, err := nexusiq.New(iqURL, iqUser, iqPassword)
	if err != nil {
		panic(err)
	}

	type args struct {
		iq    nexusiq.IQ
		iqApp string
		token string
		pull  GithubPullRequest
	}
	tests := []struct {
		name    string
		args    args
		wantErr bool
	}{
		{
			"real data",
			args{
				iq:    iq,
				token: token,
				iqApp: "APP",
				pull:  pull,
			},
			false,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if err := ProcessPullRequestForRemediations(tt.args.iq, tt.args.token, tt.args.iqApp, tt.args.pull); (err != nil) != tt.wantErr {
				t.Errorf("processPullRequestForRemediations() error = %v, wantErr %v", err, tt.wantErr)
			}
		})
	}
}
