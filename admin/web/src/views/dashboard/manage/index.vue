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
  <div class="manage">
    <div class="manage-head hd_bar">
        <div class="manage-head__name">Dubbo Go Pixiu</div>
        <div class="right-menu el-menu-demo">
            <el-dropdown class="avatar-container right-menu-item hover-effect" trigger="click">
                <div class="avatar-wrapper">
                    <div class="name">{{operatorInfo.username}}</div>
                <!-- <img :src="avatar" class="user-avatar"> -->
                <i class="el-icon-caret-bottom" />
                </div>
                <el-dropdown-menu slot="dropdown">
                    <el-dropdown-item>
                        <span>你好 - {{operatorInfo.username}}</span>
                    </el-dropdown-item>
                    <!-- <el-dropdown-item divided @click.native="handleOnOpen">
                        <span>发行系统</span>
                    </el-dropdown-item> -->
                    <el-dropdown-item divided @click.native="handlePersonInfo">
                        <span>个人中心</span>
                    </el-dropdown-item>
                    <el-dropdown-item divided @click.native="handleOnLogOut">
                        <span>退出登录</span>
                    </el-dropdown-item>
                </el-dropdown-menu>
            </el-dropdown>
        </div>
    </div>
    <div class="manage-foot">
        <div class="menu-container" :class="{'menu-container--hide': !fromType}">
            <custom-menu></custom-menu>
        </div>
        <div class="lyt-container">
            <router-view></router-view>
        </div>
    </div>
    
  </div>
</template>

<script>
import CustomMenu from '@/components/CustomMenu.vue'
import { getLocationSearchObj } from '@/utils/common.js'
import { getToken, setToken, removeToken, getLocalStorage, setLocalStorage,clearLocalStorage } from '@/utils/auth'
import {mapGetters, mapActions} from 'vuex'
export default {
  name: 'Manage',
  data() {
    return {
      fromType: true,
      operatorInfo:{}
    }
  },
  components: {
    CustomMenu
  },
  computed: {
      // ...mapGetters([
      //     'operatorInfo'
      // ])
  },
  methods: {
    ...mapActions([
        'FedLogOut',
    ]),
    init() {
      this.operatorInfo = getLocalStorage('operatorInfo')
      let params = getLocationSearchObj()
      this.fromType = params.from == 1 ? false : true // 判断是否隐藏左侧菜单 true是不隐藏 false的隐藏
    },
    handlePersonInfo() {
        this.$router.push({
            path:'/personInfo'
        })
    },
    handleOnLogOut() {
        console.log(1);
        this.FedLogOut().then(res => {
            this.$router.push({path: '/login'});
        });
    },
    handleOnOpen() {
        let port = '30026';
        const key = getSessionStorage('tskey')
        window.open('http://' + document.location.hostname + '/#/login?key=' + key, '_blank')  // NOSONAR
    },
  },
  created() {
    this.init()
  }
}
</script>

<style lang="less" scoped>
.manage {
  height: 100vh;
  &-head {
    display: flex;
    height: 60px;
    justify-content: space-between;
    align-items: center;
    padding:0 60px;
    background-color: rgba(68, 126, 217, 100);
    &__name{
        font-size: 28px;
        font-weight: 600;
        color: #FFF;
    }
    &__info{
        
    }
  }
}
.manage-foot{
    display: flex;
    height: calc(100vh - 60px);
}
.menu-container {
  flex-shrink: 0;
  width: 150px;
  height: 100%;
  background: #FFF;
  border-right: 1px solid #dcdfe6;
  box-sizing: border-box;
}
.menu-container--hide {
  width: 0;
}
.lyt-layout {
  flex: 1;
  display: flex;
  box-sizing: border-box;
}
.lyt-container {
  flex: 1;
  height: 100%;
  overflow-y: auto;
  overflow-x: hidden;
}
.lyt-container--inherit {
  max-width: inherit;
}
.spe{
    margin-right: 20px;
    cursor: pointer;
}
.right-menu {
    float: right;
    height: 100%;
    line-height: 50px;

    &:focus {
      outline: none;
    }

    .right-menu-item {
      display: inline-block;
      height: 100%;
      font-size: 14px;
      color: #5a5e66;
      vertical-align: text-bottom;

      &.hover-effect {
        cursor: pointer;
        transition: background .3s;

        &:hover {
          background: rgba(0, 0, 0, .025)
        }
      }
    }
}
/deep/ .avatar-wrapper{
    display: flex;
    align-items: center;
    color: #fff;
    font-size: 18px;
}
</style>
