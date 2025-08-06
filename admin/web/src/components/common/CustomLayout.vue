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
  <div class="page">
    <div class="page-head" v-if="fromType">
      <div class="page-head__title">
        <span v-if="!subTitle">{{ title }}</span>
        <el-breadcrumb separator="/" v-else>
          <el-breadcrumb-item>
            <div @click="goBack" class="previous-route">
              {{ title && title.length > 4 ? `${title.substring(0, 4)}...` : title }}
            </div>
          </el-breadcrumb-item>
          <el-breadcrumb-item> {{ subTitle }}</el-breadcrumb-item>
          <el-breadcrumb-item v-if="thirdTitle">{{
            thirdTitle
          }}</el-breadcrumb-item>
        </el-breadcrumb>
      </div>
      <div class="page-head__toolbox">
        <slot name="toolbox"></slot>
      </div>
    </div>
    <el-scrollbar class="page-content" :class="{ 'page-content--have-fixed-bottom': fixedBottom,
      'page-content--fromType' : !fromType }">
      <slot></slot>
      <!-- <div class="page-logo"></div> -->
    </el-scrollbar>
  </div>
</template>

<script>
import { getLocationSearchObj } from '@/utils/common.js'

export default {
  namne: 'CustomLayout',
  props: {
    title: String,
    subTitle: String,
    thirdTitle: String,
    fixedBottom: Boolean
  },
  data () {
    return {
      fromType: false
    }
  },
  methods: {
    init() {
      let params = getLocationSearchObj()
      // this.fromType = params.from == 1 ? false : true // 判断是否隐藏左侧菜单 true是不隐藏  false的隐藏
    },
    goBack () {
      if (this.$router) {
        // this.$router.go(-1)
        this.$router.back()
      }
      this.$emit('back')
    },
    watchScroll () {
      let last_known_scroll_position = 0
      let ticking = false
      let direction = 1
      if (!this.$el.querySelector) {
        return
      }
      let elem = this.$el.querySelector('.page-content')
      let handler = e => {
        if (elem.scrollTop > last_known_scroll_position) {
          direction = 1
        } else {
          direction = -1
        }
        last_known_scroll_position = elem.scrollTop
        if (!ticking) {
          window.requestAnimationFrame(() => {
            this.$emit('scroll', {
              scrollTop: last_known_scroll_position,
              clientHeight: elem.clientHeight,
              scrollHeight: elem.scrollHeight,
              direction: direction
            })
            ticking = false
          })
        }
        ticking = true
      }
      elem.addEventListener('scroll', handler, {
        passive: true
      })
      this.$on('hook:beforeDestroy', () => {
        elem.removeEventListener('scroll', handler, {
          passive: true
        })
      })
    },
    setPaginationStyle () {
      if (!this.$el.querySelector) {
        return
      }
      let elem = this.$el.querySelector('.page-content')
      let handler = options => {
        if (
          options.scrollTop + options.clientHeight >
            options.scrollHeight - 16
        ) {
          let pagination = this.$el.querySelector(
            '.custom-panel--sticky .custom-pagination'
          )
          if (pagination) {
            pagination.classList.add('custom-pagination--no-shadow')
          }
        } else {
          let pagination = this.$el.querySelector(
            '.custom-panel--sticky .custom-pagination'
          )
          if (pagination) {
            pagination.classList.remove('custom-pagination--no-shadow')
          }
        }
      }
      this.$on('scroll', handler)
      this.$parent.$on('hook:updated', () => {
        handler({
          scrollTop: elem.scrollTop,
          clientHeight: elem.clientHeight,
          scrollHeight: elem.scrollHeight
        })
      })
    }
  },
  mounted () {
    this.watchScroll()
    this.setPaginationStyle()
  },
  created() {
    this.init()
  }
}

</script>

<style scoped lang="less">
  .page {
    min-height: calc(100vh - 60px);
    height: calc(100vh - 60px);
  }

  .page-head {
    display: flex;
    background-color: #fff;
    height: 60px;
    padding: 0 24px;
    z-index: 2;
    border-bottom: 1px solid #dcdfe6;
    position: relative;
    box-sizing: border-box;
  }

  .page-head__title {
    color: #303133;
    font-size: 16px;
    line-height: 60px;
  }

  .page-head__toolbox {
    text-align: right;
    flex-grow: 1;
    padding-top: 14px;
  }

  .page-content {
    padding: 12px;
    background-color: #f2f2f2;
    box-sizing: border-box;
    height: calc(100vh - 60px);
    overflow: hidden;
  }
  .page-content--fromType { // 当from为1时  去掉padding和 改变height
    height: calc(100vh - 60px);
  }

  .page-content--have-fixed-bottom {
    /* padding-bottom: 72px; */
    max-height: calc(100vh - 116px);
  }
  .page-content--fromType.page-content--have-fixed-bottom {
    max-height: calc(100vh - 56px) !important;
  }

  .page-head__toolbox i.iconfont {
    font-size: 22px;
    vertical-align: middle;
    cursor: pointer;
  }

  .previous-route {
    font-weight: 700;
    text-decoration: none;
    transition: color 0.2s cubic-bezier(0.645, 0.045, 0.355, 1);
    color: #303133;
    cursor: pointer;
    display: inline-block;
  }

  .previous-route:hover {
    color: #409eff;
  }

  /deep/ .el-scrollbar__wrap {
    overflow-x: hidden;
    overflow-y: auto;
    margin-right: 0 !important;
  }

  /deep/ .el-breadcrumb {
    line-height: 60px;
  }

  /deep/.el-scrollbar__view {
    height: 100%;
  }

  .page-logo {
    height: 52px;
  }
</style>
<style>
.page-content--fromType .fixed-bottom-buttons {
  left: 0;
}
</style>
