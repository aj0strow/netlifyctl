# `netlifyctl`

Forked from the official non-official client. It works on **Circle CI** if you cross compile and store the binary in your project. Below are steps to reproduce the hack.

1. Fetch golang source code.

```
$ go get github.com/aj0strow/netlifyctl
$ cd netlifyctl
```

2. Cross compile for circle ci operating system from within the source code project folder. 

```
$ GOOS=linux GOARCH=amd64 go build .
```

3. You should have a `netlifyctl` binary now. Copy it to your project bin folder.

```
/project/bin/netlifyctl
```

4. Run the following command in circle yml with the full project id (not the name, the uuid) and a valid access token.

```
$ bin/netlifyctl deploy --site-id {uuid} --path dist --access-token {token}
````

Hope this helps!
