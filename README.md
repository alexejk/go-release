# Release tools for Go projects

The main idea is to provide a way to make new releases, version bumps and publishing of created artifacts with one command.
To support this a descriptor file is used which defined behavior and properties in a configuration format, such as YAML or TOML.

## Sample Project File

```
project:
  version:
    file: build.properties
    property: APP_VERSION
    increment: minor

  build:
    command: make build-in-docker

  git:

    tag:
      format: $version
      push: true

    message:
      release: "[release] Release version"
      development: "[release] Next development version"

publishing:

  artifactory:
    url:
    repository:

  s3upload:
    bucket: my-tools
    path: /release-tool/
    artifacts:
    - build/releaser-darwin
    - build/release-linux

  github:
    draft: true
    changelog:
      file: CHANGELOG.md
      format: Release $version

    releaseName: $version
    artifacts:
    - build/releaser-$version-darwin.tar.gz
    
```
------

## Release Flow:

- Read configuration
- Get current version
- Bump version to release
- Git commit
- Build
- Git tag
- Set version to next development
