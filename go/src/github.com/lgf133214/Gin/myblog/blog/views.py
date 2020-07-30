from django.shortcuts import render, reverse
from django.http import *

from blog.models import *
from django.views.generic.base import View
from django.db.models import Q
from django.contrib.auth import login, logout
from django.contrib.auth.hashers import make_password, check_password
from random import Random
from django.core.mail import send_mail
from .models import EmailVerifyRecord
from myblog.settings import DEFAULT_FROM_EMAIL, DOMAIN
import json
from .custom_function import change_info, posts_page_divide, verify_istrue
from django.utils import timezone
from django.conf import settings
from django.views.decorators.csrf import csrf_exempt
import os
import uuid
import datetime as dt


# Create your views here.

def get_day_data():
    if DayNumber.objects.filter(day=timezone.now().date()):
        visitor_day_count = DayNumber.objects.get(day=timezone.now().date()).day_visit_ip.count()
    else:
        visitor_day_count = DayNumber.objects.create()
        visitor_day_count.save()
        visitor_day_count = DayNumber.objects.get(day=timezone.now().date()).day_visit_ip.count()
    return visitor_day_count


def get_categories():
    return BlogCategory.objects.all()


def get_tags():
    return Tags.objects.all()


def get_recommends():
    return Recommend.objects.all()


def index(request):
    change_info(request)
    banner_list = Banner.objects.all()
    post_list = Post.objects.order_by('-pub_date').all()[:4]
    ctx = {
        'banner_list': banner_list,
        'post_list': post_list,
        'recommend_list': get_recommends(),
        'blogcategory_list': get_categories(),
        'tag_list': get_tags(),
        'visitor_day_count': get_day_data(),
    }
    return render(request, 'index.html', ctx)


# 生成随机字符串
def make_random_str(randomlength=8):
    str = ''
    chars = 'AaBbCcDdEeFfGgHhIiJjKkLlMmNnOoPpQqRrSsTtUuVvWwXxYyZz0123456789'
    length = len(chars) - 1
    random = Random()
    for i in range(randomlength):
        str += chars[random.randint(0, length)]
    return str


# 发送邮件
def my_send_email(email, send_type="register"):
    email_record = EmailVerifyRecord()
    code = make_random_str(16)
    email_record.code = code
    email_record.email = email
    email_record.send_type = send_type
    email_record.save()

    if send_type == "register":
        email_title = "Feng`s Blog-注册激活链接"
        email_body = "请点击下面的链接激活你的账号: https://{0}/active/{1}".format(DOMAIN, code)

        send_mail(email_title, email_body, DEFAULT_FROM_EMAIL, [email])
    elif send_type == "forget":
        email_title = "Feng`s Blog-密码重置链接"
        email_body = "请点击下面的链接重置密码：https://{}/reset/?code={}".format(DOMAIN, code)

        send_mail(email_title, email_body, DEFAULT_FROM_EMAIL, [email])
    elif send_type == "update_email":
        email_title = "Feng`s Blog-邮箱修改链接"
        email_body = "请点击下面的链接完成邮箱的修改: https://{}/reset_email/?code={}".format(DOMAIN, code)

        send_mail(email_title, email_body, DEFAULT_FROM_EMAIL, [email])


def reset_email(request):
    code = request.GET.get('code')
    if code:
        email_record = EmailVerifyRecord.objects.filter(code=code)
        if email_record:
            user = BlogUser.objects.filter(email=email_record[0].email)
            user = user.get()
            user.is_active = True
            user.save()
            email_record.delete()
            return render(request, 'tips.html', {'msg': '激活成功', 'location': '/', 'recommend_list': get_recommends(),
                                                 'blogcategory_list': get_categories(),

                                                 'visitor_day_count': get_day_data(), })
        else:
            return HttpResponseNotFound()
    else:
        return HttpResponseNotFound()


