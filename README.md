# refdoc
The CLI tool for generate reference documentation

## LT;DR
### 1. create config file [(./refdoc.yaml)](./refdoc.yaml)
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

### 2, run command `refdoc`

### 3. markdown file generated [(./refdoc/README.md)](./refdoc/README.md)
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

#### 1. short expression
```yaml
refs:
  - ++ https://abcdefg.hij
  - https://zxyvuts.rqp
```

firstly `++` expresses a number of star mark ⭐. description is default; means blank.
