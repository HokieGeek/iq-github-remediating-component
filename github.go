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

	nexusiq "github.com/sonatype-nexus-community/gonexus/iq"
)

// GithubPullRequest defines the structure of a Github pull request
type GithubPullRequest struct {
	Action      string      `json:"action"`
	Number      int64       `json:"number"`
	PullRequest pullRequest `json:"pull_request"`
	Repository  repo        `json:"repository"`
	Sender      githubUser  `json:"sender"`
}

type pullRequest struct {
	URL                 string     `json:"url"`
	ID                  int64      `json:"id"`
	NodeID              string     `json:"node_id"`
	HTMLURL             string     `json:"html_url"`
	DiffURL             string     `json:"diff_url"`
	PatchURL            string     `json:"patch_url"`
	IssueURL            string     `json:"issue_url"`
	Number              int64      `json:"number"`
	State               string     `json:"state"`
	Locked              bool       `json:"locked"`
	Title               string     `json:"title"`
	User                githubUser `json:"user"`
	Body                string     `json:"body"`
	CreatedAt           string     `json:"created_at"`
	UpdatedAt           string     `json:"updated_at"`
	ClosedAt            string     `json:"closed_at"`
	MergedAt            string     `json:"merged_at"`
	MergeCommitSHA      string     `json:"merge_commit_sha"`
	Assignee            string     `json:"assignee"`
	Assignees           []string   `json:"assignees"`
	RequestedReviewers  []string   `json:"requested_reviewers"`
	RequestedTeams      []string   `json:"requested_teams"`
	Labels              []string   `json:"labels"`
	Milestone           string     `json:"milestone"`
	CommitsURL          string     `json:"commits_url"`
	ReviewCommentsURL   string     `json:"review_comments_url"`
	ReviewCommentURL    string     `json:"review_comment_url"`
	CommentsURL         string     `json:"comments_url"`
	StatusesURL         string     `json:"statuses_url"`
	Head                gitref     `json:"head"`
	Base                gitref     `json:"base"`
	Links               links      `json:"_links"`
	AuthorAssociation   string     `json:"author_association"`
	Draft               bool       `json:"draft"`
	Merged              bool       `json:"merged"`
	Mergeable           bool       `json:"mergeable"`
	Rebaseable          bool       `json:"rebaseable"`
	MergeableState      string     `json:"mergeable_state"`
	MergedBy            string     `json:"merged_by"`
	Comments            int64      `json:"comments"`
	ReviewComments      int64      `json:"review_comments"`
	MaintainerCanModify bool       `json:"maintainer_can_modify"`
	Commits             int64      `json:"commits"`
	Additions           int64      `json:"additions"`
	Deletions           int64      `json:"deletions"`
	ChangedFiles        int64      `json:"changed_files"`
}

type gitref struct {
	Label string     `json:"label"`
	Ref   string     `json:"ref"`
	SHA   string     `json:"sha"`
	User  githubUser `json:"user"`
	Repo  repo       `json:"repo"`
}

