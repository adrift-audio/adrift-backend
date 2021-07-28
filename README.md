## adrift-backend

Backend for the [Adrift](https://github.com/adrift-audio/adrift-desktop) project

Stack: [Golang](https://golang.org), [Fiber](https://docs.gofiber.io), [MongoDB](https://www.mongodb.com), [Redis](https://github.com/go-redis/redis), [JWT](https://github.com/dgrijalva/jwt-go), [Argon2](https://github.com/alexedwards/argon2id)

### Deploy

Golang v1.15.X is required

```shell script
git clone https://github.com/adrift-audio/adrift-backend
cd ./adrift-backend
```

### Environment

The `.env` file is required, check [.env.example](.env.example) for details

### Launch

```shell script
go run ./
```

Alternatively, use [AIR](https://github.com/cosmtrek/air)

### License

[MIT](./LICENSE.md)
