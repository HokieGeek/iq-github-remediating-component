package main

import (
	"bytes"
	"encoding/json"
	"errors"
	"fmt"
	"io"
	"io/ioutil"
	"log"
	"net/http"
	"time"
)

// GitlabMergeRequest defines the structure of a Gitlab merge request
type GitlabMergeRequest struct {
	ID                          int64                `json:"id"`
	Iid                         int64                `json:"iid"`
	ProjectID                   int64                `json:"project_id"`
	Title                       string               `json:"title"`
	Description                 string               `json:"description"`
	State                       string               `json:"state"`
	CreatedAt                   string               `json:"created_at"`
	UpdatedAt                   string               `json:"updated_at"`
	MergedBy                    string               `json:"merged_by"`
	MergedAt                    string               `json:"merged_at"`
	ClosedBy                    string               `json:"closed_by"`
	ClosedAt                    string               `json:"closed_at"`
	TargetBranch                string               `json:"target_branch"`
	SourceBranch                string               `json:"source_branch"`
	UserNotesCount              int64                `json:"user_notes_count"`
	Upvotes                     int64                `json:"upvotes"`
	Downvotes                   int64                `json:"downvotes"`
	Assignee                    string               `json:"assignee"`
	Author                      author               `json:"author"`
	Assignees                   []string             `json:"assignees"`
	SourceProjectID             int64                `json:"source_project_id"`
	TargetProjectID             int64                `json:"target_project_id"`
	Labels                      []string             `json:"labels"`
	WorkInProgress              bool                 `json:"work_in_progress"`
	Milestone                   string               `json:"milestone"`
	MergeWhenPipelineSucceeds   bool                 `json:"merge_when_pipeline_succeeds"`
	MergeStatus                 string               `json:"merge_status"`
	SHA                         string               `json:"sha"`
	MergeCommitSHA              string               `json:"merge_commit_sha"`
	DiscussionLocked            bool                 `json:"discussion_locked"`
	ShouldRemoveSourceBranch    bool                 `json:"should_remove_source_branch"`
	ForceRemoveSourceBranch     bool                 `json:"force_remove_source_branch"`
	Reference                   string               `json:"reference"`
	WebURL                      string               `json:"web_url"`
	TimeStats                   timeStats            `json:"time_stats"`
	Squash                      bool                 `json:"squash"`
	TaskCompletionStatus        taskCompletionStatus `json:"task_completion_status"`
	Subscribed                  bool                 `json:"subscribed"`
	ChangesCount                string               `json:"changes_count"`
	LatestBuildStartedAt        string               `json:"latest_build_started_at"`
	LatestBuildFinishedAt       string               `json:"latest_build_finished_at"`
	FirstDeployedToProductionAt string               `json:"first_deployed_to_production_at"`
	Pipeline                    string               `json:"pipeline"`
	HeadPipeline                string               `json:"head_pipeline"`
	DiffRefs                    diffRefs             `json:"diff_refs"`
	MergeError                  string               `json:"merge_error"`
	User                        user                 `json:"user"`
	Changes                     []change             `json:"changes,omitempty"`
	ApprovalsBeforeMerge        int64                `json:"approvals_before_merge"`
}

type author struct {
	ID        int64  `json:"id,omitempty"`
	Name      string `json:"name,omitempty"`
	Username  string `json:"username,omitempty"`
	State     string `json:"state,omitempty"`
	AvatarURL string `json:"avatar_url,omitempty"`
	WebURL    string `json:"web_url,omitempty"`
	Email     string `json:"email,omitempty"`
}

type diffRefs struct {
	BaseSHA  string `json:"base_sha"`
	HeadSHA  string `json:"head_sha"`
	StartSHA string `json:"start_sha"`
}

type taskCompletionStatus struct {
	Count          int64 `json:"count"`
	CompletedCount int64 `json:"completed_count"`
}

type timeStats struct {
	TimeEstimate        int64  `json:"time_estimate"`
	TotalTimeSpent      int64  `json:"total_time_spent"`
	HumanTimeEstimate   string `json:"human_time_estimate"`
	HumanTotalTimeSpent string `json:"human_total_time_spent"`
}

