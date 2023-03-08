---
id: ibzb5
title: Router Group
file_version: 1.1.2
app_version: 1.3.8
---

## Introduction

Implemented two versions of router group. The first one is in the main branch, the other is in group\_v2 branch.

The only difference between the two versions is the router struct's location.

## Main features

In the former commit, this framework only has the Engine and router implemented. Where Engine struct only has one field: router. But now I want to add more features into this framework like groups which is useful when we need to group together some routes that share the same prefix or the same authorization.

So what should a group have? Prefix, middleware, router. How to get the ability of routers? See before we have our Engine struct with holds a router field. If we have a pointer to the engine we are able to use the engine's router.

## Interface

Let's take a look at the main branch version.

<br/>

`RouterGroup`<swm-token data-swm-token=":gu/gu.go:10:2:2:`type RouterGroup struct {`"/> and `Engine`<swm-token data-swm-token=":gu/gu.go:16:2:2:`type Engine struct {`"/> struct
<!-- NOTE-swimm-snippet: the lines below link your snippet to Swimm -->
### 📄 gu/gu.go
```go
10     type RouterGroup struct {
11     	prefix      string
12     	middlewares []HandlerFunc
13     	engine      *Engine
14     }
15     
16     type Engine struct {
17     	router *router
18     	*RouterGroup
19     	groups []*RouterGroup
20     }
21     
```

<br/>

Few methods of RouterGroup
<!-- NOTE-swimm-snippet: the lines below link your snippet to Swimm -->
### 📄 gu/gu.go
```go
29     func (group *RouterGroup) Group(prefix string) *RouterGroup {
30     	engine := group.engine
31     	newGroup := &RouterGroup{
32     		prefix: group.prefix + prefix,
33     		engine: engine,
34     	}
35     	engine.groups = append(engine.groups, newGroup)
36     	return newGroup
37     }
38     
39     func (group *RouterGroup) addRoute(method string, comp string, handler HandlerFunc) {
40     	pattern := group.prefix + comp
41     	log.Printf("Route %4s - %s", method, pattern)
42     	group.engine.router.addRoute(method, pattern, handler)
43     }
44     
45     func (group *RouterGroup) GET(pattern string, handler HandlerFunc) {
46     	group.addRoute("GET", pattern, handler)
47     }
48     
49     func (group *RouterGroup) POST(pattern string, handler HandlerFunc) {
50     	group.addRoute("POST", pattern, handler)
51     }
52     
```

<br/>

So when someone wants to use our framework, what one has to do is like below. Because our Engine has embedded RouterGroup, we case use engine struct to call Group() method without having to implicitly transfer it to Group struct.

<br/>


<!-- NOTE-swimm-snippet: the lines below link your snippet to Swimm -->
### 📄 main.go
```go
8      func main() {
9      	e := gu.New()
10     	e.GET("/", func(c *gu.Context) {
11     		c.HTML(http.StatusOK, "<h1>gu-library</h1>")
12     	})
13     
14     	v1 := e.Group("/v1")
15     	{
16     		v1.GET("/", func(c *gu.Context) {
17     			c.HTML(http.StatusOK, "<h1>Hello Gee</h1>")
18     		})
19     
20     		v1.GET("/hello", func(c *gu.Context) {
21     			c.String(http.StatusOK, "hello %s, you're at %s\n", c.Query("name"), c.Path)
22     		})
23     	}
24     
25     	v2 := e.Group("/v2")
26     	{
27     		v2.GET("/hello/:name", func(c *gu.Context) {
28     			c.String(http.StatusOK, "hello %s, you're at %s\n", c.Param("name"), c.Path)
29     		})
30     		v2.POST("/login", func(c *gu.Context) {
31     			c.JSON(http.StatusOK, gu.H{
32     				"username": c.PostForm("username"),
33     				"password": c.PostForm("password"),
34     			})
35     		})
36     	}
37     
38     	e.Run(":9999")
39     }
```

<br/>

## Version two

Now one thing doesn't feel natural though. If we take a look at the `RouterGroup`<swm-token data-swm-token=":gu/gu.go:10:2:2:`type RouterGroup struct {`"/> and `Engine`<swm-token data-swm-token=":gu/gu.go:16:2:2:`type Engine struct {`"/> struct, we'll see that they have each other embedded. Refresh our memory of why we embedded the Engine in the RouterGroup again. It's so that the router group has the ability of router( inside the engine ). Seems a bit unpractical not to embed the router directly.

