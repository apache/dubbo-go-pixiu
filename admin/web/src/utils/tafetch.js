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
/* eslint-disable*/
import axios from 'axios';
import { getToken, setLocalStorage, getLocalStorage, getSessionStorage } from '@/utils/auth';
import store from '../store/store';
import * as dlgUtils from './dialogUtils';
// 创建axios实例
const service = axios.create({
    baseURL: '/login', // api的base_url
    dataType: 'json',
    headers: {"Content-Type": "multipart/form-data; boundary=----WebKitFormBoundary7MA4YWxkTrZu0gW" },
    timeout: 30000 // 请求超时时间
});

// request拦截器
service.interceptors.request.use(config => {
    
    const tk = store.getters.token
    
    if (tk) {
        config.headers.token = tk;
        config.headers.username = getLocalStorage('operatorInfo').username;
    }
    if (config.method.toLowerCase() !== 'get' && Object.keys(config.data).length > 0) {
        let body = {};

        for (let key in config.data) {
            if (key != 'method') {
                body[key] = config.data[key];
            }
        }

        config.data = JSON.stringify(body);
    }

    return config;
}, error => {
    console.log(error); // for debug
    Promise.reject(error);
});

// respone拦截器
service.interceptors.response.use(response => {
    // console.log(response, 'response');
    if (!response.data) {
        return Promise.reject({
            code: '',
            message: '网络异常'
        });
    }
    const res = response.data;
    if (res.hasOwnProperty('code')) {
        // 验签
        // if (!sign(JSON.stringify(res))){
        //     return Promise.reject(
        //         {
        //             code: '',
        //             message: '验签异常'
        //         }
        //     );
        // }


        if (res.code == '10001') {
            console.log(res.data, '-----111')
            return Promise.resolve(res.data);
        } else if (res.code == 503) {
            dlgUtils.loginTimeout();
            return Promise.reject(
                {
                    code: '',
                    message: '网络异常'
                }
            );
        } else {
            return Promise.reject(
                {
                    code: res.code,
                    message: res.message
                }
            );
        }
    } else {
        // size: 3967
        // type: "application/x-download"
        if (res.type) {
            return Promise.resolve(res);
        } else {
            return Promise.reject(
                {
                    code: '',
                    message: '网络异常'
                }
            );
        }

    }
}, error => {
    console.log(error, 'errorerrorerror');
    return Promise.reject({
        code: error.code,
        message: error.message
    });
});
export default service;
