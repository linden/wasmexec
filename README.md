# `wasmexec`
## Getting Started
### Install
```sh
# install `wasmexec`.
go install github.com/linden/wasmexec@latest
```

### Run
```sh
# set the build target to wasm/js then run using wasmexec as the executor.
GOOS=js GOARCH=wasm go run -exec wasmexec 
```

### Test
```sh
# set the build target to wasm/js then test using wasmexec as the executor.
GOOS=js GOARCH=wasm go test -exec wasmexec -v
```


## Credits
Based-off of [`github.com/agnivade/wasmbrowsertest`](https://github.com/agnivade/wasmbrowsertest).