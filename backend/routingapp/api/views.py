from datetime import datetime, timedelta
from re import sub
from decouple import config
from django.http import JsonResponse
from django.views.decorators.csrf import csrf_exempt
from django.utils.decorators import method_decorator
from django.contrib.auth.hashers import check_password
import json
from .models import User
from .serializers import RouteSerializer
from .service_object import create_route_and_points, create_user, send_to_queue, get_user_id_by_header, generate_JWT
from rest_framework.views import APIView

@method_decorator(csrf_exempt, name='dispatch')
class RouteView(APIView):
    def post(self, request):
        try:
            data = json.loads(request.body)
            
            user_id = get_user_id_by_header(request)
            
            try:
                User.objects.get(id=user_id)
                points = data.get('points')
                route = create_route_and_points(points, user_id)
                route_serializer = RouteSerializer(route)
            except User.DoesNotExist:
                return JsonResponse({'error': 'bad request'}, status=401)

            return JsonResponse(route_serializer.data, status = 200)
        
        except Exception as e:
            return JsonResponse({'error': str(e)}, status=500)


@method_decorator(csrf_exempt, name='dispatch')
class RegistrationView(APIView):
    def post(self, request):
        try:
            userData = json.loads(request.body)
            print(request.body)
            create_user(userData['first_name'],
                        userData['last_name'],
                        userData['email'],
                        userData['password'])

            return JsonResponse({"message": f"{userData['first_name']}, your account was successfully created"}, status=201) 
        except Exception as e:
            return JsonResponse({"error": str(e)}, status=500)

@method_decorator(csrf_exempt, name='dispatch')
class EmailVerification(APIView):
    def post(self, request):
        try:
            user_id = get_user_id_by_header(request)
            try:
                user = User.objects.get(id=user_id)
                email = user.email
                send_to_queue(email)
            except User.DoesNotExist:
                return JsonResponse({'error': 'bad request'}, status=401)

            return JsonResponse({"message": f"{email}, was successfully delevered"}, status=201) 
        except Exception as e:
            return JsonResponse({"error": str(e)}, status=500)

@method_decorator(csrf_exempt, name='dispatch')
class LoginView(APIView):
    def post(self, request):
        try:
            email = request.data.get('email')
            password = request.data.get('password')

            try:
                user = User.objects.get(email=email)
            except User.DoesNotExist:
                return JsonResponse({'error': 'Invalid credentials'}, status=401)

            if check_password(password, user.password_hash):
                encoded_token = generate_JWT(user.id)
                return JsonResponse({
                    'encoded_token': str(encoded_token),
                })
            else:
                return JsonResponse({'error': 'Invalid credentials'}, status=401)

        except Exception as e:
            return JsonResponse({'error': str(e)}, status=500)
