# -*-coding:utf-8-*-
"""myblog URL Configuration

The `urlpatterns` list routes URLs to views. For more information please see:
    https://docs.djangoproject.com/en/2.2/topics/http/urls/
Examples:
Function views
    1. Add an import:  from my_app import views
    2. Add a URL to urlpatterns:  path('', views.home, name='home')
Class-based views
    1. Add an import:  from other_app.views import Home
    2. Add a URL to urlpatterns:  path('', Home.as_view(), name='home')
Including another URLconf
    1. Import the include() function: from django.urls import include, path
    2. Add a URL to urlpatterns:  path('blog/', include('blog.urls'))
"""
from django.urls import path, re_path
from blog.views import *
from django.views.static import serve


app_name = 'blog'

urlpatterns = [
    path('', index, name='index'),
    path('about.html', about, name='about'),
    path('contact.html', contact, name='contact'),
    path('search/', search, name='search'),
    path('blog/<int:bid>/', blog_detail, name='blog_detail'),
    path('login/', log_in, name='login'),
    path('logout/', LogoutView.as_view(), name='logout'),
    path('register/', RegisterView.as_view(), name='register'),
    path('passwd/', PasswdView.as_view(), name='passwd'),
    path('reset/', reset, name='reset'),
    path('reset_email/', reset_email, name='reset_email'),
    path('setpd/', setpd, name='setpd'),
    path('validate/', Validate.as_view(), name='validate'),
    path('active/<str:active_code>', ActiveView.as_view(), name='active'),
    path('comment_sub/', comment_sub, name='comment_sub'),
    path('email_val/', email_val, name='email_val'),
    path('new_email_val/', new_email_val, name='new_email_val'),
    path('category/', category, name='category'),
    path('tag/', tag, name='tag'),
    path('all/', all, name='all'),
    path('send_contact_mail/', send_contact_mail, name='send_contact_mail'),
    path('subscript/', subscript, name='subscript'),
    path('subscript_val/', subscript_val, name='subscript_val'),
    path('username_val/', username_val, name='username_val'),
    path('change_info/', change_profile, name='change_info'),
    path('change_email/', change_email, name='change_email'),
    path('profile_modify/', profile_modify, name='profile_modify'),
    path('admin/upload/<str:dir_name>', upload_image, name='upload_image'),
    re_path(r'^upload/(?P<path>.*)$', serve, {'document_root': settings.MEDIA_ROOT, }),
]
