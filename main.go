package main

import (
	"bytes"
	"flag"
	"fmt"
	"os"
	"sort"
	"strings"
	"sync"

	"github.com/wI2L/fizz/markdown"
)

func main() {

	var (
		sourcePath, cachePath, destPath string
		cacheConfig                     = CacheConfig{
			Caches: map[string]Cache{},
		}
		srcConfig = SourceConfig{
			Title:       "IA - Informative Articles",
			Description: "",
			Mark:        "⭐",
		}
		wg = sync.WaitGroup{}
	)

	flag.StringVar(&sourcePath, "source", sourcePath, "The filepath of source file")
	flag.StringVar(&cachePath, "cache", "", "The filepath of cache file")
	flag.StringVar(&destPath, "dest", "./refdoc/README.md", "The filepath of destination file")
	flag.Parse()

	wg.Add(2)
	go func() {
		defer wg.Done()
		if err := ReadSourceYaml(sourcePath, &srcConfig); err != nil {
			logger.Error(err, "cannot read source file", "path", sourcePath)
			os.Exit(1)
		}
	}()

	go func() {
		defer wg.Done()
		if cachePath == "" {
			return
		}
		if err := ReadCacheJson(cachePath, &cacheConfig); err != nil {
			if os.IsNotExist(err) {
				logger.Info("not found cache file", "path", cachePath)
			} else {
				logger.Error(err, "cannot read cache file", "path", cachePath)
			}
		}
	}()

	wg.Wait()
	if err := srcConfig.Fill(&cacheConfig); err != nil {
		logger.Error(err, "cannot fill web page's title")
	}

	wg.Add(1)
	go func() {
		defer wg.Done()
		if cachePath == "" {
			return
		}
		WriteCacheJson(cachePath, cacheConfig)
	}()
	defer wg.Wait()

	md := bytes.NewBufferString(buildMarkdown(srcConfig))
	if err := WriteDestMarkdown(destPath, md); err != nil {
		logger.Error(err, "cannot write markdown file", "path", destPath)
	}
}

func buildMarkdown(srcConfig SourceConfig) string {

	builder := &markdown.Builder{}
	if srcConfig.Title != "" {
		builder = builder.H1(srcConfig.Title)
	}
	if srcConfig.Description != "" {
		builder = builder.P(replaceLF(srcConfig.Description))
	}

	for _, category := range srcConfig.Categories {
		builder = builder.H2(category.Name)
		if category.Description != "" {
			builder.P(replaceLF(category.Description))
		}

		sort.Sort(References(category.Refs))
		lists := []interface{}{}

		for _, ref := range category.Refs {
			internal := fmt.Sprintf("%s%s</br>%s",
				strings.Repeat(srcConfig.Mark, ref.Good),
				builder.Link(ref.Link, ref.Title),
				// ref.FoundAt.Format(time.ANSIC),
				replaceLF(ref.Description),
			)
			lists = append(lists, internal)
		}

		builder.BulletedList(lists...)
		builder.BR()
	}

	artifact := builder.String()
	logger.Info("success to build markdown file",
		"size", fmt.Sprintf("%d KB", len(artifact)/1024))
	return artifact
}

func replaceLF(plain string) string {
	return strings.ReplaceAll(plain, "\n", "</br>")
}
