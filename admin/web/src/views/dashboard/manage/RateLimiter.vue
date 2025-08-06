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
    <CustomLayout >
    <div class="custom-body">
      <div>
        <CommonTitle title="概览"></CommonTitle>
      </div>
      <div class="custom-tools">
        <div class="table-head">
          <div class="custom-tools__info">限流配置</div>
          <div>
              <el-button type="primary" icon="el-icon-plus" size="mini"
                     @click="handleSave">保存</el-button>
              <el-button type="primary" icon="el-icon-plus" size="mini"
                     @click="handleDelete">删除</el-button>
          </div>
         
        </div>
        
        <div class="custom-tools__content">
          <el-form :model="form"
                 :inline="true"
                 @submit.native.prevent=""
                 class="table-form bg-gray"
                 label-width="130px">
          <el-row>
            <div style="clear: 'both'; height: 300px;width: 100%;" id="containForm" ref="containForm"/>
          </el-row>
        </el-form>
        </div>
      </div>
    </div>
  </CustomLayout>
</template>

<script>
import CommonTitle from '@/components/common/CommonTitle'
import CustomLayout from '@/components/common/CustomLayout.vue'
import * as monaco from 'monaco-editor/esm/vs/editor/editor.main.js'
import 'monaco-editor/esm/vs/basic-languages/javascript/javascript.contribution'
import { StandaloneCodeEditorServiceImpl } from 'monaco-editor/esm/vs/editor/standalone/browser/standaloneCodeServiceImpl.js'

export default {
  name: '',
  components:{
    CommonTitle,
    CustomLayout
  },
  data () {
    return {
      create: true,
      form:{

      },
      edit:false,
      dialogVisible:false,
      chooseRow:{
        timeout:'',
        description:'',
        type:'',
        path:''
      },
      loading: false,
      tableHeight: window.document.documentElement.clientHeight - 248, // table高度
      multipleSelection: [], // 多选时，已选项
      currentRow: {}, // 当前选中行
      areaDialogVisible: false, // 是否显示行政区划弹窗
      emptText: '暂无数据', // table数据为空的提示语
      isSizeInit: false, // 是否是页面初始化时设置的PageSize
      previewDialogVisible: false, // 是否显示行政区划弹窗
      currentTemplate: {},//当前选中的模板
      monacoEditor: null,
      monacoEditored: null,
      detailSource:null,
    }
  },
  mounted() {
    //获取基础信息
    this.getRateLimiterDetail()
  },
  methods: {
    initMoacoEditored(language, value) {
      this.monacoEditored = monaco.editor.create(document.getElementById('containForm'), {
        value,
        language: 'yaml',
        codeLens: true,
        selectOnLineNumbers: true,
        roundedSelection: false,
        readOnly: false,
        lineNumbersMinChars: true,
        theme: 'vs-dark',
        wordWrapColumn: 120,
        folding: false,
        showFoldingControls: 'always',
        wordWrap: 'wordWrapColumn',
        cursorStyle: 'line',
        automaticLayout: true,
      });
    },
    //映射服务列表
    getRateLimiterDetail() {
        this.$get('/config/api/plugin/ratelimit')
        .then((res) => {
            if (res.code == 10001) {
                this.create = false            
                this.$nextTick(() =>[
                    this.initMoacoEditored('yaml', res.data)
                ])
            } else {
                this.create = true            
                this.$nextTick(() =>[
                    this.initMoacoEditored('yaml', '暂未配置')
                ])
            }
        })
        .catch((err) => {
            console.log(err)
        })
    },
    handleSelectionChange(val) {
      this.multipleSelection = val
    },
    //修改用户信息
    handleSave() {
      let formData = new FormData();
      let data = this.monacoEditored.getValue()
      formData.append('content', data);

      if (this.create == true) {
        this.hanldeCreate(formData)
      } else {
        this.handleModify(formData)
      }
    },
    hanldeCreate(data) {
      this.$post('config/api/plugin/ratelimit/', data)
        .then((res) => {
          if (res.code == 10001) {
            this.$message({
              type: 'success',
              message: '创建成功！',
            })
            this.monacoEditored.dispose()
            this.getRateLimiterDetail()
          } 
        })
        .catch((err) => {
          console.log(err)
        })
    },
    handleModify(data) {
      this.$put('config/api/plugin/ratelimit/', data)
        .then((res) => {
          if (res.code == 10001) {
            this.$message({
              type: 'success',
              message: '修改成功！',
            })
            this.monacoEditored.dispose()
            this.getRateLimiterDetail()
          } 
        })
        .catch((err) => {
          console.log(err)
        })
    },
    handleDelete() {
      this.$delete('config/api/plugin/ratelimit/')
        .then((res) => {
          if (res.code == 10001) {
            this.$message({
              type: 'success',
              message: '删除成功！',
            })
            this.monacoEditored.dispose()
            this.getRateLimiterDetail()
          } 
        })
        .catch((err) => {
          console.log(err)
        })
    },
    
  },
  destroyed(){
    // this.monacoEditored.dispose();
  }
}
</script>


<style scoped lang="less">
.custom-panel{
  margin-top: 20px;
}
.dltbtn{
  color: red;
}
.custom-tools__info{
  color: rgba(16, 16, 16, 100);
  font-size: 18px;
  text-align: left;
  margin-top: 10px;
}
.custom-tools__content{
  background-color: #fff;
  margin-top: 10px;
  padding: 10px 20px;
}
.table-head{
  display: flex;
  margin-top: 10px;
  justify-content: space-between;
}
.dialog_main {
    display: flex;
    flex-wrap: wrap;

    .dialog_item {
        width: 100%;
        display: flex;
        align-items: center;
        margin-bottom: 10px;
    }

    .item_label {
        min-width: 60px;
    }
}
</style>
