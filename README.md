# Практическое приложение для курса Golang5 Pro

Разработано для школы программирования Golang5

## Микросервисная система Movix


## GRPC

### Windows

Скачайте скомпилируемые бинарные файлы и добавьте в путь поиска исполнительных файлов

## Запуск кодогенерации для GRPC

```bash
cd api
protoc --go_out=./gen --go_opt=paths=source_relative --go-grpc_out=./gen --go-grpc_opt=paths=source_relative movie.proto
```