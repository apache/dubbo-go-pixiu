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
    <div style="display: flex;flex-direction: row">
      <img src="@/assets/login_img.png" alt="" class="login-img">
      <div style="flex: 1;background-color: #fff;position: relative;">
          <div style="display: flex;flex-direction: column;justify-content: center;margin: 6.6rem 8rem 0 8rem"
                v-on:keyup.enter="handleLogin">
            <div style="margin-left: 0rem">
                <!-- <div class="login-en-text">BDI-IoT</div> -->
                <div class="login-cn-text">DubbGo Pixiu</div>
                <div class="login-text">欢迎您 ! 请{{isLogin? "登入" : "注册"}}您的账号</div>
            </div>
            <div v-if="isLogin">
                <el-form
                        class="login-form"
                        autoComplete="off"
                        :model="ruleForm"
                        :rules="rules"
                        ref="ruleForm"
                        label-position="left">
                    <el-form-item prop="userName" class="item">
                        <div class="login-id-container">
                            <!-- <img src="@/assets/login_people.svg" alt="" class="login-id-icon"> -->
                            <el-input v-model="ruleForm.username"  type="text" placeholder="请输入您的账号" prefix-icon="el-icon-user"/>
                        </div>
                        
                    </el-form-item>
                    <el-form-item prop="password" class="item">
                        <div class="login-password-container">
                            <!-- <img src="@/assets/login_password.svg" alt="" class="login-password-icon"> -->
                            <el-input
                                type="password"
                                 prefix-icon="el-icon-lock"
                                    placeholder="请输入您的密码"
                                    name="password"
                                    show-password
                                    v-model="ruleForm.password"
                                    autoComplete="off"/>
                        </div>
                        
                    </el-form-item>
                    
                </el-form>
                <div style="display: flex;flex-direction: row;position: relative;margin-top: 4.6rem;justify-content: center;cursor: pointer;">
                    <button @click="handleLogin" class="login-btn-lg">登入</button>
                    
                    <div class="regist">
                        <small>没有账号？</small>
                        <el-button type="text" class="signup" @click="handleRegist">&nbsp;注册</el-button>
                    </div>  
                </div>
            </div>
            <div v-else>
                <el-form :model="registerParam" :rules="rules" ref="registerForm" label-width="0px" class="ms-content login-form">
                    <el-form-item prop="username">
                        <el-input v-model="registerParam.username" onkeyup="value=value.replace(/[^\w\.\/]/ig,'')" autocomplete="new-password" placeholder="用户名" prefix-icon="el-icon-user">
                        </el-input>
                    </el-form-item>
                    <el-form-item prop="password">
                        <el-input type="password"  autocomplete="new-password" placeholder="密码" v-model="registerParam.password" prefix-icon="el-icon-lock">
                        </el-input>
                    </el-form-item>
                    <el-form-item prop="r_password">
                        <el-input type="password"  autocomplete="new-password" placeholder="确认密码" v-model="registerParam.r_password" prefix-icon="el-icon-lock">
                        </el-input>
                    </el-form-item>
                    <!-- <div class="login-btn">
                        <el-button type="primary" @click="submitRegisterForm('registerForm')">注册</el-button>
                    </div>
                    <el-link type="primary" @click="isLogin = true" style="text-align: center;">去登陆 ></el-link> -->
                </el-form>
                <div style="display: flex;flex-direction: row;position: relative;margin-top: 4.6rem;justify-content: center;cursor: pointer;">
                    <button class="login-btn-lg" @click="submitRegisterForm()">注册</button>
                    
                    <div class="regist">
                        <el-button type="text" class="signup" @click="isLogin = true" style="text-align: center;">去登陆 ></el-button>
                    </div>  
                </div>
            </div>
            <div class="login-footer">
          </div>
        </div>
      </div>
    </div>
  </div>
</template>

