package main

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"log"
	"net/http"
	"net/http/httputil"
	"time"

	tektonv1alpha1 "github.com/tektoncd/pipeline/pkg/apis/pipeline/v1alpha1"
	tekton "github.com/tektoncd/pipeline/pkg/client/clientset/versioned"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/client-go/rest"
)

type PushEvent struct {
	Ref        string        `json:"ref"`
	Before     string        `json:"before"`
	After      string        `json:"after"`
	Created    bool          `json:"created"`
	Deleted    bool          `json:"deleted"`
	Forced     bool          `json:"forced"`
	BaseRef    interface{}   `json:"base_ref"`
	Compare    string        `json:"compare"`
	Commits    []interface{} `json:"commits"`
	HeadCommit interface{}   `json:"head_commit"`
	Repository struct {
		ID       int    `json:"id"`
		NodeID   string `json:"node_id"`
		Name     string `json:"name"`
		FullName string `json:"full_name"`
		Owner    struct {
			Name              string `json:"name"`
			Email             string `json:"email"`
			Login             string `json:"login"`
			ID                int    `json:"id"`
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
		} `json:"owner"`
		Private          bool        `json:"private"`
		HTMLURL          string      `json:"html_url"`
		Description      interface{} `json:"description"`
		Fork             bool        `json:"fork"`
		URL              string      `json:"url"`
		ForksURL         string      `json:"forks_url"`
		KeysURL          string      `json:"keys_url"`
		CollaboratorsURL string      `json:"collaborators_url"`
		TeamsURL         string      `json:"teams_url"`
		HooksURL         string      `json:"hooks_url"`
		IssueEventsURL   string      `json:"issue_events_url"`
		EventsURL        string      `json:"events_url"`
		AssigneesURL     string      `json:"assignees_url"`
		BranchesURL      string      `json:"branches_url"`
		TagsURL          string      `json:"tags_url"`
		BlobsURL         string      `json:"blobs_url"`
		GitTagsURL       string      `json:"git_tags_url"`
		GitRefsURL       string      `json:"git_refs_url"`
		TreesURL         string      `json:"trees_url"`
		StatusesURL      string      `json:"statuses_url"`
		LanguagesURL     string      `json:"languages_url"`
		StargazersURL    string      `json:"stargazers_url"`
		ContributorsURL  string      `json:"contributors_url"`
		SubscribersURL   string      `json:"subscribers_url"`
		SubscriptionURL  string      `json:"subscription_url"`
		CommitsURL       string      `json:"commits_url"`
		GitCommitsURL    string      `json:"git_commits_url"`
		CommentsURL      string      `json:"comments_url"`
		IssueCommentURL  string      `json:"issue_comment_url"`
		ContentsURL      string      `json:"contents_url"`
		CompareURL       string      `json:"compare_url"`
		MergesURL        string      `json:"merges_url"`
		ArchiveURL       string      `json:"archive_url"`
		DownloadsURL     string      `json:"downloads_url"`
		IssuesURL        string      `json:"issues_url"`
		PullsURL         string      `json:"pulls_url"`
		MilestonesURL    string      `json:"milestones_url"`
		NotificationsURL string      `json:"notifications_url"`
		LabelsURL        string      `json:"labels_url"`
		ReleasesURL      string      `json:"releases_url"`
		DeploymentsURL   string      `json:"deployments_url"`
		CreatedAt        int         `json:"created_at"`
		UpdatedAt        time.Time   `json:"updated_at"`
		PushedAt         int         `json:"pushed_at"`
		GitURL           string      `json:"git_url"`
		SSHURL           string      `json:"ssh_url"`
		CloneURL         string      `json:"clone_url"`
		SvnURL           string      `json:"svn_url"`
		Homepage         interface{} `json:"homepage"`
		Size             int         `json:"size"`
		StargazersCount  int         `json:"stargazers_count"`
		WatchersCount    int         `json:"watchers_count"`
		Language         interface{} `json:"language"`
		HasIssues        bool        `json:"has_issues"`
		HasProjects      bool        `json:"has_projects"`
		HasDownloads     bool        `json:"has_downloads"`
		HasWiki          bool        `json:"has_wiki"`
		HasPages         bool        `json:"has_pages"`
		ForksCount       int         `json:"forks_count"`
		MirrorURL        interface{} `json:"mirror_url"`
		Archived         bool        `json:"archived"`
		OpenIssuesCount  int         `json:"open_issues_count"`
		License          interface{} `json:"license"`
		Forks            int         `json:"forks"`
		OpenIssues       int         `json:"open_issues"`
		Watchers         int         `json:"watchers"`
		DefaultBranch    string      `json:"default_branch"`
		Stargazers       int         `json:"stargazers"`
		MasterBranch     string      `json:"master_branch"`
	} `json:"repository"`
	Pusher struct {
		Name  string `json:"name"`
		Email string `json:"email"`
	} `json:"pusher"`
	Sender struct {
		Login             string `json:"login"`
		ID                int    `json:"id"`
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
	} `json:"sender"`
}

