# logfile-open
wrap the os.OpenFile to respect the USR1 signal

```
go get github.com/6543/logfile-open@latest
```

```go
readWriteCloser, err := logfile.OpenFile("/tmp/some_file.log", 0o660)
```
