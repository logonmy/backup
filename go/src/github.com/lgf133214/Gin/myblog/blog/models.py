from django.db import models
from django.contrib.auth.models import AbstractUser
from django.utils import timezone


# Create your models here.


# 访问网站的ip地址和次数
class Userip(models.Model):
    ip = models.CharField(verbose_name='IP地址', max_length=30)  # ip地址
    count = models.IntegerField(verbose_name='访问次数', default=0)  # 该ip访问次数
    last_visit = models.DateTimeField(verbose_name="上次访问时间", default=timezone.now)

    class Meta:
        verbose_name = '访问用户信息'
        verbose_name_plural = verbose_name

    def __str__(self):
        return self.ip


# 网站总访问次数
class VisitNumber(models.Model):
    count = models.IntegerField(verbose_name='网站访问总次数', default=0)  # 网站访问总次数

    class Meta:
        verbose_name = '网站访问总次数'
        verbose_name_plural = verbose_name

    def __str__(self):
        return str(self.count)


# 单日访问量统计
class DayNumber(models.Model):
    day = models.DateField(verbose_name='日期', default=timezone.now)
    count = models.IntegerField(verbose_name='网站访问次数', default=0)  # 网站访问总次数
    day_visit_ip = models.ManyToManyField(Userip, verbose_name="ip地址")

    class Meta:
        verbose_name = '网站日访问量统计'
        verbose_name_plural = verbose_name

    def __str__(self):
        return str(self.day)


class BlogUser(AbstractUser):
    nick_name = models.CharField('昵称', max_length=20, default='游客', blank=True)
    age = models.IntegerField('年龄', default=0)
    sex = models.CharField('性别', default='保密', max_length=3)
    cover = models.ImageField('头像', upload_to='images/user_icon/',
                              default='/static/images/user_icon_default.ico')
    views = models.IntegerField('热度', default=0)
    quote = models.CharField('格言', max_length=150, default='', blank=True)
    quote_author = models.CharField('格言出处', max_length=20, default='', blank=True)
    personal_profile = models.CharField('个人简介', max_length=500, default='', blank=True)

    class Meta:
        verbose_name = '用户'
        verbose_name_plural = verbose_name
        ordering = ['-id']

    def __str__(self):
        return self.username


class EmailSubscription(models.Model):
    email = models.EmailField(max_length=50, verbose_name="订阅者邮箱")
    sub_time = models.DateTimeField(verbose_name="订阅时间", default=timezone.now)
    user = models.ForeignKey(BlogUser, verbose_name='用户', on_delete=models.CASCADE, default=None)
    is_active = models.BooleanField("是否active", default=False)

    class Meta:
        verbose_name = '订阅者邮箱'
        verbose_name_plural = '订阅者邮箱'

    def __str__(self):
        return self.email


class EmailVerifyRecord(models.Model):
    code = models.CharField(verbose_name='验证码', max_length=50, default='')
    email = models.EmailField(max_length=50, verbose_name="邮箱")
    send_type = models.CharField(verbose_name="验证码类型",
                                 choices=(("register", "注册"), ("forget", "找回密码"), ("update_email", "修改邮箱")),
                                 max_length=30)
    send_time = models.DateTimeField(verbose_name="发送时间", default=timezone.now)

    class Meta:
        verbose_name = "邮箱验证码"
        # 复数
        verbose_name_plural = verbose_name

    def __str__(self):
        return '{0}({1})'.format(self.code, self.email)


class BlogCategory(models.Model):
    name = models.CharField('分类名称', max_length=20, default='')

    class Meta:
        verbose_name = '博客分类'
        verbose_name_plural = '博客分类'

    def __str__(self):
        return self.name


class Banner(models.Model):
    title = models.CharField('标题', max_length=50)
    cover = models.ImageField('轮播图', upload_to='images/banner',
                              default='/static/images/banner_default.png')
    link_url = models.CharField('图片链接', max_length=100)
    content = models.CharField('概述', max_length=200)
    is_active = models.BooleanField('是否active', default=False)
    category = models.ForeignKey(BlogCategory, verbose_name='分类', default=None, on_delete=models.DO_NOTHING, blank=True)

    def __str__(self):
        return self.title

    class Meta:
        verbose_name = '轮播图'
        verbose_name_plural = '轮播图'


class Tags(models.Model):
    name = models.CharField('标签名称', max_length=20, default='')

    class Meta:
        verbose_name = '标签'
        verbose_name_plural = '标签'

    def __str__(self):
        return self.name


class Post(models.Model):
    user = models.ForeignKey(BlogUser, verbose_name='用户', on_delete=models.CASCADE, default=None)
    category = models.ForeignKey(BlogCategory, verbose_name='博客分类', default=None, on_delete=models.DO_NOTHING)
    tags = models.ManyToManyField(Tags, verbose_name='标签')
    title = models.CharField('标题', max_length=50)
    content = models.TextField('内容')
    pub_date = models.DateTimeField('发布日期', default=timezone.now)
    update_date = models.DateTimeField('修改时间', default=timezone.now)
    cover = models.ImageField('博客封面', upload_to='images/post-cover', default='/static/images/post_cover_default.png')
    views = models.ManyToManyField(Userip, verbose_name='浏览数', blank=True)

    def __str__(self):
        return self.title

    class Meta:
        verbose_name = '博客'
        verbose_name_plural = '博客'


class Comment(models.Model):
    post = models.ForeignKey(Post, verbose_name='博客', on_delete=models.CASCADE)
    user = models.ForeignKey(BlogUser, verbose_name='用户', on_delete=models.CASCADE, default=None)
    pub_date = models.DateTimeField('发布时间', default=timezone.now)
    content = models.TextField('内容')

    def __str__(self):
        return self.content

    class Meta:
        verbose_name = '评论'
        verbose_name_plural = '评论'


class Recommend(models.Model):
    title = models.CharField('标题', max_length=50)
    link = models.CharField('链接', max_length=50, default='')

    def __str__(self):
        return self.title

    class Meta:
        verbose_name = '推荐'
        verbose_name_plural = '推荐'
