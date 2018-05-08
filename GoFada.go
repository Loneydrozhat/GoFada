package main

// Usage
// ./GoFada

//
// Imports
//

import (
    "fmt"
    "net/http"
    "golang.org/x/net/html"
    "strings"
    "io/ioutil"
    "sort"
    "math/rand"
    "time"
)

//
// Types
//

type visit struct {
    using_depth int // Depth when link was visited
    visited bool // Whether link was visited or not
}

//
// Variables
//

var (
    base_url = "localhost" // Base url to start searching from
    scope_filter = "localhost" // Will limit search to urls starting with this
    depth = 2 // Maximum search depth (dom based)
    visited map[string]visit // Keeps visited urls/status
    wordlist = "none" // Path to wordlist for brute force path discovery
    discovered map[string]bool // Urls discovered via brute force, true => 2xx, false => non-2xx
    throttle = 0 // Throttle between consecutive http gets
    jitter = 0 // Throttle jitter
)

//
// Costants
//

const str_puke = "\n     (}\n    /Y\\`;,\n    /^\\  ;:,\n:::GoFada:::" // Banner

//
// Input handlers
//

// Reads option from command line
func read_option(option *int) {
    fmt.Scanf("%d", option)
}

//
// Operations
//

// Picks random number between min and max (to jitter)
func random(min, max int) int {
    
    rand.Seed(time.Now().Unix())
    return rand.Intn(max - min) + min // Intn panics if arg <= 0 
}

// Delay using throttle and jitter, handle Intn panic condition
func delay() {
    time.Sleep(time.Duration(throttle * int(time.Millisecond)))
    if jitter > 0 { // O/w Intn panics
        time.Sleep(time.Duration(random(0, jitter) * int(time.Millisecond)))
    }
}

// Restart visited urls map
func restart_visited() {
    visited = make(map[string]visit)
    v := visit{0, false}
    visited[base_url] = v
}

// Restart discovered urls map
func restart_discovered() {
    discovered = make(map[string]bool)
}

// Set GoFada parameters
func set_params() {

    keep_setting_params := true // Flag to control set_params loop
    for keep_setting_params {
        fmt.Printf("\n:::current params:::\n(1) base_url: %s\n(2) scope_filter: %s\n(3) depth: %d\n(4) wordlist: %s\n(5) throttle: %d ms\n(6) jitter: at most %d ms\n(*) back\n\ngfd|params>", base_url, scope_filter, depth, wordlist, throttle, jitter)
        param := -1
        read_option(&param)
        switch param {
            case 1:
                fmt.Printf("\n:::set new base_url:::\ngfd|set params|base_url>")
                fmt.Scanf("%s", &base_url)
                restart_visited()
            case 2:
                fmt.Printf("\n:::set new scope_filter:::\ngfd|set params|scope_filter>")
                fmt.Scanf("%s", &scope_filter)
            case 3:
                fmt.Printf("\n:::set new depth:::\ngfd|set params|depth>")
                fmt.Scanf("%d", &depth)
            case 4:
                fmt.Printf("\n:::set new wordlist:::\ngfd|set params|wordlist>")
                fmt.Scanf("%s", &wordlist)
            case 5:
                fmt.Printf("\n:::set new throttle:::\ngfd|set params|throttle>")
                fmt.Scanf("%d", &throttle)
            case 6:
                fmt.Printf("\n:::set new jitter:::\ngfd|set params|jitter>")
                fmt.Scanf("%d", &jitter);
            default:
                keep_setting_params = false
                fmt.Println("\n:::all set:::")
        }
    }
}

// Extracts links from given http response
func extract_links(resp *http.Response) []string {
    
    links := []string{}
    tokenizer := html.NewTokenizer(resp.Body)
    
    for {
    
        temp_token := tokenizer.Next()

        switch temp_token {
            case html.ErrorToken:
                // End of document, done
                return links
            case html.StartTagToken:
                token := tokenizer.Token()
                is_anchor := token.Data == "a"
                if is_anchor {
                    for _, a := range token.Attr {
                        if a.Key == "href" {
                            links = append(links, a.Val)
                            break
                        }
                    }
                }
        }
    }
    
}

