<template>
  <div class="app-container">
    <div style="margin-left:25px;margin-right:20px">
      <el-table :data="perimisions">
        <el-table-column align="center" property="p_type" label="策略" min-width="10"></el-table-column>
        <el-table-column align="center" property="v0.String" label="角色" min-width="20"></el-table-column>
        <el-table-column align="left" property="v1.String" label="权限" min-width="60"></el-table-column>
        <el-table-column align="center" property="v2.String" label="请求方式" min-width="20"></el-table-column>
        <el-table-column align="center" label="操作" min-width="60">
          <template slot-scope="scope">
            <el-button-group>
              <el-tooltip class="item" effect="dark" content="新增" placement="top">
                <el-button type="primary" size="mini" @click="dialogVisible=true">新增</el-button>
              </el-tooltip>
              <el-popconfirm
                :hideIcon="true"
                title="是否要删除策略?"
                @confirm="delper(scope.row)"
              >
                <el-button type="danger" slot="reference" size="mini" >删除</el-button>
              </el-popconfirm>
            </el-button-group>
          </template>
        </el-table-column>
      </el-table>
      <div style="margin-top: 10px;float:right;height: 70px;">
        <el-pagination
         :page-size="hostquery.limit"
          @current-change="handleCurrentChangerun"
          background
          layout="total,prev, pager, next"
          :total="pagecount"
        ></el-pagination>
      </div>
      <!-- 对话框 -->
      <el-dialog
      title="新增策略"
      :visible.sync="dialogVisible"
      width="30%"
      :before-close="handleClose">
      <el-form ref="ruleForm" label-position="right" label-width="80px" :model="changeform">
        <el-form-item label="策略">
          <el-input v-model="changeform.p_type"></el-input>
        </el-form-item>
        <el-form-item label="角色">
          <el-input v-model="changeform.v0"></el-input>
        </el-form-item>
        <el-form-item label="权限">
          <el-input v-model="changeform.v1"></el-input>
        </el-form-item>
        <el-form-item label="请求方式">
          <el-input v-model="changeform.v2"></el-input>
        </el-form-item>
      </el-form>
      <span slot="footer" class="dialog-footer">
        <el-button @click="handleClose()">取 消</el-button>
        <el-button type="primary" @click="submitdata()">确 定</el-button>
      </span>
    </el-dialog>
    </div>
  </div>
</template>

<script>
import {getperimision,set_permission,del_permission} from '@/api/perimision'
import { Message } from "element-ui";

export default {
  name: "index",
  data() {
    return {
      dialogVisible: false,
      perimisions: [],
      pagecount: 0,
      hostquery: {
        offset: 0,
        limit: 15
      },
      changeform: {
        p_type: "",
        v0:"",
        v1:"",
        v2:""
      }
    }
  },
   created() {
    this.startgetdata();
  },
  methods: {
    startgetdata() {
      getperimision(this.hostquery).then(resp => {
        this.perimisions = resp.data;
        this.pagecount = resp.count;
      });
    },
    handleCurrentChangerun(page) {
      this.hostquery.offset = (page - 1) * this.hostquery.limit;
      this.startgetdata();
    },
    handleClose(done) {
      this.$confirm('确认关闭？')
        .then(_ => {
          this.changeform = {
            p_type: "",
            v0:"",
            v1:"",
            v2:""
          }
          this.dialogVisible=false;
          done();
        })
        .catch(_ => {});
    },
    submitdata() {
      this.dialogVisible = !this.dialogVisible;
      set_permission(this.changeform).then(res =>{
            if (res.code === 0) {
              Message.success("新增策略成功")
            }else{
              Message.error(`${res.msg}`)
            }
            this.startgetdata();
      });
      this.changeform = {
        p_type: "",
        v0:"",
        v1:"",
        v2:""
      };
    },
    delper(row) {
      del_permission(row).then(resp => {
        if (resp.code === 0){
          Message.success("删除成功");
          this.startgetdata();
        }else{
          Message.warning(`删除失败 ${resp.msg}`)
        }
      })
    }
  }
}
</script>

<style scoped>

</style>