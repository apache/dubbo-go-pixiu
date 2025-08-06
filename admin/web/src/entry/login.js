/*
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
 */
import Vue from 'vue'
import App from '@/App.vue'
import router from '@/router/router'
import store from '@/store/store'
import '@/plugins/element.js'
import '@/styles/icons/iconfont.css'
import '@/styles/common.scss'
import '@/styles/normalize.css'
import '@/permission'
import '@/utils/directives'
import {loadding,updatePublic} from '@/utils/dialogUtils'
import '@/element-variables.scss'
import ElementUI from 'element-ui'
Vue.config.productionTip = false
import Router from 'vue-router'
 
const originalPush = Router.prototype.push
Router.prototype.push = function push(location) {
  return originalPush.call(this, location).catch(err => err)
}
Vue.use(ElementUI, {
  size: 'small'
})
Vue.mixin(loadding)
Vue.mixin(updatePublic)
new Vue({
  router,
  store,
  render: h => h(App)
}).$mount('#app')
