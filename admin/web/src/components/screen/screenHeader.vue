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
    <div class="cheader">
        <div class="cheader-left">
            <label></label>
        </div>
        <div class="cheader-center">
            
        </div>
        <div class="cheader-right">
            <!-- <div class="cheader-right__user">admin</div>
            <div class="cheader-right__day">3天</div> -->
            <div class="cheader-right__time">{{currentTime}}</div>
        </div>
    </div>
</template>
<script>
export default {
    name:'screenHeader',
    data(){
        return{
            currentTime:''
        }
    },
    created() {
        this.timeFormate(new Date()) 
        // 定时刷新时间
        this.timer = setInterval(()=> {
            this.timeFormate(new Date()) // 修改数据date
        }, 100)
    },
    methods:{
        timeFormate(date) {
            var year = date.getFullYear();
            var month = date.getMonth()+1;
            month=month<10?"0"+month:month;
            var day = date.getDate();
            day=day<10?"0"+day:day;
            var week = date.getDay();
            week="日一二三四五六".charAt(week);
            week="星期"+week;
            var hour = date.getHours();
            hour=hour<10?"0"+hour:hour;
            var minute = date.getMinutes();
            minute=minute<10?"0"+minute:minute;
            var second = date.getSeconds();
            second=second<10?"0"+second:second;
            var result = year+"-"+month+"-"+day+" "+week+" "+hour+":"+minute+":"+second;
            this.currentTime = result;
		},
        
    },
    destroyed() {
        if (this.timer) { // 注意在vue实例销毁前，清除我们的定时器
            clearInterval(this.timer);
        }
    }
}
</script>
<style scoped lang="scss">
    .cheader{
        display: flex;
        height: 48px;
        line-height: 48px;
        padding: 0 32px;
        &-left{
            color: #fff;
            font-size: 26px;
            label{
                color: rgb(16, 131, 254);
            }
            width: 300px;
        }
        &-center{
            color: rgb(26, 192, 255);
            font-size: 32px;
            flex: 1;
            text-align: center;
        }
        &-right{
            width: 300px;
            display: flex;
            color: rgb(25, 174, 14);
            &__user{

            }
            &__day{

            }
            &__time{
                font-size: 20px;
            }
        }
    }
</style>