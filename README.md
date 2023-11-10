# httpc
[![PkgGoDev](https://pkg.go.dev/badge/github.com/bikbah/httpc)](https://pkg.go.dev/github.com/bikbah/httpc)

golang package to DRY making http client calls

## Usage

1. Install
    ```
    go get -v github.com/bikbah/httpc
    ```

2. Use:

    ```go
    package main
    
    import (
    	"context"
    	"fmt"
    
    	"github.com/bikbah/httpc"
    )
    
    func main() {
    	client := httpc.Must("https://httpbin.org")
    
    	req := httpc.GET("/get").
    		WithQuery(httpc.String("query", "query_value")).
    		WithHeader(httpc.String("Authorization", "Bearer some_token"))
    
    	var res struct {
    		Args    map[string]any    `json:"args"`
    		Headers map[string]string `json:"headers"`
    		URL     string            `json:"url"`
    		JSON    map[string]any    `json:"json"`
    	}
    
    	err := client.Do(context.Background(), req, &res, nil)
    	fmt.Printf("%+v %v\n", res, err)
    	// Output: {
    	//   Args:map[query:query_value]
    	//   Headers:map[
    	//     Accept-Encoding:gzip
    	//     Authorization:Bearer some_token
    	//     Host:httpbin.org
    	//     User-Agent:Go-http-client/2.0
    	//   ]
    	//   URL:https://httpbin.org/get?query=query_value
    	//   JSON:map[]
    	// }
    	// <nil>
    
    	reqBody := struct {
    		Key string `json:"key"`
    	}{
    		Key: "value",
    	}
    
    	res = struct {
    		Args    map[string]any    `json:"args"`
    		Headers map[string]string `json:"headers"`
    		URL     string            `json:"url"`
    		JSON    map[string]any    `json:"json"`
    	}{}
    
    	req = httpc.POST("/post").
    		WithJSONBody(reqBody)
    
    	err = client.Do(context.Background(), req, &res, nil)
    	fmt.Printf("%+v %v\n", res, err)
    	// Output: {
    	//   Args:map[]
    	//   Headers:map[
    	//     Accept-Encoding:gzip
    	//     Content-Length:16
    	//     Content-Type:application/json
    	//     Host:httpbin.org
    	//     User-Agent:Go-http-client/2.0
    	//   ]
    	//   URL:https://httpbin.org/post
    	//   JSON:map[key:value]
    	// }
    	// <nil>
    }
    ```
