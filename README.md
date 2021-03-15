# Proxy server

### Запуск
```
sudo docker build -t proxy . 
sudo docker run -p 8080:8080 -p 8000:8000 -t proxy
```

### Порты
#### 8080 - прокси сервер
#### 8000 - repeater

### Зпросы
```
http://127.0.0.1:8000/requests – список запросов
http://127.0.0.1:8000/requests/id – вывод 1 запроса
http://127.0.0.1:8000/repeat/id – повторная отправка запроса
http://127.0.0.1:8000/scan/id – сканирование запрос
```
