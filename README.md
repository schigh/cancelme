# cancelme

api skeleton demonstrating context cancellation

## How To

- clone repo
- get deps: `go mod download`
- run main `go run main.go`
    - you can pass a port param `-p`, e.g. `-p 8000`.  By default it runs on `8080`

## URL params

### all requests
- `depth` - (integer) - the depth of nested calls where a context is passed
- `stagepause` - the `time.Duration` that each stage will sleep between nested calls.

### synchronous requests
- `timeout` - the total request timeout, set as a `time.Duration`
- `cancels` - (must be `true` if set) - the context will cancel sometime during the request
- `canceldepth` - (integer) - the depth in the call stack where the context will cancel. `cancels` must be `true` for this to work

### asynchronous requests
- `async` - (must be `true` if set) - dispatches all parts of the call stack asynchronously
- `asyncdies` - (integer) - the Nth dispatched goroutine that will return an error and cause its error group to cancel.  Set this number >= `depth` to make all goroutines succeed

## sample outputs

```
http://localhost:8080/?depth=10&stagepause=3000000&canceldepth=2&cancels=true

2019/05/24 22:04:39 depth: 10
2019/05/24 22:04:39 cancels: true
2019/05/24 22:04:39 canceldepth: 2
2019/05/24 22:04:39 stage pause: 3ms
2019/05/24 22:04:39 dispatching baz 9
2019/05/24 22:04:39 dispatching child baz 8
2019/05/24 22:04:39 dispatching child baz 7
2019/05/24 22:04:39 context canceled
```

```
http://localhost:8080/?depth=10&stagepause=300000000&timeout=1500000000

2019/05/24 22:06:32 timeout: 1.5s
2019/05/24 22:06:32 depth: 10
2019/05/24 22:06:32 stage pause: 300ms
2019/05/24 22:06:32 dispatching baz 9
2019/05/24 22:06:32 dispatching child baz 8
2019/05/24 22:06:32 dispatching child baz 7
2019/05/24 22:06:32 dispatching child baz 6
2019/05/24 22:06:33 dispatching child baz 5
2019/05/24 22:06:33 dispatching child baz 4
2019/05/24 22:06:33 context deadline exceeded
```

```
http://localhost:8080/?depth=10&stagepause=300000000&async=true&asyncdies=2

2019/05/24 22:08:41 async: true
2019/05/24 22:08:41 asyncdies: 2
2019/05/24 22:08:41 depth: 10
2019/05/24 22:08:41 stage pause: 300ms
2019/05/24 22:08:41 async #5 succeeded after 6.514771ms
2019/05/24 22:08:41 async #0 succeeded after 27.980978ms
2019/05/24 22:08:41 async #3 succeeded after 32.210818ms
2019/05/24 22:08:41 async #2 dies after 58.054053ms
2019/05/24 22:08:41 async #8 terminated by context after 97.713919ms: context canceled
2019/05/24 22:08:41 async #7 terminated by context after 115.232382ms: context canceled
2019/05/24 22:08:41 async #1 terminated by context after 130.395949ms: context canceled
2019/05/24 22:08:41 async #9 terminated by context after 241.168034ms: context canceled
2019/05/24 22:08:41 async #4 terminated by context after 249.547137ms: context canceled
2019/05/24 22:08:41 async #6 terminated by context after 278.306873ms: context canceled
2019/05/24 22:08:41 dies
```


