# Generated by Django 2.2.6 on 2019-10-28 16:55

from django.db import migrations, models


class Migration(migrations.Migration):

    dependencies = [
        ('blog', '0003_auto_20191025_2234'),
    ]

    operations = [
        migrations.AlterField(
            model_name='banner',
            name='cover',
            field=models.ImageField(default='/static/images/banner_default.png', upload_to='images/banner', verbose_name='轮播图'),
        ),
    ]
