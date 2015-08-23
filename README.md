# Hook notifier

An app to accept webhook notifications, translate and send them to configured services.

## Features:

* Configurable service payloads `./templates`
* Configurable service URLs
* Templates are Go [html/templates](https://golang.org/pkg/html/template/#pkg-overview)


## Usage:

1. Add a service template in `./templates/` like the example template.
2. Add a service URL as an ENV variable with the same name as upper case.
3. Post to http://127.0.0.1:3000/service with the payload.

This will translate your payload into the template and POST to the service URL.


## Payload:

```json
{
  "app": "test-convox",
  "release": "v7",
  "url": "http://test-convox.example.com"
}
```

## Example:

`./templates/slack.tmpl`

```json
{
  "channel": "#notifications",
  "username": "webhookbot",
  "text": "Deployed {{.App}} {{.Release}} to {{.URL}}",
  "icon_emoji": ":ghost:"
}
```

`.env`
```shell
SLACK_URL=http://example.slack.com/the/webhook/endpoint
```

Run `$ ./notifier`

Send the webhook payload:
```shell
$ curl -X "POST" "http://127.0.0.1:3000/slack" \
  -d $'{   "app": "test-convox",
  "release": "v7",
  "url": "http://test-convox.example.com"
}'
```

Notifier's output:
```shell
2015/08/23 13:21:09 Received Webhook for slack
2015/08/23 13:21:09 Sending Notification to http://example.slack.com/the/webhook/endpoint
```


`http://example.slack.com/the/webhook/endpoint` was sent:
```json
{
"channel": "#notifications",
"username": "webhookbot",
"text": "Deployed test-convox v7 to http://test-convox.example.com",
"icon_emoji": ":ghost:"
}
```
