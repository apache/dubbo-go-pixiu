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
        <CommonTitle title="映射服务"></CommonTitle>
      </div>
      <div class="custom-tools">
        <div class="table-head">
          <div class="custom-tools__info">Resource信息</div>
          <el-button type="primary" icon="el-icon-plus" size="mini"
                     @click="handleChange">修改</el-button>
        </div>
        
        <div class="custom-tools__content">
          <el-form :model="form"
                 :inline="true"
                 @submit.native.prevent=""
                 class="table-form bg-gray"
                 label-width="130px">
          <el-row>
            <div style="clear: 'both'; height: 100px;width: 100%;" id="containForm" ref="containForm"/>
          </el-row>
        </el-form>
        </div>
      </div>

      <div class="table-head">
        <div class="custom-tools__info">方法映射</div>
        <el-button type="primary" icon="el-icon-plus" size="mini"
                     @click="handleAdd">新增</el-button>
      </div>
      <div class="custom-panel">
        <el-table v-loading="loading"
                  size="normal"
                  :data="tableData"
                  :empty-text="emptText"
                  row-class-name="custom-table-tr--hover"
                  header-row-class-name="custom-table-header"
                  @selection-change="handleSelectionChange"
                  style="width: 100%">
          <el-table-column 
                           prop="id"
                           label="ID">
          </el-table-column>
          <el-table-column 
                           class-name="custom-popper--overflow"
                           prop="httpVerb"
                           label="方法">
            <template slot-scope="scope">
                <span></span>
            </template>
          </el-table-column>
          <el-table-column class-name="custom-popper--overflow"
                           label="timeout">
            <template slot-scope="scope">
                <span>{{scope.row.timeout}}</span>
            </template>
          </el-table-column>
          <el-table-column 
                           prop="resourcePath"
                           label="路径">
          </el-table-column>
          <el-table-column
                           prop="type"
                           label="类型">
          </el-table-column>
          <el-table-column
                           prop="description"
                           label="描述">
          </el-table-column>
          <el-table-column
                           label="操作">
            <template slot-scope="scope">

              <el-button type="text"
                         @click="handleLook(scope.row)">查看</el-button>
              <el-button type="text" class="dltbtn"
                         @click="deleteRow(scope.row)">删除</el-button>
            </template>
          </el-table-column>
        </el-table>

        <div class="custom-pagination"
             style="float:right">
          
          <div style="display: flex;">
            <div>
              总共<span>{{ pagination.total }}</span>条记录<i class="custom-pagination__interval"></i>每页显示
              <el-select v-model="pagination.pageSize"
                         placeholder="请选择"
                         style="width:70px;">
                <el-option v-for="item in pagination.pageSizes"
                           :key="item"
                           :label="item"
                           :value="item">
                </el-option>
              </el-select>
              条记录
            </div>
            <el-pagination :current-page.sync="pagination.pageIndex"
                           :page-size.sync="pagination.pageSize"
                           layout="prev, pager, next"
                           :total="pagination.total">
            </el-pagination>
          </div>
        </div>
      </div>
    </div>
            <el-dialog
                title="新增"
                :visible.sync="createDialogVisible"
                width="640px"
                :before-close="handleClose">
            <div class="dialog_main">
                <div style="clear: 'both'; height: 300px;width: 100%;" id="createContainer" ref="createContainer"/>
            </div>
            <span slot="footer" class="dialog-footer">
                <el-button size="mini" @click="handleClose">取 消</el-button>
                <el-button type="primary"  size="mini" @click="makeCreate">确 定</el-button>
            </span>
        </el-dialog>

            <el-dialog
                title="查看修改"
                :visible.sync="updateDialogVisible"
                width="640px"
                :before-close="handleClose">
            <div class="dialog_main">
                <div style="clear: 'both'; height: 300px;width: 100%;" id="updateContainer" ref="updateContainer"/>
            </div>
            <span slot="footer" class="dialog-footer">
                <el-button size="mini" @click="handleClose">取 消</el-button>
                <el-button type="primary"  size="mini" @click="makeModify">确 定</el-button>
            </span>
        </el-dialog>
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
      edit:false,
      createDialogVisible:false,
      updateDialogVisible:false,
      form:{

      },
      chooseRow:{
        timeout:'',
        description:'',
        type:'',
        path:''
      },
      pagination: { // 分页字段
        pageIndex: 1, // 当前页
        pageSize: 20, // 每页显示多少条
        total: 0, // 总条数
        pageSizes: [20, 50, 100]
      },
      tableData: [], // table列表数据
      loading: false,
      tableHeight: window.document.documentElement.clientHeight - 248, // table高度
      multipleSelection: [], // 多选时，已选项
      currentRow: {}, // 当前选中行
      areaDialogVisible: false, // 是否显示行政区划弹窗
      emptText: '暂无数据', // table数据为空的提示语
      isSizeInit: false, // 是否是页面初始化时设置的PageSize
      previewDialogVisible: false, // 是否显示行政区划弹窗
      currentTemplate: {},//当前选中的模板
      updateDetailSource:null,
      createDetailSource:null,
      createMonacoEditor: null,
      updateMonacoEditor: null,
      monacoEditored: null,
    }
  },
  mounted(){
    console.log(this.$route)
    this.getMethodList()
    this.getResourceDetail()
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
    initCreateMoacoEditor(language, value) {
      this.createMonacoEditor = monaco.editor.create(document.getElementById('createContainer'), {
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
    initUpdateMoacoEditor(language, value) {
      this.updateMonacoEditor = monaco.editor.create(document.getElementById('updateContainer'), {
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
    handleSelectionChange(val) {
      this.multipleSelection = val
    },
    handleAdd() {
      this.areaDialogVisible = true
    },
    getMethodDeatil() {
      this.$get('/config/api/resource/method/detail', {
        resourceId: 1,
        methodId: 2
      })
        .then((res) => {
          if (res.code == 10001) {
              let data = JSON.parse(res.data)
              console.log(data)
          } else {
            
          }
        })
        .catch((err) => {
          console.log(err)
        })
    },
    //删除
    deleteRow(row) {
      this.$delete('/config/api/resource/method', {
        resourceId: this.$route.query.resourceId,
        methodId: row.id
      })
        .then((res) => {
          if (res.code == 10001) {
            this.$message({
              type: 'success',
              message: '删除成功！',
            })
            this.getMethodList()
            console.log(res)
          } else {
            
          }
        })
        .catch((err) => {
          console.log(err)
        })
    },
     //映射服务列表
    getMethodList() {
      this.$get('/config/api/resource/method/list', {
        resourceId: this.$route.query.resourceId
      })
        .then((res) => {
          if (res.code == 10001) {
             this.tableData = JSON.parse(res.data)
             console.log(this.tableData)
          } else {
            
          }
        })
        .catch((err) => {
          console.log(err)
        })
    },
    getResourceDetail() {
      this.$get('/config/api/resource/detail', {
        resourceId: this.$route.query.resourceId
      })
        .then((res) => {
          if (res.code == 10001) {
            this.$nextTick(() =>[
              this.initMoacoEditored('yaml', res.data)
            ])
          } else {
            this.$nextTick(() =>[
              this.initMoacoEditored('yaml', '')
            ])
          }
        })
        .catch((err) => {
          console.log(err)
        })
    },
    //修改resource detail 信息
    handleChange() {
      let formData = new FormData();
      let data = this.monacoEditored.getValue()
      formData.append('content', data);
      this.$put('/config/api/resource?resourceId=' + this.$route.query.resourceId, formData)
        .then((res) => {
          if (res.code == 10001) {
            this.$message({
              type: 'success',
              message: '修改成功！',
            })
            this.monacoEditored.dispose()
            this.getResourceDetail()
          } 
        })
        .catch((err) => {
          console.log(err)
        })
    },
    makeCreate() {
      let formData = new FormData();
      let data = this.createMonacoEditor.getValue()
      formData.append('content', data);
      
      this.$post('/config/api/resource/method?resourceId=' + this.$route.query.resourceId, formData)
        .then((res) => {
          if (res.code == 10001) {
            this.handleClose()
            this.getMethodList()
            console.log(res)
          } else {
            
          }
        })
        .catch((err) => {
          console.log(err)
        })
    },
    makeModify() {
      let formData = new FormData();
      let data = this.updateMonacoEditor.getValue()
      formData.append('content', data);
      
      this.$put('/config/api/resource/method?resourceId=' + this.$route.query.resourceId, formData)
        .then((res) => {
          if (res.code == 10001) {
            this.handleClose()
            this.getMethodList()
            console.log(res)
          } else {
            
          }
        })
        .catch((err) => {
          console.log(err)
        })
    },
    handleClose() {
      this.createDialogVisible = false
      this.updateDialogVisible = false
      if (this.createMonacoEditor != null) {
        this.createMonacoEditor.dispose();
      }
      if (this.updateMonacoEditor != null) {
        this.updateMonacoEditor.dispose();
      }
    },
    handleLook(row) {

      this.edit = true
      this.$get('/config/api/resource/method/detail', {
        resourceId: this.$route.query.resourceId,
        methodId: row.id,
      })
        .then((res) => {
          if (res.code == 10001) {
              let data = res.data
              this.modifyDetailSource = data
              this.updateDialogVisible = true
              this.$nextTick(() =>[
                this.initUpdateMoacoEditor('yaml', data)
              ])
          } else {
            
          }
        })
        .catch((err) => {
          console.log(err)
        })
    },
    //新增
    handleAdd() {
      this.createDialogVisible = true
      this.edit = false
      this.$nextTick(() =>[
        this.initCreateMoacoEditor('yaml', '')
      ])
    },
  }
}
</script>


<style scoped lang="less">
.dltbtn{
  color: red;
}
.custom-panel{
  margin-top: 20px;
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
</style>
