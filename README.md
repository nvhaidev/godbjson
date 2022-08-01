# GOBDJson

GOBDJson is a Golang library.

## Installation

Use the package manager [go](https://pkg.go.dev/cmd/go#hdr-Download_and_install_packages_and_dependencies) to install foobar.

```bash
go get -u github.com/nvhaidev/godbjson
```

## Usage

```python
package main

import (
	"github.com/nvhaidev/godbjson"
)

func main() {
	db := godbjson.NewDB("db.json")
	db.Create(map[string]interface{}{
		"name":        "Hai",
		"age":         "17",
		"description": "Hai dep trai",
	})
	db.FindById("5d8f8f8f-8f8f-8f8f-8f8f-8f8f8f8f8f8f")
	
}
```

## Contributing
Pull requests are welcome. For major changes, please open an issue first to discuss what you would like to change.

Please make sure to update tests as appropriate.

## License
[MIT](https://choosealicense.com/licenses/mit/)