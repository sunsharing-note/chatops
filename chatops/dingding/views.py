#!/usr/bin/env python
# -*- coding:utf-8 -*- 
# __author__ = 'bigc'
# @Time    : 2020/4/12 21:44
# @Email   : luocs1@lenovo.com
# views部分代码

from django.http import HttpResponse, JsonResponse
import json,re
from chatops.settings import *
import hmac
import hashlib
import base64,requests


# 机器人的app_secret
app_secret = "fkyymqRK_Jr-nmKX4tk9ogMn9OmT87HeDmeJC6Gt30JeRwLQcR94PEf_UwX0EiKi"
#app_secret="D_3RLRTqq76GGePa7XlVWVuaO3fnIlRbkpSCAhrnpTnrIK81Vp1vvBP5-65zWN0E"
zabbix_url=ZABBIX_URL
zabbix_user=ZABBIX_USER
zabbix_password=ZABBIX_PASSWORD
headers = {"Content-Type": "application/json"}
# Create your views here.
def auth():
    data = {"jsonrpc": "2.0",
            "method": "user.login",
            "params": {
                "user": zabbix_user,
                "password": zabbix_password
            },
            "id": 1,
            "auth": None
            }
    auth = requests.post(url=zabbix_url, headers=headers, json=data)
    auth = json.loads(auth.content)['result']
    return auth

def logout(auth):
    data = {
        "jsonrpc": "2.0",
        "method": "user.logout",
        "params": [],
        "id": 1,
        "auth": auth
    }
    logout = requests.post(url=zabbix_url, headers=headers, json=data)
    logout=json.loads(logout.content)['result']
    return logout
def gethost(auth):
    #print(content)
    data={
        "jsonrpc": "2.0",
        "method": "host.get",
        "params": {
            "output":"extend"
        },
        "auth": auth,
        "id": 1
    }

    gethost=requests.post(url=zabbix_url,headers=headers,json=data)
    gethost=json.loads(gethost.content)['result']
    return gethost
def getalert(auth):
    data={
        "jsonrpc": "2.0",
        "method": "event.get",
        "params": {
            "output": "extend",
            "select_acknowledges": "extend",
            "selectTags": "extend",
            "sortfield": ["clock", "eventid"],
            "sortorder": "DESC"
        },
        "auth": auth,
        "id": 1
    }
    getalert=requests.post(url=zabbix_url,headers=headers,json=data)
    getalert=json.loads(getalert.content)['result']
    return getalert

login=auth()
gethost=gethost(login)
getalert=getalert(login)
logout=logout(login)

def robot(request):
    if request.method == "POST":
        HTTP_SIGN = request.META.get("HTTP_SIGN")
        HTTP_TIMESTAMP = request.META.get("HTTP_TIMESTAMP")
        res = json.loads(request.body)
        print(res)
        print(zabbix_url)
        # 用户输入钉钉的信息
        content = res.get("text").get("content")
        print(content)
        string_to_sign = '{}\n{}'.format(HTTP_TIMESTAMP, app_secret)
        string_to_sign_enc = string_to_sign.encode('utf-8')
        hmac_code = hmac.new(app_secret.encode("utf-8"), string_to_sign_enc, digestmod=hashlib.sha256).digest()
        sign = base64.b64encode(hmac_code).decode("utf-8")
        print(sign)
        print(HTTP_SIGN)
        # 验证签名是否为钉钉服务器发来的

        if sign == HTTP_SIGN:
            if "test" in content:
                return JsonResponse(
                    {"msgtype": "text",
                     "text": {
                         "content": "test success"
                     }
                     }
                )
            if "登录zabbix" in content:

                return JsonResponse(
                    {"msgtype": "text",
                     "text": {
                         "content": login
                     }
                     }
                )
            if "注销zabbix" in content:
                return JsonResponse(
                    {"msgtype": "text",
                     "text": {
                         "content": logout
                     }
                     }
                )
            if "获取zabbix主机" in content:
                return JsonResponse(
                    {"msgtype": "text",
                     "text": {
                         "content": gethost
                     }
                     }
                )
            if "获取告警事件" in content:
                return JsonResponse(
                    {"msgtype": "text",
                     "text": {
                         "content": getalert
                     }
                     }
                )
            return JsonResponse(
                {"msgtype": "text",
                 "text": {
                     "content": "谢谢使用此机器人，{}".format(content)
                 }
                 }
            )
        return JsonResponse({"error": "你没有权限访问此接口"})
    if request.method == "GET":
        return HttpResponse("hello")