type PushEventTrigger struct {
	client            *tekton.Clientset
	namespace         string
	pipelineRunPrefix string
	repoName          string
	pipelineName      string
}

func GetConfig() (*tekton.Clientset, error) {
	config, err := rest.InClusterConfig()
	if err != nil {
		return &tekton.Clientset{}, err
	}

	client, err := tekton.NewForConfig(config)
	if err != nil {
		return &tekton.Clientset{}, err
	}

	return client, nil
}

func (pet *PushEventTrigger) createPipelineRun(timestamp int) error {
	runName := fmt.Sprintf("%s-%d", pet.pipelineRunPrefix, timestamp)

	log.Printf("Create PipelineRun: %s in namespace %s\n", runName, pet.namespace)

	pipelineRun := &tektonv1alpha1.PipelineRun{
		ObjectMeta: metav1.ObjectMeta{
			Name:      runName,
			Namespace: pet.namespace,
		},
		Spec: tektonv1alpha1.PipelineRunSpec{
			PipelineRef: tektonv1alpha1.PipelineRef{
				Name: pet.pipelineName,
			},
			Resources: []tektonv1alpha1.PipelineResourceBinding{
				tektonv1alpha1.PipelineResourceBinding{
					Name: "git-source",
					ResourceRef: tektonv1alpha1.PipelineResourceRef{
						Name: "todo-app-git",
					},
				},
			},
			Params: []tektonv1alpha1.Param{
				tektonv1alpha1.Param{
					Name:  "pathToYamlFile",
					Value: "knative/todo-app.yaml",
				},
				tektonv1alpha1.Param{
					Name:  "imageUrl",
					Value: "johscheuer/todo-app-web",
				},
				tektonv1alpha1.Param{
					Name:  "pathToContext",
					Value: "dir:///workspace/git-source/",
				},
				tektonv1alpha1.Param{
					Name:  "pathToDockerFile",
					Value: "Dockerfile",
				},
				tektonv1alpha1.Param{
					Name:  "imageTag",
					Value: fmt.Sprintf("tekton-%d", timestamp),
				},
			},
			ServiceAccount: "build-bot",
		},
	}

	_, err := pet.client.TektonV1alpha1().PipelineRuns(pet.namespace).Create(pipelineRun)

	return err
}

func (pet *PushEventTrigger) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	w.WriteHeader(http.StatusOK)

	if reqBytes, err := httputil.DumpRequest(r, true); err == nil {
		log.Printf("Message Dumper received a message: %+v", string(reqBytes))
	} else {
		log.Printf("Error dumping the request: %+v :: %+v", err, r)
		return
	}

	var pushEvent PushEvent

	body, err := ioutil.ReadAll(r.Body)
	if err != nil {
		log.Println(err)
		return
	}

	if err := json.Unmarshal(body, &pushEvent); err != nil {
		log.Println(err)
		return
	}

	if pushEvent.Repository.Name != pet.repoName {
		log.Printf("Push event from %s repo no action will be done", pushEvent.Repository.Name)
		return
	}

	err = pet.createPipelineRun(pushEvent.Repository.PushedAt)
	if err != nil {
		log.Println(err)
	}
}

// curl -vX POST http://localhost:8080 -d @./fixtures/dummy.json --header "Content-Type: application/json"
// ToDo we should write tests :D
func main() {
	// ToDo this should be flags
	pet := &PushEventTrigger{
		namespace:         "todo-app",
		pipelineRunPrefix: "todo-app",
		pipelineName:      "build-and-deploy-pipeline",
		repoName:          "todo-app-web",
	}

	log.Printf("Listening for repo: %s", pet.repoName)
	client, err := GetConfig()
	if err != nil {
		log.Fatal(err)
	}

	pet.client = client

	log.Println("Starting server")
	http.ListenAndServe(":8080", pet)
}
