FROM python:3.10.12
ENV PYTHONUNBUFFERED=1

# Установите рабочую директорию
WORKDIR /app

# Скопируйте файл зависимостей
COPY requirements.txt requirements.txt
RUN pip install -r requirements.txt

# Скопируйте весь проект в контейнер
COPY . .

# Установите переменные окружения
ENV DJANGO_SETTINGS_MODULE=routingapp.settings


