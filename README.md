# SteveYi Short Link

A Short Link System, Written by Golang and PostgreSQL.

## Usage

GET /{code}

```
code: You will get it when you create a short link.
```

POST /api/v1/create/

```
admin: open admin features  
token: contact admin to get it.(Need enabled admin features)  
link: need redicert link, (like as https://steveyi.net/)  
g-recaptcha-response: Google Recaptcha  
custom: open custom code features  
customcode: custom code.(Need enabled custom features)  
```
