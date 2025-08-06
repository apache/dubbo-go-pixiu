<!--
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
 -->
<template>
  <div class="app-container">
    <el-row :gutter="20">
      <el-col :span="6" :xs="24">
        <el-card class="box-card">
          <div slot="header" class="clearfix">
            <span>个人信息</span>
          </div>
          <div>
            <div class="text-center">
              <userAvatar :user="user" />
            </div>
            <ul class="list-group list-group-striped">
              <li class="list-group-item">
                用户名称
                <div class="pull-right">{{ userInfo.username }}</div>
              </li>
              <li class="list-group-item">
                用户ID
                <div class="pull-right">{{ userInfo.userId }}</div>
              </li>
              <!-- <li class="list-group-item">
                手机号码
                <div class="pull-right">{{ operatorInfo.phone }}</div>
              </li>
              <li class="list-group-item">
                用户邮箱
                <div class="pull-right">{{ operatorInfo.email }}</div>
              </li>
              <li class="list-group-item">
                所属部门
                <div class="pull-right" v-if="operatorInfo.branch_name">{{ operatorInfo.branch_name }} </div>
              </li> -->
              <li class="list-group-item">
                所属角色
                <div class="pull-right">{{userInfo.role == 1 ? 'admin' : 'admin'}}</div> 
              </li>
              <!-- <li class="list-group-item">
                <svg-icon icon-class="date" />创建日期
                <div class="pull-right">{{ user.createTime }}</div>
              </li> -->
            </ul>
          </div>
        </el-card>
      </el-col>
      <el-col :span="18" :xs="24">
        <el-card>
          <div slot="header" class="clearfix">
            <span>基本资料</span>
          </div>
          <el-tabs v-model="activeTab" :before-leave="beforeLeave">
            <el-tab-pane label="修改密码" name="userinfo">
              <userInfo :user="user" />
            </el-tab-pane>
            <!-- <el-tab-pane label="修改密码" name="resetPwd">
              <resetPwd :user="user" />
            </el-tab-pane> -->
          </el-tabs>
        </el-card>
      </el-col>
    </el-row>
  </div>
</template>

<script>
import userAvatar from "./userAvatar";
import userInfo from "./userInfo";
import resetPwd from "./resetPwd";
import {mapGetters, mapActions} from 'vuex'
import api from '@/api'
import { getToken, setToken, removeToken, getLocalStorage, setLocalStorage,clearLocalStorage } from '@/utils/auth'
import fetch from '@/utils/tafetch';
export default {
  name: "Profile",
  components: { userAvatar, userInfo, resetPwd },
  data() {
    return {
      user: {},
      roleGroup: {},
      postGroup: {},
      activeTab: "userinfo",
      roleMap:new Map([
          [1, "本部"],
          [2, "片区"],
          [3, "网点"],
      ]),
      role:'',
      userInfo:{}
    };
  },
    computed: {
        ...mapGetters([
            'operatorInfo'
        ])
    },
    watch:{
      
    },
  created() {
    this.getUser();
    const role = this.operatorInfo.department
    this.role = this.roleMap.get(role)
  },
  methods: {
    getUser() {
      let formData = new FormData();
        fetch({
          url: api['getInfo'].url || '',
          method: 'post',
          data: {}
      }).then(res => {
          console.log(res, '=====>1221')
          this.userInfo = res    
      }).catch(error => {
          console.log(error);
          this.$msgbox({
              message: error.message,
              title: '失败',
              customClass: 'my_msgBox singelBtn',
              dangerouslyUseHTMLString: true,
              confirmButtonText: '确定',
              type: 'error'
          })
      });
    },
    beforeLeave() {
      if(this.operatorInfo.phone == '' || !this.operatorInfo.phone){
        this.$message({
          type:"warning",
          message:"请先保存基本资料里的手机号码，不得为空！"
        })
        return false
      }
    },
  }
};
</script>
<style type="text/less" lang="less" scoped>
  .app-container{
    padding: 10px;
  }
</style>
<style type="text/less" lang="less">

.el-row{
  margin-left: 0 !important;
  margin-right: 0 !important;
}
.text-center {
  text-align: center
}
.list-group-striped > .list-group-item {
	border-left: 0;
	border-right: 0;
	border-radius: 0;
	padding-left: 0;
	padding-right: 0;
}

.list-group {
	padding-left: 0px;
	list-style: none;
}

.list-group-item {
	border-bottom: 1px solid #e7eaec;
	border-top: 1px solid #e7eaec;
	margin-bottom: -1px;
	padding: 11px 0px;
	font-size: 13px;
}
.pull-right {
	float: right !important;
}
</style>