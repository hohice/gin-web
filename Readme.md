# ginS
ginS enable an ApiServer with monitor and trace.
# Usage
## start server
ginS serv [-a addr] [-p port]
## Build
### Builder
```
cd build/builder-docker && make
```
### Build(local)
```
make all
```
### Build(local)
```
make build
```
## Test
### All
```
make test
```
### Unit test 
```
make unit-test
```
### E2E test 
```
make e2e-test
```
# Road Map
- [x] RestFul API Server