from django.urls import path
from . import views

urlpatterns = [
    path('route', views.RouteView.as_view(), name='route'),
    path('register', views.RegistrationView.as_view(), name="register"),
    path('login', views.LoginView.as_view(), name='login'),
    path('verification', views.EmailVerification.as_view(), name='email_verifications'),
]