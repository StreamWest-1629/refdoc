package main

import (
	"fmt"
	"net/http"
	"net/url"
	"path/filepath"
	"regexp"
	"time"

	"github.com/pkg/errors"
	"golang.org/x/net/html"
)

type (
	SourceConfig struct {
		Title       string     `yaml:"title"`
		Description string     `yaml:"desc"`
		Mark        string     `yaml:"mark"`
		Categories  []Category `yaml:"categories"`
	}

	Category struct {
		Name        string      `yaml:"name"`
		Description string      `yaml:"desc"`
		Refs        []Reference `yaml:"refs"`
	}

	Reference struct {
		Title       string    `yaml:"title"`
		Link        string    `yaml:"link"`
		Good        int       `yaml:"-"`
		Description string    `yaml:"desc"`
		FoundAt     time.Time `yaml:"-"`
	}

	CacheConfig struct {
		ResetCache bool             `json:"-"`
		Caches     map[string]Cache `json:"caches"`
	}

	Cache struct {
		URL     string    `json:"url"`
		Title   string    `json:"title"`
		FoundAt time.Time `json:"foundAt"`
	}
)

var (
	regexpLinkExpr = `^\s*([+]*)\s*(https?://[\w!?/+\-_~;.,*&@#$%()'[\]]+)\s*$`
	regexpLink     = regexp.MustCompile(regexpLinkExpr)
)

func (r *Reference) parseLink(expr string) error {
	matched := regexpLink.FindStringSubmatch(expr)
	if matched == nil {
		return errors.New(fmt.Sprintf("unmatched regexp: (want: %s", regexpLinkExpr))
	}

	r.Good, r.Link = 0, ""
	r.Good = len(matched[1])
	r.Link = matched[2]
	if r.Link == "" {
		return errors.Errorf("invalid link: but regexp matched: %v", matched)
	}
	return nil
}

func (r *Reference) UnmarshalYAML(unmarshal func(interface{}) error) error {

	link := ""
	if err := unmarshal(&link); err == nil {
		return r.parseLink(link)
	}

	ref := map[string]string{}
	if err := unmarshal(&ref); err != nil {
		return err
	}

	r.Description = ref["desc"]
	return r.parseLink(ref["link"])
}

func (r *SourceConfig) Fill(cache *CacheConfig) error {

	err := error(nil)
	num, finished := 0, 0
	ticker := time.NewTicker(time.Second)
	for i := range r.Categories {
		if r.Categories[i].Name == "" {
			return errors.New("category's name (categories[*].name) is required")
		}
		for j := range r.Categories[i].Refs {
			num++
			go func(i, j int) {
				defer func() { finished++ }()
				err = r.Categories[i].Refs[j].Fill(cache)
			}(i, j)
		}
	}

	for {
		<-ticker.C
		if finished < num {
			logger.Info("reading web pages",
				"status", fmt.Sprintf("%06d / %06d", finished, num))
		} else {
			return err
		}
	}
}

func (ref *Reference) Fill(cache *CacheConfig) error {
	cached, exist := cache.Caches[ref.Link]

	if !exist || cached.Title == "" || cache.ResetCache {
		title, err := getHtmlTitle(ref.Link)

		if title == "" && err != nil {
			return err
		} else if err != nil {
			logger.Error(err, "failed to get html title", "link", ref.Link)
		}

		ref.Title = title
		if exist {
			ref.FoundAt = cached.FoundAt
		} else {
			ref.FoundAt = time.Now().UTC()
		}

		cache.Caches[ref.Link] = Cache{
			URL:     ref.Link,
			Title:   ref.Title,
			FoundAt: ref.FoundAt,
		}
	} else {
		*ref = Reference{
			Title:       cached.Title,
			FoundAt:     cached.FoundAt,
			Link:        ref.Link,
			Description: ref.Description,
			Good:        ref.Good,
		}
	}

	return nil
}

func getHtmlTitle(link string) (string, error) {

	response, err := http.Get(link)
	if err != nil {
		title, urlErr := getHtmlTitleFromUrl(link)
		if urlErr != nil {
			return "", urlErr
		}
		return title, errors.WithStack(err)
	}
	defer response.Body.Close()

	doc, err := html.Parse(response.Body)
	if err != nil {
		title, urlErr := getHtmlTitleFromUrl(link)
		if urlErr != nil {
			return "", urlErr
		}
		return title, errors.WithStack(err)
	}

	if title, exist := searchHtmlTitle(doc); exist {
		return title, nil
	}

	return getHtmlTitleFromUrl(link)
}

func getHtmlTitleFromUrl(link string) (string, error) {
	parsed, err := url.Parse(link)
	if err != nil {
		return "", errors.WithStack(err)
	}
	return filepath.Join(parsed.Host, parsed.RawPath), nil
}

func searchHtmlTitle(node *html.Node) (string, bool) {
	if node.Type == html.ElementNode && node.Data == "title" {
		return node.FirstChild.Data, true
	}

	for child := node.FirstChild; child != nil; child = child.NextSibling {
		if result, exist := searchHtmlTitle(child); exist {
			return result, true
		}
	}

	return "", false
}

type References []Reference

func (a References) Len() int      { return len(a) }
func (a References) Swap(i, j int) { a[i], a[j] = a[j], a[i] }
func (a References) Less(i, j int) bool {
	if a[i].Good != a[j].Good {
		return a[i].Good > a[j].Good
	} else if !a[i].FoundAt.Equal(a[j].FoundAt) {
		return a[i].FoundAt.After(a[j].FoundAt)
	} else {
		return a[i].Link < a[j].Link
	}
}