<script>
    import api from '@/api'
    import {mapActions} from 'vuex';
    import {loading} from "../utils/dialogUtils";
    import {getSecuCode} from "../utils/dialogUtils";
    import { getToken, setToken, removeToken, getLocalStorage, setLocalStorage,clearLocalStorage } from '@/utils/auth'
    import fetch from '@/utils/tafetch';
    export default {
        data() {
            let validatePwd = (rule, value, callback) => {

                if (value === "" || value === undefined) {
                    callback(new Error("请输入密码"));
                } else {
                    callback();
                }
            };
            let validatePass = (rule, value, callback) => {
                if (value === '') {
                    callback(new Error('请再次输入密码'));
                } else if (value !== this.registerParam.password) {
                    callback(new Error('两次输入密码不一致!'));
                } else {
                    callback();
                }
            }
            return {
                isLogin:true,
                ruleForm: {
                    username: "",
                    password: "",
                    checked: true
                },
                registerParam:{
                    username: "",
                    password: "",
                    r_password:""
                },
                rules: {
                    username: [
                        {required: true, message: "请输入登录名", trigger: "blur"}
                    ],
                    password: [{validator: validatePwd, trigger: "blur"}],
                    r_password: [
                        { required: true, message: '请输入确认密码', trigger: 'blur' },
                        { validator: validatePass, trigger: 'blur' }
                    ],
                },
                code: {
                    isVCode: false,//是否显示验证码
                    src: ''
                },
                rand: 0,
                error12: 0,
                loginMsg: '',
                isShowPwd: false, // 是否显示密码
                loading: false, // 登录loading
            };
        },
        methods: {
            ...mapActions([
                'Login'
            ]),
            handleRegist() {
                for(var key in this.registerParam) {
                    this.registerParam[key] = ''
                }
                this.isLogin = false
            },  
            submitRegisterForm() {
                this.$refs["registerForm"].validate(valid => {
                    if (valid) {
                        let formData = new FormData();
                        formData.append('username', this.registerParam.username);
                        formData.append('password', this.registerParam.password);
                         fetch({
                            url: api['register'].url || '',
                            method: 'post',
                            data: formData
                        }).then(res => {
                                console.log(res, '=====>1')
                                this.$msgbox({
                                    message: "注册成功",
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
                    
                        // let formData = new FormData();
                        // formData.append('username', this.registerParam.username);
                        // formData.append('password', this.registerParam.password);
                        // let params = {
                        //     username: this.registerParam.username,
                        //     password: this.registerParam.password,
                        // }
                        // this.$post('/config/api/register', formData)
                        // .then((res) => {
                        //     if (res.code == 10001) {
                        //         console.log(res)
                        //     } else {
                                
                        //     }
                        // }) .catch((err) => {
                        //     console.log(err)
                        // })
                    } else {
                        return false;
                    }
                });
            },
            handleLogin() {
                
                this.$refs["ruleForm"].validate(valid => {
                    if (valid) {
                        let formData = new FormData();
                        formData.append('username', this.ruleForm.username);
                        formData.append('password', this.ruleForm.password);
                        this.Login(formData).then(res => {
                            let _this = this
                            // this.loginMsg = '';
                            console.log(res, "-----111")
                            const getData = getLocalStorage('operatorInfo')
                            console.log(getData, "-----1")
                           if(getData.token && getData.token != ''){
                                console.log('----token')
                                _this.$router.push('/manage');
                           }else{
                                // this.loginMsg = getLocalStorage('operatorInfo').msg;
                                return;
                           }
                        }, error => {
                            console.log(error)
                            this.loginMsg = error.msg;
                        });
                    } else {
                        return false;
                    }
                });
            },
            handleOnRandom() {
                this.rand += 1;
            }
        },
        mounted() {
            // this.init();		//初始化
			// this.animate();	//动画效果
        }
    };
</script>

<style type="text/scss" lang="scss">
    @import "../styles/mixin";
    @import "../styles/common";


    $dark_gray: #889aa4;
    $light_gray: #eee;
    .item{
        margin-bottom: 3rem;
    }
    .login-footer {
        position: absolute;
        left: 0;
        right: 0;
        bottom: 0;
        height: 32px;
        text-align: center;
        color: #143992;
    }
    .login-img{
    width:50vw ;
    height: 100vh;
    margin: -0.4rem 0 0 -0.6rem;
    }
    .login-en-text {
    font-size: 3.84rem;
    font-family: Source Sans Pro;
    font-weight: bold;
    letter-spacing: 1rem;
    color: #143992;
    }

    .login-cn-text {
    font-size: 2rem;
    font-weight: bold;
    letter-spacing: 0.4rem;
    color:#143992;
    text-align: center;
    }

    .login-text {
    font-weight: 400;
    font-size: 1.21rem;
    color: #A6A7AD;
    text-align: center;
    margin: 2.8rem 0 2.8rem 0rem;
    }

    .login-id-container {
    display: flex;
    flex-direction: row;
    border-bottom: 0.16rem solid #E9E9F0;
    
    }

    .login-id-icon {
    width: 1rem;
    height: 1rem;
    margin: 0.36rem 0.6rem 0.2rem 0.6rem
    }

    .login-password-container {
    display: flex;
    flex-direction: row;
    border-bottom: 0.16rem solid #E9E9F0
    }

    .login-password-icon {
    width: 2rem;
    height: 1.8rem;
    margin: 0.16rem 0.1rem 0rem 0.1rem;
    }

    .login-input {
    width: 26rem;
    height: 2rem;
    border: none;
    font-size: 1.08rem;
    color: #43425D;
    display: flex;
    flex-direction: row;
    align-items: center;
    }

    .login-btn-lg {
    width: 12.49rem;
    height: 3.37rem;
    background:#143992;
    color: white;
    border-radius: 1.6rem;
    display: flex;
    justify-content: center;
    align-items: center;
    font-size: 1.21rem;
    cursor: pointer;
    }
    .regist{
        height: 3.37rem;
        line-height: 3.37rem;
        position: absolute;
        right: 0;
        small {
            color: #aaa;
            font-size: 1.11rem;
        }
        .signup {
            font-size: 1.21rem;
            color: #85b4f2;
        }
    }
    .login-btn-rg {
    border: 0.06rem solid #143992;
    border-radius: 1.6rem;
    width: 12.49rem;
    height: 3.37rem;
    display: flex;
    justify-content: center;
    align-items: center;
    font-size: 1.21rem;
    color: #143992;
    margin-left: 3.66rem; 
    background: white;
    }
    .login-container {
        // @include relative;
        margin: 0;
        overflow: hidden;
        background: linear-gradient(to bottom, #19778c, #095f88);
        background-size:1400% 300%;
        animation: dynamics 6s ease infinite;
        -webkit-animation: dynamics 6s ease infinite;
        -moz-animation: dynamics 6s ease infinite;
        font-size: 14px;
        color: #ffffff;
        min-height: 700px;
        height: 100%;

        .login-header {
            width: 100%;
            height: 50%;
            // background: url('../assets/login_bg.png') no-repeat center #fff;
            background-size: 100% 100%;
        }

        .card-box {
            // @include fxied-center;
            text-align: center;
            width: 60%;
            transform: translate(-50%, -57%);

            h1 {
                color: #f5f5f5;
                margin-bottom: 20px;
                text-shadow: 2px 2px 5px #0c0c0c;
            }

            .login-form {
                padding: 51px 48px 48px 48px;
                background: rgba(255, 255, 255, 0.2);
                border-radius: 5px;
                box-shadow: 1px 1px 3px #999;
                width: 35%;
                margin: 0 auto;
            }

            .login_msg {
                background-color: #fef2f2;
                color: #6C6C6C;
                line-height: 16px;
                padding: 6px 10px;
                background: #fef2f2;
                border: 1px solid #ffb4a8;
                margin-bottom: 22px;
                overflow: hidden;
                text-align: left;
                display: flex;

                i {

                    padding-right: 10px;
                    color: #f40;
                }

                p {
                    // @include f_left;
                    white-space: normal;
                    word-wrap: break-word;
                    padding: 0;
                    margin: 0;
                }
            }
        }

        .item {
            .el-form-item__content {
                display: flex;
                flex-flow: row;

                .el-input__inner {
                    color: #999;
                }
            }

        }

        .vCode {
            display: flex;
            flex-flow: row;
            justify-content: space-between;

            .item {
                width: 60%;
            }

            #code {
                display: block;
                color: #ffffff;
                font-size: 20px;
                padding: 5px 35px 10px 35px;
                margin-left: 5%;
                height: 27px;
                cursor: pointer;
            }
        }

        input {
            border: 0;
            -webkit-appearance: none;
            color: $light_gray;
            height: 100%;
        }

        .el-input {
            display: inline-block;
        }

        .tips {
            font-size: 14px;
            color: #fff;
            margin-bottom: 0.13333rem;
        }

        .svg-container {
            padding: 0.08rem 0.0666rem 0.08rem 0.2rem;
            color: $dark_gray;
            vertical-align: middle;
            display: inline-block;

            &_login {
                font-size: 20px;
            }
        }

        .title {
            font-size: 26px;
            color: $light_gray;
            margin: 0 auto 0.5333rem auto;
            text-align: center;
            font-weight: bold;
        }

        .el-form-item {
            border: 1px solid #cbcbcb;
            background: #fff;
            border-radius: 5px;
            color: #9f9f9f;
        }

        .show-pwd {
            position: absolute;
            right: 0.1333rem;
            top: 0.09333rem;
            font-size: 16px;
            color: $dark_gray;
            cursor: pointer;
        }

        .login-footer {
            position: absolute;
            left: 0;
            right: 0;
            bottom: 0;
            height: 32px;
            text-align: center;
            color: #143992;
            // background: url('../assets/login-ft.png') no-repeat center;
        }

        .SoftwareDownload{
            display: flex;
            flex-direction: row;
            justify-content: flex-start;
            .SoftwareDownloadTitle{
                margin-right: 10px;
            }
            .SoftwareDownloadClick{
                flex: 1;
                display: flex;
                flex-direction: row;
                justify-content: flex-start;

            }
            span{
                color: $three-color;
                margin-right: 8px;
                cursor: pointer;
            }
            .line{
                font-weight: 500;
                color: black;
            }
        }
    }
</style>
        