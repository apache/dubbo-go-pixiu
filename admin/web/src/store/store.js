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
import Vue from 'vue'
import Vuex from 'vuex'
import { getToken, setToken, removeToken, getLocalStorage, setLocalStorage,clearLocalStorage } from '@/utils/auth'
import fetch from '@/utils/tafetch'
import api from '@/api'
Vue.use(Vuex)

export default new Vuex.Store({
  state: {
    id: '123456',
    openedTab: ['index'],
    activeTab: '',
    token:getToken(),
    operatorInfo:getLocalStorage('info') || {},//系统操作员
    customerInfo:{},//ETC客户
    modelType:'',//当前用户正在操作的模块
    cancelVehicle: false,//当前用户是否从新增车辆中返回的
    vehicleInfo:{},//车辆信息
    changeCardInfo:{},//补领换卡功能中的新卡信息以及公务卡信息
    appointInfo:{},//预约信息
    application_order_no:'',//申领单号
  },
  getters:{
    token: (state) => state.token,
    operatorInfo: (state) => state.operatorInfo,
    customerInfo: (state) => state.customerInfo,
    modelType: (state) => state.modelType,
    cancelVehicle: (state) => state.cancelVehicle,
    vehicleInfo:(state) => state.vehicleInfo,
    changeCardInfo:(state) => state.changeCardInfo,
    appointInfo:(state) => state.appointInfo,
    application_order_no:(state) => state.application_order_no,
  },
  mutations: {
    addTab (state, componentName) {
      state.openedTab.push(componentName)
    },
    changeTab (state, componentName) {
      state.activeTab = componentName
    },
    deductTab (state, componentName) {
      let index = state.openedTab.indexOf(componentName)
      state.openedTab.splice(index, 1)
    },
    SET_TOKEN: (state, token) => {
      state.token = token
    },
    SET_APPLICATION_ORDER_NO: (state, application_order_no) => {
      state.application_order_no = application_order_no
    },
    SET_OPERATORINFO: (state, operatorInfo) => {
      state.operatorInfo = operatorInfo
    },
    SET_CUSTOMERINFO: (state, customerInfo) => {
      state.customerInfo = customerInfo
    },
    SET_APPOINTINFO: (state, appointInfo) => {
      state.appointInfo = appointInfo
    },
    SET_CHANGECARDINFO: (state, changeCardInfo) => {
      state.changeCardInfo = changeCardInfo
    },
    SET_MODEL_TYPE: (state, modelType) => {
      state.modelType = modelType
    },
    SET_CANCEL_VEHICLE: (state, cancelVehicle) => {
      state.cancelVehicle = cancelVehicle;
    },
    SET_VEHICLEINFO: (state, vehicleInfo) => {
      console.log(vehicleInfo, '-------');
      state.vehicleInfo = vehicleInfo
    },
  },
  actions: {
    Login({ commit }, userInfo) {
      console.log(userInfo)
      return new Promise((resolve, reject) => {
        
        fetch({
          url: api['login'].url || '',
          method: 'post',
          data: userInfo
          
        }).then(res => {
          console.log(res, "======>21")
          const data = res
            // console.log(data.data.token)
          setToken(data.token)
          commit('SET_TOKEN', data.token);
          setLocalStorage('expireTime', new Date().getTime() + 1000*60*60*24*7)
          setLocalStorage('operatorInfo',data);
          commit('SET_OPERATORINFO',data);
          
          
          resolve()
        }).catch(error => {
          setLocalStorage('expireTime', 0)
          reject(error)
        })
      })
    },
    // 设置ETC客户
    setCustomerInfo({ commit },data) {
      commit('SET_CUSTOMERINFO',data)
    },
    setAppointInfo({ commit },data) {
      commit('SET_APPOINTINFO',data)
    },
    setModelType({ commit },data) {
      commit('SET_MODEL_TYPE',data)
    },
    setCancelVehicle({ commit },data) {
      commit('SET_CANCEL_VEHICLE',data)
    },
    setChangeCardInfo({ commit },data) {
      commit('SET_CHANGECARDINFO',data)
      console.log(this.getters.changeCardInfo, 'changeCardInfo')
    },
    setVehicleInfo({ commit },data){
      commit('SET_VEHICLEINFO',data);
    },
    setApplicationOrderNo({ commit },data){
      commit('SET_APPLICATION_ORDER_NO',data);
    },

    // 登出
    LogOut({ commit, state }) {
      return new Promise((resolve, reject) => {
        logout(state.ticket).then(() => {

          commit('SET_TOKEN', '')
          commit('SET_OPERATORINFO', {})
          commit('SET_CUSTOMERINFO',{})
          commit('SET_VEHICLEINFO',{})
          commit('SET_APPOINTINFO',{})

          removeToken();
          clearLocalStorage('operatorInfo');
          clearLocalStorage('expireTime');
          resolve()
        }).catch(error => {
          reject(error)
        })
      })
    },
    // 前端 登出
    FedLogOut({ commit }) {
      return new Promise(resolve => {
        commit('SET_TOKEN', '')
        commit('SET_OPERATORINFO', {})
        commit('SET_CUSTOMERINFO',{})
        commit('SET_VEHICLEINFO',{})
        commit('SET_APPOINTINFO',{})

        removeToken();
        clearLocalStorage('operatorInfo');
        clearLocalStorage('expireTime');

        resolve()
      })
    }

  }
})
