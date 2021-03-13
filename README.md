command to get available options

`make help`

##### Build Project

`make build`

##### Run project

`make start` will build and run project with default options



#### Available options

`./ck-order-delivery-system server --ops 2 --order-file-path ./data/orders.json`

`--ops` number of orders need to be processed per second

`--order-file-path` json file path can be passed here to read and process orders.Default path is `./data/orders.json`

`./ck-order-delivery-system server -h` to list available options

##### Running Tests
`make test`

#### Confessions
Libraries used:
- github.com/spf13/cobra
- github.com/stretchr/testify

**Due to time concern there is good room for optimization and add more test coverage.**

- can improve error reporting
- Split code into context vised modules 
- more abstract