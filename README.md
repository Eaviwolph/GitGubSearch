# GitHubSearch
A GitHub Search server implemented in Go (backend) and html/css/js (frontend)  

__githubSearch.go__ contains a main that start the server and a command parser  
__gitAPISearch/gitAPISearch.go__ contains all searching logic and can be reused to get Data from GitHub
Html code is located in __template__  
All resources linked to templates are located in __web__

The program can be started by running the command  
`go run githubSearch.go`  
