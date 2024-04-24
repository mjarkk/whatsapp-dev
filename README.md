# WhatsApp-Dev

![Screenshot](/screenshot.jpeg?raw=true "Screnshot")

Maildev but for the WhatsApp business API.

This API is mainly designed for 2 things:

1. Locally test whatsapp API code without all the hassle of the WhatsApp's testing phone number (handy if you work with multiple people)
2. Staging / Testing envourments where you freely want to test your whatsapp integration code but don't want to be scared of accedentially sending messages to real people or be limited to the testing API and have to change your api key every day.

## Limitations / TODO

- Sending something other than text messages like images, videos, sticker, etc..
- Send status updates via the webhook
- Templates
  - Support Website, Phone number and Promo offer action buttons (currently only quick reply is supported)
  - Support Media header (currently only text is supported)

## Local setup and run

```
go run . --webhook-url http://your-app-webhook.local:8080/webhook
```

## Docker setup and run

```sh
docker build -t whatsapp-dev .

# IMPORTANT! If you don't do this docker will create a folder instaid of a file
touch db.sqlite

docker run \
  -it \
  --rm \
  -p 80:1090 \
  -v `pwd`/db.sqlite:/usr/src/app/db.sqlite \
  whatsapp-dev \
  app --webhook-url http://your-app-webhook.local:8080/api/webhook
```

Now visit http://localhost to see whatsapp-dev

## From WhatsApp business API to this?

Replace `https://graph.facebook.com` with a instance of WhatsApp-Dev and assuming you have setup your instance of WhatsApp-Dev the same as the real api this should be all.

Note that if your instance of WhatsApp-Dev has different credentials then the real API you also need to change those.

It's recomment to use diffrent credentials for this and the real api.

## Going from this to the real api?

Here is a rough outline of what you need todo:

- On [developers.facebook.com](https://developers.fracebook.com) create a new app
  - App kind: **Other**
  - App type: **Company**
- Add WhatsApp as product
- You will probably hit a roadblock here and need to verivy your facebook business
- Add a (real) phone number under Whatsapp > Api Setup
- Setup a webhook (that works. important!) under Whatsapp > Configuration
- Toggle the live mode within your API
- Add a creditcard to your facebook bussness here: [business.facebook.com](https://business.facebook.com/billing_hub/accounts/details)
  - Note the selected org MUST the same as the org used to created your facebook app earlier
  - Also add your creditcard under Accounts > Whatsapp Business-accounts

Now you _should_ beable to send the hello_world template using your real phone number (note that this cost money 💸)

You can add new templates on this website [business.facebook.com](https://business.facebook.com/wa/manage/home/)
