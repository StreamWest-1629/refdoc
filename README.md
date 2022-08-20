# refdoc
The CLI tool for generate reference documentation

## TL;DR
### 1. create config file (current: [./refdoc.yaml])
```yaml
title: Informative Articles
desc: |
  This document presents pages that were useful in the development.
  このドキュメントは開発を進めるうえで有用だったページを紹介しています．
mark: 🥰
categories:
  - name: Github Action
    refs:
      - link: + https://docs.github.com/ja/actions/creating-actions/creating-a-docker-container-action
        desc: This is Github official document

  - name: Golang
    refs:
      - link: https://abhinavg.net/posts/flexible-yaml/
      - ++ https://pkg.go.dev/github.com/wi2L/fizz/markdown
```

### 2, run action `refdoc` (testing: [./.github/workflows/check_it.yaml])
```yaml
steps:
  - uses: actions/checkout@v3
  - uses: streamwest-1629/refdoc@0.1.0
    with:
      refdoc: refdoc.yaml
      cache: refdoc/cache.json   # default is refdoc/cache.json
      markdown: refdoc/README.md # default is refdoc/README.md
```
### 3. markdown file generated! (current: [./refdoc/README.md])
```md
# Informative Articles

This document presents pages that were useful in the development.</br>このドキュメントは開発を進めるうえで有用だったページを紹介しています．</br>

## Github Action

* 🥰[Docker コンテナのアクションを作成する - GitHub Docs](https://docs.github.com/ja/actions/creating-actions/creating-a-docker-container-action)</br>This is Github official document

   
## Golang

* 🥰🥰[markdown package - github.com/wi2L/fizz/markdown - Go Packages](https://pkg.go.dev/github.com/wi2L/fizz/markdown)</br>
* [How to write flexible YAML shapes in Go](https://abhinavg.net/posts/flexible-yaml/)</br>
```

## Configuration in Yaml file
### title
The title for generated markdown file (default: `IA - Informative Articles`).

### desc
The description for generated markdown file.

### mark
This is star mark, you can add articles more useful than others (default: ⭐)

### categories[*].name
The category name. this value is required.

### categories[*].refs
The link or description for informative articles. You can choice short expression and long expression to add it.

By the default, title is gotten from article page.

#### 1. short expression
```yaml
refs:
  - ++ https://abcdefg.hij
  - https://zxyvuts.rqp
```

firstly `++` expresses a number of star mark ⭐. description is default; means blank.

#### 2. full expression
```yaml
refs:
  - link: ++ https://abcdefg.hij
    desc: |
      This is description. 
      Write how this reference is helpful in your words.
```

## Configurations of Github actions
| input name | required | default | description
| :--: | :--: | :--: | :---
| refdoc | **Yes** | - | Filepath of your configuration file for refdoc (in this repository, [./refdoc.yaml]).
| cache | No | `./refdoc/cache.json` | Filepath of cache file in refdoc (in this repository, [./refdoc/cache.json]).</br>If it was empty, refdoc won't use cache, but cache file containing `foundAt` each your found references. We highly recommend for using cache file.
| markdown | No | `./refdoc/README.md` | Filepath of refdoc's generated markdown file (in this repository, [./refdoc/README.md]).

[./refdoc.yaml]: ./refdoc.yaml
[./.github/workflows/check_it.yaml]: ./.github/workflows/check_it.yaml
[./refdoc/README.md]: ./refdoc/README.md
[./refdoc/cache.json]: ./refdoc/cache.json