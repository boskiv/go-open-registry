# Go Open Registry

[![Build Status](https://travis-ci.org/boskiv/go-open-registry.svg?branch=master)](https://travis-ci.org/boskiv/go-open-registry)
[![Maintainability](https://api.codeclimate.com/v1/badges/cd4770aade4ad722f9ca/maintainability)](https://codeclimate.com/github/boskiv/go-open-registry/maintainability)
[![Test Coverage](https://api.codeclimate.com/v1/badges/cd4770aade4ad722f9ca/test_coverage)](https://codeclimate.com/github/boskiv/go-open-registry/test_coverage)


Crates.io cargo registry Golang implementation using amazing Gin web framework
https://github.com/gin-gonic/gin 

Based on https://doc.rust-lang.org/cargo/reference/registries.html

## Limitations

* Yank feature not implemented yet (https://doc.rust-lang.org/cargo/reference/registries.html#yank)

## Configure with environment variables

A shortlist, description bellow

* [GIN_MODE](#gin_mode)
* [PORT](#port)
* [CARGO_API_URL](#cargo_api_url)
* [GIT_REPO_URL](#git_repo_url)
* [GIT_REPO_PATH](#git_repo_path)
* [GIT_REPO_USERNAME](#git_repo_username)
* [GIT_REPO_PASSWORD](#git_repo_path)
* [GIT_REPO_EMAIL](#git_repo_email)
* [MONGODB_URI](#mongodb_uri)
* [MONGO_CONNECTION_TIMEOUT](#mongo_connection_timeout)
* [STORAGE_TYPE](#storage_type)
* [LOCAL_STORAGE_PATH](#local_storage_path)
* [ARTIFACTORY_URL](#artifactory_url)
* [ARTIFACTORY_LOGIN](#artifactory_login)
* [ARTIFACTORY_PASSWORD](#artifactory_password)
* [ARTIFACTORY_REPO_NAME](#artifactory_repo_name)
* [AWS_ACCESS_KEY_ID](#aws_access_key_id)
* [AWS_SECRET_ACCESS_KEY](#aws_secret_access_key)
* [AWS_DEFAULT_REGION](#aws_default_region)
* [AWS_S3_BUCKET_NAME](#aws_s3_bucket_name)
* [AWS_S3_ENDPOINT](#aws_s3_endpoint)

### Generic
##### `STORAGE_TYPE`
* **Description**
    > Storage type is used to choose backend storage to upload/download binary crates files <br> 
    It can be one of choice: <br>
    .. local - local filesystem folder <br>
    .. s3 - AWS S3 bucket, or S3 API compatible system (Minio, Ceph) <br>
    .. artifactory - JFrog Artifactory system
* **Example** 
    > `STORAGE_TYPE=s3`
* **Default**
    > `STORAGE_TYPE=local`
### GIN

##### `GIN_MODE`
* **Description**
    > Gin Application working mode <br> 
    https://github.com/gin-gonic/gin/blob/master/mode.go#L52 <br>
    There is some log verbosity applied based on the mode. <br>
    Debug is more verbose.
* **Example** 
    > `GIN_MODE=release`
* **Default**
    > `GIN_MODE=debug`

##### `PORT`
* **Description**
    > HTTP port to application listen on<br> 
* **Example** 
    > `PORT=8000`
* **Default**
    > `PORT=8000`

### Cargo config
Cargo config file stored in a git registry generated automatically from environment variables on start
If you change variables, the file will be updated and pushed to the repository

https://doc.rust-lang.org/cargo/reference/registries.html#index-format

`dl` url made from `api` with concatenation `/api/v1/crates` in [internal/config/config.go:60](https://github.com/boskiv/go-open-registry/blob/master/internal/config/config.go#L60)

##### `CARGO_API_URL`
* **Description**
    > This is the base URL for the web API 
* **Example** 
    > `CARGO_API_URL=http://my-registry-api:3000`
* **Default**
    > `CARGO_API_URL=http://localhost:8000`

### Git
##### `GIT_REPO_URL`
* **Description**
    > Git repository to store index files. <br>
    Cargo use it to search particular package and version. <br>
    [Cargo configuration example](#cargo-configuration-example)
* **Example** 
    > `GIT_REPO_URL=https://github.com/boskiv/open-registry-index`
* **Default**
    > Empty

##### `GIT_REPO_PATH`
* **Description**
    > Temporary directory to clone repo. <br>
    Application use it to commit and push cargo package information. <br>
* **Example** 
    > `GIT_REPO_PATH=/data/gitRepo`
* **Default**
    > `GIT_REPO_PATH=tmpGit`

##### `GIT_REPO_USERNAME`
* **Description**
    > Login to work with git repo. <br>
* **Example** 
    > `GIT_REPO_USERNAME=boskiv`
* **Default**
    > Empty

##### `GIT_REPO_PASSWORD`
* **Description**
    > Password to work with git repo. <br>
* **Example** 
    > `GIT_REPO_PASSWORD=123123123`
* **Default**
    > Empty

##### `GIT_REPO_EMAIL`
* **Description**
    > Email used in commit signature. <br>
* **Example** 
    > `GIT_REPO_EMAIL=crates@company.org`
* **Default**
    > Empty
### Mongo
Mongo database used to store package name and version info.

Multiple field index used to control version uniqueness.

https://docs.mongodb.com/manual/core/index-unique

```
Keys: bson.M{
  "name": 1,
  "version": 1,
},
Options: options.Index().SetUnique(true)
```

##### `MONGODB_URI`
* **Description**
    > URI formats for defining connections between applications and MongoDB instances in the official MongoDB drivers. <br>
    https://docs.mongodb.com/manual/reference/connection-string
* **Example** 
    > `MONGODB_URI=mongodb://mongo:27017`
* **Default**
    > `MONGODB_URI=mongodb://127.0.0.1:27017`

##### `MONGO_CONNECTION_TIMEOUT`
* **Description**
    > Timeout to check mongo availability <br>
    Application wil exit with code 1, if timeout fires.
* **Example** 
    > `MONGO_CONNECTION_TIMEOUT=15`
* **Default**
    > `MONGO_CONNECTION_TIMEOUT=5`
### Local storage

Be sure that use set `STORAGE_TYPE=local` to use or remove that env, so this will apply by default.

##### `LOCAL_STORAGE_PATH`
* **Description**
    > Filesystem path to store uploaded crates and get it for download. <br>
    If the path does not exist, the application will try to create it. <br>
    You can use docker mounted volumes, shared volumes and network filesystem volumes, to sure about data persistence. <br>
    However, it's not production recommended storage because of maintenance, backup, and support issues. <br> 
    So use it for **testing only** <br>
* **Example** 
    > `LOCAL_STORAGE_PATH=/data/crates`
* **Default**
    > `LOCAL_STORAGE_PATH=upload`
### Artifactory storage

Be sure that use set `STORAGE_TYPE=artifactory` to use it.

##### `ARTIFACTORY_URL`
* **Description**
    > Path to artifactory API. <br> 
    Should include schema (http/https) and full path to API. <br> 
    You can find this in the description at [Set Me Up](https://www.jfrog.com/confluence/display/RTF/Using+Artifactory#UsingArtifactory-SetMeUp) dialog
* **Example**
    > `ARTIFACTORY_URL=http://localhost:8081/artifactory`
* **Default**
    > Empty

##### `ARTIFACTORY_LOGIN`
* **Description**
    > Login to artifactory with write access for binary repository. <br> 
* **Example**
    > `ARTIFACTORY_LOGIN=bot`
* **Default**
    > Empty
    
##### `ARTIFACTORY_PASSWORD`
* **Description**
    > Password to access artifactory with login provided in `ARTIFACTORY_LOGIN`. <br> 
* **Example**
    > `ARTIFACTORY_PASSWORD=password`
* **Default**
    > Empty

##### `ARTIFACTORY_REPO_NAME`
* **Description**
    > Generic binary repository name created in artifactory. <br>
    https://www.jfrog.com/confluence/display/RTF/Configuring+Repositories
* **Example**
    > `ARTIFACTORY_REPO_NAME=crates`
* **Default**
    > Empty
### S3 storage

Most of variables are common used.

https://docs.aws.amazon.com/cli/latest/userguide/cli-configure-envvars.html

##### `AWS_ACCESS_KEY_ID`
* **Description**
    > Specifies an AWS access key associated with an IAM user or role. <br> 
* **Example**
    > `AWS_ACCESS_KEY_ID=AKIAIOSFODNN7EXAMPLE`
* **Default**
    > Empty

##### `AWS_SECRET_ACCESS_KEY`
* **Description**
    > Specifies the secret key associated with the access key. This is essentially the "password" for the access key. 
* **Example**
    > `AWS_SECRET_ACCESS_KEY=wJalrXUtnFEMI/K7MDENG/bPxRfiCYEXAMPLEKEY`
* **Default**
    > Empty

##### `AWS_DEFAULT_REGION`
* **Description**
    > Specifies the [AWS Region](https://docs.aws.amazon.com/en_us/cli/latest/userguide/cli-chap-configure.html#cli-quick-configuration-region) to send the request to. 
* **Example**
    > `AWS_DEFAULT_REGION=us-west-2`
* **Default**
    > Empty

##### `AWS_S3_BUCKET_NAME`
* **Description**
    > Bucket name to store crate files. It will be created automatically if not exist.
* **Example**
    > `AWS_S3_BUCKET_NAME=crates`
* **Default**
    > Empty

##### `AWS_S3_USE_SSL`
* **Description**
    > Should application use SSL protocol when connect to S3 API  
* **Example**
    > `AWS_S3_USE_SSL=false`
* **Default**
    > AWS_S3_USE_SSL=true

##### `AWS_S3_ENDPOINT`
* **Description**
    > S3 API Endpoint to upload file. <br>
    Application use minio sdk to work with S3 like APIs
    you should put here url without schema f.e. play.minio.io, but you should to put schema in AWS_S3_USE_SSL as boolean param  
* **Example**
    > `AWS_S3_ENDPOINT=localhost:9000`
* **Default**
    > AWS_S3_ENDPOINT=s3
## Run
### Prerequisite

Make a accessible git repository with config.json

https://doc.rust-lang.org/cargo/reference/registries.html#index-format

```json
{
    "dl": "http://localhost:8000/api/v1/crates",
    "api": "http://localhost:8000"
}
```

> The keys are: <br>
  dl: This is the URL for downloading crates listed in the index.<br>
  api: This is the base URL for the web API. 

### Developer 

Run this configuration just to have third-party services

* mongo
* minio
* artifactory

#### Steps

* `docker-compose run -f compose/developer.yaml`
* `go build cmd/go-open-registry`
* `./go-open-registry`

### Test with local storage

* `docker-compose up -d -f compose/local.yaml`
* `open http://localhost:3000` for Gogs and create an account and repository.
* put config.json file to the repository
* `use http://localhost:8000` to Configure cargo

### Test with artifactory

* `docker-compose up -d -f compose/artifactory.yaml`
* `open http://localhost:8081` for Artifactory, create a login `bot` and password `password`.
* create a repository named crates, and give user bot a write access to this repo.
* `open http://localhost:3000` for Gogs and create an account and repository.
* put config.json file to the repository
* `use http://localhost:8000` to Configure cargo

### Test with s3 storage

* `docker-compose up -d -f compose/minio.yaml`
* `open http://localhost:9000` for Minio and login with `minio` as user and `minio123` as password.
* bucket `crates` will be created automatically
* `open http://localhost:3000` for Gogs and create an account and repository.
* put config.json file to the repository
* `use http://localhost:8000` to Configure cargo


# Cargo configuration example
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