// Crawl starting from url using depth first search
func crawl(url string, current_depth int) {

    fmt.Println(url)

    // Mark URL visited using current_depth
    v := visit{current_depth, true}
    visited[url] = v
    
    // Get page body
    resp, err := http.Get(url)
    delay()
    
    if err != nil {
        fmt.Println("\n:::http.Get(base_url) returned err not nil:::\n\n:::can't keep going, maybe check params?:::")
        return
    }
    
    // Extract links from page body
    links := extract_links(resp)
    resp.Body.Close()
    
    // Add links to map if they didn't already exist
    for _, link := range links {
        if _, exists := visited[link]; !exists { // Exists? If not do
            
            v := visit{999999, false} // Not visited, depth -> infinity
            visited[link] = v
        }   
    }
    
    // Check end depth, if so return else keep crawling
    if current_depth < depth {
        for _, link := range links {
            if visited[link].visited { // Visited 
            
                if visited[link].using_depth > current_depth + 1 { // Visited with greater depth => visit again with lower depth

                    if strings.HasPrefix(link, scope_filter) {
                        crawl(link, current_depth + 1)
                    } else if strings.HasPrefix(link, "/") { // relative links
                        crawl(url + link[1:], current_depth + 1)
                    }    
                }
                
            } else { // Not visited => visit using current_depth + 1
                
                if strings.HasPrefix(link, scope_filter) {
                    crawl(link, current_depth + 1)
                } else if strings.HasPrefix(link, "/") { // relative links
                    crawl(url + link[1:], current_depth + 1)
                }
            }
        }
    } 
}

// Shows visited urls (sorted))
func show_visited() {

    count_visited := 0
    
    fmt.Println("\n:::showing visited links:::\n")
    
    var keys []string
    for k := range visited {
        keys = append(keys, k)
    }
    sort.Strings(keys)

    for _, link := range keys {
        
        if visited[link].visited == true {
            count_visited++
            fmt.Println(link)   
        }        
    }
    
    fmt.Println("\n:::total", count_visited, "\b:::")

}

// Employs brute force to discover urls using given wordlist
func discover() {

    if wordlist == "none" {
    
        fmt.Println("\n:::can't go on with wordlist set to 'none':::\n\n:::maybe check wordlist path?:::")
        return
    }
    
    fmt.Printf("\n:::requesting %sxxx using wordlist at %s:::\n\n", base_url, wordlist)
    
    // Read wordlist
    list, err := ioutil.ReadFile(wordlist)
    if err != nil {
        
        fmt.Println("\n:::ioutil.ReadFile(wordlist) returned err not nil:::\n\n:::can't keep going, maybe check wordlist path?:::")
    }
    words := strings.Split(string(list), "\n")
    
    for _, word := range words {
        
        // Build URL
        url := base_url + word

        if url == base_url {
            continue
        }
        
        fmt.Println(url)
        
        resp, _ := http.Get(url)
        
        if resp.StatusCode >= 200 && resp.StatusCode < 299 { // Http 200 OK family since http.Get follows redirects

            discovered[word] = true
        } else {
            
            discovered[word] = false
        }
    }
    
    fmt.Println("\n:::done brute forcing:::")
}

// Shows discovered urls (from brute force search) - sorted
func show_discovered() { 
    
    count_discovered := 0
    
    fmt.Println("\n:::showing brute forced links:::\n")
    
    var keys []string
    for k := range discovered {
        keys = append(keys, k)
    }
    sort.Strings(keys)

    for _, link := range keys {
        
        if discovered[link] == true {
            count_discovered++
            fmt.Println(base_url + link)
        }
    }
    
    fmt.Println("\n:::total", count_discovered, "\b:::")
}

//
// Main
//

// Main function
func main() {

    keep_going := true
    fmt.Println(str_puke)
    
    for keep_going {
        fmt.Printf("\n:::pick an option:::\n(1) set params\n(2) crawl\n(3) discover\n(4) show crawled\n(5) show brute forced\n(0) quit\n(*) show puking stickman\n\ngfd>")
        option := -1
        read_option(&option)
        switch option {
            case 0:
                keep_going = false
            case 1:
                set_params()
            case 2:
                restart_visited()
                fmt.Println("\n:::crawling", base_url, "using scope_filter", scope_filter, "\b, depth", depth, "\b:::\n")
                if strings.HasPrefix(base_url, scope_filter) {
                    crawl(base_url, 0)
                } else {
                    fmt.Println("\n:::can't crawl anything, maybe check params?:::")
                }
                fmt.Println("\n:::done crawling:::")
            case 3:
                restart_discovered()
                discover()
            case 4:
                show_visited()
            case 5:
                show_discovered()
            default:
                fmt.Println(str_puke)
        }
    }
    
    fmt.Println("\n:::bye:::\n")
}
