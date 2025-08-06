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
import router from '@/router/router'
import store from '@/store/store'
import { Message } from 'element-ui'
import { getToken,getLocalStorage } from '@/utils/auth'
import NProgress from 'nprogress' // Progress 进度条
import 'nprogress/nprogress.css'// Progress 进度条样式

const whiteList = ['/login'] // 不重定向白名单
router.beforeEach(( to, from, next) => {
    NProgress.start()
    // next()
    if(getToken()){
        if(to.path == '/login'){
            next({path:'/'})
            NProgress.done();
        }else{
            console.log(Object.keys(getLocalStorage('operatorInfo')))
            if(Object.keys(getLocalStorage('operatorInfo')).length > 0){
                next()
            }else{
                store.dispatch('FedLogOut').then(() => {
                    next({ path: '/login' })
                })
            }
        }
    }else{
        if (whiteList.indexOf(to.path) !== -1) {
            next()
        } else {
            next('/login')
            NProgress.done()
        }
    }

    router.afterEach(() => {
        NProgress.done() // 结束Progress
    })
})