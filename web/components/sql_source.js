var SqlSource = Vue.extend({
    template: `<div>
        <el-button type="primary" plain @click="createShow()">新增链接</el-button>
        
        <el-table :data="sql_source_list">
            <el-table-column property="create_dt" label="链接名称"></el-table-column>
           
            <el-table-column property="duration" label="主机"></el-table-column>
            <el-table-column property="" label="操作">
                <template slot-scope="scope">
                    <el-popover trigger="hover" placement="left">
                        <div>{{scope.row.body}}</div>
                        <div slot="reference" class="name-wrapper">
                            <el-tag size="medium">详情</el-tag>
                        </div>
                    </el-popover>
                </template>
            </el-table-column>
        </el-table>

        <!--设置弹窗-->
        <el-dialog :title="setConfigTitle" :visible.sync="setConfigShow" :close-on-click-modal="false">
            <el-form :model="form">
                <el-form-item label="链接名*">
                    <el-input v-model="form.title"></el-input>
                </el-form-item>
                <el-form-item label="主机*">
                    <el-input v-model="form.source.hostname"></el-input>
                </el-form-item>
                <el-form-item label="端口*">
                    <el-input v-model="form.source.port"></el-input>
                </el-form-item>
                <el-form-item label="用户名">
                    <el-input v-model="form.source.username"></el-input>
                </el-form-item>
                <el-form-item label="密码">
                    <el-input v-model="form.source.password"></el-input>
                </el-form-item>
            </el-form>
            <div slot="footer" class="dialog-footer">
                <el-button @click="setConfigShow = false">取 消</el-button>
                <el-button type="primary" @click="submitForm()">确 定</el-button>
            </div>
        </el-dialog>
    </div>`,

    name: "SqlSource",
    props: {
        reload_list:false, // 重新加载列表
    },
    data(){
        return {
            sql_source_list:[],
            page:{
                index: 1,
                size: 10,
                total: 0
            },
            listParam:{
                page: 1,
                size: 20,
            },
            form:{}, // 表单数据

        }
    },

    mounted(){
        console.log("sql_source:reload_list", this.reload_list)
        if (this.reload_list){
            this.getList()
        }
    },

    // 具体方法
    methods:{
        // 列表
        getList(){
            api.innerGet("/setting/sql_source_list", this.listParam, (res)=>{
                console.log("sql_source:sql_source_list 响应", this.reload_list)
                if (res.code != "000000"){
                    return this.$message.error(res.message);
                }
                for (i in res.data.list){
                    res.data.list[i].status = res.data.list[i].status.toString()
                }
                this.sql_source_list = res.data.list;
                this.page = res.data.page;
            })
        },
        handleSizeChange(val) {
            console.log(`每页 ${val} 条`);
        },
        handleCurrentChange(val) {
            console.log(`当前页: ${val}`);
            this.listParam.page = val
            this.getList()
        },

        // 添加弹窗
        createShow(){
            this.setConfigShow = true
            this.setConfigTitle = '添加任务'
            // 应该有个地方定义空结构体
            this.initFormData()
        },
        submitForm(){
            let body = this.form
            api.innerPost("/setting/sql_source_set", body, (res) =>{
                console.log("sql源设置响应",res)
            })
        },
        initFormData(){
            this.form = {
                id: 0,
                title:"",
                source:{
                    hostname: "",
                    port: "",
                    username: "",
                    password: ""
                }
            }
        },
    }
})

Vue.component("SqlSource", SqlSource);