class ActiveView(View):
    def get(self, request, active_code):
        records = EmailVerifyRecord.objects.filter(code=active_code)
        if records:
            email = records[0].email
            user = BlogUser.objects.get(email=email)
            user.is_active = True
            name = user.username
            user.save()
            records[0].delete()
            return render(request, "login.html", {'uname': name, 'active_msg': '激活成功', })
        return HttpResponseNotFound()


def log_in(request):
    return render(request, 'login.html')


def reset(request):
    code = request.GET.get('code')
    email = EmailVerifyRecord.objects.filter(code=code, send_type='forget')
    if email:
        return render(request, 'reset.html', {'email': email.get().email, 'visitor_day_count': get_day_data(),
                                              'blogcategory_list': get_categories(), })
    else:
        return HttpResponseNotFound()


def setpd(request):
    email_url = request.POST.get('email')
    new_pd = request.POST.get('password')
    if EmailVerifyRecord.objects.filter(email=email_url, send_type='forget'):
        EmailVerifyRecord.objects.filter(email=email_url, send_type='forget').delete()
        user = BlogUser.objects.get(email=email_url)
        user.password = make_password(new_pd)
        user.is_active = True
        user.save()

        return HttpResponse(json.dumps({'msg': '修改成功！'}))
    else:
        return HttpResponseNotFound()


def about(request):
    return render(request, 'about.html',
                  {'visitor_day_count': get_day_data(), 'blogcategory_list': get_categories(), })


class PasswdView(View):
    def get(self, request):
        return render(request, 'passwd.html',
                      {'visitor_day_count': get_day_data(), 'blogcategory_list': get_categories(), })

    def post(self, request):
        email = request.POST.get('email')

        tmp = BlogUser.objects.filter(email=email)
        if tmp:
            tmp = tmp[0]
            tmp.is_active = False
            tmp.save()
            my_send_email(email, send_type="forget")
            return HttpResponse(json.dumps({'msg': '重置链接已发送至您的邮箱，请注意查收', 'status': 1}))
        else:
            return HttpResponse(json.dumps({'msg': '邮箱输入错误，请检查后重试', 'status': 0}))


class RegisterView(View):
    def get(self, request):
        return render(request, 'register.html')

    def post(self, request):
        username = request.POST.get('username')
        if BlogUser.objects.filter(username=username):
            return HttpResponse(json.dumps({'id': '0', 'error_msg': '该用户名已被使用', 'tip': '名字早被人取了，哈哈😄'}))
        password = request.POST.get('password')
        email = request.POST.get('email')
        tmp = BlogUser.objects.filter(email=email)
        if tmp:
            tmp = tmp.get()
            if tmp.is_active:
                return HttpResponse(json.dumps({'id': '0', 'error_msg': '该邮箱已被注册且已激活', 'tip': '世界上有两个你？！！🐂🍺'}))
            else:
                return HttpResponse(json.dumps({'id': '0', 'error_msg': '该邮箱已被注册但未激活', 'tip': '快去邮箱看看吧📫'}))

        status = verify_istrue(email)

        if status:
            my_send_email(email)

            user = BlogUser()
            user.username = username
            user.password = make_password(password)
            user.email = email
            user.is_active = False

            user.save()
            return HttpResponse(json.dumps({'error_msg': '快去邮箱看看吧📫', 'id': '1'}))

        else:
            return HttpResponse(json.dumps({'id': '0', 'error_msg': '发送失败', 'tip': '请检查输入的邮箱是否正确'}))


class Validate(View):
    def get(self, request):
        return render(request, 'login.html', {})

    def post(self, request):
        username = request.POST.get('username')
        password = request.POST.get('password')
        user = BlogUser.objects.filter(username=username)
        if not user:
            user = BlogUser.objects.filter(email=username)
        if user:
            if check_password(password, user.get().password):
                if user.get().is_active:
                    login(request, user.get())
                    return HttpResponse(json.dumps({'id': '2'}))
                else:
                    return HttpResponse(json.dumps({"error_msg": '还没激活呢, 快去邮箱看看吧', 'id': '1', 'tip': '不知道要激活吗？！！！'}))
            else:
                return HttpResponse(json.dumps({"error_msg": '用户名或密码错误', 'id': '0', 'tip': '是不是不知道账号密码！？？'}))

        else:
            return HttpResponse(json.dumps({"error_msg": '用户名或密码错误', 'id': '0', 'tip': '是不是不知道账号密码！？？？'}))


