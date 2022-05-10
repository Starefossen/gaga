package labels

import (
	"bytes"
	"encoding/json"
	"flag"
	"fmt"
	"github.com/varneberg/gaga/flags"
	"io/ioutil"
	"log"
	"net/http"
	"os"
	"strings"
	"time"
)

var (
	ghRepoOwner                  = os.Getenv("GITHUB_REPOSITORY_OWNER")
	ghRef                        = os.Getenv("GITHUB_REF")
	ghRefName                    = os.Getenv("GITHUB_REF_NAME")
	ghRepo                       = os.Getenv("GITHUB_REPOSITORY")
	ghToken                      = os.Getenv("GITHUB_TOKEN")
	ghEvent                      = os.Getenv("GITHUB_EVENT_NAME")
	ghActor                      = os.Getenv("GITHUB_ACTOR")
	ghWorkflow                   = os.Getenv("GITHUB_WORKFLOW")
	ghActionsIDTokenRequestURL   = os.Getenv("ACTIONS_ID_TOKEN_REQUEST_URL")
	ghActionsIDTokenRequestToken = os.Getenv("ACTIONS_ID_TOKEN_REQUEST_TOKEN")
	ghAPIURL                     = os.Getenv("GITHUB_API_URL")
)

// Returns URL to the active pull request
func GetPrUrl() string {
	prNumber := strings.Split(ghRefName, "/")[0]
	return ghAPIURL + "/repos/" + ghRepo + "/issues/" + prNumber + "/labels"
}

func getRepoUrl() string {
	// https://api.github.com/repos/OWNER/REPO/labels
	return ghAPIURL + "/repos/" + ghRepo + "/labels"
}

// colors
// orange : #D93F0B
type Label struct {
	Name        []string `json:"labels"`
	Description string   `json:"description,omitempty"`
	Color       string   `json:"color,omitempty"`
}

func parseLabel(label Label) []byte {
	rb, err := json.Marshal(label)
	if err != nil {
		log.Fatalln(err)
	}
	fmt.Println(string(rb))
	return rb
}

// Function for sending requests to the github API
func APIRequest(requestType string, url string, requestBody []byte) {
	timeout := time.Duration(5 * time.Second)
	client := &http.Client{
		Timeout: timeout,
	}
	//request, err := http.NewRequest("POST", url, bytes.NewBuffer(requestBody))
	request, err := http.NewRequest(requestType, url, bytes.NewBuffer(requestBody))
	request.Header.Add("Accept", "application/vnd.github.v3+json")
	request.Header.Add("Authorization", "token "+ghToken)
	request.Header.Set("Content-Type", "application/json")
	if err != nil {
		log.Fatalln(err)
	}
	resp, err := client.Do(request)
	if err != nil {
		log.Fatalln(err)
	}
	defer resp.Body.Close()

	body, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		log.Fatalln(err)
		os.Exit(2)
	}
	fmt.Printf(string(body))
	fmt.Println()
}

// Adds labels to current pull request
func addLabelPR(label Label) {
	requestBody := parseLabel(label)
	url := GetPrUrl()
	APIRequest("POST", url, requestBody)
}

func addLabelRepo(label Label) {

}

func LabelHandler(args []string) {
	var labelName flags.FlagSlice
	labelFlag := flag.NewFlagSet("label", flag.ExitOnError)
	labelFlag.Var(&labelName, "n", "Name of the labels")
	labelNewName := labelFlag.String("N", "", "Name new labels to add")
	labelDesc := labelFlag.String("d", "", "Description of labels, enclosed with \"\"")
	labelColor := labelFlag.String("c", "", "Color of labels")
	labelFlag.Parse(args)
	newLabel := Label{
		Name:        labelName,
		Description: *labelDesc,
		Color:       *labelColor,
	}
	fmt.Println(*labelNewName)
	//fmt.Println(newLabel)
	addLabelPR(newLabel)
}