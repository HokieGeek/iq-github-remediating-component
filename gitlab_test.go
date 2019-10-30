package main

import "testing"

func Test_addMergeRequestComment(t *testing.T) {
	/*
		token := os.Args[1]
		repoID := os.Args[2]
		mrid, err := strconv.ParseInt(os.Args[3], 10, 64)
		if err != nil {
			panic(err)
		}

		file := "Gemfile"
		pos := int64(47)
		comment := "```suggestion:-0+0\ngem 'doorkeeper', '~> 4.5'\n```"

		url := fmt.Sprintf("https://gitlab.com/api/v4/projects/%s/merge_requests/%d", repoID, mrid)
		resp, err := req("GET", url, token, nil)
		if err != nil {
			panic(err)
		}

		fmt.Printf("%#v\n", resp)

		defer resp.Body.Close()
		body, err := ioutil.ReadAll(resp.Body)
		if err != nil {
			panic(err)
		}

		var mr GitlabMergeRequest
		if err := json.Unmarshal(body, &mr); err != nil {
			panic(err)
		}

		addMergeRequestComment(token, mr, pos, file, comment)
	*/
	type args struct {
		token   string
		mr      GitlabMergeRequest
		pos     int64
		path    string
		comment string
	}
	tests := []struct {
		name string
		args args
	}{
		// TODO: Add test cases.
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			addMergeRequestComment(tt.args.token, tt.args.mr, tt.args.pos, tt.args.path, tt.args.comment)
		})
	}
}
