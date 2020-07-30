# -*-coding:utf-8-*-
"""
@File    : cus_tags.py
@Time    : 2019/10/20 14:00
@Author  : 李高峰
@Email   : ligaofeng.own@gmail.com
"""
from django import template
import re
from blog.models import *

register = template.Library()


@register.filter(name='get_range')
def get_range(value):
    return range(value)


@register.filter(name='img_strip')
def img_strip(value):
    imgs = re.findall('<img(.*?)>', value)
    for img in imgs:
        value = value.replace(img, "")
    return value


@register.filter(name='img_size')
def img_size(value, width):
    imgs = re.findall('<img(.*?)>', value)
    for img in imgs:
        value = value.replace(img, img + " width='{}'".format(width))
    return value


@register.filter(name='add_class_img')
def add_class_img(value):
    imgs = re.findall('<img(.*?)>', value)
    for img in imgs:
        value = value.replace(img, img + " class='{}'".format('img'))
    return value


@register.filter(name='html_tags_strip')
def html_tags_strip(value):
    res = re.findall('<(.*?)>', value)
    for i in res:
        value = value.replace(i, "")
    return value


@register.filter(name='views_count')
def views_count(value):
    post = Post.objects.get(id=value)
    num = post.views.all()
    return num.count()


@register.filter(name='cover_url')
def cover_url(value):
    default = ['/static/images/post_cover_default.png', '/static/images/banner_default.png',
               '/static/images/user_icon_default.ico']
    if str(value) in default:
        return value
    else:
        return '/upload/' + str(value)
