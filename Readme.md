# GoWiki

A simple wiki website following https://go.dev/doc/articles/wiki/

Basically a demo of the `net/http` of Go.

## URLs

### FrontPage

`/` is redirected to `/view/FrontPage`

### Interfaces

- `/view/xxx` - view page xxx (jump to `edit/xxx` if not found)
- `/edit/xxx` - edit page xxx
- `/save/xxx` - save the page xxx and jump to `/view/xxx`

## Syntax

`[PageName]` can be used to refer to PageName

## build

```bash
$ go build -o bin/
$ ./wiki
```