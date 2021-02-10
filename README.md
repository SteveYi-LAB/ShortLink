# SteveYi Short Link

A Short Link System, Written by Golang and SQLing.

## Usage

GET /{code}

```
code: You will get it when you create a short link.
```

POST /api/create/

```
admin: open admin features  
token: contact admin to get it.(Need open admin features)  
link: need redicert link, (like https://steveyi.net/)  
g-recaptcha-response: Google Recaptcha  
custom: open custom code features  
customcode: custom code.(Need open custom features)  
```