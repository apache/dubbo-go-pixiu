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
<div>
  <el-form ref="form" :model="user" :rules="rules" label-width="80px">
    <el-form-item label="用户名称" prop="username">
      <el-input v-model="user.username" disabled maxlength="12" />
    </el-form-item> 
    <el-form-item label="旧密码" prop="oldPassword" >
      <el-input v-model="user.oldPassword" maxlength="12" />
    </el-form-item>
    <el-form-item label="新密码" prop="newPassword">
      <el-input v-model="user.newPassword" maxlength="12" />
    </el-form-item>
    <el-form-item>
      <el-button type="primary" size="mini" @click="submit">保存</el-button>
      <!-- <el-button type="danger" size="mini" @click="close">关闭</el-button> -->
    </el-form-item>
  </el-form>
  <div>
    <el-dialog
        title="短信验证码"
        :visible.sync="dialogVisible"
        width="300px"
        :before-close="handleClose">
        <div class="login">
            <!-- 手机号 -->
            <InputGroup
                placeholder="手机号"
                v-model="phone" 
                :error="errors.phone"
                :disableInput="myDisable"
            />
            <!-- 输入验证码 -->
            <InputGroup
                v-model="verifyCode"
                placeholder="验证码"
                :error="errors.code"
                :btnTitle="btnTitle"
                :disabled="disabled"
                @btnClick="getVerifyCode"
            />
            <!-- 登录按钮 -->
        　　　　　　　　
        <div class="login_btn">
              <button @click="handleLogin" :disabled="isClick">确 定</button>
        </div>
      </div>
    </el-dialog>
  </div>
</div>
</template>

<script>
import InputGroup from '@/components/InputGroup'
import api from '@/api'
import fetch from '@/utils/fetch'
import {mapGetters, mapActions} from 'vuex'
import { getToken, setToken, removeToken, getLocalStorage, setLocalStorage,clearLocalStorage } from '@/utils/auth'
export default {
  components:{InputGroup},
  props: {
    
  },
  data() {
    return {
      myDisable:true,
      phone:"", //手机号
      verifyCode:"", //验证码
      btnTitle:"获取验证码",
      disabled:false,  //是否可点击
      errors:{}, //验证提示信息
      dialogVisible:false,
      user:{  
        username:"",
        oldPassword:"",
        newPassword:''
      },
      // 表单校验
      rules: {
        username: [
          { required: true, message: "用户名称不能为空", trigger: "blur" }
        ],
        oldPassword: [
          { required: true, message: "旧密码不能为空", trigger: "blur" }
        ],
        newPassword: [
          { required: true, message: "新密码不能为空", trigger: "blur" }
        ]
      }
    };
  },
  created() {
    this.user.username = getLocalStorage('operatorInfo').username
  },
  computed: {
      ...mapGetters([
          'operatorInfo'
      ]),
      
      //手机号和验证码都不能为空
      isClick(){
        if(!this.phone || !this.verifyCode) {
          return true
        } else {
          return false
        }
                  
      }
  },
  methods: {
    ...mapActions([
        'setOperatorInfo',
        'setToken',
    ]),
    submit() {
      this.$refs["form"].validate(valid => {
        if (valid) {
          let formData = new FormData();
          formData.append('username', this.user.username);
          formData.append('oldPassword', this.user.oldPassword);
          formData.append('newPassword', this.user.newPassword);
            fetch({
              url: api['getPasswordEdit'].url || '',
              method: 'post',
              data: formData
          }).then(res => {
                  console.log(res, '=====>1')
                  this.$msgbox({
                      message: "密码修改成功",
                      title: '成功',
                      customClass: 'my_msgBox singelBtn',
                      dangerouslyUseHTMLString: true,
                      confirmButtonText: '确定',
                      type: 'success'
                  })
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
        }
      });
    },
    getVerifyCode(){
    //获取验证码
      if(this.validatePhone()) {
         let params = {
            url: api['operatorSendSms'].url,
            method: 'post',
            data: {
                opId: this.operatorInfo.op_code,
                phone: this.phone,
                smsType: 3
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
      }
    },
    handleLogin() {
        let params = {
            url: api['operatorUpdateData'].url,
            method: 'post',
            data: {
                code: this.verifyCode,
                email: this.user.email,
                opId: this.operatorInfo.op_code,
                phone: this.phone,
                sex: this.user.sex,
                address: "",
            }
        }
        fetch(params).then(res => {
            this.$msgbox({
                message: '数据保存成功',
                title: '成功',
                customClass: 'my_msgBox singelBtn',
                dangerouslyUseHTMLString: true,
                confirmButtonText: '确定',
                type: 'success'
            })
            let data = this.operatorInfo;
            Object.assign(data, this.user)
            this.setOperatorInfo(data)
            this.dialogVisible = false
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
    handleClose() {
      this.dialogVisible = false
    }, 
    validatePhone(){
      //判断输入的手机号是否合法
      if(!this.phone) {
        this.errors = {
        phone:"手机号码不能为空"
      }
        // return false
      } else if(!/^1[345678]\d{9}$/.test(this.phone)) {
        this.errors = {
        phone:"请输入正确的手机号"
      }
        // return false
      } else {
        this.errors ={}
        return true
      }
    },  
    close() {
    //   this.$store.dispatch("tagsView/delView", this.$route);
    //   this.$router.push({ path: "/index" });
    }
  }
};
</script>
<style lang="less" scoped>
/deep/ .el-form-item__content{
  width: 300px;
}
.input_group input{
    margin-top: 20px;
    height: 28px;
    width: 240px;
}
/deep/ .el-dialog__body{
  padding: 10px 20px 30px ;
}
.login_btn button {
  margin-top: 20px;
  width: 248px;
  background-color:dodgerblue;
  height: 32px;
  color: #fff;
  border-radius: 4px;
  border: none;
  cursor: pointer;
}
</style>