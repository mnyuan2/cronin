var MyEnv = Vue.extend({
    template: `<div>
        <el-button type="primary" plain @click="initForm(true)" style="margin-left: 20px">新增环境</el-button>
        
        <el-table :data="list">
            <el-table-column property="name" label="key"></el-table-column>
            <el-table-column property="title" label="环境名称"></el-table-column>
            <el-table-column property="default" label="是否默认">
                <template slot-scope="scope">
                    <el-switch :value="scope.row.default" active-value="2" inactive-value="1" validate-event="true" @change="setContent(scope.row,$event)"></el-switch>
                </template>
            </el-table-column>
            <el-table-column property="status_name" label="状态"></el-table-column>
            <el-table-column property="update_dt" label="更新时间"></el-table-column>
            <el-table-column label="操作">
                <template slot-scope="{row}">
                    <el-button type="text" @click="initForm(true, row)">编辑</el-button>
                    <el-button type="text" @click="changeStatus(row, 2)" v-if="row.status!=2">激活</el-button>
                    <el-button type="text" @click="changeStatus(row, 1)" v-if="row.status==2">停用</el-button>
                    <el-button type="text" @click="deleteEnv(row.id)" v-if="row.status!=2">删除</el-button>
                    
                    
                </template>
                
            </el-table-column>
        </el-table>

        <!--设置弹窗-->
        <el-dialog :title="form.box.title" :visible.sync="form.box.show" :close-on-click-modal="false" append-to-body="true" width="400px">
            <el-form :model="form.data" ref="form.data" :rules="form.rule" label-position="left" label-width="80px" size="small">
                <el-form-item label="key" prop="name">
                    <el-input v-model="form.data.name" :disabled="form.data.id>0"></el-input>
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

    name: "MyEnv",
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
                    name: [
                        { required: true, message: '请输入key', trigger: 'blur'},
                        { min: 1, max: 8, message: '长度在 1 到 8 个字符', trigger: 'blur' }
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
                if (!res.status){
                    console.log("setting/env_list 错误", res)
                    return this.$message.error(res.message);
                }
                for (let i in res.data.list){
                    if (res.data.list[i].default){
                        res.data.list[i].default = res.data.list[i].default.toString()
                    }
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

        submitForm(){
            let body = copyJSON(this.form.data)
            this.$refs['form.data'].validate((valid) => {
                if (!valid) {
                    console.log("submit", valid)
                    return false;
                }
                body.default = Number(body.default)
                api.innerPost("/setting/env_set", body, (res) =>{
                    console.log("sql源设置响应",res)
                    if (!res.status){
                        return this.$message.error(res.message)
                    }
                    this.initForm(false)
                    this.getList()
                })
            });
        },

        setContent(row,newValue){
            if (newValue != 2){// 只能启用不能关闭
                return this.$message.warning("请选择其它环境，当前默认将会自动取消！")
            }
            this.$confirm(newValue==2? '确认启用默认，之前的默认将会自动取消！': '确认关闭默认','提示',{
                type:'warning',
            }).then(()=>{
                api.innerPost("/setting/env_set_content", {id: row.id, default:Number(newValue)}, (res)=>{
                    if (!res.status){
                        console.log("setting/env_set_content 错误", res)
                        return this.$message.error(res.message)
                    }
                    row.default = newValue.toString()
                    this.getList()
                })
            }).catch(()=>{
                // row.default = (Number(newValue == 2)+1).toString()
            })
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
                name: "",
            }
            if ( typeof data === 'object' && data["id"] != undefined && data.id > 0){
                let _data = copyJSON(data)
                this.form.box.title = '编辑环境'
                this.form.data = _data
                console.log("编辑源",_data)
            }
        },

        // 改变状态
        changeStatus(row, newStatus){
            this.$confirm(newStatus==1? '环境停用前，请确保不存在进行中的任务！': '确认环境激活', '提示',{
                type: 'warning',
            }).then(()=>{
                // 确认操作
                api.innerPost("/setting/env_change_status", {id:row.id,status:Number(newStatus)}, (res)=>{
                    if (!res.status){
                        console.log("setting/env_change_status 错误", res)
                        return this.$message.error(res.message)
                    }
                    row.status = newStatus.toString()
                    row.status_name = newStatus == 1 ? '停用' :'激活'
                    return this.$message.success(res.message)
                })
            }).catch(()=>{
                // 取消操作
            })
        },

        // 删除连接
        deleteEnv(id){
            if (id<=0){
                return this.$message.warning('参数异常，操作取消')
            }
            this.$confirm('环境删除后，下面所有任务将会连带删除！', '提示',{
                type: 'warning',
            }).then(()=>{
                api.innerPost("/setting/env_del", {id:id}, (res) =>{
                    if (!res.status){
                        console.log("setting/env_del 错误", res)
                        return this.$message.error(res.message)
                    }
                    this.getList()
                })
            }).catch(()=>{
                // 取消操作
            })
        },

        close(){
            this.$emit('update:visible', false) // 向外传递关闭表示
        }
    }
})

Vue.component("MyEnv", MyEnv);