type user struct {
	CanMerge bool `json:"can_merge"`
}

type gitlabDiscussionRequest struct {
	ID              int64    `json:"id"`
	MergeRequestIID int64    `json:"merge_request_iid"`
	Body            string   `json:"body"`
	Position        position `json:"position"`
}

type position struct {
	PositionType string `json:"position_type"`
	BaseSHA      string `json:"base_sha,omitempty"`  // Base commit SHA in the source branch
	StartSHA     string `json:"start_sha,omitempty"` // SHA referencing commit in target branch
	HeadSHA      string `json:"head_sha,omitempty"`  // SHA referencing HEAD of this merge request
	OldPath      string `json:"old_path,omitempty"`
	NewPath      string `json:"new_path,omitempty"`
	OldLine      int64  `json:"old_line,omitempty"`
	NewLine      int64  `json:"new_line,omitempty"`
}

type change struct {
	OldPath     string `json:"old_path"`
	NewPath     string `json:"new_path"`
	AMode       string `json:"a_mode"`
	BMode       string `json:"b_mode"`
	NewFile     bool   `json:"new_file"`
	RenamedFile bool   `json:"renamed_file"`
	DeletedFile bool   `json:"deleted_file"`
	Diff        string `json:"diff"`
}

type gitlabMergeRequestWebhookEvent struct {
	ObjectKind       string           `json:"object_kind"`
	User             eventUser        `json:"user"`
	Project          project          `json:"project"`
	Repository       repository       `json:"repository"`
	ObjectAttributes objectAttributes `json:"object_attributes"`
	Labels           []label          `json:"labels"`
	Changes          changes          `json:"changes"`
}

type changes struct {
	UpdatedByID updatedByID `json:"updated_by_id"`
	UpdatedAt   updatedAt   `json:"updated_at"`
	Labels      labels      `json:"labels"`
}

type labels struct {
	Previous []label `json:"previous"`
	Current  []label `json:"current"`
}

type label struct {
	ID          int64  `json:"id"`
	Title       string `json:"title"`
	Color       string `json:"color"`
	ProjectID   int64  `json:"project_id"`
	CreatedAt   string `json:"created_at"`
	UpdatedAt   string `json:"updated_at"`
	Template    bool   `json:"template"`
	Description string `json:"description"`
	Type        string `json:"type"`
	GroupID     int64  `json:"group_id"`
}

type updatedAt struct {
	Previous string `json:"previous"`
	Current  string `json:"current"`
}

type updatedByID struct {
	Previous int64 `json:"previous"`
	Current  int64 `json:"current"`
}

type objectAttributes struct {
	ID              int64     `json:"id"`
	TargetBranch    string    `json:"target_branch"`
	SourceBranch    string    `json:"source_branch"`
	SourceProjectID int64     `json:"source_project_id"`
	AuthorID        int64     `json:"author_id"`
	AssigneeID      int64     `json:"assignee_id"`
	Title           string    `json:"title"`
	CreatedAt       string    `json:"created_at"`
	UpdatedAt       string    `json:"updated_at"`
	MilestoneID     int64     `json:"milestone_id"`
	State           string    `json:"state"`
	MergeStatus     string    `json:"merge_status"`
	TargetProjectID int64     `json:"target_project_id"`
	Iid             int64     `json:"iid"`
	Description     string    `json:"description"`
	Source          project   `json:"source"`
	Target          project   `json:"target"`
	LastCommit      commit    `json:"last_commit"`
	WorkInProgress  bool      `json:"work_in_progress"`
	URL             string    `json:"url"`
	Action          string    `json:"action"`
	Assignee        eventUser `json:"assignee"`
}

type eventUser struct {
	Name      string `json:"name"`
	Username  string `json:"username"`
	AvatarURL string `json:"avatar_url"`
}

type commit struct {
	ID        string `json:"id"`
	Message   string `json:"message"`
	Timestamp string `json:"timestamp"`
	URL       string `json:"url"`
	Author    author `json:"author"`
}