So in version two, I moved the router from Engine to RouterGroup. Make sense if we think about their own duty. The Engine can ServeHttp, Run server and manage groups. The TouterGroup can manage its prefix, and its middleware, and add path: handlers to routes (the trie struct I implemented before).

Here's version two of `RouterGroup`<swm-token data-swm-token=":gu/gu.go:10:2:2:`type RouterGroup struct {`"/> and `Engine`<swm-token data-swm-token=":gu/gu.go:16:2:2:`type Engine struct {`"/> struct which you can also see in branch group\_v2.

```go
type RouterGroup struct {
	prefix      string
	middlewares []HandlerFunc
	*router
}

type Engine struct {
	*RouterGroup
	groups []*RouterGroup
}
```

Meanwhile, some of their methods are changed too. Note that this version does not require any change in `📄 main.go` file.

```go
func New() *Engine {
	group := &RouterGroup{router: newRouter()}
	engine := &Engine{
		RouterGroup: group,
		groups:      []*RouterGroup{group},
	}
	return engine
}

func (group *RouterGroup) Group(prefix string) *RouterGroup {
	newGroup := &RouterGroup{
		prefix: group.prefix + prefix,
		router: group.router,
	}
	return newGroup
}

func (group *RouterGroup) addRoute(method string, comp string, handler HandlerFunc) {
	pattern := group.prefix + comp
	log.Printf("Route %4s - %s", method, pattern)
	group.router.addRoute(method, pattern, handler)
}

func (group *RouterGroup) GET(pattern string, handler HandlerFunc) {
	group.addRoute("GET", pattern, handler)
}

func (group *RouterGroup) POST(pattern string, handler HandlerFunc) {
	group.addRoute("POST", pattern, handler)
}

func (e *Engine) Group(prefix string) *RouterGroup {
	group := e.RouterGroup.Group(prefix)
	e.groups = append(e.groups, group)
	return group
}
```

Though I still remain one question. Since I removed the embedded \*Engine pointer inside RouterGroup, the Engine struct has to implement its own Group() method now which, if you remember in version 1 we call the method Group() through Engine struct because it has RouterGroup embedded.

## Design decisions

In version 2, the Engine struct has its own Group() method in which it calls method Group() of RouterGroup. This is usually how we do in other-oriented languages. However, this doesn't feel quite like Go. version 1 feels quite "Go". Why? I think it's cause I don't see embedding the same thing as inheritance.

The reason Engine has its own Group() method is to save its groups field, adding the new group to this slice. For now, this field doesn't really have any usefulness. If we remove this field, we can actually get rid of Group() method. Let's see what it looks like after.

```go
type RouterGroup struct {
	prefix      string
	middlewares []HandlerFunc
	*router
}

type Engine struct {
	*RouterGroup
	//groups []*RouterGroup
}

func New() *Engine {
	group := &RouterGroup{router: newRouter()}
	engine := &Engine{
		RouterGroup: group,
		//groups:      []*RouterGroup{group},
	}
	return engine
}

func (group *RouterGroup) Group(prefix string) *RouterGroup {
	newGroup := &RouterGroup{
		prefix: group.prefix + prefix,
		router: group.router,
	}
	return newGroup
}

func (group *RouterGroup) addRoute(method string, comp string, handler HandlerFunc) {
	pattern := group.prefix + comp
	log.Printf("Route %4s - %s", method, pattern)
	group.router.addRoute(method, pattern, handler)
}

func (group *RouterGroup) GET(pattern string, handler HandlerFunc) {
	group.addRoute("GET", pattern, handler)
}

func (group *RouterGroup) POST(pattern string, handler HandlerFunc) {
	group.addRoute("POST", pattern, handler)
}

//func (e *Engine) Group(prefix string) *RouterGroup {
//	group := e.RouterGroup.Group(prefix)
//	e.groups = append(e.groups, group)
//	return group
//}
```

After this change, tests still pass and the "Gu" library still works. We don't have to change a word in `📄 main.go`. I think if we don't use it, we don't need it.

## When to embed

I was curious enough to do some research on when should embed inside a struct. Because I think everything finally traces back to the basics.

