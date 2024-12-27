# defect-detect

Detect vulnerable components using sboms

## Pre-requisites

* Create Github fine grained token with `"Contents" repository permissions (read)`
* Go Installed on Machine
* Docker or other container oci environment


## Installation

* Clone repo and install tool

    ```bash
    git clone https://github.com/dmdhrumilmistry/defect-detect
    cd defect-detect
    go install -v ./...
    ```

* Create config file

    ```bash
    cp .env.sample .env
    export GITHUB_TOKEN="your-github-token" # this can be also added in config file
    ```

* Start container env (mongodb)

    ```bash
    docker compose up -d 
    ```

* Start backend

    ```bash
    defect-detect
    ```

## Usage

### Import SBOM and Analyze components

* Import Sbom into DB

    * Using File

        ```bash
        curl -X POST -F "sbom=@example-sbom.json" http://localhost:8080/api/v1/sbom
        ```

    * Import Github Repo

        ```bash
        curl -X POST -H "application/json" -d '{"owner":"dmdhrumilmistry", "repo_name":"pyhtools"}' http://localhost:8080/api/v1/sbom/githubImport
        ```

    * Example Output

        ```
        {"id":"676f0bac3da126bf929f246c","message":"SBOM uploaded successfully"}
        ```

* Analyze components

    ```bash
    curl -X POST "http://localhost:8080/api/v1/component?sbom_id=676f0bac3da126bf929f246c"

    # Output
    # {"ids":["676f0c4ff986a31a1ab2ecf5", "...snip..."],"message":"Components created successfully from Sbom"}
    ```
