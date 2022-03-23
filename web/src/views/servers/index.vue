<template>
    <div class="app-container">
        <div style="margin-left:25px;margin-right:20px">
            <el-row>
                  <el-button type="primary" >新增</el-button>
            </el-row>
            <el-row :gutter="12">
                <el-col :span="8" v-for="(item) in hosts" :key="item.id" style="margin-top:20px">
                    <el-card shadow="hover">
                        <div slot="header" class="clearfix">
                            <el-row type="flex" justify="space-between">
                                <el-col >
                                    <span>服务器名称： {{item.hostname}}</span> <br>
                                    <span>服务器IP： {{item.addr}}</span>
                                </el-col>
                                <el-col>
                                    <el-button style="float:right" type="primary" plain @click.stop="opendrawer(item)" >命令行界面</el-button>
                                </el-col>
                            </el-row>
                        </div>
                    </el-card>
                </el-col>  
            </el-row>
        </div>
        <el-drawer
            title="命令行"
            :visible.sync="drawer"
            v-if="drawer"
            direction="rtl"
            size="40%"
            :before-close="handleClose">
            <el-divider></el-divider>
            <div style="height: 100%;background: black;margin-top:-20px">
                <my-terminal  :terminal="terminal"
                            :url="url"></my-terminal>
            </div>
        </el-drawer>
        <el-dialog
        title="新增服务器"
        :visible.sync="dialogVisible"
        width="30%"
        :before-close="handleClosedig">
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
</template>

<script>
import { gethost } from "@/api/hosts";
import Console from './console'

export default {
  name: "servers",
  data() {
      return {
          hosts:[],
          hostquery: {
           offset: 0,
           limit: 15
          },
          changeform:{
              addr:"",
              hostname:"",
              username:"",
              password:""
          },
          drawer: false,
          dialogVisible: false,
          terminal: {
                    pid: 1,
                    name: 'terminal',
                    cols: 400,
                    rows: 400
                },
        //    url : `ws://127.0.0.1:8080/ws?addr=114.215.84.163&username=root&password=4477123Wl!&cols=188&rows=50`
          url:""
      }
    },
    created() {
        this.startgethosts();
    },
    components: {
            'my-terminal': Console
    },
    methods:{
        startgethosts() {
            gethost(this.hostquery).then(resp => {
                this.hosts = resp.data
            })
        },
        opendrawer(row) {
            this.url = `${process.env.VUE_APP_BASE_API_WS}/ws?addr=${row.addr}&username=${row.username}&password=${row.hashpassword}&cols=188&rows=50`
            this.drawer = true;
        },
        handleClose(done) {
            this.$confirm('确认关闭？')
            .then(_ => {
                this.drawer = false;
                // this.url="";
                done();
            })
            .catch(_ => {});
        },
        handleClosedig(done) {
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
    }
}
</script>