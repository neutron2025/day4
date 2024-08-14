### Generate ECDSA256 Private Key if not existed

- 1、Install OpenSSL in your computer. Remember Add to PATH for Windows user
- 2、Run the below code (detailed command can be found in Makefile)
```
make generate-ecdsa

go run .
```

### Generating Swagger API Document
- Add comments to your API source code, See [Declarative Comments Format](https://github.com/swaggo/swag#declarative-comments-format).
- Download Swag by using:
```
go install github.com/swaggo/swag/cmd/swag@latest
```
- Run the Swag in your Go project root folder which contains main.go file, Swag will parse comments and generate required files (docs folder and docs/doc.go). All the packages have been imported to main.go for dependency parsing when init document.
```
swag init
```
- (Optional) Use swag fmt format the SWAG comment. (Please upgrade to the latest version)
```
swag fmt
```
- from https://github1s.com/pcminh0505/gofiber-casbin/blob/HEAD/config/restful_rbac_model.conf#L1-L14