type repo struct {
	ID               int64      `json:"id"`
	NodeID           string     `json:"node_id"`
	Name             string     `json:"name"`
	FullName         string     `json:"full_name"`
	Private          bool       `json:"private"`
	Owner            githubUser `json:"owner"`
	HTMLURL          string     `json:"html_url"`
	Description      string     `json:"description"`
	Fork             bool       `json:"fork"`
	URL              string     `json:"url"`
	ForksURL         string     `json:"forks_url"`
	KeysURL          string     `json:"keys_url"`
	CollaboratorsURL string     `json:"collaborators_url"`
	TeamsURL         string     `json:"teams_url"`
	HooksURL         string     `json:"hooks_url"`
	IssueEventsURL   string     `json:"issue_events_url"`
	EventsURL        string     `json:"events_url"`
	AssigneesURL     string     `json:"assignees_url"`
	BranchesURL      string     `json:"branches_url"`
	TagsURL          string     `json:"tags_url"`
	BlobsURL         string     `json:"blobs_url"`
	GitTagsURL       string     `json:"git_tags_url"`
	GitRefsURL       string     `json:"git_refs_url"`
	TreesURL         string     `json:"trees_url"`
	StatusesURL      string     `json:"statuses_url"`
	LanguagesURL     string     `json:"languages_url"`
	StargazersURL    string     `json:"stargazers_url"`
	ContributorsURL  string     `json:"contributors_url"`
	SubscribersURL   string     `json:"subscribers_url"`
	SubscriptionURL  string     `json:"subscription_url"`
	CommitsURL       string     `json:"commits_url"`
	GitCommitsURL    string     `json:"git_commits_url"`
	CommentsURL      string     `json:"comments_url"`
	IssueCommentURL  string     `json:"issue_comment_url"`
	ContentsURL      string     `json:"contents_url"`
	CompareURL       string     `json:"compare_url"`
	MergesURL        string     `json:"merges_url"`
	ArchiveURL       string     `json:"archive_url"`
	DownloadsURL     string     `json:"downloads_url"`
	IssuesURL        string     `json:"issues_url"`
	PullsURL         string     `json:"pulls_url"`
	MilestonesURL    string     `json:"milestones_url"`
	NotificationsURL string     `json:"notifications_url"`
	LabelsURL        string     `json:"labels_url"`
	ReleasesURL      string     `json:"releases_url"`
	DeploymentsURL   string     `json:"deployments_url"`
	CreatedAt        string     `json:"created_at"`
	UpdatedAt        string     `json:"updated_at"`
	PushedAt         string     `json:"pushed_at"`
	GitURL           string     `json:"git_url"`
	SSHURL           string     `json:"ssh_url"`
	CloneURL         string     `json:"clone_url"`
	SvnURL           string     `json:"svn_url"`
	Homepage         string     `json:"homepage"`
	Size             int64      `json:"size"`
	StargazersCount  int64      `json:"stargazers_count"`
	WatchersCount    int64      `json:"watchers_count"`
	Language         string     `json:"language"`
	HasIssues        bool       `json:"has_issues"`
	HasProjects      bool       `json:"has_projects"`
	HasDownloads     bool       `json:"has_downloads"`
	HasWiki          bool       `json:"has_wiki"`
	HasPages         bool       `json:"has_pages"`
	ForksCount       int64      `json:"forks_count"`
	MirrorURL        string     `json:"mirror_url"`
	Archived         bool       `json:"archived"`
	Disabled         bool       `json:"disabled"`
	OpenIssuesCount  int64      `json:"open_issues_count"`
	License          license    `json:"license"`
	Forks            int64      `json:"forks"`
	OpenIssues       int64      `json:"open_issues"`
	Watchers         int64      `json:"watchers"`
	DefaultBranch    string     `json:"default_branch"`
}

type license struct {
	Key    string `json:"key"`
	Name   string `json:"name"`
	SpdxID string `json:"spdx_id"`
	URL    string `json:"url"`
	NodeID string `json:"node_id"`
}

type githubUser struct {
	Login             string `json:"login"`
	ID                int64  `json:"id"`
	NodeID            string `json:"node_id"`
	AvatarURL         string `json:"avatar_url"`
	GravatarID        string `json:"gravatar_id"`
	URL               string `json:"url"`
	HTMLURL           string `json:"html_url"`
	FollowersURL      string `json:"followers_url"`
	FollowingURL      string `json:"following_url"`
	GistsURL          string `json:"gists_url"`
	StarredURL        string `json:"starred_url"`
	SubscriptionsURL  string `json:"subscriptions_url"`
	OrganizationsURL  string `json:"organizations_url"`
	ReposURL          string `json:"repos_url"`
	EventsURL         string `json:"events_url"`
	ReceivedEventsURL string `json:"received_events_url"`
	Type              string `json:"type"`
	SiteAdmin         bool   `json:"site_admin"`
}

type links struct {
	Self           href `json:"self"`
	HTML           href `json:"html"`
	Issue          href `json:"issue"`
	Comments       href `json:"comments"`
	ReviewComments href `json:"review_comments"`
	ReviewComment  href `json:"review_comment"`
	Commits        href `json:"commits"`
	Statuses       href `json:"statuses"`
}

type href struct {
	Href string `json:"href"`
}

// GET /repos/:owner/:repo/pulls/:pull_number/files
type githubPullRequestFile struct {
	SHA         string `json:"sha"`
	Filename    string `json:"filename"`
	Status      string `json:"status"`
	Additions   int64  `json:"additions"`
	Deletions   int64  `json:"deletions"`
	Changes     int64  `json:"changes"`
	BlobURL     string `json:"blob_url"`
	RawURL      string `json:"raw_url"`
	ContentsURL string `json:"contents_url"`
	Patch       string `json:"patch"`
}

// POST /repos/:owner/:repo/pulls/:pull_number/comments
type githubPullRequestCommentSinglelineRequest struct {
	CommitID string `json:"commit_id"`
	Path     string `json:"path"`
	Position int64  `json:"position"`
	Side     string `json:"side"`
	Body     string `json:"body"`
}

func getGitHubEventType(requestHeaders map[string]string) (string, error) {
	eventType, ok := requestHeaders["X-GitHub-Event"]
	if !ok {
		return "", errors.New("error: did not receive a github event")
	}
	return eventType, nil
}

