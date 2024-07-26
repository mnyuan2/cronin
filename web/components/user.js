var MyUser = Vue.extend({
    template: `<el-main>
        <div class="top">
          <el-button type="text" @click="setBox(true)" v-if="user.id==detail.id || $auth_tag.user_edit">编辑</el-button>
          <el-button type="text" @click="passwordBox(true)" v-if="user.id==detail.id || user.role_ids.filter(x => x == 1).length">修改密码</el-button>
          <el-button v-if="detail.status == 2" type="text" @click="changeStatus(1)" v-if="user.role_ids.filter(x => x == 1).length">停用</el-button>
          <el-button v-if="detail.status == 1" type="text" @click="changeStatus(2)" v-if="user.role_ids.filter(x => x == 1).length">激活</el-button>
        </div>
        
        <el-divider></el-divider>
        <el-descriptions title="" column="1" :colon="false" label-style="width: 60px">
            <el-descriptions-item label="账号">
                <span v-if="detail.account==''">未设置</span>
                {{detail.account}}
                <el-link type="primary" @click="accountBox(true)" style="margin-left: 20px" v-if="$auth_tag.user_account">设置账号</el-link>
            </el-descriptions-item>
            <el-descriptions-item label="用户名">{{detail.username}}</el-descriptions-item>
            <el-descriptions-item label="手机号">{{detail.mobile}}</el-descriptions-item>
            <el-descriptions-item label="排序">{{detail.sort}}</el-descriptions-item>
            <el-descriptions-item label="状态">{{detail.status_name}}</el-descriptions-item>
            <el-descriptions-item label="角色">{{detail.roles}}</el-descriptions-item>
            <el-descriptions-item label="创建时间">{{detail.create_dt}}</el-descriptions-item>
            <el-descriptions-item label="更新时间">{{detail.update_dt}}</el-descriptions-item>
        </el-descriptions>

        <!--设置-->
        <el-dialog title="编辑用户信息" :visible.sync="set_box.show" :close-on-click-modal="false" append-to-body="true" width="400px">
            <el-form :model="set_box.form" label-position="left" label-width="100px">
                <el-form-item label="用户名*">
                    <el-input v-model="set_box.form.username"></el-input>
                </el-form-item>
                <el-form-item label="手机号">
                    <el-input v-model="set_box.form.mobile"></el-input>
                </el-form-item>
                <el-form-item label="角色*" v-if="user.role_ids.filter(x => x == 1).length">
                    <el-select v-model="set_box.form.role_ids" multiple filterable clearable>
                        <el-option v-for="item in dic_role" :key="item.id" :label="item.name" :value="item.id"></el-option>
                    </el-select>
                </el-form-item>
                <el-form-item label="序号">
                    <el-input v-model="set_box.form.sort"></el-input>
                </el-form-item>
            </el-form>
            <div slot="footer" class="dialog-footer">
                <el-button @click="setBox(false)">取 消</el-button>
                <el-button type="primary" @click="setSubmit">确 定</el-button>
            </div>
        </el-dialog>
        <!--重置密码-->
        <el-dialog title="重置密码" :visible.sync="password_box.show" :close-on-click-modal="false" append-to-body="true" width="400px">
            <el-form :model="password_box.form" label-position="left" label-width="100px">
                <el-form-item label="新密码">
                    <el-input v-model="password_box.form.password" show-password></el-input>
                </el-form-item>
                <el-form-item label="确认新密码">
                    <el-input v-model="password_box.form.re_password" show-password></el-input>
                </el-form-item>
            </el-form>
            <div slot="footer" class="dialog-footer">
                <el-button @click="passwordBox(false)">取 消</el-button>
                <el-button type="primary" @click="passwordSubmit">确 定</el-button>
            </div>
        </el-dialog>
        <!--重置账号-->
        <el-dialog title="设置账号" :visible.sync="account_box.show" :close-on-click-modal="false" append-to-body="true" width="400px">
            <el-form :model="account_box.form" label-position="left" label-width="100px">
                <el-form-item label="原账号">
                    <el-input v-model="detail.account" :disabled="true"></el-input>
                </el-form-item>
                <el-form-item label="新账号">
                    <el-input v-model="account_box.form.account"></el-input>
                </el-form-item>
            </el-form>
            <div slot="footer" class="dialog-footer">
                <el-button @click="accountBox(false)">取 消</el-button>
                <el-button type="primary" @click="accountSubmit">确 定</el-button>
            </div>
        </el-dialog>
    </el-main>`,

    name: "MyUser",
    props: {
        data_id:Number
    },
    data(){
        return {
            user:{},
            dic_role: [],
            detail:{},
            set_box:{
                show: false,
                form:{}
            },
            password_box:{
                show: false,
                form:{}
            },
            account_box:{
                show: false,
                form:{}
            }

        }
    },
    // 模块初始化
    created(){
        setDocumentTitle('用户信息')
        this.user = cache.getUser()
        let searchParams = getHashParams(window.location.href)
        let data_id = searchParams["data_id"]
        if (!this.data_id && data_id){
            this.data_id = data_id
        }
    },
    // 模块初始化
    mounted(){
        this.getDic()
        this.getDetail()
        console.log("user",this.user)
    },

    // 具体方法
    methods:{
        // 查询详情
        getDetail(){
            if (!this.data_id){
                return
            }
            api.innerGet("/user/detail", {id: this.data_id}, (res)=>{
                if (!res.status){
                    return this.$message.error(res.message);
                }
                res.data.roles = ''
                res.data.role_ids.forEach((id) => {
                    this.dic_role.forEach((item)=>{
                        if (item.id == id){
                            res.data.roles += item.name +" "
                        }
                    })
                })
                this.detail = res.data;
            })
        },

        setSubmit(){
            let body  = copyJSON(this.set_box.form)
            if (body.username == ""){
                return this.$message.warning("请输入用户名");
            }
            if (!body.role_ids || body.role_ids.length <= 0){
                return this.$message.warning("请选择角色");
            }
            body.sort = Number(body.sort)

            // 附加权限标记
            let path = "/user/set"
            if (body.id != this.user.id){
                path += "?auth_type=set"
            }

            api.innerPost(path, body, (res) =>{
                if (!res.status){
                    return this.$message.error(res.message)
                }
                this.setBox(false)
                this.getDetail()
                // api.dicList([Enum.dicUser],()=>{}, true) // 存在变化，更新缓存
            })
        },
        setBox(show=false){
            this.set_box.show = Boolean(show)
            if (this.set_box.show){
                this.set_box.form = {
                    id: this.detail.id,
                    username: this.detail.username,
                    mobile: this.detail.mobile,
                    sort: this.detail.sort,
                    role_ids: this.detail.role_ids,
                }
            }else{
                this.set_box.form = {}
            }
        },
        changeStatus(newStatus){
            this.$confirm(newStatus==1? '确认停用': '确认激活', '提示',{
                type: 'warning',
            }).then(()=>{
                // 确认操作
                api.innerPost("/user/change_status", {id:this.detail.id, status:Number(newStatus)}, (res)=>{
                    if (!res.status){
                        return this.$message.error(res.message)
                    }
                    return this.$message.success(res.message)
                    this.getDetail()
                })
            }).catch(()=>{
                // 取消操作
            })
        },
        passwordBox(show=false){
            this.password_box.show = Boolean(show)
            if (this.password_box.show){
                this.password_box.form = {
                    id: this.detail.id,
                    password: '',
                    re_password: ''
                }
            }else{
                this.password_box.form = {}
            }
        },
        passwordSubmit(){
            let body = {id:this.detail.id, password:this.password_box.form.password}
            if (this.password_box.form.password != this.password_box.form.re_password){
                return this.$message.warning("两次密码不一致")
            }
            api.innerPost("/user/change_password", body, (res)=>{
                if (!res.status){
                    return this.$message.error(res.message)
                }
                this.passwordBox(false)
                return this.$message.success(res.message)
            })
        },

        accountBox(show=false){
            this.account_box.show = Boolean(show)
            if (this.account_box.show){
                this.account_box.form = {
                    id: this.detail.id,
                    account: '',
                }
            }else{
                this.account_box.form = {}
            }
        },
        accountSubmit(){
            let body  = copyJSON(this.account_box.form)

            api.innerPost("/user/change_account", body, (res) =>{
                if (!res.status){
                    return this.$message.error(res.message)
                }
                this.accountBox(false)
                this.getDetail()
            })
        },
        getDic(){
            let types = [
                Enum.dicRole,
            ]
            api.dicList(types,(res) =>{
                this.dic_role = res[Enum.dicRole]
            })
        },

        close(){
            this.$emit('update:visible', false) // 向外传递关闭表示
        }
    }
})

Vue.component("MyUser", MyUser);