var MySqlSource = Vue.extend({
    template: `<el-main>
        <el-button type="text" @click="initForm(true)">新增链接</el-button>
        <el-table :data="sql_source_list">
            <el-table-column property="title" label="链接名称"></el-table-column>
            <el-table-column property="create_dt" label="创建时间"></el-table-column>
            <el-table-column property="update_dt" label="更新时间"></el-table-column>
            <el-table-column label="操作">
                <template slot-scope="scope">
                    <el-button plain @click="initForm(true, scope.row)">编辑</el-button>
                    <el-button plain @click="deleteSqlSource(scope.row.id)">删除</el-button>
                </template>
                
            </el-table-column>
        </el-table>

        <!--设置弹窗-->
        <el-dialog :title="form.box.title" :visible.sync="form.box.show" :close-on-click-modal="false" append-to-body="true" width="400px">
            <el-form :model="form.data" label-position="left" label-width="80px" size="small">
                <el-form-item label="链接名*">
                    <el-input v-model="form.data.title"></el-input>
                </el-form-item>
                <el-form-item label="主机*">
                    <el-input v-model="form.data.source.hostname"></el-input>
                </el-form-item>
                <el-form-item label="端口*">
                    <el-input v-model="form.data.source.port"></el-input>
                </el-form-item>
                <el-form-item label="用户名">
                    <el-input v-model="form.data.source.username"></el-input>
                </el-form-item>
                <el-form-item label="密码">
                    <el-input v-model="form.data.source.password" show-password="true"></el-input>
                </el-form-item>
                <el-form-item label="选中库名">
                    <el-input v-model="form.data.source.database"></el-input>
                </el-form-item>
            </el-form>
            <div slot="footer" class="dialog-footer">
                <el-button @click="pingForm()" style="float: left;">连接测试</el-button>
                <el-button @click="initForm(false,'-')">取 消</el-button>
                <el-button type="primary" @click="submitForm()">确 定</el-button>
            </div>
        </el-dialog>
    </el-main>`,

    name: "MySqlSource",
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
            form:{}, // 表单

        }
    },
    // 模块初始化
    created(){
        this.initForm(false,"-")
    },
    // 模块初始化
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
                if (!res.status){
                    return this.$message.error(res.message);
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
        // 删除连接
        deleteSqlSource(id){
            if (id<0){
                return this.$message.warning('参数异常，操作取消')
            }
            api.innerPost("/setting/sql_source_change_status", {id:id, status: 9}, (res) =>{
                console.log("sql源设置响应",res)
                if (!res.status){
                    return this.$message.error(res.message)
                }
                this.getList()
            })
        },
        submitForm(){
            let body = this.form.data
            api.innerPost("/setting/sql_source_set", body, (res) =>{
                console.log("sql源设置响应",res)
                if (!res.status){
                    return this.$message.error(res.message)
                }
                this.initForm(false)
                this.getList()
            })
        },
        // sql连接连接
        pingForm(){
            let body = this.form.data.source
            api.innerPost("/setting/sql_source_ping", body, (res) =>{
                if (!res.status){
                    console.log("setting/sql_source_ping 错误", res)
                    return this.$message.error(res.message)
                }
                return this.$message.success('连接成功');
            })
        },
        // 初始化表单数据
        initForm(show, data){
            this.form = {
                box:{
                    show: show == true,
                    title: "添加sql链接",
                },
                data: {
                    id: 0,
                    title:"",
                    source:{
                        hostname: "",
                        port: "",
                        database:"",
                        username: "",
                        password: ""
                    }
                }
            }
            if ( typeof data === 'object' && data["id"] != undefined && data["source"] != undefined
                && data.id > 0 && typeof data.source === 'object'){
                this.form.box.title = '编辑sql连接'
                this.form.data = data
                console.log("编辑源",data)
            }
        },
        close(){
            this.$emit('update:visible', false) // 向外传递关闭表示
        }
    }
})

Vue.component("MySqlSource", MySqlSource);