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
  <CustomLayout>
    <div class="custom-body">
      <div>
        <CommonTitle title="OPA 策略配置"></CommonTitle>
      </div>
      <div class="custom-tools">
        <div class="table-head">
          <div class="custom-tools__info">OPA 策略</div>
          <div>
            <el-button size="mini" @click="handleSync">同步</el-button>
            <el-button type="primary" size="mini" @click="handleSave">保存</el-button>
            <el-popconfirm
              title="确定要删除该策略吗？"
              @confirm="handleDelete">
              <el-button slot="reference" type="danger" size="mini">删除</el-button>
            </el-popconfirm>
          </div>
        </div>
        <div class="custom-tools__content">
          <el-form :model="form"
                   :inline="true"
                   @submit.native.prevent=""
                   class="table-form bg-gray"
                   label-width="130px">
            <el-row>
              <el-form-item label="policy_id">
                <el-input v-model="form.policy_id"
                          clearable
                          placeholder="pixiu-authz"></el-input>
              </el-form-item>

            </el-row>
            <el-row>
              <div style="clear: both; height: 300px;width: 100%;" id="policyEditor" ref="policyEditor"/>
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
import { getLocalStorage, setLocalStorage } from '@/utils/auth'

const POLICY_CACHE_KEY = 'opaPolicyLast'
const DEFAULT_POLICY_ID = 'pixiu-authz'

const DEFAULT_POLICY = `package pixiu.authz

default allow := false

allow if {
  # write your logic here
}
`
const CONTROL_CHAR_REGEX = /[\u0000-\u0008\u000B\u000C\u000E-\u001F\u007F\u2028\u2029\uFEFF]/g
let regoRegistered = false

function normalizePolicyText(value) {
  if (!value) {
    return ''
  }
  return value
    .replace(/\r\n/g, '\n')
    .replace(/\r/g, '\n')
    .replace(CONTROL_CHAR_REGEX, '')
}

function registerRegoLanguage(monaco) {
  if (regoRegistered) {
    return
  }
  regoRegistered = true
  monaco.languages.register({ id: 'rego' })
  monaco.languages.setMonarchTokensProvider('rego', {
    keywords: [
      'package', 'default', 'import', 'as', 'with', 'not', 'some',
      'else', 'if', 'in', 'true', 'false', 'null'
    ],
    operators: [
      '=', ':=', '==', '!=', '<', '>', '<=', '>=', '+', '-', '*', '/',
      '%', 'and', 'or'
    ],
    tokenizer: {
      root: [
        [/[a-zA-Z_][\w\-]*/, {
          cases: {
            '@keywords': 'keyword',
            '@default': 'identifier'
          }
        }],
        [/[{}()[\]]/, '@brackets'],
        [/[=><!]=?/, 'operator'],
        [/\"([^\"\\]|\\.)*$/, 'string.invalid'],
        [/\"/, { token: 'string.quote', bracket: '@open', next: '@string' }],
        [/#[^\r\n]*/, 'comment'],
        [/\d+(\.\d+)?/, 'number'],
        [/[;,.]/, 'delimiter']
      ],
      string: [
        [/[^\\"]+/, 'string'],
        [/\\./, 'string.escape'],
        [/\"/, { token: 'string.quote', bracket: '@close', next: '@pop' }]
      ]
    }
  })
}

export default {
  name: 'OPAPolicyConfig',
  components: {
    CommonTitle,
    CustomLayout
  },
  data () {
    return {
      form: {
        policy_id: DEFAULT_POLICY_ID
      },
      monacoEditor: null
    }
  },
  mounted () {
    this.$nextTick(() => {
      this.initPolicyEditor()
      this.handleSync(false)
      window.addEventListener('resize', this.handleEditorResize)
    })
  },
  methods: {
    handleEditorResize() {
      if (this.monacoEditor) {
        this.monacoEditor.layout()
      }
    },
    initPolicyEditor() {
      registerRegoLanguage(monaco)
      let cached = getLocalStorage(POLICY_CACHE_KEY)
      let value = cached && cached !== '' ? cached : DEFAULT_POLICY
      this.monacoEditor = monaco.editor.create(document.getElementById('policyEditor'), {
        value,
        language: 'rego',
        codeLens: true,
        selectOnLineNumbers: true,
        roundedSelection: false,
        readOnly: false,
        lineNumbers: 'on',
        theme: 'vs-dark',
        wordWrapColumn: 120,
        folding: false,
        showFoldingControls: 'always',
        wordWrap: 'wordWrapColumn',
        cursorStyle: 'line',
        automaticLayout: true
      })
      const model = this.monacoEditor.getModel()
      if (model) {
        model.setEOL(monaco.editor.EndOfLineSequence.LF)
      }
      monaco.editor.remeasureFonts()
      this.monacoEditor.layout()
    },
    handleSync(showMessage = true) {
      this.$get('/config/api/opa/policy', {
        policy_id: this.form.policy_id || DEFAULT_POLICY_ID
      })
        .then((res) => {
          if (res) {
            let content = ''
            if (typeof res === 'object') {
              if (res.code == 10001) {
                content = res.data || ''
              } else if (res.result && typeof res.result.raw === 'string') {
                content = res.result.raw
              } else if (typeof res.data === 'string') {
                content = res.data
              }
            } else if (typeof res === 'string') {
              content = res
            }
            if (content.trim() === '') {
              this.setEditorValue(DEFAULT_POLICY)
              if (showMessage) {
                this.$message({
                  type: 'warning',
                  message: '当前无策略，请编写并部署',
                })
              }
            } else {
              this.setEditorValue(content)
            }
          }
        })
        .catch((err) => {
          console.error(err)
          this.$message({
            type: 'error',
            message: '同步失败，请稍后重试',
          })
        })
    },
    handleDelete() {
      this.$delete('/config/api/opa/policy', {
        policy_id: this.form.policy_id || DEFAULT_POLICY_ID
      })
        .then((res) => {
          if (res.code == 10001) {
            this.setEditorValue(DEFAULT_POLICY)
            this.$message({
              type: 'success',
              message: '删除成功',
            })
          }
        })
        .catch((err) => {
          console.log(err)
        })
    },
    setEditorValue(value) {
      if (this.monacoEditor) {
        this.monacoEditor.setValue(value)
        const model = this.monacoEditor.getModel()
        if (model) {
          model.setEOL(monaco.editor.EndOfLineSequence.LF)
        }
        this.$nextTick(() => {
          this.monacoEditor.layout()
        })
      }
    },
    handleSave() {
      let formData = new FormData()
      let policy = this.monacoEditor ? this.monacoEditor.getValue() : ''
      if (!policy || policy.trim() === '') {
        this.$message({
          type: 'warning',
          message: '策略内容不能为空',
        })
        return
      }
      policy = normalizePolicyText(policy)
      formData.append('content', policy)
      formData.append('policy_id', this.form.policy_id || DEFAULT_POLICY_ID)
      this.$put('/config/api/opa/policy', formData)
        .then((res) => {
          if (res.code == 10001) {
            setLocalStorage(POLICY_CACHE_KEY, policy)
            this.$message({
              type: 'success',
              message: '保存成功',
            })
          }
        })
        .catch((err) => {
          console.log(err)
        })
    }
  },
  destroyed() {
    if (this.monacoEditor) {
      this.monacoEditor.dispose()
    }
    window.removeEventListener('resize', this.handleEditorResize)
  }
}
</script>

<style scoped lang="less">
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
<style lang="less">
#policyEditor {
  height: 100%;
  width: 100%;
  overflow: hidden;
  font-family: Consolas, "Courier New", monospace;
}
</style>
