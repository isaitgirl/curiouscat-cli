Build:
```
go mod tidy
go build -o curiouscat-cli main.go
```

Execução:
```shell
./curiouscat-cli -username foo -limit 100
```

Se informar `-limit 0`, a CLI vai chamar a API até encerrar todos posts.