class LogoutView(View):
    def get(self, request):
        logout(request)
        return HttpResponseRedirect(reverse("myblog:index"))


def contact(request):
    return render(request, 'contact.html',
                  {'visitor_day_count': get_day_data(), 'blogcategory_list': get_categories(), })


def search(request):
    kw = request.POST.get('keyword') if not request.GET.get('kw') else request.GET.get('kw')
    if kw:
        post_list = Post.objects.filter(Q(title__icontains=kw) | Q(content__icontains=kw)).order_by('-pub_date')
        blogs = posts_page_divide(request, post_list)
        ctx = {
            'posts': blogs,
            'kw': kw,
            'visitor_day_count': get_day_data(),
            'blogcategory_list': get_categories(),
        }
        return render(request, 'search-result.html', ctx)
    else:
        return HttpResponseNotFound()


def blog_detail(request, bid):
    post = Post.objects.get(id=bid)
    if 'HTTP_X_FORWARDED_FOR' in request.META:  # 获取ip
        client_ip = request.META['HTTP_X_FORWARDED_FOR']
        client_ip = client_ip.split(",")[0]  # 所以这里是真实的ip
    else:
        client_ip = request.META['REMOTE_ADDR']  # 这里获得代理ip
    if not post.views.filter(ip=str(client_ip)):
        post.views.add(Userip.objects.get(ip=str(client_ip)))
    post.save()

    ctx = {
        'post': post,
        'visitor_day_count': get_day_data(),
        'blogcategory_list': get_categories(),
        'tag_list': get_tags(),
        'recommend_list': get_recommends(),
    }
    return render(request, 'blog-details.html', ctx)


def comment_sub(request):
    if request.user.is_authenticated:
        comment = request.POST.get('comment')
        id = request.POST.get('post_id')
        com_obj = Comment()
        com_obj.content = comment
        com_obj.user = request.user
        com_obj.post = Post.objects.get(id=id)
        com_obj.post.save()
        com_obj.save()
        return render(request, 'tips.html',
                      {'visitor_day_count': get_day_data(), 'blogcategory_list': get_categories(), 'msg': '留言成功'})
    else:
        comment = request.POST.get('comment')
        id = request.POST.get('post_id')
        email = request.POST.get('email')
        com_obj = Comment()
        com_obj.content = comment
        com_obj.user = BlogUser.objects.get(email=email)
        com_obj.post = Post.objects.get(id=id)
        com_obj.save()
        com_obj.post.save()
        return render(request, 'tips.html',
                      {'visitor_day_count': get_day_data(), 'blogcategory_list': get_categories(), 'msg': '留言成功'})


def email_val(request):
    email = request.POST.get('email')
    user = BlogUser.objects.filter(email=email)
    if user:
        if user.get().is_active:
            return HttpResponse(json.dumps({'msg': "", 'status': 1}))
        else:
            return HttpResponse(json.dumps({'msg': "邮箱未激活，先去激活吧", 'status': 0}))
    else:
        return HttpResponse(json.dumps({'msg': "邮箱输入错误，请在检查后重试", 'status': 0}))


def send_contact_mail(request):
    if request.method == 'POST':
        name = request.POST.get('name')
        email = request.POST.get('email')
        message = request.POST.get('message')
        email_title = " - ".join([email, name, ' 反馈建议'])
        send_mail(email_title, message, DEFAULT_FROM_EMAIL, [DEFAULT_FROM_EMAIL])
        return HttpResponse(json.dumps({'msg': '反馈成功，感谢您的耐心指正'}))
    else:
        return HttpResponseNotFound()


