
Generated via:

    go run ./genreflect -out genreflect/go18stdreflect/reflect.go -env GOOS=js -env GOARCH=wasm -exclude vendor/.* -exclude syscall -exclude syscall/js -exclude log/syslog std