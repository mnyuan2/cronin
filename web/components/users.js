var MyUsers = Vue.extend({
    template: `<el-main>
        <el-button type="text" @click="addBox(true)" v-if="$auth_tag.user_set">添加用户</el-button>
        <el-table :data="list">
            <el-table-column property="sort" label="序号"></el-table-column>
            <el-table-column property="username" label="用户名"></el-table-column>
            <el-table-column property="mobile" label="手机号"></el-table-column>
            <el-table-column property="status_name" label="状态"></el-table-column>
            <el-table-column property="update_dt" label="更新时间"></el-table-column>
            <el-table-column label="操作">
                <template slot-scope="scope">
                    <el-button plain size="small" @click="detailBox(true, scope.row)">查看</el-button>
                    <el-button plain size="small" @click="detailDelete(scope.row.id)">删除</el-button>
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
                <el-form-item label="密码">
                    <el-input v-model="form.data.password" placeholder="默认值 123456"></el-input>
                </el-form-item>
                <el-form-item label="角色*">
                    <el-select v-model="form.data.role_ids" multiple filterable clearable>
                        <el-option v-for="item in dic_role" :key="item.id" :label="item.name" :value="item.id"></el-option>
                    </el-select>
                </el-form-item>
            </el-form>
            <div slot="footer" class="dialog-footer">
                <el-button @click="addBox(false)">取 消</el-button>
                <el-button type="primary" @click="addSubmit">确 定</el-button>
            </div>
        </el-dialog>
        
        <el-drawer title="用户详情" :visible.sync="detail_box.show" direction="rtl" size="40%" wrapperClosable="false">
            <my-user v-if="detail_box.show" :data_id="detail_box.data_id" @close="detailClose"></my-user>
        </el-drawer>
    </el-main>`,

    name: "MyUsers",
    data(){
        return {
            dic_role:[],
            list:[],
            page:{},
            listParam:{
                page: 1,
                size: 20,
            },
            form:{}, // 表单
            detail_box:{
                show: false,
                data_id: 0
            }

        }
    },
    // 模块初始化
    created(){
        setDocumentTitle('用户')
        this.addBox(false)
    },
    // 模块初始化
    mounted(){
        console.log("sql_source mounted")
        this.getList()
        this.getDic()
    },

    // 具体方法
    methods:{
        getDic(){
            let types = [
                Enum.dicRole,
            ]
            api.dicList(types,(res) =>{
                this.dic_role = res[Enum.dicRole]
            })
        },
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

        addSubmit(){
            let body = copyJSON(this.form.data)
            if (body.username == ""){
                return this.$message.warning("请输入用户名");
            }
            if (body.password == ''){
                body.password = '123456'
            }
            body.sort = Number(body.sort)
            let path = "/user/set"
            if (body.id != this.$user.id){
                path += "?auth_type=set"
            }

            api.innerPost(path, body, (res) =>{
                if (!res.status){
                    return this.$message.error(res.message)
                }
                this.addBox(false)
                this.getList()
                // api.dicList([Enum.dicUser],()=>{}, true) // 存在变化，更新缓存
                this.detailBox(true, {id: res.data.id})
            })
        },
        // 初始化表单数据
        addBox(show){
            this.form = {
                box:{
                    show: Boolean(show),
                    title: "新增用户",
                },
                data: {
                    id: 0,
                    username:"",
                    mobile: "",
                    sort: "1",
                    password: '',
                }
            }
        },
        detailBox(show, data){
            this.detail_box.show = show == true
            if ( typeof data === 'object' && data["id"] != undefined &&  data.id > 0 ){
                this.detail_box.data_id = data.id
            }else{
                this.detail_box.data_id = 0
            }
        },
        detailClose(e){
          if (e.is_change){
              this.getList()
          }
        },
        detailDelete(row){
            this.$confirm('确认删除用户', '提示',{
                type: 'warning',
            }).then(()=>{
                // 确认操作
                api.innerPost("/user/change_status", {id:row.id, status:9}, (res)=>{
                    if (!res.status){
                        return this.$message.error(res.message)
                    }
                    return this.$message.success(res.message)
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

Vue.component("MyUsers", MyUsers);