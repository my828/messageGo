package handlers

import (
	"encoding/json"
	"fmt"
	"io"
	"log"
	"net/http"
	"net/url"
	"path"
	"strconv"
	"strings"

	"golang.org/x/net/html"
)

//PreviewImage represents a preview image for a page
type PreviewImage struct {
	URL       string `json:"url,omitempty"`
	SecureURL string `json:"secureURL,omitempty"`
	Type      string `json:"type,omitempty"`
	Width     int    `json:"width,omitempty"`
	Height    int    `json:"height,omitempty"`
	Alt       string `json:"alt,omitempty"`
}

//PageSummary represents summary properties for a web page
type PageSummary struct {
	Type        string          `json:"type,omitempty"`
	URL         string          `json:"url,omitempty"`
	Title       string          `json:"title,omitempty"`
	SiteName    string          `json:"siteName,omitempty"`
	Description string          `json:"description,omitempty"`
	Author      string          `json:"author,omitempty"`
	Keywords    []string        `json:"keywords,omitempty"`
	Icon        *PreviewImage   `json:"icon,omitempty"`
	Images      []*PreviewImage `json:"images,omitempty"`
}

const headerCORS = "Access-Control-Allow-Origin"
const corsAnyOrigin = "*"

//SummaryHandler handles requests for the page summary API.
//This API expects one query string parameter named `url`,
//which should contain a URL to a web page. It responds with
//a JSON-encoded PageSummary struct containing the page summary
//meta-data.
func SummaryHandler(w http.ResponseWriter, r *http.Request) {
	/*TODO: add code and additional functions to do the following:
	- Add an HTTP header to the response with the name
	 `Access-Control-Allow-Origin` and a value of `*`. This will
	  allow cross-origin AJAX requests to your server.
	- Get the `url` query string parameter value from the request.
	  If not supplied, respond with an http.StatusBadRequest error.
	- Call fetchHTML() to fetch the requested URL. See comments in that
	  function for more details.
	- Call extractSummary() to extract the page summary meta-data,
	  as directed in the assignment. See comments in that function
	  for more details
	- Close the response HTML stream so that you don't leak resources.
	- Finally, respond with a JSON-encoded version of the PageSummary
	  struct. That way the client can easily parse the JSON back into
	  an object. Remember to tell the client that the response content
	  type is JSON.

	Helpful Links:
	https://golang.org/pkg/net/http/#Request.FormValue
	https://golang.org/pkg/net/http/#Error
	https://golang.org/pkg/encoding/json/#NewEncoder
	*/
	w.Header().Add(headerCORS, corsAnyOrigin)
	w.Header().Add("Content-Type", "application/json")
	URL := r.FormValue("url")
	if URL == "" {
		http.Error(w, fmt.Sprintf("Bad Request! Try something else! %v", http.StatusBadRequest), 400)
	}
	html, err := fetchHTML(URL)
	if err != nil {
		http.Error(w, fmt.Sprintf("Bad Request! Try something else! %v", err), 400)
	}

	pageSummary, err := extractSummary(URL, html)
	if err != nil {
		http.Error(w, fmt.Sprintf("Bad Request! %v", err), 500)
	}
	enc := json.NewEncoder(w)
	if enc.Encode(pageSummary); err != nil {
		fmt.Printf("error fetching HTML %v", err)
	}
	defer html.Close()
}

//fetchHTML fetches `pageURL` and returns the body stream or an error.
//Errors are returned if the response status code is an error (>=400),
//or if the content type indicates the URL is not an HTML page.
func fetchHTML(pageURL string) (io.ReadCloser, error) {
	/*TODO: Do an HTTP GET for the page URL. If the response status
	code is >= 400, return a nil stream and an error. If the response
	content type does not indicate that the content is a web page, return
	a nil stream and an error. Otherwise return the response body and
	no (nil) error.

	To test your implementation of this function, run the TestFetchHTML
	test in summary_test.go. You can do that directly in Visual Studio Code,
	or at the command line by running:
		go test -run TestFetchHTML

	Helpful Links:
	https://golang.org/pkg/net/http/#Get
	*/
	resp, err := http.Get(pageURL)
	if err != nil {
		return nil, err
	}
	if resp.StatusCode >= 400 {
		return nil, fmt.Errorf("Status code error: %d", resp.StatusCode)
	}
	ctype := resp.Header.Get("Content-Type")
	if !strings.HasPrefix(ctype, "text/html") {
		return nil, fmt.Errorf("Content type error: %s", ctype)
	}
	return resp.Body, nil
}

// convert relative path to absolute path
func makeAbsPath(path, pageURL string) string {
	parsedURL, err := url.Parse(path)
	if err != nil {
		fmt.Errorf("error parsing path: %v\n", err)
	}
	base, err := url.Parse(pageURL)
	if err != nil {
		fmt.Errorf("error parsing URL: %v\n", err)
	}
	return base.ResolveReference(parsedURL).String()
}

