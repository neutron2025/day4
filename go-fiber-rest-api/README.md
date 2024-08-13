参考文章：https://dev.to/koddr/build-a-restful-api-on-go-fiber-postgresql-jwt-and-swagger-docs-in-isolated-docker-containers-475j#toc


生成迁移文件

```migrate create -ext sql -dir platform/migrations init```

迁移文件里面自己填


postgresql使用容器 或者也可用
```
docker run --rm -d \
		--name dev-postgres \
		--network dev-network \
		-e POSTGRES_USER=postgres \
		-e POSTGRES_PASSWORD=password \
		-e POSTGRES_DB=postgres \
		-v ${HOME}/dev-postgres/data/:/var/lib/postgresql/data \
		-p 5432:5432 \
		postgres
```

生成到数据库
```migrate -verbose -path="./platform/migrations" -database "postgres://postgres:password@localhost/postgres?sslmode=disable" up```
