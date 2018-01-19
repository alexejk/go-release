# Release tools for Go projects

The main idea is to provide a way to make new releases, version bumps and publishing of created artifacts with one command.
To support this a descriptor file is used which defined behavior and properties in a configuration format, such as YAML or TOML.

```
release:
  version:
    file: build.properties
    property: APP_VERSION
    format:
      release: 1.0.0
      development: 1.0.1-dev

  build:
    command: make build-in-docker

  git:
    tag:
      format: $version
      push: true

    commit:
      releaseFormat: [release] Release version $version
      developmentFormat: [release] Next development version $version

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
    releaseName: $version
    artifacts:
    - build/releaser-$version-darwin.tar.gz


    
```