def category(request):
    if request.method == "GET":
        id = request.GET.get('id')
        if id:
            cat = BlogCategory.objects.filter(id=int(id))
            if cat:
                posts = Post.objects.filter(category_id=int(id)).order_by('-pub_date')
                blogs = posts_page_divide(request, posts)
                return render(request, 'blog_list.html',
                              {'desc': "Category:" + cat.get().name, 'posts': blogs, 'categories': cat,
                               'visitor_day_count': get_day_data(), 'blogcategory_list': get_categories(),
                               'tag_list': get_tags(), 'recommend_list': get_recommends(), })
            else:
                return HttpResponseNotFound()
        else:
            return HttpResponseNotFound()
    else:
        return HttpResponseNotFound()


def subscript_val(request):
    email = request.POST.get('email')
    user = BlogUser.objects.filter(email=email)
    if user:
        if user[0].is_active:
            email_obj = EmailSubscription.objects.filter(email=email)
            if email_obj:
                return HttpResponse(json.dumps({'msg': "这个邮箱已经订阅过了哦", 'status': 0}))
            else:
                return HttpResponse(json.dumps({'msg': "", 'status': 1}))
        else:
            return HttpResponse(json.dumps({'msg': "邮箱未激活，先去激活吧", 'status': 0}))
    else:
        return HttpResponse(json.dumps({'msg': "邮箱输入错误，请在检查后重试", 'status': 0}))


def subscript(request):
    email = request.POST.get('email')
    if email:
        user = BlogUser.objects.get(email=email)
        if not EmailSubscription.objects.filter(email=email, user=user):
            email_obj = EmailSubscription.objects.create(email=email, user=user)
            email_obj.save()
            email_title = "Feng`s Blog-订阅成功提醒"
            email_body = "Hello，就在刚才，你已经开通了我的邮箱通知，如果有关于本站重要的消息我会第一时间进行通知的，同时，本站的解释权归本人所有，如有需要，可以更改本站的一切数据，对于由此造成的不便，我表示我也没有办法😄，毕竟数据库在我这里🙂。。。"

            send_mail(email_title, email_body, DEFAULT_FROM_EMAIL, [email])

            return render(request, 'tips.html',
                          {'visitor_day_count': get_day_data(), 'blogcategory_list': get_categories(),
                           'msg': '订阅成功'})
    return HttpResponseNotFound()


def tag(request):
    if request.method == 'GET':
        id = request.GET.get('id')
        if id:
            tag = Tags.objects.filter(id=int(id))
            if tag:
                posts = Post.objects.filter(tags=tag.get()).order_by('-pub_date')
                blogs = posts_page_divide(request, posts)
                return render(request, "blog_list.html", {'desc': "Tag:" + tag.get().name, 'posts': blogs,
                                                          'visitor_day_count': get_day_data(),
                                                          'blogcategory_list': get_categories(),
                                                          'tag_list': get_tags(), 'recommend_list': get_recommends(), })
            else:
                return HttpResponseNotFound()
        else:
            return HttpResponseNotFound()
    else:
        return HttpResponseNotFound()


def all(request):
    posts = Post.objects.all().order_by('-pub_date')
    blogs = posts_page_divide(request, posts)
    return render(request, 'blog_list.html',
                  {'desc': '全部文章', "posts": blogs, 'visitor_day_count': get_day_data(),
                   'blogcategory_list': get_categories(), 'tag_list': get_tags(),
                   'recommend_list': get_recommends(), })


def change_profile(request):
    if request.user.is_authenticated:
        return render(request, 'change_info.html', )
    else:
        return HttpResponseNotFound()


def username_val(request):
    name = request.POST.get('username')
    if BlogUser.objects.filter(username=name):
        return HttpResponse(json.dumps({'msg': '用户名已被使用', 'status': 0}))
    else:
        return HttpResponse(json.dumps({'msg': '', 'status': 1}))


