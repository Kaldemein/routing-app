from celery import shared_task
import pika
from .models import User

@shared_task
def process_email(email):
    # email processing
    print(f"Processing email: {email}")
    user = User.objects.get(email=email)
    user.verified = True
    user.save()
    print(f"Updated user {user} to verified.")


def callback(ch, method, properties, body):
    email = body.decode()
    print(f"Received email for processing: {email}")
    process_email.delay(email)
    ch.basic_ack(delivery_tag=method.delivery_tag)

def start_listening():
    print('Connecting to RabbitMQ...')
    
    connection = pika.BlockingConnection(pika.ConnectionParameters('rabbitmq'))
    channel = connection.channel()
    channel.queue_declare(queue='verification_queue', durable=True)

    print('Declaring queue and starting to consume...')
    channel.basic_consume(queue='verification_queue', on_message_callback=callback)

    print('Waiting for messages. To exit press CTRL+C')
    channel.start_consuming()
