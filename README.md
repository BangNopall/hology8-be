<div align="center">

<h1 align="center">Hology 8 Backend (Legacy Hology7)</h1>
</div>

## Getting Started

This application is using [Go v1.22.3](https://tip.golang.org/doc/go1.22) and [PostgreSQL](https://www.postgresql.org/), make sure you already installed the required dependency for this application.

## 📄 API Docs

For API Docs, you can look up in the following link.

[API Swagger Documentation](https://awsdocs.dvnnfrr.my.id/docs/index.html)

## 📝 Convention 

Please look and follow the convention that you can see in [here](./CONVENTION.md)

## 🗂️ Directory Structure

```zsh
hology-be
├── cmd                 # executeables file
│  ├── app              # application entry
├── config              # config files (firebase admin, etc)
├── data            
│  ├── seeders          # seeders data for db
│  │   ├── dev          # for development purposes
│  │   ├── prod         # for production also
├── deploy              # for deployment purposes
├── docs                # application documentation
├── domain              # entities structure, dtos, contracts
├── internal            # private application and libraries
│  ├── app              # application functionality
│  ├── infra            # application external systems
│  ├── middlewares      # application middlewares
├── pkg                 # external reusable packges and libraries            
├── tests               # integration testing files 
├── web                 # static html files
```

## Penggunaan Nama Branch
Untuk nama branch harus menggunakan standar seperti berikut :
- ``(nama)-(tipe): (deskripsi)``

contoh : ``nopal-feat/navbar``

Setiap membuat fitur baru atau memperbaiki fitur harus menggunakan branch baru seperti SOP diatas dan melakukan pull request ke akuu (BangNopall) yaa

## Cara Commit / Push ke Github untuk Update Progresan
```c
- git add .
- git commit -m "feat: adding navbar section"
- git push origin nopal-feat/navbar
```
Harus menggunakan [`conventional commits`](https://gist.github.com/qoomon/5dfcdf8eec66a051ecd85625518cfd13)! <br/>
dengan format ``tipe: deskripsi``

<h3>Berikut Tipenya</h3>

- API or UI relevant changes
    - `feat` Commits, that add or remove a new feature to the API or UI
    - `fix` Commits, that fix an API or UI bug of a preceded `feat` commit
- `refactor` Commits, that rewrite/restructure your code, however do not change any API or UI behaviour
    - `perf` Commits are special `refactor` commits, that improve performance
- `style` Commits, that do not affect the meaning (white-space, formatting, missing semi-colons, etc)
- `test` Commits, that add missing tests or correcting existing tests
- `docs` Commits, that affect documentation only
- `build` Commits, that affect build components like build tools, dependencies, project version, ci pipelines, ...
- `ops` Commits, that affect operational components like infrastructure, deployment, backup, recovery, ...
- `chore` Miscellaneous commits e.g. modifying `.gitignore`

## Cara Pull / Mengambil Data Terbaru dari Github ke Lokal
```c
- git pull
- git pull origin (branch) // untuk pull dari branch lain
```

## 🧰 Libraries Used

NOTES : if you change a library or add another library, please update the following list

Library | Usage | 
--- | --- | 
[OAuth2](https://github.com/golang/oauth2) | Auth |
[JWT](https://github.com/golang-jwt/jwt) | Auth
[Uuid](https://github.com/google/uuid) | UUID
[Logger](https://github.com/sirupsen/logrus) | Logging
[Gorm](https://github.com/go-gorm/gorm) | ORM
[Gin](https://github.com/gin-gonic/gin) | HTTP Framework
[GoMail](https://github.com/go-gomail/gomail) | Mail
[Firebase](https://github.com/firebase/firebase-admin-go) | File Storage
[Redis](https://github.com/redis/go-redis) | In-Memory Cache
[Swaggo](https://github.com/swaggo/swag) | API Docs

## 🛠️ Tech Stacks

[![My Skills](https://skillicons.dev/icons?i=nginx,docker,golang,aws,postgres,redis,linux,ubuntu,)](https://skillicons.dev)