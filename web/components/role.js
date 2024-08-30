var MyRole = Vue.extend({
    template: `<el-row class="my-role">
<!--
        先规划一下页面布局（整体仿造tapd，还有我们自己的手机端），
        # 各个环境名称（如果没有对应任务就不在哪上）（默认只有第一个标签页打开）
            流水线和任务混合在一起（sql写法要研究一下）
                点击跳转到详情页（新的标签页）
-->
        <el-col :span="4" class="sidebar">
            <el-row>
                <span class="h3">角色列表</span>
                <el-button type="text" @click="setBoxOpen" style="margin-left: 25px;" v-if="$auth_tag.role_set">添加角色</el-button>
            </el-row>
            <el-row>
                <ul>
                    <li v-for="(role_v,role_i) in roleList" @click="selectedRole(role_v)" style="cursor: pointer;padding: 4px 6px;" :class="role_detail.id==role_v.id? 'active': ''">
                        {{role_v.name}}
                    </li>
                </ul>
            </el-row>
        </el-col>
        <el-col :span="20" class="body" :style="auth.style">
            <el-card shadow="never">
                <el-row>
                    <span style="margin-right: 20px">{{role_detail.name}}</span>
                    <el-button size="mini" @click="setBoxEdit()" v-if="$auth_tag.role_set"><i class="el-icon-edit"></i></el-button> 
                    <el-button size="mini" v-if="$auth_tag.role_set"><i class="el-icon-delete"></i></el-button>
                </el-row>
                <el-row>备注: {{role_detail.remark}}</el-row>
            </el-card>
            
            <el-tree
              :data="auth.list"
              show-checkbox
              node-key="id"
              default-expand-all
              ref="tree"
              @check-change="handleCheckChange"
              highlight-current
              :props="auth.props">
            </el-tree>
            <div class="footer">
                <el-button type="primary" @click="roleAuthSubmit()" size="medium" v-if="$auth_tag.auth_set">保 存</el-button>
            </div>
        </el-col>
        
        <el-dialog :title="setBox.title" :visible.sync="setBox.show" width="30%" :show-close="false" :close-on-click-modal="false">
            <el-form>
                <el-form-item label="名称">
                    <el-input v-model="setBox.form.name"></el-input>
                </el-form-item>
                <el-form-item label="备注">
                    <el-input v-model="setBox.form.remark"></el-input>
                </el-form-item>
            </el-form>
            <span slot="footer" class="dialog-footer">
                <el-button @click="setBox.show = false">取 消</el-button>
                <el-button type="primary" @click="setBoxSubmit()">确 定</el-button>
            </span>
        </el-dialog>
    </el-row>`,

    name: "MyRole",
    props: {
        data_id:Number
    },
    data(){
        return {
            group_list: [],
            // 角色盒子
            setBox:{
                show:false,
                title: "",
                form:{},
            },
            roleList:[],
            role_detail:{},
            // 节点信息
            auth:{
                list: [],
                props: {
                    children: 'children',
                    label: function(data, node){
                        return data.name
                    }
                },
                style: 'height: 500px',
                selected: {}, // 选中部分
            },

        }
    },
    // 模块初始化
    created(){
        setDocumentTitle('角色&权限')
    },
    // 模块初始化
    mounted(){
        this.auth.style = 'height: '+(window.innerHeight-128)+'px';
        this.getAuthList()
        this.getRoleList()
    },
    // 具体方法
    methods:{
        setBoxOpen(){
            this.setBox.show = true
            this.setBox.title = '新增'
        },
        setBoxEdit(){
            this.setBox.show = true
            this.setBox.title = '编辑'
            this.setBox.form = copyJSON(this.role_detail)
        },
        setBoxSubmit() {
            this.setBox.show = false
            let body = copyJSON(this.setBox.form)
            api.innerPost("/role/set", body, res => {
                if (!res.status){
                    return this.$message.error(res.message);
                }
                this.$message.success(this.setBox.title+' 成功')
                if (body.id){
                    this.role_detail.name = body.name
                    this.role_detail.remark = body.remark
                }
                this.getRoleList()
            })
        },
        // 获得角色列表
        getRoleList(){
            api.innerGet("/role/list", {}, res => {
                if (!res.status){
                    return this.$message.error(res.message);
                }
                this.roleList = res.data.list
                if (!this.role_detail.id && res.data.list.length>0){
                    this.selectedRole(res.data.list[0])
                }
            })
        },
        // 选中角色
        selectedRole(row){
            // console.log("角色选择", row)
            this.role_detail = copyJSON(row)
            this.$refs.tree.setCheckedKeys(this.role_detail.auth_ids); // 要确保dom渲染完成后，才能调用此方法
        },
        // 获得权限列表
        getAuthList(){
            api.innerGet("/role/auth_list", {}, res => {
                if (!res.status){
                    return this.$message.error(res.message);
                }
                let trace = arrayToTree(res.data.list, 'name', 'group', 'children')
                this.auth.list = trace
            })
        },
        /**
         * 选中权限
         * @param data object 数据行
         * @param checked bool 是否选择
         * @param indeterminate
         */
        handleCheckChange(data, checked, indeterminate){
            if (checked){
                this.auth.selected[data.id] = data.id
            }else{
                delete this.auth.selected[data.id]
            }
            // console.log("选中权限",data, checked, indeterminate, this.auth.selected);
        },
        // 权限提交
        roleAuthSubmit(){
            let body = {
                id: this.role_detail.id,
                auth_ids: [],
            }
            for (let i in this.auth.selected){
                body.auth_ids.push(Number(i))
            }

            api.innerPost("/role/auth_set", body, res =>{
                if (!res.status){
                    return this.$message.error(res.message);
                }
                this.$message.success('保存成功')
            })
        }
    }
})

Vue.component("MyRole", MyRole);