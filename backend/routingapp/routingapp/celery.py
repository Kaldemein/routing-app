from __future__ import absolute_import, unicode_literals
import os
from celery import Celery

# Устанавливаем переменную окружения для настройки Django
os.environ.setdefault('DJANGO_SETTINGS_MODULE', 'routingapp.settings')

app = Celery('routingapp')

# Загрузка конфигурации Celery из настроек Django
app.config_from_object('django.conf:settings', namespace='CELERY')

# Автоматическое обнаружение задач в приложениях Django
app.autodiscover_tasks()