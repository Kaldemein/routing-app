from datetime import datetime, timedelta
from re import sub
import requests
import pika
from rest_framework.views import APIView
import jwt
from .models import Point, Route, User
from decouple import config
import polyline


def create_route_and_points(points, user_id):

    #saving points in a list
    waypoints = [[point['lon'], point['lat']] for point in points]
    print(waypoints)

    #turn list into string
    string_waypoints = ''
    for i, point in enumerate(waypoints):
        lon = point[0]
        lat = point[1]
        string_waypoints+=f"{lon},{lat}"
        if i != len(waypoints)-1:
            string_waypoints+=";"

    #OSRM API string
    url = f"http://osrm:5000/route/v1/driving/{string_waypoints}"

    #request to OSRM API
    try:
        response = requests.post(url)
        response.raise_for_status() 
    except requests.exceptions.RequestException as e:
        print("Error:", e)

    #getting data from OSRM response
    osrm_data = response.json()
    encoded_coordinates = osrm_data['routes'][0]['geometry']
    decoded_coordinates = polyline.decode(encoded_coordinates)

    route_start_at = datetime.now()
    route_duration = osrm_data['routes'][0]['duration']
    route_ends_at = route_start_at + timedelta(seconds=route_duration)

    
    #creating route in DB
    user_instance = User.objects.get(id=user_id)  

    route = Route(
        start_at = route_start_at,
        ends_at = route_ends_at,
        user=user_instance,
        osrm_coordinates = decoded_coordinates
    )
    route.save()

    #creating points in DB
    arrival_time = datetime.now()
    for i, waypoint in enumerate(waypoints):
        if i>0:
            legs_duration = osrm_data['routes'][0]['legs'][i-1]['duration'] #get duration between two waypoints
            arrival_time+=timedelta(seconds=legs_duration)
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
    