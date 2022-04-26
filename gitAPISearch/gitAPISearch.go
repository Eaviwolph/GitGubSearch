package gitAPISearch

import (
	"encoding/json"
	"io/ioutil"
	"log"
	"net/http"
	"sync"
)

//Structure Definition.
type getSearchResult struct {
	/*
	 *Structure that will be encoded and sent to the client as JSON.
	 */
	Repos []gitRepos
}

type gitRepos struct {
	/*
	 *Structure that contains name and languages of the repository.
	 */
	Name           string
	Full_name      string
	Languages      []string
	LanguagesCount []int64
}

type tmpGitRepos struct {
	/*
	 *Sub Structure that is build with the github query and is
	 *contained in the gitGlob Structure.
	 */
	Name          string
	Full_name     string
	Languages_url string
}

type gitGlob struct {
	/*
	 *Structure that is build with the github query.
	 */
	Items []tmpGitRepos
}

/*
 *Shared Variables used by the multi threading.
 */
var threadNum = 10
var threads = make([][]gitRepos, threadNum)

/*
 *Function that transform the tmpGitRepos structure to the gitRepos
 *structure by requesting GitHub all languages used by the repository.
 */
func tmpGitToGitRepos(tmpGit tmpGitRepos) gitRepos {
	var gitRepos = gitRepos{
		Name:           tmpGit.Name,
		Full_name:      tmpGit.Full_name,
		Languages:      make([]string, 0),
		LanguagesCount: make([]int64, 0),
	} //Struct instantiation.

	var resp, err = http.Get(tmpGit.Languages_url)
	if err != nil {
		return gitRepos
	}
	bytes, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		return gitRepos
	}
	var d interface{}
	json.Unmarshal(bytes, &d)
	var msg = d.(map[string]interface{})
	for k, v := range msg {
		if _, ok := v.(float64); !ok {
			/*
			 *If the value is not a float64, then it's a string and
			 *so we have a GitHub quota exceeded.
			 */
			log.Println("Quota Exceeded")
			gitRepos.Full_name = "!GitHub Quota Exceeded!"
			gitRepos.Name = "!GitHub Quota Exceeded!"
			return gitRepos
		}
		gitRepos.Languages = append(gitRepos.Languages, k)
		gitRepos.LanguagesCount = append(gitRepos.LanguagesCount, int64(v.(float64)))
	}
	return gitRepos
}

/*
 *Function that return the gitRepos structure by requesting GitHub
 *and converting the tmpGitRepos structure to the gitRepos structure.
 */
func gitSearch(filter string) []gitRepos {
	url := "https://api.github.com/search/repositories?q=" + filter + "&page=1&per_page=20&sort=update"
	var resp, err = http.Get(url)
	if err != nil {
		log.Println(err)
		var list = make([]gitRepos, 0)
		return list
	}
	var data gitGlob
	if err := json.NewDecoder(resp.Body).Decode(&data); err != nil { //Decode the body of the response into the data object
		log.Println(err)
		var list = make([]gitRepos, 0)
		return list
	}

	var list = make([]gitRepos, 0)
	var i = 0 //Counter for the number of started threads
	var wg = sync.WaitGroup{}
	for _, v := range data.Items {
		if i >= len(threads) {
			/*If the number of started threads is greater than the number of threads,
			 *we wait for all the threads to finish.
			 *We copy all the threads into 'list' to avoid concurrency issues.
			 *We reset the counter to 0 to start the next batch of threads.
			 */
			wg.Wait()
			for _, v := range threads {
				list = append(list, v...)
			}
			i = 0
			threads = make([][]gitRepos, threadNum)
		}

		wg.Add(1) //Start a new thread
		go func(tmpGit tmpGitRepos, x int) {
			defer wg.Done()
			threads[x] = append(threads[x], tmpGitToGitRepos(tmpGit))
		}(v, i)
		i++
	}
	if i > 0 {
		/*If there are still threads running, we wait for them to finish
		 *and copy all the threads into 'list'.
		 */
		wg.Wait()
		var j = 0
		for j < i {
			list = append(list, threads[j]...)
			j++
		}
		i = 0
		threads = make([][]gitRepos, threadNum)
	}
	return list
}

/*
 *That call the function gitSearch and send the result to the client
 *if the url request contains a query.
 */
func GetSearch(w http.ResponseWriter, r *http.Request) {
	keys, errQuery := r.URL.Query()["search"]
	var key = ""
	if !errQuery || len(keys[0]) < 1 {
		log.Println("Url Param 'search' is missing or invalid")
		return
	}
	key = keys[0]
	var sr = getSearchResult{}
	sr.Repos = gitSearch(key)
	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(sr)
}
