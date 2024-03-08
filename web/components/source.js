var MySource = Vue.extend({
    template: `<el-main>
        <el-menu :default-active="list.labelIndex" class="el-menu-demo" mode="horizontal" @select="handleClickTypeLabel">
            <el-menu-item index="11" :disabled="list.request">mysql</el-menu-item>
            <el-menu-item index="12" :disabled="list.request">jenkins</el-menu-item>
            <el-menu-item index="13" :disabled="list.request">git</el-menu-item>
            <div style="float: right">
                <el-button type="text" @click="initForm(true)">新增链接</el-button>
            </div>
        </el-menu>
        
        <el-table :data="list.items">
            <el-table-column property="title" label="链接名称"></el-table-column>
            <el-table-column property="create_dt" label="创建时间"></el-table-column>
            <el-table-column property="update_dt" label="更新时间"></el-table-column>
            <el-table-column label="操作">
                <template slot-scope="scope">
                    <el-button plain @click="initForm(true, scope.row)">编辑</el-button>
                    <el-button plain @click="deleteSource(scope.row.id)">删除</el-button>
                </template>
            </el-table-column>
        </el-table>

        <!--设置弹窗-->
        <el-dialog :title="form.box.title" :visible.sync="form.box.show" :close-on-click-modal="false" append-to-body="true" width="600px">
            <el-form :model="form.data" label-position="left" label-width="80px" size="small">
                <el-form-item label="链接名*">
                    <el-input v-model="form.data.title"></el-input>
                </el-form-item>
                
                <el-tabs type="border-card" v-model="form.data.type">
                    <el-tab-pane label="sql" name="11">
                        <el-form-item label="主机*">
                            <el-input v-model="form.data.source.sql.hostname"></el-input>
                        </el-form-item>
                        <el-form-item label="端口*">
                            <el-input v-model="form.data.source.sql.port"></el-input>
                        </el-form-item>
                        <el-form-item label="用户名">
                            <el-input v-model="form.data.source.sql.username"></el-input>
                        </el-form-item>
                        <el-form-item label="密码">
                            <el-input v-model="form.data.source.sql.password" show-password="true"></el-input>
                        </el-form-item>
                        <el-form-item label="选中库名">
                            <el-input v-model="form.data.source.sql.database"></el-input>
                        </el-form-item>
                    </el-tab-pane>
                    
                    <el-tab-pane label="jenkins" name="12">
                        <el-form-item label="地址*">
                            <el-input v-model="form.data.source.jenkins.hostname" placeholder="http://ip:prod 或 https://hostname"></el-input>
                        </el-form-item>
                        <el-form-item label="用户名">
                            <el-input v-model="form.data.source.jenkins.username" placeholder="登录账户"></el-input>
                        </el-form-item>
                        <el-form-item label="密码">
                            <el-input v-model="form.data.source.jenkins.password" show-password="true" placeholder="api token"></el-input>
                        </el-form-item>
                    </el-tab-pane>
                    <el-tab-pane label="git" name="13">
                        <el-form-item label="类型">
                            <el-select v-model="form.data.source.git.type">
                                <el-option label="gitee" value="gitee"></el-option>
                                <!--<el-option label="github" value="github"></el-option>-->
                            </el-select>
                        </el-form-item>
                        <el-form-item label="授权码">
                            <el-input v-model="form.data.source.git.access_token" placeholder="gitee.com / 个人设置 / 私密令牌 可进行获取"></el-input>
                        </el-form-item>
                    </el-tab-pane>
                </el-tabs>
                
            </el-form>
            <div slot="footer" class="dialog-footer">
                <el-button @click="pingForm()" style="float: left;">连接测试</el-button>
                <el-button @click="initForm(false,'-')">取 消</el-button>
                <el-button type="primary" @click="submitForm()">确 定</el-button>
            </div>
        </el-dialog>
    </el-main>`,

    name: "MySource",
    data(){
        return {
            list:{
                labelIndex: '11',
                items: [],
                page:{
                    index: 1,
                    size: 10,
                    total: 0
                },
                param:{
                    type: 11,
                    page: 1,
                    size: 20,
                },
                request: false,
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
        this.getList()
    },

    // 具体方法
    methods:{
        handleClickTypeLabel(tab, event) {
            this.list.param.type = tab
            this.getList()
        },
        // 列表
        getList(){
            api.innerGet("/setting/source_list", this.list.param, (res)=>{
                if (!res.status){
                    return this.$message.error(res.message);
                }
                this.list.items = res.data.list;
                this.list.page = res.data.page;
            })
        },
        handleSizeChange(val) {
            console.log(`每页 ${val} 条`);
        },
        handleCurrentChange(val) {
            this.list.param.page = val
            this.getList()
        },
        // 删除连接
        deleteSource(id){
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
        // 提交表单
        submitForm(){
            let body = this.form.data
            body.type = Number(body.type)
            api.innerPost("/setting/source_set", body, (res) =>{
                console.log("源设置响应",res)
                if (!res.status){
                    return this.$message.error(res.message)
                }
                this.initForm(false)
                this.getList()

                api.dicDel([body.type])
            })
        },
        // sql连接连接
        pingForm(){
            let body = JSON.parse(JSON.stringify(this.form.data));
            body.type = Number(body.type)
            api.innerPost("/setting/source_ping", body, (res) =>{
                if (!res.status){
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
                    type: this.list.param.type,
                    source:{
                        sql:{
                            hostname: "",
                            port: "",
                            database:"",
                            username: "",
                            password: ""
                        },
                        jenkins:{
                            hostname: "",
                            database:"",
                            username: ""
                        },
                        git:{
                            type: 'gitee',
                            access_token: ''
                        }
                    }
                }
            }
            if ( typeof data === 'object' && data["id"] != undefined && data["source"] != undefined
                && data.id > 0 && typeof data.source === 'object'){
                this.form.box.title = '编辑连接'
                this.form.data = data
                console.log("编辑源",data)
            }
            this.form.data.type = this.form.data.type.toString()
        },
        close(){
            this.$emit('update:visible', false) // 向外传递关闭表示
        }
    }
})

Vue.component("MySource", MySource);