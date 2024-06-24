from datetime import datetime, timedelta
import requests

from .models import Point, Route, User
from decouple import config

def create_route_and_points(points, user_id):
    waypoints = [[point['lon'], point['lat']] for point in points]

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
    user.save()
