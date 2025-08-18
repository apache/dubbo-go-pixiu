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
    this.init()
  },
  methods: {
    async init() {
      await this.getData()
      this.getNowData()
    },
    async getData() {
      setTimeout(function(){
        console.log(111)
      }, 5000)
      clearTimeout()
    },
    getNowData() {
      console.log(222)
    }
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
</style>
