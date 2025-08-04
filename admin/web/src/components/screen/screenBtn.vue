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
    <div class="cbtn">
        <div class="cbtn-top">
            <div class="cbtn-top__viewTop">
                <div class="txt">机车电压</div>
                <div class="result"><label>??? </label>V</div>
            </div>
             <div class="cbtn-top__viewTop">
                <div class="txt">电机电流</div>
                <div class="result"><label>??? </label>A</div>
            </div>
             <div class="cbtn-top__viewTop">
                <div class="txt">机车速度</div>
                <div class="result"><label>??? </label>M/S</div>
            </div>
             <div class="cbtn-top__viewTop">
                <div class="txt">当前位置</div>
                <div class="result"><label>??? </label>M</div>
            </div>
             <div class="cbtn-top__viewTop">
                <div class="txt">累计里程</div>
                <div class="result"><label>??? </label>M</div>
            </div>
            <div class="cbtn-top__viewTxt">
                停止运行
            </div>
        </div>
        <div class="cbtn-btn">
            <div class="left">
                <div>
                    <div class="first">
                        <el-switch
                        style="display: block"
                        v-model="runByHand"
                        active-color="#13ce66"
                        inactive-color="#ff4949"
                        active-text="自动运行"
                        inactive-text="手动运行"
                        @change="handleSwitch">
                        </el-switch>
                    </div>
                    <div class="second">
                        <el-button round  @click="handleWhistle">鸣笛</el-button>
                        <el-button round  @click="handleSand">撒沙</el-button>
                        <el-button round >复位</el-button>
                        <el-button round >主驾前</el-button>
                        <el-button round >从驾前</el-button>
                        <div style="margin-left: 12px;display:flex;">
                            <div>
                                <div>
                                    <el-button icon="el-icon-plus" circle size="mini"></el-button>
                                <!-- </div>
                                <div> -->
                                    <el-button icon="el-icon-minus" circle size="mini"></el-button>
                                </div>
                            </div>
                            <div>       
                                <label style="margin-left: 12px;">??? M/S</label>
                            </div>
                            
                        </div>
                    </div>
                    <div class="third">
                        <div class="third-left">
                            <el-button plain type="info">任务选择</el-button>
                            <el-button round>主驾前</el-button>
                            <el-button round>从驾前</el-button>
                            <el-button round>任务执行</el-button>
                            <el-button round>任务停止</el-button>
                        </div>
                    </div>
                </div>
                <div>
                    <div class="third-open">
                        <div style="margin-top:10px;">   
                            <el-button round>停止</el-button>
                        </div>
                        <div style="margin-top:10px;">
                            <img src="@/assets/open.png" class="open-img"/>
                        </div>
                        <div style="margin-top:10px;font-size:22px;color:rgb(212, 0, 9);">
                            <label>急停</label>
                        </div>
                    </div> 
                </div>
            </div>
            <div class="right">
                <div class="right-row">
                    <div class="right-row_flex">
                        <span></span>鸣笛
                    </div>
                    <div class="right-row_flex">
                         <span></span>1号抱闸
                    </div>
                </div>
                 <div class="right-row">
                    <div class="right-row_flex">
                        <span></span>前灯
                    </div>
                    <div class="right-row_flex">
                         <span></span>2号抱闸
                    </div>
                </div>
                 <div class="right-row">
                    <div class="right-row_flex">
                        <span></span>后灯
                    </div>
                    <div class="right-row_flex">
                         <span></span>松闸
                    </div>
                </div>
            </div>
        </div>
    </div>
</template>
<script>
import fetch from "@/utils/fetch";
import api from '@/api'
export default {
    name:'screenHeader',
    data(){
        return{
            runByHand: false,
        }
    },
    methods:{
        handleSand() {
            let data = {
                userName: 'cjg'+Math.floor(Math.random() * 1000),  // NOSONAR
                userPwd: '1234',  // NOSONAR
                code:Math.floor(Math.random() * 100),  // NOSONAR
                registered_lon:'1',
                registered_lat:'2',
                name:'tony'+ Math.random()  // NOSONAR
            }
            this.$post('/users/sand', data).then(res => {
                
                this.$message({
                    showClose: true,
                    type: 'success',
                    message: res.msg,
                });
            }).catch(err => {
                console.log(err)
            })
        },
        handleSwitch(e) {
            console.log(e)
        },
        handleWhistle() {
            //重新加载当前车辆信息，存入store
            // 议题发布，新增，编辑
            let data = {
                userName: 'admin',  // NOSONAR
                userPwd: 'admin'  // NOSONAR
            }
            this.$get('/users/whistle', data).then(res => {
                this.$message({
                    showClose: true,
                    type: 'success',
                    message: res.msg,
                });
            }).catch(err => {
                console.log(err)
            })
        }
    }
}
</script>
<style scoped lang="scss">
    /deep/ .el-switch.is-checked .el-switch__core::after{
        margin-left: -27px;
    }
    /deep/ .el-switch{
        height: 30px;
    }  
    /deep/ .el-switch__label span{
        font-size: 22px;
    }
    /deep/ .is-active span{
        font-size: 24px;
    }
    /deep/ .el-switch__core:after{
        width: 26px;
        height: 26px;
    }
    /deep/ .el-switch__core{
        width: 80px !important;
        height: 30px;
        border-radius: 15px;
    }
    .cbtn{
        color: #fff;
        padding: 10px 30px;
        &-top{
            display: flex;
            border:1px solid rgb(161, 161, 88);
            padding: 6px 32px;
            &__viewTop{
                flex: 1;
                text-align: center;
                .txt{
                    font-size: 18px;
                }
                .result{
                    font-size: 18px;
                    color: rgb(33, 255, 255);
                }
            }
            &__viewTxt{
                font-size: 24px;
                color: rgb(194, 0, 6);
                height: 50px;
                line-height: 50px;
            }
        }
        &-btn{
            // padding: 20px 0;
            display: flex;
            /deep/ .el-switch__label--left{
                color: rgb(251, 47, 57) !important;
            }
            /deep/ .el-switch__label--right{
                color: rgb(31, 200, 83);
            }
            .left{
                flex: 1;
                border:1px solid rgb(161, 161, 88);
                padding:0 22px;
                display: flex;
                justify-content: space-between;
                .first{
                    margin-top: 10px;
                }
                .second{
                    margin-top: 10px;
                    display: flex;
                }
                .third{
                    margin-top: 10px;
                    &-left{

                    }
                    
                }
                .third-open{
                    text-align: center;
                    .open-img{
                        cursor: pointer;
                    }
                }
            }
            .right{
                width: 300px;
                height: 300px;
                border:1px solid rgb(161, 161, 88);
                padding: 10px;
                &-row{
                    display: flex;
                    align-items: center;
                    margin-top: 10px;
                    span:first-child{
                        width: 40px;
                        height: 25px;
                        box-sizing: border-box;
                        display: inline-block;
                        border-color: rgb(255, 73, 73);
                        background-color: rgb(255, 73, 73);
                    }
                    &_flex{
                        flex: 1;
                            height: 26px;
                            line-height: 26px;
                            display: flex;
                    }
                }
            }
            
        }
    }
</style>