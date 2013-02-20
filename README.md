ironme
=========

A supercharged development webserver. 

## Introduction

This tool is born after I realize how much I miss the Rails asset pipeline when working on a project not using Rails (especially single page app).  
When working on simple project with only HTML, CSS and Javascript there is no real problem, but as soon as you start adding technology like SASS or CoffeeScript in the mix it become much more difficult.  
I really dislike my workflow in that situation, relying on grunt tasks, livereload, ruby/python/node.js scripts, so I decided to write this tool which hopefully should have everything I need when working on that kind of project.  

## Installation

**Dependencies :**

* [Go](http://golang.org/doc/install)
* [libsass](https://github.com/hcatlin/libsass), if you want SASS support
* [CoffeeScript](http://coffeescript.org/#installation), if you want CoffeeScript support

Right the only way to install is via Go tool :

```bash
$ go get github.com/stevedomin/ironme
```

## Usage

```bash
$ ironme # Current directory served on default port (4000)
$ ironme ~/mypath -p 3000 # ~/mypath served on port 3000
```

## TODO

* Handle 
* Write tests
* Flag to disable filter (--no-sass, --no-coffee)
* Live reload v1 (watch files and reload page on change)
* Cache compiled files
* Read config from ./.ironme/config
* Store server log in ./.ironme/logs
* Live reload v2 (watch files and refresh page on change)
