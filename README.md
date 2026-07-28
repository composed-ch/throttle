# throttle

Limit requests to certain time intervals, both globally and per entity.

## Rationale

Consider a website with a contact form, which is mostly used by spammers, and occasionally by real users. There are two issues:

1. The server must not be overwhelmed by spam.
2. The real user must not be delayed by spam countermeasures.

This throttle library achieves both by defining two intervals:

1. A _global_ interval, after which a new request can be sent.
2. A _entity_-scoped interval, after which a particular user can send a new request.

This is achieved by defining a short global and a much longer by-entity interval.

The first token, both globally and per entity, is spawned immediately. So there is no waiting time for the initial user, and no waiting time beyond the global interval for each subsequent user.

The user only needs to await the _longer_ of the two intervals, nor both durations cumulated.

## Example

Define a `throttle` allowing requests coming in every second, but only once per minute for a single entity:

```go
t := throttle.New[*net.IP, *http.Request](time.Second, time.Minute)
```

The `throttle` identifies the entities by their IPv4 addresses and throttles HTTP requests.

Now consider a legitimate user in the following pseude-HTTP handling code:

```go
// demo boilerplate
ip := net.IPv4(192, 168, 1, 100)
req, _ := http.NewRequest(http.MethodGet, "/index.html", bytes.NewBufferString(""))

// http handling code
r, err := t.Await(&ip, &req, 5*time.Second)
if err == throttle.TimeoutError {
    fmt.Fprintf(os.Stderr, "timeout")
}
fmt.Println("serve request", r)
```

The request will be served without any delay, but then Alice needs to wait for a minute for her second request.

Check out the test cases in `throttle_test.go` for further understanding.