// //extractSummary tokenizes the `htmlStream` and populates a PageSummary
// //struct with the page's summary meta-data.
func extractSummary(pageURL string, htmlStream io.ReadCloser) (*PageSummary, error) {
	// 	/*TODO: tokenize the `htmlStream` and extract the page summary meta-data
	// 	according to the assignment description.

	// 	To test your implementation of this function, run the TestExtractSummary
	// 	test in summary_test.go. You can do that directly in Visual Studio Code,
	// 	or at the command line by running:
	// 		go test -run TestExtractSummary

	// 	Helpful Links:
	// 	https://drstearns.github.io/tutorials/tokenizing/
	// 	http://ogp.me/
	// 	https://developers.facebook.com/docs/reference/opengraph/
	// 	https://golang.org/pkg/net/url/#URL.ResolveReference
	// 	*/

	tokenizer := html.NewTokenizer(htmlStream)
	images := []*PreviewImage{}
	ps := new(PageSummary)
	for {
		tokenType := tokenizer.Next()
		token := tokenizer.Token()
		tag := token.Data
		if tokenType == html.ErrorToken {
			err := tokenizer.Err()
			if err == io.EOF || tag == "/head" {
				break
			}
			log.Fatalf("error tokenizing HTML: %v", tokenizer.Err())
		}

		if tokenType == html.StartTagToken || tokenType == html.SelfClosingTagToken {
			if "meta" == tag || "link" == tag {
				for _, attr := range token.Attr {
					proKey := attr.Key
					proVal := attr.Val
					if proKey == "property" || proKey == "name" {
						for _, attr := range token.Attr {
							//ps.Type = "this is test"
							//fmt.Println(ps.Type)
							conKey := attr.Key
							conVal := attr.Val
							//fmt.Println("I am in content: " + conKey)
							if conKey == "content" {
								switch proVal {
								case "og:type":
									ps.Type = conVal
								case "og:url":
									ps.URL = conVal
								case "og:title":
									ps.Title = conVal
								case "og:site_name":
									ps.SiteName = conVal
								case "og:description":
									ps.Description = conVal
								case "description":
									if ps.Description == "" {
										ps.Description = conVal
									}
								case "og:image":
									image := new(PreviewImage)
									images = append(images, image)
									ps.Images = images
									abURL := conVal
									if path.IsAbs(abURL) {
										abURL = makeAbsPath(conVal, pageURL)
									}
									images[len(images)-1].URL = abURL
								case "og:image:secure_url":
									abURL := conVal
									if path.IsAbs(abURL) {
										abURL = makeAbsPath(conVal, pageURL)
									}
									images[len(images)-1].SecureURL = abURL
								case "og:image:type":
									images[len(images)-1].Type = conVal
								case "og:image:height":
									height, err := strconv.Atoi(conVal)
									if err != nil {
										return nil, fmt.Errorf("Cannot find height %v", err)
									}
									images[len(images)-1].Height = height
								case "og:image:width":
									width, err := strconv.Atoi(conVal)
									if err != nil {
										return nil, fmt.Errorf("Cannot find width %v", err)
									}
									images[len(images)-1].Width = width
								case "og:image:alt":
									images[len(images)-1].Alt = conVal
								case "author":
									for _, attr := range token.Attr {
										if attr.Key == "content" {
											ps.Author = attr.Val
											break
										}
									}
								case "keywords":
									for _, attr := range token.Attr {
										if attr.Key == "content" {
											keywords := strings.Split(attr.Val, ",")
											for _, keyword := range keywords {
												ps.Keywords = append(ps.Keywords, strings.Trim(string(keyword), " "))
											}
										}
									}
								}
							}
						}
					} else if proKey == "rel" && proVal == "icon" {
						icon := new(PreviewImage)
						for _, attr := range token.Attr {
							val := attr.Val
							switch attr.Key {
							case "href":
								abURL := val
								if path.IsAbs(abURL) {
									abURL = makeAbsPath(val, pageURL)
								}
								icon.URL = abURL
							case "sizes":
								if val != "any" {
									size := strings.Split(val, "x")
									height, err := strconv.Atoi(size[0])
									width, err := strconv.Atoi(size[1])
									if err != nil {
										return nil, fmt.Errorf("Cannot find height or width: %v", err)
									}
									icon.Height = height
									icon.Width = width
								}
							case "type":
								icon.Type = val
							}
						}
						ps.Icon = icon
					}
				}
			} else if tag == "title" {
				if ps.Title == "" {
					tokenType = tokenizer.Next()
					if tokenType == html.TextToken {
						ps.Title = tokenizer.Token().Data
					}
				}
			}
		}
	}
	return ps, nil
}
