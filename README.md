# Homework Helper

Homework Helper — удобный сервис для проверки домашних заданий.  
Проект позволяет учителям загружать задания, а ученикам прикреплять решения.

## Установка

```bash
docker-compose build homework-helper
docker-compose up homework-helper


```

## Конфигурация
Пример файла .env:

```bash
DB_HOST="db"
DB_NAME="postgres"
DB_PORT=5432
DB_USERNAME="postgres"
DB_PASSWORD="qwerty"
SSL_MODE="disable"

SALT="gKnre432mj"
PRIVATE_KEY="-----BEGIN EC PRIVATE KEY-----
MHcCAQEEIN8J7BqIV6EDOAh+jnuVnyMJYDT2UtYXOf8mEAPLBNVRoAoGCCqGSM49
AwEHoUQDQgAEzi4jzzjCD/OOgsxFnvSxuigPVrzVB3UopAVQKDsoYFwrnRMtGaFc
S2IZDDh1NXvz5tpdR5c9unPvH+fXrlibyw==
-----END EC PRIVATE KEY-----"
PUBLIC_KEY="-----BEGIN PUBLIC KEY-----
MFkwEwYHKoZIzj0CAQYIKoZIzj0DAQcDQgAEzi4jzzjCD/OOgsxFnvSxuigPVrzV
B3UopAVQKDsoYFwrnRMtGaFcS2IZDDh1NXvz5tpdR5c9unPvH+fXrlibyw==
-----END PUBLIC KEY-----"

SYMMETRICK_KEY="5ebKcIHM21lmSFqQrOABR3CY5NtMDDF0"

EMAIL="alexmikheev05@gmail.com"
EMAIL_PASSWORD="injk runc eypk dcvh"