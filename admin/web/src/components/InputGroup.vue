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
    <div class="text_group">
        <div class="input_group" :class="{'is-invalid': error}">
            <!-- 输入框 -->
            <input 
                :type="type"
                :placeholder="placeholder"
                :value="value"
                :name="name"
                :disabled = "disableInput"
                @input="$emit('input',$event.target.value)"
            >
            <!-- 输入框后面的内容 -->
            <button v-if="btnTitle" @click="$emit('btnClick')" :disabled="disabled">{{btnTitle}}</button>
        </div>
        <!-- 验证提示 -->
        <div v-if="error" class="invalid-feedback">{{error}}</div>
    </div>
</template>
<script>
    export default {
        name:"inputGroup",
        props:{
            type: {
                type: String,
                default: "text"
            },
            disableInput:Boolean,
            placeholder:String,
            value:String,
            name:String,
            disabled:Boolean,
            btnTitle:String,  //input框中的文字
            error:String  //验证不正确提示
        }
    }
</script>
<style lang="less" scoped>
.is-invalid{

}
.invalid-feedback{
    color: crimson;
}
.input_group input{
    margin-top: 12px;
    height: 28px;
    width: 240px;
}
.input_group button{
    background-color: #fff;
    border: none;
    margin-left: -84px;
    cursor: pointer;
}
</style>