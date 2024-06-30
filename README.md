# WhatsApp-Dev

![Screenshot](/screenshot.jpeg?raw=true "Screenshot")

Maildev, but for the WhatsApp business API.

This API is mainly designed for 2 things:

1. Locally test WhatsApp API code without all the hassle of the WhatsApp's testing phone number (handy if you work with multiple people)
2. Staging / Testing environment where you can freely test your WhatsApp integration code and don't have to be afraid of accidentally sending messages to real people or be limited to the testing API and have to change your api key every single day.

## Local setup and run

```
go run . --webhook-url http://your-app-webhook.local:8080/webhook
```

Replace `https://graph.facebook.com` with an instance of WhatsApp-Dev `http://localhost:1090` and, assuming you have setup your instance of WhatsApp-Dev the same as the real api, this should be 👌.

## Docker setup and run

```sh
# IMPORTANT! If you don't do this, docker will create a folder instead of a file
touch db.sqlite

docker run \
  -it \
  --rm \
  -p 1090:1090 \
  -v `pwd`/db.sqlite:/usr/src/app/db.sqlite \
  ghcr.io/mjarkk/whatsapp-dev:latest \
  app --webhook-url http://your-app-webhook.local:8080/api/webhook
```

Replace `https://graph.facebook.com` with an instance of WhatsApp-Dev `http://localhost:1090` and, assuming you have setup your instance of WhatsApp-Dev the same as the real api, this should be 👌.

Now visit http://localhost:1090 to see whatsapp-dev

## Options:

| name                          | flag                          | env                        | default                           |
| ----------------------------- | ----------------------------- | -------------------------- | --------------------------------- |
| webhook url (required)        | `--webhook-url` `-w`          | `WEBHOOK_URL`              |                                   |
| webhook verify token          | `--webhook-verify-token` `-t` | `WEBHOOK_VERIFY_TOKEN`     | _Randomly generated_              |
| Secrets seed                  | `--secrets-seed` `-s`         | `SECRETS_SEED`             | `fallback-secrets-seed`           |
| HTTP address                  | `--http-addr` `-a`            | `HTTP_ADDR`                | `:3000`                           |
| HTTP server username          | `--http-username` `-u`        | `HTTP_USERNAME`            | _No auth required if not defined_ |
| HTTP server password          | `--http-password` `-p`        | `HTTP_PASSWORD`            | _No auth required if not defined_ |
| Mocked phone number           | `--whatsapp-phone-number`     | `WHATSAPP_PHONE_NUMBER`    | _Randomly generated_              |
| Mocked phone number id        | `--whatsapp-phone-number-id`  | `WHATSAPP_PHONE_NUMBER_ID` | _Randomly generated_              |
| Facebook Graph token          | `--facebook-graph-token`      | `FACEBOOK_GRAPH_TOKEN`     | _Randomly generated_              |
| Facebook developer app secret | `--facebook-app-secret`       | `FACEBOOK_APP_SECRET`      | _Randomly generated_              |

_Note that all randomly generated values are generated using the secrets seed. If you don't change your seed, all randomly generated values will stay the same when restarting the service_

## Limitations / TODO

- Sending something other than text messages like images, videos, stickers, etc..
- Send status updates via the webhook
- Templates
  - Support Website, Phone number and Promo offer action buttons (currently only quick reply is supported)
  - Support Media header (currently only text is supported)

## From WhatsApp business API to this?

Replace `https://graph.facebook.com` with an instance of WhatsApp-Dev and, assuming you have setup your instance of WhatsApp-Dev the same as the real api, this should be 👌.

Note that if your instance of WhatsApp-Dev has different credentials then the real API you also need to change those.

It's recommended to use different credentials for this and the real api.

## Going from this to the real api?

Here is a rough outline of what you need todo:

- On [developers.facebook.com](https://developers.fracebook.com) create a new app
  - App kind: **Other**
  - App type: **Company**
- Add WhatsApp as product
- You will probably hit a roadblock here and need to verify your Facebook business
- Add a (real) phone number under WhatsApp > Api Setup
- Setup a webhook under WhatsApp > Configuration
  - Note: the webhook url must point to a working service
- Toggle the live mode within your API
- Add a credit card to your Facebook business here: [business.facebook.com](https://business.facebook.com/billing_hub/accounts/details)
  - Note: the selected org MUST the same as the org used to created your Facebook app earlier
  - Also add your credit card under Accounts > WhatsApp business accounts

Now you _should_ be able to send the _hello_world_ template using your real phone number (note that this cost money 💸)

You can add new templates on this website [business.facebook.com](https://business.facebook.com/wa/manage/home/)