def new_email_val(request):
    email = request.POST.get('email')
    if email:
        obj = BlogUser.objects.filter(email=email)
        if not obj:
            status = verify_istrue(email)
            if status:
                my_send_email(email, send_type='update_email')
                request.user.is_active = False
                request.user.email = email
                request.user.save()
                logout(request)
                return HttpResponse(json.dumps({'status': 1, 'msg': '发送成功，请尽快完成验证'}))
            else:
                return HttpResponse(json.dumps({'status': 0, 'msg': '发送失败，请检查输入的邮箱是否正确'}))
        else:
            return HttpResponse(json.dumps({'status': 0, 'msg': '该邮箱已被使用'}))


def change_email(request):
    if request.user.is_authenticated:
        return render(request, 'change_email.html',
                      {'visitor_day_count': get_day_data(), 'blogcategory_list': get_categories(), })
    else:
        return HttpResponseNotFound()


def profile_modify(request):
    if request.method == 'POST':
        if request.user.is_authenticated:
            ctx = {}
            ctx['nick_name'] = request.POST.get('nickname')
            ctx['username'] = request.POST.get('username')
            ctx['age'] = request.POST.get('age')
            ctx['cover'] = request.FILES.get('cover')
            ctx['sex'] = request.POST.get('sex')
            ctx['quote'] = request.POST.get('quote')
            ctx['quote_author'] = request.POST.get('quote_author')
            ctx['personal_profile'] = request.POST.get('personal_profile')
            user = request.user
            if ctx['nick_name']:
                user.nick_name = ctx['nick_name']
            if ctx['username']:
                user.username = ctx['username']
            if ctx['age']:
                user.age = int(ctx['age'])
            if ctx['sex']:
                user.sex = ctx['sex']
            if ctx['quote']:
                user.quote = ctx['quote']
            if ctx['quote_author']:
                user.quote_author = ctx['quote_author']
            if ctx['personal_profile']:
                user.personal_profile = ctx['personal_profile']
            if ctx['cover']:
                user.cover = ctx['cover']
            user.save()
            return render(request, 'tips.html',
                          {'visitor_day_count': get_day_data(), 'blogcategory_list': get_categories(),
                           'msg': '修改成功',
                           'location': '/'})

    else:
        return HttpResponseNotFound()


@csrf_exempt
def upload_image(request, dir_name):
    ##################
    #  kindeditor图片上传返回数据格式说明：
    # {"error": 1, "message": "出错信息"}
    # {"error": 0, "url": "图片地址"}
    ##################
    result = {"error": 1, "message": "上传出错"}
    files = request.FILES.get("imgFile", None)
    if files:
        result = image_upload(files, dir_name)
    return HttpResponse(json.dumps(result), content_type="application/json")


# 目录创建
def upload_generation_dir(dir_name):
    today = dt.datetime.today()
    url_part = dir_name + '/%d/%d/' % (today.year, today.month)
    dir_name = os.path.join(dir_name, str(today.year), str(today.month))
    print("*********", os.path.join(settings.MEDIA_ROOT, dir_name))
    if not os.path.exists(os.path.join(settings.MEDIA_ROOT, dir_name)):
        os.makedirs(os.path.join(settings.MEDIA_ROOT, dir_name))
    return dir_name, url_part


# 图片上传
def image_upload(files, dir_name):
    # 允许上传文件类型
    allow_suffix = ['jpg', 'png', 'jpeg', 'gif', 'bmp']
    file_suffix = files.name.split(".")[-1]
    if file_suffix not in allow_suffix:
        return {"error": 1, "message": "图片格式不正确"}
    relative_path_file, url_part = upload_generation_dir(dir_name)
    path = os.path.join(settings.MEDIA_ROOT, relative_path_file)
    print("&&&&path", path)
    if not os.path.exists(path):  # 如果目录不存在创建目录
        os.makedirs(path)
    file_name = str(uuid.uuid1()) + "." + file_suffix
    path_file = os.path.join(path, file_name)
    file_url = settings.MEDIA_URL + url_part + file_name
    open(path_file, 'wb').write(files.file.read())
    return {"error": 0, "url": file_url}


def page_not_found(request, exception):
    return render(request, '404.html',
                  {'visitor_day_count': get_day_data(), 'blogcategory_list': get_categories()})
