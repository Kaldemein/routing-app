# api/management/commands/listen_rabbitmq.py

from django.core.management.base import BaseCommand
from api.tasks import start_listening

class Command(BaseCommand):
    help = 'Starts the RabbitMQ listener'

    def handle(self, *args, **kwargs):
        self.stdout.write(self.style.SUCCESS('Starting RabbitMQ listener...'))
        start_listening()
