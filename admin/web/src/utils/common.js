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

// 获取地址栏url参数并对象化
export function getLocationSearchObj () {
  var search = location.search.substring(1)
  if (!search && location.href.lastIndexOf('?') > -1) {
    search = location.href.substring(location.href.lastIndexOf('?') + 1)
  }
  var obj = {}
  if (search.length > 0) {
    var arr = [], item
    arr = search.split('&')
    for (var i = arr.length; --i >= 0;) {
      item = arr[i].split('=')
      obj[item[0]] = item[1]
    }
  }
  return obj
}

// 校验身份证
export function judgeIdCard (val) {
  let bool = false
  let reg = /(^[1-9]\d{5}(18|19|([23]\d))\d{2}((0[1-9])|(10|11|12))(([0-2][1-9])|10|20|30|31)\d{3}[0-9Xx]$)|(^[1-9]\d{5}\d{2}((0[1-9])|(10|11|12))(([0-2][1-9])|10|20|30|31)\d{3}$)/
  if (reg.test(val)) {
    bool = true
  }
  return bool
}

/* 
 * 只允许输入数字和. 用于处理一些金额 xxxx.xx
 * @param {Number | String} val 需要限制的内容
 * @param {Number} max 整数部分最大值 比如四位数：9999
 * @param {Number} num 是否有小数点，如果有小数点则为小数的最大数 比如两位小数: 99
 */
export function numberLimit (val, max, num) {
  val = val.toString()
  val = val.replace(/。/g, '.')
  // 没有小数
  if (!num) {
    val = val.replace(/[^0-9]/g, '')
    val = val > max ? max : val * 1
  } else { // 有小数时限制
    val = val.replace(/[^0-9.]/g, '')
    let list = val.split('.')
    let a = list[0]
    if (list.length > 1) {
      let b = list[1]
      a = a > max ? max : a * 1
      b = b > num ? num : b.toString().slice(0, 2)
      if (a >= max) {
        b = '00'
      }
      val = (a || 0) + '.' + b
    } else {
      a = a > max ? max : (a === '' ? a : a * 1)
      val = a
    }
  }
  return val
}

// 经纬度转换成三角函数中度分表形式。
function rad (d){
  return d * Math.PI / 180.0
}

// 参考地址https://blog.csdn.net/zzjiadw/article/details/7031610
/* 
 * 计算点球上两点距离
 * @params {Number} lat1 第一个点纬度
 * @params {Number} lng1 第一个点纬度
 * @params {Number} lat2 第二个点纬度
 * @params {Number} lng2 第二个点纬度
 * @return 返回两点之间的距离，单位km
*/
export function getDistance (lat1, lng1, lat2, lng2) {
  var radLat1 = rad(lat1)
  var radLat2 = rad(lat2)
  var a = radLat1 - radLat2
  var b = rad(lng1) - rad(lng2)
  var s = 2 * Math.asin(Math.sqrt(Math.pow(Math.sin(a / 2), 2) +
  Math.cos(radLat1) * Math.cos(radLat2) * Math.pow(Math.sin(b / 2), 2)))
  s = s * 6378.137 // 地球半径 6378.137
  s = Math.round(s * 10000) / 10000 //输出为公里
  // s=s.toFixed(4)
  return s
}

/*
 * 将数字进行逗号拼接,每3位加一个逗号，支持小数(小数不做逗号处理) 123456.32 => 123,456.32
 * @params {Number | String} 需要处理的数值
 * @return 返回已拼接的数值
*/
export function numJoint (val) {
  if (val) {
    val = val.toString()
    // 是否包含小数
    let idx = val.indexOf('.')
    let point = ''
    if (idx > 0) {
      point = val.slice(idx)
      val = val.slice(0, idx)
      val = val.replace(/(\d)(?=(?:\d{3})+$)/g, '$1,')  // NOSONAR
      val = val + point
    } else {
      val = val.replace(/(\d)(?=(?:\d{3})+$)/g, '$1,')  // NOSONAR
    }
  }
  return val
}
