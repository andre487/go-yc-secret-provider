# Go Yandex Cloud secret provider

API for Yandex Cloud LockBox

## Usage

Mode: `YC_SECRET_MODE={prod,dev}`

* dev - using CLI util `yc` for retrieving IAM token
* prod - using metadata server for retrieving IAM token. Server address: `YC_METADATA_SERVICE` or `169.254.169.254` by default

```go
package main

import (
	"fmt"

	ycSecretProvider "github.com/andre487/go-yc-secret-provider"
)

func main() {
	fmt.Println(ycSecretProvider.GetLockBoxTextValue("d6qbv0lnihrdt4mmer19", "token"))
	fmt.Println(ycSecretProvider.GetLockBoxBinaryValue("d6qmn9f60sspf916ncu1", "content"))
}
```
