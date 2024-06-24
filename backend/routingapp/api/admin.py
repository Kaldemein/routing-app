from django.contrib import admin
from .models import User, Route, Point

# Register your models here.
admin.site.register(User)
admin.site.register(Route)
admin.site.register(Point)