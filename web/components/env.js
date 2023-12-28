var Env = Vue.extend({
    template: `<div>
        <el-button type="primary" plain @click="initForm(true)" style="margin-left: 20px">新增环境</el-button>
        
        <el-table :data="list">
            <el-table-column property="key" label="key"></el-table-column>
            <el-table-column property="title" label="环境名称"></el-table-column>
            <el-table-column property="status_name" label="状态"></el-table-column>
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
            <el-form :model="form.data" ref="form.data" :rules="form.rule" label-position="left" label-width="80px" size="small">
                <el-form-item label="key" prop="key">
                    <el-input v-model="form.data.key"></el-input>
                </el-form-item>
                <el-form-item label="名称" prop="title">
                    <el-input v-model="form.data.title"></el-input>
                </el-form-item>
            </el-form>
            <div slot="footer" class="dialog-footer">
                <el-button @click="initForm(false,'-')">取 消</el-button>
                <el-button type="primary" @click="submitForm()">确 定</el-button>
            </div>
        </el-dialog>
    </div>`,

    name: "Env",
    props: {
        reload_list:false, // 重新加载列表
    },
    data(){
        return {
            list:[],
            page:{
                index: 1,
                size: 10,
                total: 0
            },
            listParam:{
                page: 1,
                size: 20,
            },
            form:{
                box:{},
                data:{},
                rule:{
                    title: [
                        { required: true, message: '请输入名称', trigger: 'blur'},
                        { min: 1, max: 60, message: '长度在 1 到 60 个字符', trigger: 'blur' }
                    ],
                    key: [
                        { required: true, message: '请输入key', trigger: 'blur'},
                        { min: 1, max: 60, message: '长度在 1 到 60 个字符', trigger: 'blur' }
                    ],
                }
            }, // 表单

        }
    },
    // 模块初始化
    created(){
        this.initForm(false,"-")
    },
    // 模块初始化
    mounted(){
        if (this.reload_list){
            this.getList()
        }

    },

    // 具体方法
    methods:{
        // 列表
        getList(){
            api.innerGet("/setting/env_list", this.listParam, (res)=>{
                console.log("env_list 响应", this.reload_list)
                if (!res.status){
                    return this.$message.error(res.message);
                }
                this.list = res.data.list;
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
            api.innerPost("/setting/env_change_status", {id:id, status: 9}, (res) =>{
                console.log("sql源设置响应",res)
                if (!res.status){
                    return this.$message.error(res.message)
                }
                this.getList()
            })
        },
        submitForm(){
            this.$refs['form.data'].validate((valid) => {
                if (!valid) {
                    return false;
                }
                api.innerPost("/setting/env_set", body, (res) =>{
                    console.log("sql源设置响应",res)
                    if (!res.status){
                        return this.$message.error(res.message)
                    }
                    this.initForm(false)
                    this.getList()
                })
            });
            let body = this.form.data

        },

        // 初始化表单数据
        initForm(show, data){
            // if (this.$refs['form.data'] != undefined){
            //     this.$refs['form.data'].resetFields() // 重置验证, 慎用 容易数据错位。
            // }
            this.form.box = {
                show: show == true,
                title: "添加环境",
            }
            this.form.data = {
                id: 0,
                title: "",
                key: "",
            }
            if ( typeof data === 'object' && data["id"] != undefined && data.id > 0){
                this.form.box.title = '编辑环境'
                this.form.data = data
                console.log("编辑源",data)
            }
        },
        close(){
            this.$emit('update:visible', false) // 向外传递关闭表示
        }
    }
})

Vue.component("Env", Env);