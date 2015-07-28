# cachegen

Uses "go generate" to convert the expiration cache template to
a cache that works for your package.

### Installation

Run

```bash
go get github.com/zabawaba99/cachegen
```

### Usage

In order to generate a new cache you must specify the type
key and value to cache in the comments in your code:

```go
//go:generate cachegen -key-type int -value-type MyStruct
```

This will create a new `MyStructCache` type that will allow you
to `Add`, `Get` and `Expire` MyStruct objects using an int. The file
will be called `mystruct_cache.go`

**Note:** The cache type and all it's functions are currently exported.

## Contributing

1. Fork it
2. Create your feature branch (`git checkout -b new-feature`)
3. Commit your changes (`git commit -am 'awesome things with tests'`)
4. Push to the branch (`git push origin new-feature`)
5. Create new Pull Request
