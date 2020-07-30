from django.apps import apps
from django.contrib.admin.sites import AlreadyRegistered
from django.contrib import admin
from blog.models import *


# Register your models here.

# 有自定义显示，实用

@admin.register(BlogUser)
class BlogUserAdmin(admin.ModelAdmin):
    list_per_page = 50
    actions_on_top = True
    actions_on_bottom = True
    actions_selection_counter = True
    empty_value_display = ' -空白- '
    list_editable = ['is_active', 'is_superuser', 'is_staff']
    list_display = (
        'username', 'nick_name', 'age', 'sex', 'email', 'is_active', 'is_staff', 'is_superuser', 'last_login')


@admin.register(Post)
class PostAdmin(admin.ModelAdmin):
    list_per_page = 50
    actions_on_top = True
    actions_on_bottom = True
    actions_selection_counter = True
    empty_value_display = ' -空白- '

    list_display = ('title', 'user', 'content', 'pub_date', 'update_date', 'category', )

    def tags(self, obj):
        return obj.Tags.all().get().name

    class Media:
        js = ('/static/js/editor/kindeditor-all-min.js',
              '/static/js/editor/lang/zh-CN.js',
              '/static/js/editor/config.js')


@admin.register(BlogCategory)
class BlogCategoryAdmin(admin.ModelAdmin):
    list_per_page = 50
    actions_on_top = True
    actions_on_bottom = True
    actions_selection_counter = True
    empty_value_display = ' -空白- '
    list_display = ('name',)


@admin.register(Userip)
class UseripAdmin(admin.ModelAdmin):
    list_per_page = 50
    actions_on_top = True
    actions_on_bottom = True
    actions_selection_counter = True
    empty_value_display = ' -空白- '
    list_display = ('ip', 'count', 'last_visit')


@admin.register(DayNumber)
class DayNumberAdmin(admin.ModelAdmin):
    list_per_page = 50
    actions_on_top = True
    actions_on_bottom = True
    actions_selection_counter = True
    empty_value_display = ' -空白- '
    list_display = ('day', 'count')


@admin.register(VisitNumber)
class VisitNumberAdmin(admin.ModelAdmin):
    list_per_page = 50
    actions_on_top = True
    actions_on_bottom = True
    actions_selection_counter = True
    empty_value_display = ' -空白- '
    list_display = ('count',)


@admin.register(EmailVerifyRecord)
class EmailVerifyRecordAdmin(admin.ModelAdmin):
    list_per_page = 50
    actions_on_top = True
    actions_on_bottom = True
    actions_selection_counter = True
    empty_value_display = ' -空白- '
    list_display = ('email', 'code', 'send_type', 'send_time')


@admin.register(Recommend)
class RecommendAdmin(admin.ModelAdmin):
    list_per_page = 50
    actions_on_top = True
    actions_on_bottom = True
    actions_selection_counter = True
    empty_value_display = ' -空白- '
    list_display = ('title', 'link')


@admin.register(EmailSubscription)
class EmailSubscriptionAdmin(admin.ModelAdmin):
    list_per_page = 50
    actions_on_top = True
    actions_on_bottom = True
    actions_selection_counter = True
    empty_value_display = ' -空白- '
    list_display = ('user', 'email', 'sub_time')


@admin.register(Comment)
class CommentAdmin(admin.ModelAdmin):
    list_per_page = 50
    actions_on_top = True
    actions_on_bottom = True
    actions_selection_counter = True
    empty_value_display = ' -空白- '
    list_display = ('post', 'user', 'content', 'pub_date')


@admin.register(Banner)
class BannerAdmin(admin.ModelAdmin):
    list_per_page = 50
    actions_on_top = True
    actions_on_bottom = True
    actions_selection_counter = True
    empty_value_display = ' -空白- '
    list_display = ('title', 'content', 'is_active', 'category')


# 只注册， 效果不好
app_models = apps.get_app_config("blog").get_models()  # 获取app:crm下所有的model,得到一个生成器
# 遍历注册model
for model in app_models:
    try:
        admin.site.register(model)
    except AlreadyRegistered:
        pass

admin.site.site_header = 'Feng`s blog 后台管理'
admin.site.site_title = 'Feng`s blog'
