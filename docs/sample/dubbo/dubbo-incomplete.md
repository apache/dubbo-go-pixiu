# Dubbo incomplete

## Args will nil

> fix in dubbogo 1.5.4

When the response is `User`, see below code:

```go
type User struct {
	ID   string
	Name string
	Age  int32
	Time time.Time
}
```

Although User struct has Time value, generic invoke will return nil. [the simple response](dubbo-query.md#simple) time field is disappear. 

So I suggest you can use string to time type for a short time.

[Previous](dubbo.md)
