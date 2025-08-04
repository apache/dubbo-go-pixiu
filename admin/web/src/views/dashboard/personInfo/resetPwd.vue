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
  <el-form ref="form" :model="user" :rules="rules" label-width="80px">
    <el-form-item label="手机号" class="speInput">
      <el-input v-model="operatorInfo.phone" disabled placeholder="请输入手机号" />
    </el-form-item>
    <el-form-item label="验证码" prop="code" class="speCode">
      <el-input v-model="user.code" placeholder="请输入验证码" />
      <button @click="btnClick" class="ownBtn"  autoComplete="off" :disabled="disabled">{{btnTitle}}</button>
    </el-form-item>
    <el-form-item label="确认密码" prop="newPassword" class="speInput">
      <el-input v-model="user.newPassword"
        :type="passw"
         autoComplete="off"
         placeholder="请确认密码" 
         @blur="onBlur">
        <i slot="suffix" :class="icon" @click="showPass"></i>
      </el-input>
    </el-form-item>
    <el-form-item>
      <el-button type="primary" size="mini" @click="submit">保存</el-button>
      <!-- <el-button type="danger" size="mini" @click="close">关闭</el-button> -->
    </el-form-item>
  </el-form>
</template>

<script>
import api from '@/api'
import fetch from '@/utils/fetch'
import {mapGetters, mapActions} from 'vuex'
import { getToken, setLocalStorage, sign, decrypt } from '@/utils/utils';
export default {
  data() {
    return {
      test: "1test",
      user: {
        newPassword: "",
        code: ""
      },
      //用于显示或隐藏添加修改表单
      add:false,
      //用于改变Input类型
      passw:"password",
      //用于更换Input中的图标
      icon:"el-input__icon el-icon-view",
      adduser:{
          id:null,
          name:null,
          password:null,
          dept_id:null
      },
      phone:"", //手机号
      verifyCode:"", //验证码
      btnTitle:"获取验证码",
      disabled:false,  //是否可点击
      errors:{}, //验证提示信息
      dialogVisible:false,
      // 表单校验
      rules: {
        code: [
          { required: true, message: "验证码不能为空", trigger: "blur" }
        ],
        newPassword: [
          { required: true, message: "新密码不能为空", trigger: "blur" },
          { min: 6, max: 20, message: "长度在 6 到 20 个字符", trigger: "blur" }
        ],
      }
    };
  },
  computed: {
      ...mapGetters([
          'operatorInfo'
      ]),
  },
  methods: {
     //密码失焦事件
    onBlur(){
         this.passw = "password";
         this.icon = "el-input__icon el-icon-view";
    },
    //密码的隐藏和显示
    showPass(){
　　　　　　　　　　//点击图标是密码隐藏或显示
        if( this.passw=="text"){
            this.passw="password";
            //更换图标
            this.icon="el-input__icon el-icon-view";
        }else {
            this.passw="text";
            this.icon="el-input__icon el-icon-open";
        };
    },
    btnClick() {
      console.log(111)
      let params = {
            url: api['operatorSendSms'].url,
            method: 'post',
            data: {
                opId: this.operatorInfo.op_code,
                phone: this.operatorInfo.phone,
                smsType: 1
            }
        }
        fetch(params).then(res => {
            this.$msgbox({
                message: '发送成功',
                title: '成功',
                customClass: 'my_msgBox singelBtn',
                dangerouslyUseHTMLString: true,
                confirmButtonText: '确定',
                type: 'success'
            })
            this.validateBtn()
        }).catch(error => {
            this.$msgbox({
                message:  error.message,
                title: '失败',
                customClass: 'my_msgBox singelBtn',
                dangerouslyUseHTMLString: true,
                confirmButtonText: '确定',
                type: 'error'
            })
        })
    },
    validateBtn(){
      //倒计时
      let time = 60;
      let timer = setInterval(() => {
      if(time == 0) {
        clearInterval(timer);
        this.disabled = false;
        this.btnTitle = "获取验证码";
      } else {
        this.btnTitle =time + '秒后重试';
        this.disabled = true;
        time--
      }
      },1000)
    },
    submit() {
      this.$refs["form"].validate(valid => {
        if (valid) {
          let params = {
              url: api['operatorUpdatePass'].url,
              method: 'post',
              data: {
                  opId: this.operatorInfo.op_code,
                  newPassword: this.user.newPassword,
                  code: this.user.code
              }
          }
          fetch(params).then(res => {
              this.$msgbox({
                  message: '保存数据成功',
                  title: '成功',
                  customClass: 'my_msgBox singelBtn',
                  dangerouslyUseHTMLString: true,
                  confirmButtonText: '确定',
                  type: 'success'
              })
          }).catch(error => {
              this.$msgbox({
                  message:  error.message,
                  title: '失败',
                  customClass: 'my_msgBox singelBtn',
                  dangerouslyUseHTMLString: true,
                  confirmButtonText: '确定',
                  type: 'error'
              })
          })
        }
      });
    },
    close() {
      // this.$store.dispatch("tagsView/delView", this.$route);
      // this.$router.push({ path: "/index" });
    }
  }
};
</script>
<style lang="less" scoped>
/deep/ .el-form-item{
  width: 600px;
}
.speInput /deep/ .el-form-item__content{
  width: 300px;
}
.speCode /deep/ .el-form-item__content{
  width: 300px;
  display: flex;
}
.speCode /deep/ .el-input__inner{
  width: 200px;
}
.ownBtn{
  background-color: #fff;
  border: none;
  width: 120px;
  cursor: pointer;
}
</style>