Take a look at the example from [gobyexample](https://gobyexample.com/struct-embedding).

> An embedding looks like a field without a name.
> 
> We can access the base’s fields directly on `co`, e.g. `co.num`.
> 
> Since `container` embeds `base`, the methods of `base` also become methods of a `container`.

I find a blog on [embedding](https://eli.thegreenplace.net/2020/embedding-in-go-part-1-structs-in-structs/) in go. It shows some examples of embedding in Go's standard libraries. I make a summary of the purpose of using embedding.

1.  More convenient to call/access methods
    
2.  To gain new behavior
    
3.  Implement interfaces
    

<br/>

## WHat about Gin

Let's take a look at how Gin implements this.

```go
// Engine is the framework's instance, it contains the muxer, middleware and configuration settings.
// Create an instance of Engine, by using New() or Default()
type Engine struct {
	RouterGroup

	RedirectTrailingSlash bool
	RedirectFixedPath bool
	HandleMethodNotAllowed bool
	ForwardedByClientIP bool
	AppEngine bool
	UseRawPath bool
	UnescapePathValues bool
	RemoveExtraSlash bool
	RemoteIPHeaders []string
	TrustedPlatform string
	MaxMultipartMemory int64
	UseH2C bool
	ContextWithFallback bool

	delims           render.Delims
	secureJSONPrefix string
	HTMLRender       render.HTMLRender
	FuncMap          template.FuncMap
	allNoRoute       HandlersChain
	allNoMethod      HandlersChain
	noRoute          HandlersChain
	noMethod         HandlersChain
	pool             sync.Pool
	trees            methodTrees
	maxParams        uint16
	maxSections      uint16
	trustedProxies   []string
	trustedCIDRs     []*net.IPNet
}
```

```go
// RouterGroup is used internally to configure router, a RouterGroup is associated with
// a prefix and an array of handlers (middleware).
type RouterGroup struct {
	Handlers HandlersChain
	basePath string
	engine   *Engine
	root     bool
}
```

We can see that it's very similar to our version 1. The Engine type has a RouterGroup embedded. The RouterGroup has a field engine that holds a pointer to the engine.

```go
func (group *RouterGroup) handle(httpMethod, relativePath string, handlers HandlersChain) IRoutes {
	absolutePath := group.calculateAbsolutePath(relativePath)
	handlers = group.combineHandlers(handlers)
	group.engine.addRoute(httpMethod, absolutePath, handlers)
	return group.returnObj()
}

// Handle registers a new request handle and middleware with the given path and method.
// The last handler should be the real handler, the other ones should be middleware that can and should be shared among different routes.
// See the example code in GitHub.
//
// For GET, POST, PUT, PATCH and DELETE requests the respective shortcut
// functions can be used.
//
// This function is intended for bulk loading and to allow the usage of less
// frequently used, non-standardized or custom methods (e.g. for internal
// communication with a proxy).
func (group *RouterGroup) Handle(httpMethod, relativePath string, handlers ...HandlerFunc) IRoutes {
	if matched := regEnLetter.MatchString(httpMethod); !matched {
		panic("http method " + httpMethod + " is not valid")
	}
	return group.handle(httpMethod, relativePath, handlers)
}

// POST is a shortcut for router.Handle("POST", path, handlers).
func (group *RouterGroup) POST(relativePath string, handlers ...HandlerFunc) IRoutes {
	return group.handle(http.MethodPost, relativePath, handlers)
}
```

Whenever we New() an Engine and call its GET() method, it actually calls the GET() method of RouterGroup and leads to an internal method name handle(). In which it calls the engine's addRoute() method.

```
func (engine *Engine) addRoute(method, path string, handlers HandlersChain) {
	assert1(path[0] == '/', "path must begin with '/'")
	assert1(method != "", "HTTP method can not be empty")
	assert1(len(handlers) > 0, "there must be at least one handler")

	debugPrintRoute(method, path, handlers)

	root := engine.trees.get(method)
	if root == nil {
		root = new(node)
		root.fullPath = "/"
		engine.trees = append(engine.trees, methodTree{method: method, root: root})
	}
	root.addRoute(path, handlers)

	// Update maxParams
	// ...
}
```

As we can see above, Gin's Engine also implements a tree structure to save path and handler information ( along with a bunch other information). At this point, I think RouterGroup is just an "outer source" and the Engine is still the main character that does most of the heavy work.

<br/>

<br/>

This file was generated by Swimm. [Click here to view it in the app](https://app.swimm.io/repos/Z2l0aHViJTNBJTNBR3UlM0ElM0FsdWxpYW42NjY=/docs/ibzb5).
