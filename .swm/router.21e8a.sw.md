---
id: 21e8a
title: Router
file_version: 1.1.2
app_version: 1.3.8
---

## Introduction

In order to support dynamic routes, I changed the structure of the router from a simple map ( method-pattern: handler ) to a two-field-struct. The roots contain routes of each method. The handlers hold the same information as before.

<br/>

router struct
<!-- NOTE-swimm-snippet: the lines below link your snippet to Swimm -->
### 📄 gu/router.go
```go
9      type router struct {
10     	roots    map[string]*node       // roots key eg, roots['GET'] roots['POST']
11     	handlers map[string]HandlerFunc // handlers key eg, handlers['GET-/p/:lang/doc'], handlers['POST-/p/book']
12     }
13     
```

<br/>

## Main features

The value of the roots map is represented by a tree node that fits into a particular data structure called trie. As you can see in the node struct, each node holds its part and points to its children nodes.

If we read the tree from the root node to its leaf node, combine each part together. The result will be a complete pattern of a request URL path.

<br/>

`isWild`<swm-token data-swm-token=":gu/trie.go:9:1:1:`	isWild   bool`"/>means whether the match is a regular expression like /:name.
<!-- NOTE-swimm-snippet: the lines below link your snippet to Swimm -->
### 📄 gu/trie.go
```go
5      type node struct {
6      	pattern  string
7      	part     string
8      	children []*node
9      	isWild   bool
10     }
```

<br/>

<div align="center"><img src="https://firebasestorage.googleapis.com/v0/b/swimmio-content/o/repositories%2FZ2l0aHViJTNBJTNBR3UlM0ElM0FsdWxpYW42NjY%3D%2Fe0de1189-e56b-412b-a553-079228db94a4.png?alt=media&token=9bb04f36-7346-48b8-a5f2-612344bcd05f" style="width:'50%'"/></div>

<br/>

The picture above shows an example of a tree map from GET method. Parden my hand draws. So we have to figure out how to insert and search nodes. The code is below.

<br/>

Search and insert methods are both recursion.
<!-- NOTE-swimm-snippet: the lines below link your snippet to Swimm -->
### 📄 gu/trie.go
```go
31     func (n *node) insert(pattern string, parts []string, height int) {
32     	if len(parts) == height {
33     		n.pattern = pattern
34     		return
35     	}
36     	part := parts[height]
37     	child := n.matchChild(part)
38     	if child == nil {
39     		child = &node{
40     			part:   part,
41     			isWild: part[0] == ':' || part[0] == '*',
42     		}
43     		n.children = append(n.children, child)
44     	}
45     	child.insert(pattern, parts, height+1)
46     }
47     
48     func (n *node) search(parts []string, height int) *node {
49     	if len(parts) == height || strings.HasPrefix(n.part, "*") {
50     		if n.pattern == "" {
51     			return nil
52     		}
53     		return n
54     	}
55     	part := parts[height]
56     	children := n.matchChildren(part)
57     
58     	for _, child := range children {
59     		result := child.search(parts, height+1)
60     		if result != nil {
61     			return result
62     		}
63     	}
64     	return nil
65     }
```

<br/>

## Interface

<br/>

router methods

