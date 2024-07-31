from datetime import datetime, timedelta
from re import sub
import requests
import pika
from rest_framework.views import APIView
import jwt
from .models import Point, Route, User
from decouple import config

def create_route_and_points(points, user_id):

    waypoints = [[point['lon'], point['lat']] for point in points]
    print(waypoints)

    response = requests.post(
            'https://api.openrouteservice.org/v2/directions/driving-car/geojson',
            json={'coordinates': waypoints},
            headers={
                'Authorization': config('ORS_TOKEN'),
                'Content-Type': 'application/json',
            }
        )

    ors_data = response.json()
    ors_coordinates  = ors_data['features'][0]['geometry']['coordinates']
    ors_segments = ors_data['features'][0]['properties']['segments']
    
    # CREATING ROUTE
    route_start_at = datetime.now()
    route_duration_seconds = ors_data['features'][0]['properties']['summary']['duration']
    route_ends_at = route_start_at + timedelta(seconds = route_duration_seconds)

    user_instance = User.objects.get(id=user_id)  

    route = Route(
        start_at = route_start_at,
        ends_at = route_ends_at,
        user=user_instance,
        ors_coordinates = ors_coordinates
    )
    route.save()

    #CREATING POINTS
    arrival_time = datetime.now()
    for i, waypoint in enumerate(waypoints):
        if i>0:
            segment_duration = ors_segments[i-1]['duration']
            arrival_time+=timedelta(seconds=segment_duration)
        point = Point(
            lon = waypoint[0],
            lat = waypoint[1],
            route = route,
            arrival_time = arrival_time
        )
        point.save()
    
    return route 

def create_user(first_name, last_name, email, password):
    
    user = User(first_name = first_name,
                last_name = last_name,
                email = email)
    user.set_password(password)
    print(user.password_hash)
    user.save()
    
    
def generate_JWT(user_id):
    SECRET_KEY = config('SECRET_KEY')
    payload = {
                'user_id': user_id,
                'exp': datetime.now() + timedelta(hours=1),
    }   
    encoded_token = jwt.encode(payload, SECRET_KEY, algorithm='HS256')
    return encoded_token


def send_to_queue(self, email):
        connection = pika.BlockingConnection(pika.ConnectionParameters('rabbitmq'))
        channel = connection.channel()
        channel.queue_declare(queue='email_queue', durable=True)
        channel.basic_publish(exchange='',
                              routing_key='email_queue',
                              body=email,
                              properties=pika.BasicProperties(
                                 delivery_mode=2,  # make message persistent
                              ))
        connection.close()
    
def get_user_id_by_header(request):
    auth_headers = request.META['HTTP_AUTHORIZATION']
    encoded_token = sub('Bearer ', '', auth_headers) 
    SECRET_KEY = config('SECRET_KEY')
    decoded_token = jwt.decode(encoded_token, SECRET_KEY, algorithms=["HS256"])
     
    user_id = decoded_token.get('user_id')
    
    return user_id
    