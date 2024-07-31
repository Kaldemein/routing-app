FROM python:3.10.12
ENV PYTHONUNBUFFERED=1

WORKDIR /app

COPY requirements.txt requirements.txt
RUN pip install -r requirements.txt

COPY . .


ENV DJANGO_SETTINGS_MODULE=routingapp.settings


