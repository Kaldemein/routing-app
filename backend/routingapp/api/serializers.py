from rest_framework import serializers
from .models import Route, Point, User


class PointSerializer(serializers.ModelSerializer):
    arrival_time = serializers.DateTimeField(format='%Y-%m-%d %H:%M:%S', input_formats=['%Y-%m-%d %H:%M:%S'])

    class Meta:
            model = Point
            fields = ['lon', 'lat', 'arrival_time']

class UserSerializer(serializers.ModelSerializer):
    class Meta:
            model = User
            fields = ['id', 'first_name', 'last_name', 'email',]

class RouteSerializer(serializers.ModelSerializer):
    points = PointSerializer(many=True, read_only=True)
    user = UserSerializer() 

    start_at = serializers.DateTimeField(format='%Y-%m-%d %H:%M:%S', input_formats=['%Y-%m-%d %H:%M:%S'])
    ends_at = serializers.DateTimeField(format='%Y-%m-%d %H:%M:%S', input_formats=['%Y-%m-%d %H:%M:%S'])

    
    class Meta:
        model = Route
        fields = ['id', 'start_at', 'ends_at', 'created_at', 'user', 'ors_coordinates', 'points']