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
const config = {
    'base': {
        url: '/api/base',
        method: ''
    },
    'whistle':{
        url: '/users/whistle',
        method: ''
    },
    'resetPwd':{
        url: '/operator/resetPassword',
        method: ''
    },
    'findHisList':{
        url: '/electronicAccountHis/findHisList',
        method: ''
    },
    'cancleSign':{
        url: '/electronicAccountHis/cancleSign',
        method: ''
    },
    'allChannel':{
        url: '/findAllChannel',
        method: ''
    },
    'closeChannel':{
        url: '/closeChannel',
        method: ''
    },
    'openChannel':{
        url: '/openChannel',
        method: ''
    },
    'findCustomer':{
        url: '/redis/findCustomerCar',
        method: ''
    },
    'findRedis':{
        url: '/redis/findRedisValue',
        method: ''
    },
    'removeCustomer':{
        url: '/redis/removeCustomerCar',
        method: ''
    },
    'findCardPub':{
        url: '/cardpub/findCardPub',
        method: ''
    },
    'cardLost':{
        url: '/cardpub/cardLost',
        method: ''
    },
    'cardUnLost':{
        url: '/cardpub/cardUnLost',
        method: ''
    },
    'register':{
        url: '/register',
        method: ''
    },
    'login':{
        url: '/login',
        method: ''
    },
    'getInfo':{
        url: '/user/getInfo',
        method: ''
    },
    'getPasswordEdit':{
        url: '/user/password/edit',
        method: ''
    },
}
export default config;