As we get routes, we save params into a map so we can use them in context later.
<!-- NOTE-swimm-snippet: the lines below link your snippet to Swimm -->
### 📄 gu/router.go
```go
35     func (r *router) addRoute(method string, pattern string, handler HandlerFunc) {
36     	log.Printf("Route %4s - %s", method, pattern)
37     	parts := parsePattern(pattern)
38     	key := method + "-" + pattern
39     	_, ok := r.roots[method]
40     	if !ok {
41     		r.roots[method] = &node{}
42     	}
43     	r.roots[method].insert(pattern, parts, 0)
44     	r.handlers[key] = handler
45     }
46     
47     func (r *router) getRoute(method string, path string) (*node, map[string]string) {
48     	searchParts := parsePattern(path)
49     	params := make(map[string]string)
50     	root, ok := r.roots[method]
51     
52     	if !ok {
53     		return nil, nil
54     	}
55     
56     	n := root.search(searchParts, 0)
57     
58     	if n != nil {
59     		parts := parsePattern(n.pattern)
60     		for index, part := range parts {
61     			if part[0] == ':' {
62     				params[part[1:]] = searchParts[index]
63     			}
64     			if part[0] == '*' && len(part) > 1 {
65     				params[part[1:]] = strings.Join(searchParts[index:], "/")
66     				break
67     			}
68     		}
69     		return n, params
70     	}
71     
72     	return nil, nil
73     }
74     
75     func (r *router) handle(c *Context) {
76     	n, params := r.getRoute(c.Method, c.Path)
77     	if n != nil {
78     		c.Params = params
79     		key := c.Method + "-" + n.pattern
80     		r.handlers[key](c)
81     	} else {
82     		c.String(http.StatusNotFound, "404 NOT FOUND: %s\n", c.Path)
83     	}
84     }
85     
```

<br/>

<br/>

<br/>

Add a new field to Context struct that represents params in URL
<!-- NOTE-swimm-snippet: the lines below link your snippet to Swimm -->
### 📄 gu/context.go
```go
11     type Context struct {
12     	Writer http.ResponseWriter
13     	R      *http.Request
14     
15     	Path       string
16     	Method     string
17     	StatusCode int
18     
19     	Params map[string]string
20     }
21     
```

<br/>

Add a method to get a specific param
<!-- NOTE-swimm-snippet: the lines below link your snippet to Swimm -->
### 📄 gu/context.go
```go
30     
31     func (c *Context) Param(key string) string {
32     	value, _ := c.Params[key]
33     	return value
34     }
35     
```

<br/>

## Directory structure

`📄 gu/trie.go` implements the tree node struct

`📄 gu/router.go` implements new router struct

`📄 gu/router_test.go` adds some tests for router

<br/>

tests
<!-- NOTE-swimm-snippet: the lines below link your snippet to Swimm -->
### 📄 gu/router_test.go
```go
9      func newTestRouter() *router {
10     	r := newRouter()
11     	r.addRoute("GET", "/", nil)
12     	r.addRoute("GET", "/hello/:name", nil)
13     	r.addRoute("GET", "/hello/b/c", nil)
14     	r.addRoute("GET", "/hi/:name", nil)
15     	r.addRoute("GET", "/assets/*filepath", nil)
16     	return r
17     }
18     
19     func TestParsePattern(t *testing.T) {
20     	ok := reflect.DeepEqual(parsePattern("/p/:name"), []string{"p", ":name"})
21     	ok = ok && reflect.DeepEqual(parsePattern("/p/*"), []string{"p", "*"})
22     	ok = ok && reflect.DeepEqual(parsePattern("/p/*/*name"), []string{"p", "*"})
23     	if !ok {
24     		t.Fatal("test parsePattern failed")
25     	}
26     }
27     
28     func TestGetRoute(t *testing.T) {
29     	r := newTestRouter()
30     	n, ps := r.getRoute("GET", "/hello/dj")
31     
32     	if n == nil {
33     		t.Fatal("test getRoute failed")
34     	}
35     
36     	if n.pattern != "/hello/:name" {
37     		t.Fatal("should match /hello/:name")
38     	}
39     
40     	if ps["name"] != "dj" {
41     		t.Fatal("name should be equal to 'dj")
42     	}
43     	fmt.Printf("matched path: %s, params['name']: %s\n", n.pattern, ps["name"])
44     }
45     
```

<br/>

This file was generated by Swimm. [Click here to view it in the app](https://app.swimm.io/repos/Z2l0aHViJTNBJTNBR3UlM0ElM0FsdWxpYW42NjY=/docs/21e8a).
