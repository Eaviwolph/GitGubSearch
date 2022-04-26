/*
 *This function is called when the user clicks on the button.
 *It will request a GET to the server with the query as URL parameter.
 *The response will be displayed in the div with id "results" by calling
 *displayGitInfo.
*/
function httpGet(url) {
    var xmlHttp = new XMLHttpRequest();
    var input = document.getElementById("input").value;

    xmlHttp.onreadystatechange = function () {
        if (this.readyState == 4 && this.status == 200 && this.responseText != "") {
            var myArr = JSON.parse(this.responseText);
            displayGitInfo(myArr["Repos"]);
        }
    };

    xmlHttp.open("GET", url + "?search=" + input, true);
    xmlHttp.send();
}

/*
 *Display the response from the server in the div with id "results".
*/
function displayGitInfo(map) {
    var m = Object.values(map);
    var i = 0;
    var ret = ""
    while (i < m.length) {
        var o = Object.values(m[i]);
        ret += "<div class='repo'><a class='links' href='https://github.com/" + o[1] + "'>" + o[1] + "</a><div class='languages'>";
        j = 0;
        while (j < o[2].length) {
            ret += "<p class='list'>" + o[2][j] + ": " + o[3][j] + "</p>";
            j++;
        }
        ret += "</div></div>";
        i++;
    }
    document.getElementById("results").innerHTML = ret
}