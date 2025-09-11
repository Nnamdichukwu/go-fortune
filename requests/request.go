package requests

import (
	"fmt"
	"io"
	"net/http"
	"encoding/json"
)

type Request struct {
	Owner string `json:"owner"`
	Repo  string `json:"repo"`
}

type GithubRelease struct {
	Body string `json:"body"`
	Name string `json:"name"`
}
type GithubReleaseWithRepo struct{
	Owner string 
	Repo string
	Release GithubRelease
}

func ChangeLog(r Request) (GithubRelease, error){
	url := fmt.Sprintf("https://api.github.com/repos/%s/%s/releases/latest", r.Owner, r.Repo)

	resp, err := http.Get(url)
	if err != nil {
		return GithubRelease{}, err
	}

	defer resp.Body.Close()

	if resp.StatusCode !=  http.StatusOK {
		return GithubRelease{}, fmt.Errorf("Got the status code %d", resp.StatusCode)
	}

	body, err := io.ReadAll(resp.Body) 
	if err != nil{
		return GithubRelease{}, err 
	}

	var release GithubRelease
	if  err := json.Unmarshal(body, &release); err != nil{
		return GithubRelease{}, err 
	}
	return release, nil
	
}

