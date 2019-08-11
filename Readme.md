# Go Open Registry

[![Build Status](https://travis-ci.org/boskiv/go-open-registry.svg?branch=master)](https://travis-ci.org/boskiv/go-open-registry)
[![Maintainability](https://api.codeclimate.com/v1/badges/cd4770aade4ad722f9ca/maintainability)](https://codeclimate.com/github/boskiv/go-open-registry/maintainability)
[![Test Coverage](https://api.codeclimate.com/v1/badges/cd4770aade4ad722f9ca/test_coverage)](https://codeclimate.com/github/boskiv/go-open-registry/test_coverage)


Crates.io cargo registry Golang implementation using amazing Gin web framework
https://github.com/gin-gonic/gin 

## Configure

Environment Variables available
* GIT_REPO_URL is git repository for index storage
* GIT_REPO_PATH directory to clone repo (default 'tmpGit')
* STORAGE_PATH is a directory to upload files in case of local storage (default to 'upload)
* STORAGE is a storage to store binary crate files from `cargo publish` command.
    * local (default)
    * s3 (Not Implemented)
    * artifactory (Not Implemented)
* BUCKET_NAME is a S3 bucket name to upload binary crate files in case of S3 storage
* ARTIFACTORY_REPO is a artifactory repository to upload binary crate files in case of artifactory storage
* GIT_REPO_USERNAME, GIT_REPO_EMAIL and GIT_REPO_PASSWORD is a credentials used to push info to git index repository (HINT: you can use your personal access tokens from GitHub or GitLab as password)
* GIN_MODE you can it release to run in production look more at https://github.com/gin-gonic/gin/blob/master/mode.go
* PORT is a port, GIN server application listen on (default 8000)

Mongo DB used here to store crate packages version and check if it already uploaded.

* MONGODB_URI mongo connection string (default to mongodb://127.0.0.1:27017)
* MONGO_CONNECTION_TIMEOUT connection timeout (default 5 seconds)


## Run

* `docker-compose up -d`
* `./go-open-registry` 

## Publishing package

* setup registry `.cargo/config` file with 
```toml
[registries.open-registry]
index = "https://github.com/boskiv/open-registry-index.git"
token = ""
``` 

* setup `Cargo.toml` file of your package with publish settings
```toml
[package]
...
publish = ["open-registry"]

```

* publish your package
```toml
cargo publish --registry open-registry
```

## Getting package
* setup registry `.cargo/config` file with 
```toml
[registries.open-registry]
index = "https://github.com/boskiv/open-registry-index.git"
token = ""
```

* Setup dependency in `Cargo.toml` file
```toml
...
[dependencies]
bo-helper = { version = "0.1", registry = "open-registry" }
```

* Run build
```toml
cargo update
cargo build
```


