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
        + [php](#php)
    * [Misc](#misc)
        + [Ping](#ping)
        + [Healthz](#healthz)
        + [Info](#info)
    
## Push

| Field | Type | Description |
| ----- | ---- | ----------- |
| title | string | Notification title (font size would be larger than the body) |
| body  | string | Notification content |
| category | string | Reserved field, no use yet |
| device_key | string | The key for each device |
| level (optional) | string | `'active'`, `'timeSensitive'`, or `'passive'` |
| badge (optional) | integer | The number displayed next to App icon ([Apple Developer](https://developer.apple.com/documentation/usernotifications/unnotificationcontent/1649864-badge)) |
| automaticallyCopy (optional) | string | Must be `1` |
| copy (optional) | string |  The value to be copied |
| sound (optional) | string | Value from [here](https://github.com/Finb/Bark/tree/master/Sounds) |
| icon (optional) | string | An url to the icon, available only on iOS 15 or later |
| group (optional) | string | The group of the notification |
| isArchive (optional) | string | Value must be `1`. Whether or not should be archived by the app |
| url (optional) | string | Url that will jump when click notification |

### curl

```sh
curl -X "POST" "http://127.0.0.1:8080/push" \
     -H 'Content-Type: application/json; charset=utf-8' \
     -d $'{
  "body": "Test Bark Server",
  "device_key": "ynJ5Ft4atkMkWeo2PAvFhF",
  "title": "bleem",
  "badge": 1,
  "category": "myNotificationCategory",
  "sound": "minuet.caf",
  "icon": "https://day.app/assets/images/avatar.jpg",
  "group": "test",
  "url": "https://mritd.com"
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

	json := []byte(`{"body": "Test Bark Server","device_key": "nysrshcqielvoxsa","title": "bleem", "badge": 1, "icon": "https://day.app/assets/images/avatar.jpg", "group": "test", "url": "https://mritd.com","category": "myNotificationCategory","sound": "minuet.caf"}`)
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
                "category": "myNotificationCategory",
                "sound": "minuet.caf",
                "badge": 1,
                "icon": "https://day.app/assets/images/avatar.jpg",
                "group": "test",
                "url": "https://mritd.com"
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
      .bodyString("{\"body\": \"Test Bark Server\",\"device_key\": \"nysrshcqielvoxsa\",\"title\": \"bleem\",\"url\": \"https://mritd.com\", \"group\": \"test\",\"category\": \"myNotificationCategory\",\"sound\": \"minuet.caf\"}", ContentType.APPLICATION_JSON)
      
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
    request.write("{\"device_key\":\"nysrshcqielvoxsa\",\"body\":\"Test Bark Server\",\"title\":\"bleem\",\"sound\":\"minuet.caf\",\"category\":\"myNotificationCategory\",\"url\":\"https://mritd.com\", \"group\":\"test\"}")
    request.end();
    

})((error, statusCode, headers, body) => {
    console.log('ERROR:', error); 
    console.log('STATUS:', statusCode);
    console.log('HEADERS:', JSON.stringify(headers));
    console.log('BODY:', body);
});
```

### php

```php
$curl = curl_init();
curl_setopt_array($curl, [
    CURLOPT_URL => 'http://127.0.0.1:8080/push',
    CURLOPT_CUSTOMREQUEST => 'POST',
    CURLOPT_POSTFIELDS => '{
  "body": "Test Bark Server",
  "device_key": "ynJ5Ft4atkMkWeo2PAvFhF",
  "title": "bleem",
  "badge": 1,
  "category": "myNotificationCategory",
  "sound": "minuet.caf",
  "icon": "https://day.app/assets/images/avatar.jpg",
  "group": "test",
  "url": "https://mritd.com"
}',
    CURLOPT_HTTPHEADER => [
        'Content-Type: application/json; charset=utf-8',
    ],
]);
$response = curl_exec($curl);
curl_close($curl);
echo $response;
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
