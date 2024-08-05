from django.db import models
from django.contrib.auth.hashers import make_password

class User(models.Model):
    id = models.AutoField(primary_key=True)
    first_name = models.CharField(max_length=30, blank=False)
    last_name = models.CharField(max_length=30, blank=False)
    email = models.EmailField(unique=True, blank=False)
    password_hash = models.CharField(max_length=255, blank=False)
    verified = models.BooleanField(default=False)

    def set_password(self, password):
        self.password_hash = make_password(password)

    def __str__(self):
        return f'{self.first_name} {self.last_name}'

class Route(models.Model):
    id = models.AutoField(primary_key=True)
    start_at = models.DateTimeField()
    ends_at = models.DateTimeField()
    created_at = models.DateTimeField(auto_now_add = True)
    user = models.ForeignKey(User, on_delete=models.CASCADE, related_name='routes')
    osrm_coordinates = models.JSONField()

    def __str__(self):
        return f'Route number {self.id} by user {self.user}'

class Point(models.Model):
    id = models.AutoField(primary_key=True)
    lon = models.DecimalField(max_digits=9, decimal_places=6, blank=False)
    lat = models.DecimalField(max_digits=9, decimal_places=6, blank=False)
    route = models.ForeignKey(Route, on_delete=models.CASCADE, related_name='points')
    arrival_time = models.DateTimeField()

    def __str__(self):
        return f'Point {self.lat} by user {self.lon}'
