# GoWiki

A simple wiki website following https://go.dev/doc/articles/wiki/

Basically a demo of the `net/http` of Go.

## interfaces

- `/view/xxx` - view page xxx (jump to `edit/xxx` if not found)
- `/edit/xxx` - edit page xxx
- `/save/xxx` - save the page xxx and jump to `/view/xxx`

## build

```bash
$ go build wiki.go
$ ./wiki
```