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
        <CommonTitle title="集群管理"></CommonTitle>
      </div>
      <div class="table-head">
        <div class="custom-tools__info">集群列表</div>
        <el-button
          type="primary"
          icon="el-icon-plus"
          size="mini"
          @click="handleAdd"
          >新增</el-button
        >
      </div>
      <div class="custom-panel">
        <el-table
          v-loading="loading"
          size="normal"
          :data="tableData"
          :empty-text="emptText"
          row-class-name="custom-table-tr--hover"
          header-row-class-name="custom-table-header"
          @selection-change="handleSelectionChange"
          style="width: 100%"
        >
          <el-table-column
            class-name="custom-popper--overflow"
            prop="name"
            label="名称"
          >
          </el-table-column>
          <el-table-column
            class-name="custom-popper--overflow"
            prop="type"
            label="类型"
          >
          </el-table-column>
          <el-table-column
            class-name="custom-popper--overflow"
            prop="id"
            label="ID"
          >
          </el-table-column>
          <el-table-column
            class-name="custom-popper--overflow"
            prop="address"
            label="地址"
          >
          </el-table-column>
          <el-table-column
            class-name="custom-popper--overflow"
            prop="port"
            label="端口"
          >
          </el-table-column>
          <el-table-column label="操作">
            <template slot-scope="scope">
              <el-button type="text" @click="handleLook(scope.row)"
                >查看</el-button
              >
              <el-button
                type="text"
                class="dltbtn"
                @click="deleteRow(scope.row)"
                >删除</el-button
              >
            </template>
          </el-table-column>
        </el-table>

        <div class="custom-pagination" style="float: right">
          <div style="display: flex">
            <div>
              总共<span>{{ pagination.total }}</span
              >条记录<i class="custom-pagination__interval"></i>每页显示
              <el-select
                v-model="pagination.pageSize"
                placeholder="请选择"
                style="width: 70px"
              >
                <el-option
                  v-for="item in pagination.pageSizes"
                  :key="item"
                  :label="item"
                  :value="item"
                >
                </el-option>
              </el-select>
              条记录
            </div>
            <el-pagination
              :current-page.sync="pagination.pageIndex"
              :page-size.sync="pagination.pageSize"
              layout="prev, pager, next"
              :total="pagination.total"
            >
            </el-pagination>
          </div>
        </div>
      </div>
    </div>
    <el-dialog
      title="新增"
      :visible.sync="createDialogVisible"
      width="640px"
      :before-close="handleClose"
    >
      <div class="dialog_main">
        <div
          style="clear: 'both'; height: 300px; width: 100%"
          id="createContainer"
          ref="createContainer"
        />
      </div>
      <span slot="footer" class="dialog-footer">
        <el-button size="mini" @click="handleClose">取 消</el-button>
        <el-button type="primary" size="mini" @click="makeCreate"
          >确 定</el-button
        >
      </span>
    </el-dialog>

    <el-dialog
      title="查看修改"
      :visible.sync="updateDialogVisible"
      width="640px"
      :before-close="handleClose"
    >
      <div class="dialog_main">
        <div
          style="clear: 'both'; height: 300px; width: 100%"
          id="updateContainer"
          ref="updateContainer"
        />
      </div>
      <span slot="footer" class="dialog-footer">
        <el-button size="mini" @click="handleClose">取 消</el-button>
        <el-button type="primary" size="mini" @click="makeModify"
          >确 定</el-button
        >
      </span>
    </el-dialog>
  </CustomLayout>
</template>

<script>
import CommonTitle from "@/components/common/CommonTitle";
import CustomLayout from "@/components/common/CustomLayout.vue";
import * as monaco from "monaco-editor/esm/vs/editor/editor.main.js";
import "monaco-editor/esm/vs/basic-languages/javascript/javascript.contribution";

