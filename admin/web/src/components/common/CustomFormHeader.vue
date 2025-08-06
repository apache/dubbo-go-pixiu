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
<div class="title">
    <div class="line--vertical"></div>
    <div class="name">{{ title }}</div>
    <div class="search">
        <slot name="searchContent"></slot>
    </div>
</div>
</template>

<script>
export default {
    name: 'CustomFormHeader',
    props: ['title']
}
</script>

<style scoped>
.title {
    /* height: 46px; */
    /* line-height: 46px; */
    color: rgba(48, 49, 51, 1);
    margin-bottom: 20px;
    height: 40px;
    background: rgba(248, 248, 248, 1);
}

.name {
    font-size: 15px;
    font-family: PingFangSC-Regular, PingFang SC;
    font-weight: 400;
    line-height: 40px;
}

.title div {
    float: left;
}

.line--vertical {
    width: 3px;
    height: 14px;
    background: rgba(64, 158, 255, 1);
    margin: 13px 7px 13px 12px;
}

.search {
    float: right !important;
    margin-top: 4px;
    height: 32px;
    padding-right: 16px;
}
</style>