type project struct {
	Name              string `json:"name"`
	Description       string `json:"description"`
	WebURL            string `json:"web_url"`
	AvatarURL         string `json:"avatar_url"`
	GitSSHURL         string `json:"git_ssh_url"`
	GitHTTPURL        string `json:"git_http_url"`
	Namespace         string `json:"namespace"`
	VisibilityLevel   int64  `json:"visibility_level"`
	PathWithNamespace string `json:"path_with_namespace"`
	DefaultBranch     string `json:"default_branch"`
	Homepage          string `json:"homepage"`
	URL               string `json:"url"`
	SSHURL            string `json:"ssh_url"`
	HTTPURL           string `json:"http_url"`
	ID                int64  `json:"id,omitempty"`
}

type repository struct {
	Name        string `json:"name"`
	URL         string `json:"url"`
	Description string `json:"description"`
	Homepage    string `json:"homepage"`
}

func req(method, url, token string, payload io.Reader) (*http.Response, error) {
	log.Printf("TRACE: req(%s, %s, %s, payload)", method, url, token)
	request, err := http.NewRequest(method, url, payload)
	if err != nil {
		return nil, fmt.Errorf("could not create HTTP request: %v", err)
	}

	request.Header.Set("PRIVATE-TOKEN", token)

	if payload != nil {
		request.Header.Set("Content-Type", "application/json")
	}

	client := &http.Client{
		Timeout: 60 * time.Second,
	}

	return client.Do(request)

}

func getMergeRequest(token string, projectID, mrIID int64) (GitlabMergeRequest, error) {
	url := fmt.Sprintf("https://gitlab.com/api/v4/projects/%d/merge_requests/%d", projectID, mrIID)
	resp, err := req("GET", url, token, nil)
	if err != nil {
		return GitlabMergeRequest{}, err
	}

	fmt.Printf("%#v\n", resp)

	defer resp.Body.Close()
	body, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		return GitlabMergeRequest{}, err
	}

	var mr GitlabMergeRequest
	if err := json.Unmarshal(body, &mr); err != nil {
		return GitlabMergeRequest{}, err
	}

	return mr, nil
}

func addMergeRequestComment(token string, mr GitlabMergeRequest, pos int64, path, comment string) {
	discussionReq := gitlabDiscussionRequest{
		ID:              mr.ProjectID,
		MergeRequestIID: mr.Iid,
		Body:            comment,
		Position: position{
			PositionType: "text",
			BaseSHA:      mr.DiffRefs.BaseSHA,
			HeadSHA:      mr.DiffRefs.HeadSHA,
			StartSHA:     mr.DiffRefs.StartSHA,
			OldPath:      path,
			NewPath:      path,
			OldLine:      pos,
		},
	}

	buf, err := json.Marshal(discussionReq)
	if err != nil {
		panic(err)
	}

	fmt.Println(string(buf))

	url := fmt.Sprintf("https://gitlab.com/api/v4/projects/%d/merge_requests/%d/discussions", mr.ProjectID, mr.Iid)
	resp, err := req("POST", url, token, bytes.NewBuffer(buf))
	if err != nil {
		panic(err)
	}

	fmt.Printf("%#v\n", resp)

}

func getGitlabEventType(requestHeaders map[string]string) (string, error) {
	eventType, ok := requestHeaders["X-Gitlab-Event"]
	if !ok {
		return "", errors.New("error: did not receive a gitlab event")
	}
	return eventType, nil
}

// IsValidGitlabWebhookMergeRequestEvent returns true if the given HTTP headers are for a valid merge request.
// Also returns a valid http status code.
func IsValidGitlabWebhookMergeRequestEvent(reqHeaders map[string]string) (bool, int) {
	eventType, err := getGitlabEventType(reqHeaders)
	if err != nil {
		log.Printf("could not parse request headers: %v", err)
		return false, http.StatusBadRequest
	}

	switch {
	case eventType != "Merge Request Hook":
		log.Printf("ERROR: did not receive a supported gitlab event: %s\n", eventType)
		return false, http.StatusBadRequest
	}

	return true, http.StatusOK
}
