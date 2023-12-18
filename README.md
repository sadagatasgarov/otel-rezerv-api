Qeydiyyatdan kecmek ucun
(POST) localhost:5000/api/user
```
{
  "firstName": "Admin",
  "lastName": "admin",
  "email": "admin@admin.com",
  "password": "admin_admin",
"isAdmin": true
}
```
Giriş etmək
(POST) localhost:5000/api/auth 
```
{
  "email":"admin@admin.com",
  "password":"admin_admin"
}
```
Giris eden zaman token verir onu Headers-e yazmaq lazimdirş
X-Api-Token: eyJhbGciOiJIUzI1NiIsInR5cCI6IkpXVCJ9.eyJlbWFpbCI6ImFkbWluQGFkbWluLmNvbSIsImV4cGlyZXMiOjE3MDI4NTk3MDUsImlkIjoiNjU3ZjUxYWRlNGEzNTZmMzk0Njk2MGMwIn0.roUx6JqTL5Bf5gkFX7mvrSSnjHyWkERckRQP2iGWSiY


istifadecileri siralamaq 
(GET) localhost:5000/api/v1/user/
İD-ye gore baxmaq
(GET)localhost:5000/api/v1/user/657dd7d42f6ef4be90e98b88

Hotellere baxmaq
(GET)localhost:5000/api/v1/hotel?limit=20&page=0
Otaqlara baxmaq
localhost:5000/api/v1/rooms

Otaq rezerv etmek
(POST)localhost:5000/api/v1/room/6575b0e286b1ddaf93dd6f98/book
```
{
  "numPersons":1,
  "fromDate":"2024-11-13T15:04:05Z",
  "tillDate":"2024-11-14T15:04:05Z"
}
```
Admin butun rezervlere baxa biler
(GET)localhost:5000/api/v1/admin/booking

User ise ancaq oz rezev etdiklerine baxa biler
(GET)localhost:5000/api/v1/booking/657a220c6ddfe7324b7fd90c


Proyekti ise salmaq ucun docker-compose.yml fayli
```
version: '3.1'

services:

  mongo:
    image: mongo:4.0.28-xenial
    container_name: mongo
    restart: always
    ports:
      - "27017:27017"
    environment:
      MONGO_INITDB_ROOT_USERNAME: root
      MONGO_INITDB_ROOT_PASSWORD: example
    volumes:
      - local_mongo_data:/data/db

  hotel-rezerv:
    image: sadagatasgarov/hotel-rezerv
    restart: always
    ports:
      - "5000:5000"
    environment:
      DB_NAME: hotel-rezervation
      HTTP_LISTENING_PORT: :5000
      JWT_SECRET: bunuAYRIbirYERDEsaxlamaqLAZIMDIR
      MONGODB_URL: mongodb://root:example@mongo:27017/

volumes:
  local_mongo_data:
```


