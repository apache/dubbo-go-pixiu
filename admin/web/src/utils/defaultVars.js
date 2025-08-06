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
import store from '@/store/store.js'

const defaultVars = {
    store:store,

    vehicleInfo: function () {
        return {
            //	客户编号
            customer_id: store.getters.customerInfo.customer_id,
            //	车辆号码
            vehicle_code: '浙',
            //	车牌颜色（0蓝 1黄 2黑 3白 4渐变绿 5黄绿双拼 6蓝白渐变）
            vehicle_color: '0',
            //	车型（车型 0客车,1货车）
            vehicle_type: '0',
            //	车辆用户类型(0 普通车 1集卡车 2卧铺车 8军车(交通战备车) 9警车 15紧急车 16特殊公务车)，集卡车只有车型选货车的时候显示
            vehicle_user_type: '0',
            //	座位数
            vehicle_seat: '',
            //	吨数
            vehicle_ton: '0',
            //长
            vehicle_length: '0',
            //	宽
            vehicle_width: '0',
            //	高
            vehicle_height: '0',
            //	车轮数
            vehicle_wheels: '0',
            //	车轴数
            vehicle_axles: '0',
            //	轴距
            vehicle_wheelbases: '0',
            //	车籍地（对应车籍地附件  车籍地信息city.xls）
            vehicle_city: '0100',
            //	发动机号
            vehicle_engine: '',
            //	车辆特征描述
            vehicle_specificInfo: '',
            //	车辆识别代码 VIN
            vehicle_distinguish: '',
            //	ETC通行卡功能（1记账2储值）
            card_type: '1',
            //	预约编号（如果通过预约填写的，传预约编号）
            reservation_id: null,
            //	是否发行OBU（当车型是集卡车时会传该字段。0不发行OBU；1发行OBU；）
            issue_obu: '1',
        }
    }
}

export default defaultVars