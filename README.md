# APNS/2 
test
## Example

```go
package main

import (
  "log"
  "fmt"

  "github.com/cmeyer18/apns2"
  "github.com/cmeyer18/apns2/certificate"
)

func main() {

  cert, err := certificate.FromP12File("../cert.p12", "")
  if err != nil {
    log.Fatal("Cert Error:", err)
  }

  notification := &apns2.Notification{}
  notification.DeviceToken = "11aa01229f15f0f0c52029d8cf8cd0aeaf2365fe4cebc4af26cd6d76b7919ef7"
  notification.Topic = "com.sideshow.Apns2"
  notification.Payload = []byte(`{"aps":{"alert":"Hello!"}}`) // See Payload section below

  // If you want to test push notifications for builds running directly from XCode (Development), use
  // client := apns2.NewClient(cert).Development()
  // For apps published to the app store or installed as an ad-hoc distribution use Production()

  client := apns2.NewClient(cert).Production()
  res, err := client.Push(notification)

  if err != nil {
    log.Fatal("Error:", err)
  }

  fmt.Printf("%v %v %v\n", res.StatusCode, res.ApnsID, res.Reason)
}
```

## JWT Token Example

Instead of using a `.p12` or `.pem` certificate as above, you can optionally use
APNs JWT _Provider Authentication Tokens_. First you will need a signing key (`.p8` file), Key ID and Team ID [from Apple](http://help.apple.com/xcode/mac/current/#/dev54d690a66). Once you have these details, you can create a new client:

```go
authKey, err := token.AuthKeyFromFile("../AuthKey_XXX.p8")
if err != nil {
  log.Fatal("token error:", err)
}

token := &token.Token{
  AuthKey: authKey,
  // KeyID from developer account (Certificates, Identifiers & Profiles -> Keys)
  KeyID:   "ABC123DEFG",
  // TeamID from developer account (View Account -> Membership)
  TeamID:  "DEF123GHIJ",
}
...

client := apns2.NewTokenClient(token)
res, err := client.Push(notification)
```

- You can use one APNs signing key to authenticate tokens for multiple apps.
- A signing key works for both the development and production environments.
- A signing key doesn’t expire but can be revoked.

## Notification

At a minimum, a _Notification_ needs a _DeviceToken_, a _Topic_ and a _Payload_.

```go
notification := &apns2.Notification{
  DeviceToken: "11aa01229f15f0f0c52029d8cf8cd0aeaf2365fe4cebc4af26cd6d76b7919ef7",
  Topic: "com.sideshow.Apns2",
  Payload: []byte(`{"aps":{"alert":"Hello!"}}`),
}
```

You can also set an optional _ApnsID_, _Expiration_ or _Priority_.

```go
notification.ApnsID =  "40636A2C-C093-493E-936A-2A4333C06DEA"
notification.Expiration = time.Now()
notification.Priority = apns2.PriorityLow
```

## Payload

You can use raw bytes for the `notification.Payload` as above, or you can use the payload builder package which makes it easy to construct APNs payloads.

```go
// {"aps":{"alert":"hello","badge":1},"key":"val"}

payload := payload.NewPayload().Alert("hello").Badge(1).Custom("key", "val")

notification.Payload = payload
client.Push(notification)
```

Refer to the [payload](https://godoc.org/github.com/sideshow/apns2/payload) docs for more info.

## Response, Error handling

APNS/2 draws the distinction between a valid response from Apple indicating whether or not the _Notification_ was sent or not, and an unrecoverable or unexpected _Error_;

- An `Error` is returned if a non-recoverable error occurs, i.e. if there is a problem with the underlying _http.Client_ connection or _Certificate_, the payload was not sent, or a valid _Response_ was not received.
- A `Response` is returned if the payload was successfully sent to Apple and a documented response was received. This struct will contain more information about whether or not the push notification succeeded, its _apns-id_ and if applicable, more information around why it did not succeed.

To check if a `Notification` was successfully sent;

```go
res, err := client.Push(notification)
if err != nil {
  log.Println("There was an error", err)
  return
}

if res.Sent() {
  log.Println("Sent:", res.ApnsID)
} else {
  fmt.Printf("Not Sent: %v %v %v\n", res.StatusCode, res.ApnsID, res.Reason)
}
```

## Context & Timeouts

For better control over request cancellations and timeouts APNS/2 supports
contexts. Using a context can be helpful if you want to cancel all pushes when
the parent process is cancelled, or need finer grained control over individual
push timeouts. See the [Google post](https://blog.golang.org/context) for more
information on contexts.

```go
ctx, cancel = context.WithTimeout(context.Background(), 10 * time.Second)
res, err := client.PushWithContext(ctx, notification)
defer cancel()
```

## License

The MIT License (MIT)

Copyright (c) 2016 Adam Jones

Permission is hereby granted, free of charge, to any person obtaining a copy of this software and associated documentation files (the "Software"), to deal in the Software without restriction, including without limitation the rights to use, copy, modify, merge, publish, distribute, sublicense, and/or sell copies of the Software, and to permit persons to whom the Software is furnished to do so, subject to the following conditions:

The above copyright notice and this permission notice shall be included in all copies or substantial portions of the Software.

THE SOFTWARE IS PROVIDED "AS IS", WITHOUT WARRANTY OF ANY KIND, EXPRESS OR IMPLIED, INCLUDING BUT NOT LIMITED TO THE WARRANTIES OF MERCHANTABILITY, FITNESS FOR A PARTICULAR PURPOSE AND NON INFRINGEMENT. IN NO EVENT SHALL THE AUTHORS OR COPYRIGHT HOLDERS BE LIABLE FOR ANY CLAIM, DAMAGES OR OTHER LIABILITY, WHETHER IN AN ACTION OF CONTRACT, TORT OR OTHERWISE, ARISING FROM, OUT OF OR IN CONNECTION WITH THE SOFTWARE OR THE USE OR OTHER DEALINGS IN THE SOFTWARE.
