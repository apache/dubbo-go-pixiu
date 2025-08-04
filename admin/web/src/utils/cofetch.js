/*
 * Licensed to the Apache Software Foundation (ASF) under one or more
 * contributor license agreements.  See the NOTICE file distributed with
 * this work for additional information regarding copyright ownership.
 * The ASF licenses this file to You under the Apache License, Version 2.0
 * (the "License"); you may not use this file except in compliance with
 * the License.  You may obtain a copy of the License at
 *
 *     http://www.apache.org/licenses/LICENSE-2.0
 *
 * Unless required by applicable law or agreed to in writing, software
 * distributed under the License is distributed on an "AS IS" BASIS,
 * WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
 * See the License for the specific language governing permissions and
 * limitations under the License.
 */
import axios from 'axios'
import {Message} from 'element-ui'
import {getToken, getLocalStorage, setLocalStorage} from '@/utils/auth'
import store from '@/store/store.js'
import * as dlgUtils from './dialogUtils'
import Vue from 'vue'
import ElementUI from 'element-ui'
const SUCCESS_CODE = 200;
import wsConsts from '@/utils/wsConsts'

// 创建axios实例
const service = axios.create({
    baseURL: "/login", // api的base_url
    dataType:"json",
    headers: {"Content-Type": "multipart/form-data; boundary=----WebKitFormBoundary7MA4YWxkTrZu0gW" },
    timeout: 5000 // 请求超时时间
})


// request拦截器
service.interceptors.request.use(config => {
    // console.log(store.getters)
    // if (store.getters.token) {
    //     config.headers.token = getToken();
    // }


    // if (config.method.toLowerCase() !== 'get' && Object.keys(config.data).length > 0) {
    //     let body = {};
    //     let data = {};

    //     for (let key in config.data) {
    //         if (key != 'method') {
    //             body[key] = config.data[key];
    //         }
    //     }
    //     data = body
    //     // if (getToken()) {
    //     //     data.ticket = getToken() || '';
    //     // }

    //     config.data = JSON.stringify(data)
    // }

    return config
}, error => {
    console.log(error) // for debug
    Promise.reject(error)
})
Vue.prototype.$http = service

Vue.prototype.$post = function (url, data) {
  return new Promise((resolve, reject) => {
    service.post(url, data).then(res => {
      if (res.code == 10001) {
        resolve(res)
      } else {
        resolve(res)
      }
    }).catch(err => {
      // ElementUI.Message.error('网络错误，请重试！')
      console.log(err)
      reject(err)
    })
  })
}
Vue.prototype.$put = function (url, data) {
  return new Promise((resolve, reject) => {
    service.put(url, data).then(res => {
      if (res.code == 10001) {
        resolve(res)
      } else {
        resolve(res)
      }
    }).catch(err => {
      // ElementUI.Message.error('网络错误，请重试！')
      console.log(err)
      reject(err)
    })
  })
}
Vue.prototype.$delete = function (url, params = {}) {
  return new Promise((resolve, reject) => {
    service.delete(url, {
      params
    }).then(res => {
      if (res.code == 10001) {
        resolve(res)
      } else {
        resolve(res)
      }
    }).catch(err => {
      // ElementUI.Message.error('网络错误，请重试！')
      console.log(err)
      reject(err)
    })
  })
}

Vue.prototype.$get = function (url, params = {}) {
  return new Promise((resolve, reject) => {
    service.get(url, {
      params
    }).then(res => {
      if (res.code == 10001) {
        resolve(res)
      } else {
        resolve(res)
      }
    }).catch(err => {
      // ElementUI.Message.error('网络错误，请重试！')
      console.log(err)
      reject(err)

    })
  })
}
// respone拦截器
service.interceptors.response.use(response => {
    if (!response.data) {
        return Promise.reject({code: "", message: '网络异常'});
    }
    const res = response.data;
    console.log(res)
    if (!res.hasOwnProperty('code')) {
        return res;
    }

    let expireTime = getLocalStorage('expireTime');
    if(res.code == 503){
        console.log(res)
        dlgUtils.loginTimeout()
        return Promise.reject({code: res.code, message: res.data})
    }else{
        setLocalStorage('expireTime', new Date().getTime() + 1000*60*60*24*7)
        return res;
    }

    return Promise.reject({code: res.code, message: res.msg})

}, error => {
    return Promise.reject({code: error.code, message: error.msg})
})
export default service
