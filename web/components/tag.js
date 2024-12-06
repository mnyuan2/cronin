var MyTag = Vue.extend({
    template: `<el-main>
        <el-button type="text" @click="setShow()" v-if="$auth_tag.tag_set">新增</el-button>
        <el-table :data="list">
            <el-table-column property="name" label="标记"></el-table-column>
            <el-table-column property="remark" label="备注"></el-table-column>
            <el-table-column property="create_user_name" label="创建人"></el-table-column>
            <el-table-column property="create_dt" label="创建时间"></el-table-column>
            <el-table-column property="update_user_name" label="最近编辑人"></el-table-column>
            <el-table-column property="update_dt" label="最近编辑时间"></el-table-column>
            <el-table-column label="操作">
                <template slot-scope="scope">
                    <el-button plain size="small" @click="setShow(scope.row)" v-if="$auth_tag.tag_set">编辑</el-button>
                    <el-button plain size="small" @click="changeStatus(scope.row, Enum.StatusClosed)" v-if="$auth_tag.tag_status">删除</el-button>
                </template>
            </el-table-column>
        </el-table>

        <!--设置弹窗-->
        <el-dialog :title="set_box.title" :visible.sync="set_box.show" :close-on-click-modal="false" append-to-body="true" width="400px">
            <el-form :model="set_box.form" label-position="left" label-width="100px">
                <el-form-item label="名称*">
                    <el-input v-model="set_box.form.name"></el-input>
                </el-form-item>
                <el-form-item label="备注">
                    <el-input v-model="set_box.form.remark"></el-input>
                </el-form-item>
            </el-form>
            <div slot="footer" class="dialog-footer">
                <el-button @click="set_box.show = false">取 消</el-button>
                <el-button type="primary" @click="setSubmit">确 定</el-button>
            </div>
        </el-dialog>
    </el-main>`,
    name: "MyTag",
    props: {},
    data(){
        return {
            list:[],
            set_box:{
                title: '',
                show: false,
                form:{}, // 表单
            }

        }
    },
    // 模块初始化
    created(){
        setDocumentTitle('标签')
    },
    // 模块初始化
    mounted(){
        // this.user = cache.getUser()
        this.getList()
    },

    // 具体方法
    methods:{
        // 列表
        getList(){
            api.innerGet("/tag/list", {}, (res)=>{
                if (!res.status){
                    return this.$message.error(res.message);
                }
                this.list = res.data.list;
            })
        },
        // 初始化表单数据
        setShow(row=null){
            this.set_box.show = true
            if (row == null){
                this.set_box.title= '添加任务'
                this.set_box.form = {
                    id: 0,
                    name: '',
                    remark: ''
                }
            }else{
                this.set_box.title= '编辑任务'
                this.set_box.form = row
            }
        },
        setSubmit(){
            let body = copyJSON(this.set_box.form)
            if (body.name == ""){
                return this.$message.warning("请输入名称");
            }
            if (body.name.length > 32){
                return this.$message.warning("名称长度不得超过32字符");
            }

            api.innerPost("/tag/set", body, (res) =>{
                if (!res.status){
                    return this.$message.error(res.message)
                }
                this.set_box.show = false
                this.getList()
                // api.dicList([Enum.dicUser],()=>{}, true) // 存在变化，更新缓存
            })
        },
        changeStatus(row, status){
            if (status != Enum.StatusClosed){
                return this.$message.warning('错误操作...')
            }
            this.$confirm('确认删除标签', '提示',{
                type: 'warning',
            }).then(()=>{
                api.innerPost("/user/change_status", {id:row.id, status:status}, (res)=>{
                    if (!res.status){
                        return this.$message.error(res.message)
                    }
                    return this.$message.success(res.message)
                    this.getList()
                })
            }).catch(()=>{
                // 取消操作
            })
        }

    }
})

Vue.component("MyTag", MyTag);