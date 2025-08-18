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
/* jshint esversion:6 */
class socket {
    constructor(params) {
        this._params = params
        this.init()
    }

    /* 初始化 */
    init() {
        // 重中之重，不然重连的时候会越来越快
        clearInterval(this._reconnect_timer)
        clearInterval(this._heart_timer)

        // 取出所有参数
        const params = this._params
            // 设置连接路径
        const {
            url,
            port
        } = params
        const global_params = ['heartBeat', 'heartMsg', 'reconnect', 'reconnectTime', 'reconnectTimes']

        // 定义全局变量
        Object.keys(params).forEach(key => {
            if (global_params.indexOf(key) !== -1) {
                this[key] = params[key]
            }
        })

        const ws_url = port ? url + ':' + port : url

        this._ws = new WebSocket(ws_url)

        // 默认绑定事件
        this._ws.onopen = () => {
            // 设置状态为开启
            this._alive = true
            clearInterval(this._reconnect_timer)
                // 连接后进入心跳状态
                //this.onheartbeat()
        }

        this._ws.onclose = () => {
            // 设置状态为断开
            this._alive = false

            clearInterval(this._heart_timer)

            // 自动重连开启  +  不在重连状态下
            if (this.reconnect === true) {
                /* 断开后立刻重连 */
                this.onreconnect()
            }
        }
    }

    /* 心跳事件 */
    onheartbeat(func) {
        // 在连接状态下
        setTimeout(() => {
            if (this._alive === true) {
                if (!this.first === true) {
                    this.first = true
                        /* 心跳计时器 */
                    this._heart_timer = setInterval(() => {
                        // 发送心跳信息
                        this.send('{"bissnessType":"doSale","chargeWay":"03","data":"{amount:1}"}')
                        func ? func(this) : false
                    }, this.heartBeat)
                }
            }
        }, 500)
    }

    /* 重连事件 */
    onreconnect(func) {
        /* 重连间隔计时器 */
        this._reconnect_timer = setInterval(() => {
            // 限制重连次数
            if (this.reconnectTimes <= 0) {
                // 关闭定时器
                // this._isReconnect = false
                clearInterval(this._reconnect_timer)
                    // 跳出函数之间的循环
                return
            } else {
                // 重连一次-1
                this.reconnectTimes--
            }
            // 进入初始状态
            this.init()
            func ? func(this) : false
        }, this.reconnectTime)
    }

    // 发送消息
    send(text) {
        if (text === 'undefined') return
        if (this._alive === true) {
            text = typeof text === 'string' ? text : JSON.stringify(text)
            this._ws.send(text)
        }
    }

    /**断开连接 */
    close() {
        if (this._alive === true) {
            // 关闭自动连接
            this.reconnect = false
            this._ws.close()
        }
    }

    /**接受消息 */
    onmessage(func, all = false) {
        this._ws.onmessage = data => func(!all ? data.data : data)
    }

    /**websocket连接成功事件 */
    onopen(func) {
            this._ws.onopen = event => {
                this._alive = true
                func ? func(event) : false
            }
        }
        /**websocket关闭事件 */
    onclose(func) {
            this._ws.onclose = event => {
                // 设置状态为断开
                this._alive = false

                clearInterval(this._heart_timer)

                // 自动重连开启  +  不在重连状态下
                if (this.reconnect === true) {
                    /* 断开后立刻重连 */
                    this.onreconnect()
                }
                func ? func(event) : false
            }
        }
        /**websocket错误事件 */
    onerror(func) {
        this._ws.onerror = event => {
            func ? func(event) : false
        }
    }
}

export default socket