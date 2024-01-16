var MyUsers = Vue.extend({
    template: `<el-main>
        <el-button type="text" @click="initForm(true)">添加用户</el-button>
        <el-table :data="list">
            <el-table-column property="sort" label="序号"></el-table-column>
            <el-table-column property="username" label="用户名"></el-table-column>
            <el-table-column property="mobile" label="手机号"></el-table-column>
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
            <el-form :model="form.data" label-position="left" label-width="100px">
                <el-form-item label="用户名">
                    <el-input v-model="form.data.username"></el-input>
                </el-form-item>
                <el-form-item label="手机号">
                    <el-input v-model="form.data.mobile"></el-input>
                </el-form-item>
                <el-form-item label="序号">
                    <el-input v-model="form.data.sort"></el-input>
                </el-form-item>
                
            </el-form>
            <div slot="footer" class="dialog-footer">
                <el-button @click="initForm(false,'-')">取 消</el-button>
                <el-button type="primary" @click="submitForm()">确 定</el-button>
            </div>
        </el-dialog>
    </el-main>`,

    name: "MyUsers",
    data(){
        return {
            list:[],
            page:{},
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
        console.log("sql_source mounted")
        this.getList()
    },

    // 具体方法
    methods:{
        // 列表
        getList(){
            api.innerGet("/user/list", this.listParam, (res)=>{
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

        submitForm(){
            let body = this.form.data
            if (body.username == ""){
                return this.$message.warning("请输入用户名");
            }
            body.sort = Number(body.sort)

            api.innerPost("/user/set", body, (res) =>{
                if (!res.status){
                    return this.$message.error(res.message)
                }
                this.initForm(false)
                this.getList()
            })
        },
        // 初始化表单数据
        initForm(show, data){
            this.form = {
                box:{
                    show: show == true,
                    title: "新增用户",
                },
                data: {
                    id: 0,
                    username:"",
                    mobile: "",
                    sort: "1",
                }
            }
            if ( typeof data === 'object' && data["id"] != undefined &&  data.id > 0 ){
                this.form.box.title = '编辑用户'
                this.form.data = data
                this.form.data.sort = data.sort.toString()
            }
        },
        close(){
            this.$emit('update:visible', false) // 向外传递关闭表示
        }
    }
})

Vue.component("MyUsers", MyUsers);