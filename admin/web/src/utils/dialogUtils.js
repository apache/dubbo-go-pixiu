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
import {MessageBox} from "element-ui";
import store from '../store/store'
import {Loading} from 'element-ui';
import api from '@/api'
import fetch from '@/utils/fetch'
import {mapGetters, mapActions} from 'vuex'
import {JSEncrypt} from 'jsencrypt'

export function getSecuCode(password) {
    let encrypt = new JSEncrypt();
    encrypt.setPublicKey('MIGfMA0GCSqGSIb3DQEBAQUAA4GNADCBiQKBgQCRRjobcdpzMY7UgRP4H/7936Ui09qQndfKkzpzmKzo1X30SGI6DDScYF0o8KdP/wK/AY92/V6LL1Fw0dneRd78J9QXT9vS8k3+32a2KP0K8MmBZME9LwORR9IvARpDzTmX5+c/KWjHnPpITHtCtzIodC/8c5kIHKbVcuj44ANRWQIDAQAB');

    let data = encrypt.encrypt(password);
    return data;
}

export function loginTimeout() {
    MessageBox.alert('登录信息已失效，请重新登录', '确定登出', {
        confirmButtonText: '重新登录',
        type: 'warning',
        showCancelButton: false,
        showClose: false
    }).then(() => {
        store.dispatch('FedLogOut').then(() => {
            location.reload() // reload vue-router
        })
    }).catch(e => {
        store.dispatch('FedLogOut').then(() => {
            location.reload() // reload vue-router
        })
    })
}

export const loadding = {
    data() {
        return {
            // loading content
            progressStr: '加载中, 请稍候...',
        }
    },
    methods: {
        ...mapActions([
            'setCustomerInfo',
            'setVehicleInfo'
        ]),
        /** open loading
         * @param progressStr notify content
         * @param timeout default timeout is 10s
         */
        startLoading: function (progressStr = null, timeout = 10000) {
            // console.log(this,this.progressStr)
            // debugger
            if (progressStr != null) {
                this.progressStr = progressStr
            }

            this.loadingObj = Loading.service({
                lock: true,
                text: this.progressStr,
                background: 'rgba(0, 0, 0, 0.7)'
            })
            // auto close after 10 seconds
            setTimeout(() => {
                this.endLoading()
            }, timeout)
        },
        // close loading
        endLoading: function () {
            if (this.loadingObj !== undefined && this.loadingObj !== null) {
                this.loadingObj.close();
            }
        },
        // setting loading content
        setLoadingText: function (progressStr) {
            console.log(this, this.progressStr, progressStr)
            if (this.loadingObj !== undefined && this.loadingObj !== null) {
                this.loadingObj.text = progressStr
            }
        }
    }
}

export const updatePublic = {
    computed: {
        ...mapGetters([
            'vehicleInfo',
            'customerInfo'
        ])
    },
    methods: {
        ...mapActions([
            'setCustomerInfo',
            'setVehicleInfo'
        ]),
        getCarinfo(carInfo) {
            return new Promise((resolve, reject) => {
                fetch({
                    url: api['getCarinfo'].url || '',
                    method: 'post',
                    data: {
                        method: api['getCarinfo'].method,
                        ...carInfo
                    }
                }).then(res => {
                    resolve(res)
                }).catch(error => {
                    reject(error)
                })
            })
        },
        getUserInfo(userInfo, vehicle_code) {
            let _this = this
            return new Promise((resolve, reject) => {
                fetch({
                    url: api['getUserInfo'].url || '',
                    method: 'post',
                    data: {
                        method: api['getUserInfo'].method,
                        ...userInfo
                    }
                }).then(res => {
                    this.setCustomerInfo(res[0]);
                    if (vehicle_code) {
                        // debugger
                        this.getCarinfo({
                            vehicle_code:  _this.vehicleInfo.vehicle_code,
                            vehicle_color:  _this.vehicleInfo.vehicle_color,
                            cpu_card_id: _this.vehicleInfo.cpu_card_id,
                        }).then(carInfo => {
                            if (carInfo.length > 0) {
                                this.setVehicleInfo(carInfo[0]);
                            }
                        }, (error) => {
                            this.$alert(error.message, '提示', {
                                dangerouslyUseHTMLString: true,
                                showClose: false,
                                confirmButtonText: '确定',
                                callback: action => {
                                }
                            });
                        })

                    } else {
                        this.setVehicleInfo({})
                    }

                    resolve(res)
                }).catch(error => {
                    reject(error)
                })
            })
        },
        setPublicInfo(params, cb) {

            if (!(params instanceof Object) || Object.keys(params).length == 0) return false;

            if (params.vehicle_color || params.vehicle_code) {

                this.getCarinfo(params).then(carInfo => {
                    if (carInfo.length > 0) {
                        this.setVehicleInfo(carInfo[0]);
                        cb && cb();
                    }
                    this.getUserInfo({customer_id: carInfo[0].customer_id}).then(res => {
                        if (res.length > 0) {
                            this.setCustomerInfo(res[0]);
                            cb && cb();
                        }
                    }, (error) => {
                        this.$alert(error.message, '提示', {
                            dangerouslyUseHTMLString: true,
                            showClose: false,
                            confirmButtonText: '确定',
                            callback: action => {
                            }
                        });
                    })

                }, (error) => {
                    this.$alert(error.message, '提示', {
                        dangerouslyUseHTMLString: true,
                        showClose: false,
                        confirmButtonText: '确定',
                        callback: action => {
                        }
                    });
                })
            }
        }
    }
}