func ghreq(method, url, token string, payload io.Reader) (*http.Response, error) {
	log.Printf("TRACE: req(%s, %s, %s, payload)", method, url, token)
	request, err := http.NewRequest(method, url, payload)
	if err != nil {
		return nil, fmt.Errorf("could not create HTTP request: %v", err)
	}

	request.Header.Set("Authorization", fmt.Sprintf("token %s", token))

	if payload != nil {
		request.Header.Set("Content-Type", "application/json")
	}

	client := &http.Client{
		Timeout: 60 * time.Second,
	}

	return client.Do(request)
}

func getPullRequestFiles(token string, pull GithubPullRequest) ([]changedFile, error) {
	resp, err := ghreq(http.MethodGet, fmt.Sprintf("%s/files", pull.PullRequest.URL), token, nil)
	if err != nil {
		return nil, err
	}

	if resp.StatusCode != http.StatusOK {
		return nil, fmt.Errorf("did not get OK status: %s", resp.Status)
	}

	defer resp.Body.Close()
	body, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		return nil, err
	}

	var ghfiles []githubPullRequestFile
	if err := json.Unmarshal(body, &ghfiles); err != nil {
		return nil, err
	}

	files := make([]changedFile, len(ghfiles))
	for i, f := range ghfiles {
		files[i] = changedFile{Filename: f.Filename, Patch: f.Patch}
	}

	return files, err
}

func addPullRequestComment(token string, pull GithubPullRequest, position int64, path, comment string) error {
	request := githubPullRequestCommentSinglelineRequest{
		CommitID: pull.PullRequest.Head.SHA,
		Path:     path,
		Position: position,
		Side:     "RIGHT",
		Body:     comment,
	}

	buf, err := json.Marshal(request)
	if err != nil {
		return fmt.Errorf("could not create request: %s", err)
	}

	resp, err := ghreq(http.MethodPost, pull.PullRequest.ReviewCommentsURL, token, bytes.NewBuffer(buf))
	if err != nil {
		log.Printf("ERROR: error creating comment: %s", err)
		log.Printf("TRACE: %s", buf)
		return fmt.Errorf("error creating comment: %s", err)
	}
	if resp.StatusCode != http.StatusCreated {
		log.Printf("ERROR: error creating comment. got status: %s", resp.Status)
		log.Printf("TRACE: %s", buf)
		return fmt.Errorf("error creating comment. got status: %s", resp.Status)
	}

	return nil
}

// IsValidGithubWebhookPullRequestEvent returns true if the given HTTP headers are for a valid pull request or ping.
// Also returns a valid http status code.
func IsValidGithubWebhookPullRequestEvent(reqHeaders map[string]string) (bool, int) {
	eventType, err := getGitHubEventType(reqHeaders)
	if err != nil {
		log.Printf("could not parse request headers: %v", err)
		return false, http.StatusBadRequest
	}

	switch {
	case eventType == "ping":
		return false, http.StatusOK
	case eventType != "pull_request":
		log.Printf("ERROR: did not receive a supported github event: %s\n", eventType)
		return false, http.StatusBadRequest
	}

	return true, http.StatusOK
}

// ProcessPullRequestForRemediations will take a Github pull request and add any remediations if a manifest is found
func ProcessPullRequestForRemediations(iq nexusiq.IQ, iqApp, token string, pull GithubPullRequest) error {
	log.Printf("TRACE: Received Pull Request from: %s\n", pull.Repository.HTMLURL)

	files, err := getPullRequestFiles(token, pull)
	if err != nil {
		log.Printf("ERROR: could not get files from pull request: %v\n", err)
		return fmt.Errorf("could not get files from pull request: %v", err)
	}
	log.Printf("TRACE: Got %d files from pull request\n", len(files))

	if err = addRemediationsToRequest(iq, iqApp, files, func(filename string, location changeLocation, comment string) error {
		return addPullRequestComment(token, pull, location.Position, filename, comment)
	}); err != nil {
		return fmt.Errorf("could not add remediation comments to request: %v", err)
	}

	return nil
}

// HandleGithubWebhookPullRequestEvent unmarshals a pull request event from Github and remediates if it is a new one
func HandleGithubWebhookPullRequestEvent(iq nexusiq.IQ, iqApp, token string, payload []byte) (int, error) {
	var event GithubPullRequest
	if err := json.Unmarshal(payload, &event); err != nil {
		return http.StatusBadRequest, fmt.Errorf("could not unmarshal payload as json: %v", err)
	}

	if event.Action != "opened" {
		return http.StatusNoContent, fmt.Errorf("Only processing new pull requests")
	}

	if err := ProcessPullRequestForRemediations(iq, iqApp, token, event); err != nil {
		return http.StatusInternalServerError, fmt.Errorf("error: error handling pull request: %v", err)
	}

	return http.StatusOK, nil
}
