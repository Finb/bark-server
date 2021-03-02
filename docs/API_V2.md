# API V2

**The V2 version API is switched to the standard REST request, and most of the compatibility
processing has been done for the V1 version API; users should use the new REST API when using
the V2 version.**

- [API V2](#api-v2)
    * [Push](#push)
        + [curl](#curl)
        + [golang](#golang)
        + [python](#python)
        + [java](#java)
        + [nodejs](#nodejs)
    * [Misc](#misc)
        + [Ping](#ping)
        + [Healthz](#healthz)
        + [Info](#info)
    
## Push

### curl

```sh
curl -X "POST" "http://127.0.0.1:8080/push" \
     -H 'Content-Type: application/json; charset=utf-8' \
     -d $'{
  "body": "Test Bark Server",
  "device_key": "nysrshcqielvoxsa",
  "title": "bleem",
  "ext_params": {
    "url": "https://mritd.com"
  },
  "category": "category",
  "sound": "minuet.caf"
}'
```

### golang

```go
package main

import (
	"fmt"
	"io/ioutil"
	"net/http"
	"bytes"
)

func sendPush() {
	// push (POST http://127.0.0.1:8080/push)

	json := []byte(`{"body": "Test Bark Server","device_key": "nysrshcqielvoxsa","title": "bleem","ext_params": {"url": "https://mritd.com"},"category": "category","sound": "minuet.caf"}`)
	body := bytes.NewBuffer(json)

	// Create client
	client := &http.Client{}

	// Create request
	req, err := http.NewRequest("POST", "http://127.0.0.1:8080/push", body)
	if err != nil {
		fmt.Println("Failure : ", err)
	}

	// Headers
	req.Header.Add("Content-Type", "application/json; charset=utf-8")

	// Fetch Request
	resp, err := client.Do(req)
	
	if err != nil {
		fmt.Println("Failure : ", err)
	}

	// Read Response Body
	respBody, _ := ioutil.ReadAll(resp.Body)

	// Display Results
	fmt.Println("response Status : ", resp.Status)
	fmt.Println("response Headers : ", resp.Header)
	fmt.Println("response Body : ", string(respBody))
}
```

### python

```python
# Install the Python Requests library:
# `pip install requests`

import requests
import json


def send_request():
    # push
    # POST http://127.0.0.1:8080/push

    try:
        response = requests.post(
            url="http://127.0.0.1:8080/push",
            headers={
                "Content-Type": "application/json; charset=utf-8",
            },
            data=json.dumps({
                "body": "Test Bark Server",
                "device_key": "nysrshcqielvoxsa",
                "title": "bleem",
                "ext_params": {
                    "url": "https://mritd.com"
                },
                "category": "category",
                "sound": "minuet.caf"
            })
        )
        print('Response HTTP Status Code: {status_code}'.format(
            status_code=response.status_code))
        print('Response HTTP Response Body: {content}'.format(
            content=response.content))
    except requests.exceptions.RequestException:
        print('HTTP Request failed')
```

### java

```java
import java.io.IOException;
import org.apache.http.client.fluent.*;
import org.apache.http.entity.ContentType;

public class SendRequest
{
  public static void main(String[] args) {
    sendRequest();
  }
  
  private static void sendRequest() {
    
    // push (POST )
    
    try {
      
      // Create request
      Content content = Request.Post("http://127.0.0.1:8080/push")
      
      // Add headers
      .addHeader("Content-Type", "application/json; charset=utf-8")
      
      // Add body
      .bodyString("{\"body\": \"Test Bark Server\",\"device_key\": \"nysrshcqielvoxsa\",\"title\": \"bleem\",\"ext_params\": {\"url\": \"https://mritd.com\"},\"category\": \"category\",\"sound\": \"minuet.caf\"}", ContentType.APPLICATION_JSON)
      
      // Fetch request and return content
      .execute().returnContent();
      
      // Print content
      System.out.println(content);
    }
    catch (IOException e) { System.out.println(e); }
  }
}
```

### nodejs

```node
// request push 
(function(callback) {
    'use strict';
        
    const httpTransport = require('http');
    const responseEncoding = 'utf8';
    const httpOptions = {
        hostname: '127.0.0.1',
        port: '8080',
        path: '/push',
        method: 'POST',
        headers: {"Content-Type":"application/json; charset=utf-8"}
    };
    httpOptions.headers['User-Agent'] = 'node ' + process.version;
 
    // Using Basic Auth {"username":"","password":""}
    // Paw Store Cookies option is not supported

    const request = httpTransport.request(httpOptions, (res) => {
        let responseBufs = [];
        let responseStr = '';
        
        res.on('data', (chunk) => {
            if (Buffer.isBuffer(chunk)) {
                responseBufs.push(chunk);
            }
            else {
                responseStr = responseStr + chunk;            
            }
        }).on('end', () => {
            responseStr = responseBufs.length > 0 ? 
                Buffer.concat(responseBufs).toString(responseEncoding) : responseStr;
            
            callback(null, res.statusCode, res.headers, responseStr);
        });
        
    })
    .setTimeout(0)
    .on('error', (error) => {
        callback(error);
    });
    request.write("{\"device_key\":\"nysrshcqielvoxsa\",\"body\":\"Test Bark Server\",\"title\":\"bleem\",\"sound\":\"minuet.caf\",\"category\":\"category\",\"ext_params\":{\"url\":\"https://mritd.com\"}}")
    request.end();
    

})((error, statusCode, headers, body) => {
    console.log('ERROR:', error); 
    console.log('STATUS:', statusCode);
    console.log('HEADERS:', JSON.stringify(headers));
    console.log('BODY:', body);
});
```

## Misc

### Ping

```sh
curl "http://127.0.0.1:8080/ping"
```

### Healthz

```sh
curl "http://127.0.0.1:8080/healthz"
```

### Info

```sh
curl "http://127.0.0.1:8080/info"
```
