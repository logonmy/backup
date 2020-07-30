# -*-coding:utf-8-*-
"""
@File    : custom_function.py
@Time    : 2019/10/18 22:33
@Author  : 李高峰
@Email   : ligaofeng.own@gmail.com
"""
from pure_pagination import Paginator, PageNotAnInteger
import requests

from .models import *
from django.utils import timezone


# 自定义的函数，不是视图
def change_info(request):  # 修改网站访问量和访问ip等信息
    # 每一次访问，网站总访问次数加一
    count_nums = VisitNumber.objects.filter(id=1)
    if count_nums:
        count_nums = count_nums[0]
        count_nums.count += 1
    else:
        count_nums = VisitNumber()
        count_nums.count = 1
    count_nums.save()

    # 记录访问ip和每个ip的次数
    if 'HTTP_X_FORWARDED_FOR' in request.META:  # 获取ip
        client_ip = request.META['HTTP_X_FORWARDED_FOR']
        client_ip = client_ip.split(",")[0]  # 所以这里是真实的ip
    else:
        client_ip = request.META['REMOTE_ADDR']  # 这里获得代理ip
    # print(client_ip)

    ip_exist = Userip.objects.filter(ip=str(client_ip))
    if ip_exist:  # 判断是否存在该ip
        uobj = ip_exist[0]
        uobj.count += 1
    else:
        uobj = Userip()
        uobj.ip = client_ip
        uobj.count = 1
    uobj.last_visit = timezone.now()
    uobj.save()

    # 增加今日访问次数和今日访客数
    date = timezone.now().date()
    today = DayNumber.objects.filter(day=date)
    if today:
        temp = today[0]
        temp.count += 1
    else:
        temp = DayNumber()
        temp.dayTime = date
        temp.count = 1

    temp.save()
    visitor_ip_exist = temp.day_visit_ip.filter(ip=str(client_ip))
    if not visitor_ip_exist:  # 判断是否存在该ip
        temp.day_visit_ip.add(uobj)
        temp.save()


def posts_page_divide(request, posts):
    try:
        page = request.GET.get('page', 1)
    except PageNotAnInteger:
        page = 1
    paginator = Paginator(posts, request=request)
    blogs = paginator.page(int(page))
    return blogs


def verify_istrue(email):
    res = requests.get('http://ligaofeng.xyz:8101/email_verify?email={}'.format(email)).text
    if res == '1':
        return True
    else:
        return False
