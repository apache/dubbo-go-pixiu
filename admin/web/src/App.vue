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
  <div id="app">
    <router-view/>
  </div>
</template>


<script>

export default {
  name: "app",
  components: {},
  mounted() {


  }
};
</script>

<style>
#app {
  position:relative;
  width: 100%;
  height: 100%;
  /* margin: 0 auto;
  max-width: 1560px; */
}
::-webkit-scrollbar{/*滚动条整体部分，其中的属性有width,height,background,border等（就和一个块级元素一样）（位置1）*/
  /*滚动条整体样式*/
  width: 6px;
  height: 6px;
  border-radius: 4px;
  /*高宽分别对应横竖滚动条的尺寸*/
/*background: #4494f9;*/
  -webkit-box-shadow: inset 1px 0 1px rgba(0, 0, 0, 0.2);
}
::-webkit-scrollbar-button{/*滚动条两端的按钮，可以用display:none让其不显示，也可以添加背景图片，颜色改变显示效果（位置2）*/
  display: none;
  /*background:#000;*/
  /*height: 0px;*/
}
::-webkit-scrollbar-track{/*外层轨道，可以用display:none让其不显示，也可以添加背景图片，颜色改变显示效果（位置3）*/
  /*滚动条里面轨道*/
  -webkit-box-shadow: inset 0 0 5px rgba(0, 0, 0, 0.2);
  border-radius: 4px;
  background: rgba(0, 0, 0, 0.1) ;
}
/*&::-webkit-scrollbar-track-piece{!*内层轨道，滚动条中间部分（位置4）*!*/
/*    background:#cccccc;*/
/*    border-radius:4px;*/
/*}*/
::-webkit-scrollbar-thumb{/*滚动条里面可以拖动的那部分（位置5）*/
  /*滚动条里面小方块*/
  border-radius: 3px;
  background: #4494f9;
}
</style>
