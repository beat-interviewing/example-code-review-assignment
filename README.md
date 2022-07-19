# BEAT.ly

BEAT.ly is a link shortening service for managing links and viewing click rate 
statistics. 

## Usage

### Shorten a link

Creating a short URL is as simple as issuing a `POST` request to the `/link`
endpoint. The `target` field is required and must be a valid URL. The `redirect`
field is optional and will default to a value of `302` if omitted.

```
POST /link HTTP/1.1
Host: beat.ly
Content-Type: application/json
Content-Length: 38

{
   "target": "https://example.com?foo=123#bar",
   "redirect": 307
}
```

```
HTTP/1.1 201 Created

{
    "id": "xyz",
    "url": "https://beat.ly/xyz",
    "target": "https://example.com?foo=123#bar",
    "redirect": 307
}
```

### Visit a link

Visiting a link is as simple as entering the URL in a browser of your choice.
BEAT.ly will redirect you to the `target` link using the redirect method
specified during the creation process.

```
GET /xyz HTTP/1.1
Host: beat.ly
```

```
HTTP/1.1 307 Temporary Redirect
Location: https://example.com?foo=123#bar
```

## Build

Using the `go` toolchain, building is straightforward.

    go build . -o beatly

This should produce an executable file with the name `beatly`.

During development it may be more convenient to use `go run` instead.

    go run ./...

## Run

    ./beatly