
# Intro:

Gofoto is web-based photo viewer written in Go. It starts a light-weight web server that enables you to browse & view images stored on the local filesystem in your browser. It also provides some basic opeations like resizing or rotating images. You can root `gofoto` at a file system location by supplying the fully qualitifed path as the command line argument. If no arguments are supplied, it roots at user's home directory available via the `HOME` env var.

Currently only CR2 (Canon Raw) format & JPEG encoding is supported.

# Instructions:

## Build:

```
% make build 
```

## Run:

The `run` target builds & executes the `gofoto` server rooted at the filesystem path `$HOME`.

```
% make build run 
```

To specify a different path, pass it in as command-line parameter.

```
% ./gofoto /path/to/start
```

If no path is supplied, it will set root at current user's home pointed by `$HOME` env var

## How to use:

After the server has started successfuly, point your browser to `http://localhost:8080/`

The browser view displays a list of all sub-directories, a link to go `~up~` a level or view the pictures in the current directory `$ photos $`

<br/><br/>

![](doc/browse1.png)

<br/><br/>

The photo viewer page lets you view, resize or rotate images using the link at the bottom of the page.
<br/><br/>
![](doc/view.png)