export default {
  name: "",
  components: {
    CommonTitle,
    CustomLayout,
  },
  data() {
    return {
      form: {
        name: "http_bin",
        type: "http",
        id: "backend",
        address: "http://httpbin.org", // NOSONAR
        port: 80,
      },
      edit: false,
      createDialogVisible: false,
      updateDialogVisible: false,
      dialogVisible: false,
      chooseRow: {
        timeout: "",
        description: "",
        type: "",
        path: "",
      },
      pagination: {
        pageIndex: 1,
        pageSize: 20,
        total: 0,
        pageSizes: [20, 50, 100],
      },
      tableData: [],
      loading: false,
      tableHeight: window.document.documentElement.clientHeight - 248, // default table height
      multipleSelection: [],
      currentRow: {},
      areaDialogVisible: false,
      emptText: "暂无数据",
      isSizeInit: false,
      previewDialogVisible: false,
      currentTemplate: {},
      updateDetailSource: null,
      createDetailSource: null,
      createMonacoEditor: null,
      updateMonacoEditor: null,
    };
  },
  mounted() {
    this.getClusterList();
  },
  methods: {
    initCreateMoacoEditor(language, value) {
      this.createMonacoEditor = monaco.editor.create(
        document.getElementById("createContainer"),
        {
          value,
          language: "yaml",
          tabSize: 2,
          codeLens: true,
          selectOnLineNumbers: true,
          roundedSelection: false,
          readOnly: false,
          lineNumbersMinChars: true,
          theme: "vs-dark",
          wordWrapColumn: 120,
          folding: false,
          showFoldingControls: "always",
          wordWrap: "wordWrapColumn",
          cursorStyle: "line",
          automaticLayout: true,
        }
      );
    },
    initUpdateMoacoEditor(language, value) {
      this.updateMonacoEditor = monaco.editor.create(
        document.getElementById("updateContainer"),
        {
          value,
          language: "yaml",
          tabSize: 2,
          codeLens: true,
          selectOnLineNumbers: true,
          roundedSelection: false,
          readOnly: false,
          lineNumbersMinChars: true,
          theme: "vs-dark",
          wordWrapColumn: 120,
          folding: false,
          showFoldingControls: "always",
          wordWrap: "wordWrapColumn",
          cursorStyle: "line",
          automaticLayout: true,
        }
      );
    },
    makeCreate() {
      let formData = new FormData();
      let data = this.createMonacoEditor.getValue();
      formData.append("content", data);

      this.$put("/config/api/cluster", formData)
        .then((res) => {
          if (res.code == 10001) {
            this.handleClose();
            this.getClusterList();
            console.log(res);
          } else {
          }
        })
        .catch((err) => {
          console.log(err);
        });
    },
    makeModify() {
      let formData = new FormData();
      let data = this.updateMonacoEditor.getValue();
      formData.append("content", data);

      this.$post("/config/api/cluster", formData)
        .then((res) => {
          if (res.code == 10001) {
            this.handleClose();
            this.getClusterList();
            console.log(res);
          } else {
          }
        })
        .catch((err) => {
          console.log(err);
        });
    },
    handleClose() {
      this.createDialogVisible = false;
      this.updateDialogVisible = false;
      if (this.createMonacoEditor != null) {
        this.createMonacoEditor.dispose();
      }
      if (this.updateMonacoEditor != null) {
        this.updateMonacoEditor.dispose();
      }
    },
    handleSelectionChange(val) {
      this.multipleSelection = val;
    },
    getResourceDetail() {
      this.$get("/config/api/plugin_group/detail")
        .then((res) => {
          if (res.code == 10001) {
            let data = JSON.parse(res.data);
            console.log(data);
          } else {
          }
        })
        .catch((err) => {
          console.log(err);
        });
    },
    handleAdd() {
      this.createDialogVisible = true;
      this.edit = false;
      this.$nextTick(() => [this.initCreateMoacoEditor("yaml", "")]);
    },
    getClusterList() {
      this.$get("/config/api/cluster/list", {})
        .then((res) => {
          if (res.code == 10001) {
            this.tableData = res.data;
          } else {
            this.tableData = [];
          }
        })
        .catch((err) => {
          console.log(err);
        });
    },
    makeSure() {
      let formData = new FormData();
      formData.append("groupName", this.chooseRow.groupName);
      let obj = [
        {
          name: "rate limit",
          version: "0.0.1",
          priority: 1000,
          externalLookupName: "ExternalPluginRateLimit",
        },
      ];
      formData.append("plugins", JSON.stringify(obj));
      let mobj = [];
      formData.append("methods", null);
      this.$post("/config/api/plugin_group", formData)
        .then((res) => {
          if (res.code == 10001) {
            this.dialogVisible = false;
            this.getClusterList();
            console.log(res);
          } else {
          }
        })
        .catch((err) => {
          console.log(err);
        });
    },
    // delete cluster
    deleteRow(row) {
      this.$delete("/config/api/cluster", {
        clusterId: row.id,
      })
        .then((res) => {
          if (res.code == 10001) {
            this.$message({
              type: "success",
              message: "删除成功！",
            });
            this.getClusterList();
            console.log(res);
          } else {
          }
        })
        .catch((err) => {
          console.log(err);
        });
    },
    handleLook(row) {
      this.edit = true;
      this.$get("/config/api/cluster/detail", {
        clusterId: row.id,
      })
        .then((res) => {
          if (res.code == 10001) {
            let data = res.data;
            this.modifyDetailSource = data;
            this.updateDialogVisible = true;
            this.$nextTick(() => [this.initUpdateMoacoEditor("yaml", data)]);
          } else {
          }
        })
        .catch((err) => {
          console.log(err);
        });
    },
  },
};
</script>


<style scoped lang="less">
.dltbtn {
  color: red;
}
.custom-panel {
  margin-top: 20px;
}
.custom-tools__info {
  color: rgba(16, 16, 16, 100);
  font-size: 18px;
  text-align: left;
  margin-top: 10px;
}
.custom-tools__content {
  background-color: #fff;
  margin-top: 10px;
  padding: 10px 20px;
}
.table-head {
  display: flex;
  margin-top: 10px;
  justify-content: space-between;
}
.dialog-input {
  margin-bottom: 10px;
}
</style>
