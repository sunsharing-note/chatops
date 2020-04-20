#!/usr/bin/env python
# -*- coding:utf-8 -*- 
# __author__ = 'bigc'
# @Time    : 2020/4/12 21:36
# @Email   : luocs1@lenovo.com
# views部分代码

from django.http import HttpResponse, JsonResponse


# Create your views here.
def ping(request):

    if request.method == "GET":
        return HttpResponse("pong")
