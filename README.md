# defect-detect

Detect vulnerable components using sboms

## Pre-requisites

- Create Github fine grained token with `"Contents" repository permissions (read)`
- Go Installed on Machine
- Docker or other container oci environment

## Installation

- Clone repo and install tool

  ```bash
  git clone https://github.com/dmdhrumilmistry/defect-detect
  cd defect-detect
  go install -v ./...
  ```

- Create config file

  ```bash
  cp .env.sample .env
  export GITHUB_TOKEN="your-github-token" # this can be also added in config file
  ```

- Start container env (mongodb)

  ```bash
  docker compose up -d
  ```

- Start backend

  ```bash
  defect-detect
  ```

## Usage

### Import SBOM and Analyze components

- Import Sbom into DB

  - Using File

    ```bash
    curl -X POST -F "sbom=@example-sbom.json" http://localhost:8080/api/v1/sbom
    ```

  - Import Github Repo

    ```bash
    curl -X POST -H "application/json" -d '{"owner":"dmdhrumilmistry", "repo_name":"pyhtools"}' http://localhost:8080/api/v1/sbom/githubImport
    ```

  - Example Output

    ```
    {"id":"676f0bac3da126bf929f246c","message":"SBOM uploaded successfully"}
    ```

- Create Project

  ```bash
  curl -X POST http://localhost:8080/api/v1/project -H "Content-Type: application/json" -d '{"name":"pyhtools", "description":"python hacking tools project", "sboms_to_retain": 2, "links": ["https://github.com/dmdhrumilmistry/pyhtools"], "sboms": ["676f0bac3da126bf929f246c"]}'
  ```

- Analyze components

  ```bash
  curl -X POST "http://localhost:8080/api/v1/component?sbom_id=676f0bac3da126bf929f246c"

  # Output
  # {"ids":["676f0c4ff986a31a1ab2ecf5", "...snip..."],"message":"Components created successfully from Sbom"}
  ```

- Fetch Vulnerable Components

  ```bash
  curl "http://localhost:8080/api/v1/component/vulns?sbom_ids=676f0bac3da126bf929f246c"
  ```

  > Response will be paginated

  Supported query params: `sbom_ids`, `component_names`, `component_versions`, `types`, `names`, `versions`, `purls`
  Multiple values is supported separated by `,`

  |    Query Param     | Description                                                                                                                           |
  | :----------------: | :------------------------------------------------------------------------------------------------------------------------------------ |
  |      sbom_ids      | Id of SBOM uploaded to the application                                                                                                |
  |  component_names   | Name of Component fetched from Github Repo SBOM (com.github.dmdhrumiilmistry/pyhtools) or uploaded Sbom metadata Component (pyhtools) |
  | component_versions | Version of Component fetched from Github Repo SBOM (main) or uploaded Sbom metadata Component (latest/v1.1.1)                         |
  |       types        | Type of sbom component such as package, framework, etc.                                                                               |
  |       names        | name of sbom component. It is usually dependency name                                                                                 |
  |      versions      | version of sbom component                                                                                                             |
  |        purl        | package url of sbom component                                                                                                         |
