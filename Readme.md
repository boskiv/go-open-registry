# Go Open Registry

[![Build Status](https://travis-ci.org/boskiv/go-open-registry.svg?branch=master)](https://travis-ci.org/boskiv/go-open-registry)

Crates.io cargo registry Golang implementation using amazing Gin web framework
https://github.com/gin-gonic/gin 

## Configure

Environment Variables available
* REPO is git repository for index storage
* UPLOAD_DIR is a directory to upload files in case of local storage (default to './upload)
* STORAGE is a storage to store binary crate files from `cargo publish` command.
    * local
    * s3 (Not Implemented)
    * artifactory (Not Implemented)
* BUCKET_NAME is a S3 bucket name to upload binary crate files in case of S3 storage
* ARTIFACTORY_REPO is a artifactory repository to upload binary crate files in case of artifactory storage
* GIT_USERNAME and GIT_PASSWORD is a credentials used to push info to git index repository

From GIN

## Run

`./go-open-registry` 
