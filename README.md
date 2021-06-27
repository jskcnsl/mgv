# Go modules graph view

better view for "go mod graph"

## Build

`go build -o mgv ./`

## Usage

- read from file `./mgv graph --f ./graph.txt`
- read from stdin `go mod graph | ./mgv graph`

### Which package depends on `<pkg>`

1. build a exmaple file by `go mod graph > graph.txt`
2. `./mgv dep viper --t --f ./graph.txt`

    ```shell
    github.com/jskcnsl/mg
    |  github.com/spf13/cobra@v1.1.3
    |  |  >> github.com/spf13/viper@v1.7.0 <<
    |  |  |  gopkg.in/yaml.v2@v2.2.4
    |  |  |  |  gopkg.in/check.v1@v0.0.0-20161208181325-20d25e280405
    ```

3. `./mgv ref viper --f ./graph.txt`

    ```shell
    github.com/spf13/viper@v1.7.0
    |  github.com/spf13/cobra@v1.1.3
    |  |  github.com/jskcnsl/mg
    ```
