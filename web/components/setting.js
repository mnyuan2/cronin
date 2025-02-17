var MySetting = Vue.extend({
    template: `<el-main class="my-setting">
    <el-tabs tab-position="left" v-model="setting.tab_name" :style="setting.style" stretch @tab-click="showTab">
        <el-tab-pane name="preference" class="setting-wrap" :style="setting.body_style">
            <span slot="label" style="padding-right: 30px;"><i class="el-icon-collection-tag"></i> &nbsp; 偏好设置</span>
            <el-form :model="preference.set.form" label-width="100px" size="small">
                <el-card class="box-card" shadow="always">
                    <div slot="header" class="clearfix">
                        <span class="h3">流水线默认值</span>
                    </div>
                    <div class="text item">
                        <el-form-item label="执行间隔">
                            <el-input v-model="preference.set.form.pipeline.interval"><span slot="append">s/秒</span></el-input>
                        </el-form-item>
                        <el-form-item label="任务停用">
                            <el-tooltip class="item" effect="dark" content="存在停止、错误状态任务时流水线整体停止" placement="top-start">
                                <el-radio v-model="preference.set.form.pipeline.config_disable_action" label="1">停止</el-radio>
                            </el-tooltip>
                            <el-tooltip class="item" effect="dark" content="跳过停用、错误状态任务" placement="top-start">
                                <el-radio v-model="preference.set.form.pipeline.config_disable_action" label="2">跳过</el-radio>
                            </el-tooltip>
                            <el-tooltip class="item" effect="dark" content="执行停用、错误状态任务" placement="top-start">
                                <el-radio v-model="preference.set.form.pipeline.config_disable_action" label="3">执行</el-radio>
                            </el-tooltip>
                        </el-form-item>
                    </div>
                </el-card>
                <br>
                <el-card class="box-card" shadow="always">
                    <div slot="header" class="clearfix">
                        <span class="h3">git 选项</span>
                    </div>
                    <div class="text item">
                        <el-form-item>
                            <span slot="label" style="white-space: nowrap;">
                                空间与仓库<br/>选项
                            </span>
                            <el-table :data="preference.set.form.git.owner_repo">
                                <el-table-column prop="owner" label="空间" width="180">
                                    <template slot-scope="scope">
                                        <el-input v-model="scope.row.owner" size="small" @input="e=>inputChangeArrayPush(e,scope.$index,preference.set.form.git.owner_repo, {owner:'',repos:[{name:''}]})" placeholder="空间路径"></el-input>
                                    </template>
                                </el-table-column>
                                <el-table-column prop="repos" label="仓库" width="180">
                                    <template slot-scope="scope">
                                        <el-input v-for="(v, i) in scope.row.repos" v-model="v.name" size="small" @input="e=>inputChangeArrayPush(e,i,preference.set.form.git.owner_repo[scope.$index].repos, {name:''})" placeholder="仓库路径"></el-input>
                                    </template>
                                </el-table-column>
                            </el-table>
                        </el-form-item>
                        <el-form-item label="默认空间">
                            <el-input v-model="preference.set.form.git.owner" placeholder="空间路径默认值"></el-input>
                        </el-form-item>
                        <el-form-item label="默认仓库">
                            <el-input v-model="preference.set.form.git.repo" placeholder="仓库路径默认值"></el-input>
                        </el-form-item>
                        <el-form-item label="默认分支">
                            <el-input v-model="preference.set.form.git.branch" placeholder="仓库分支默认值，通常为 master"></el-input>
                        </el-form-item>
                    </div>
                </el-card>
                <br>
                <el-card class="box-card" shadow="always">
                    <div slot="header" class="clearfix">
                        <span class="h3">杂项</span>
                    </div>
                    <div class="text item">
                        <el-form-item label="任务选择类型">
                            <el-select v-model="preference.set.form.other.config_select_type" placement="默认类型">
                                <el-option label="周期任务" value="1"></el-option>
                                <el-option label="单次任务" value="2"></el-option>
                                <el-option label="组件任务" value="5"></el-option>
                            </el-select>
                            <span class="info-2"> &nbsp; 任务选择弹窗默认标签页</span>
                        </el-form-item>
                    </div>
                </el-card>
                <div class="footer" v-if="$auth_tag.preference_set">
                    <el-button type="primary" @click="preferenceSet" size="medium">保 存</el-button>
                </div>
            </el-form>
        </el-tab-pane>
        
        <el-tab-pane name="global_variate" class="setting-wrap" :style="setting.body_style">
            <span slot="label" style="padding-right: 30px;"><i class="el-icon-key"></i> &nbsp; 全局变量</span>
            <el-card shadow="never" body-style="padding:12px">
                <el-button size="mini" round plain @click="tabChangeName()" class="active" disabled>用户变量</el-button>
                <el-button type="text" size="mini" @click="globalVariateSet()" v-if="$auth_tag.global_variate_set" style="float: right">添加变量</el-button>
            </el-card>
            
            <el-table :data="global_variate.list" style="width: 100%;margin-top: 15px;">
                <el-table-column  prop="name" label="名称"></el-table-column>
                <el-table-column prop="value" label="值"></el-table-column>
                <el-table-column prop="remark" label="描述"></el-table-column>
                <el-table-column prop="status_name" label="状态">
                    <template slot-scope="{row}">
                        <i v-show="row.register==1" class="el-badge__content is-dot el-badge__content--success" title="已注册"></i>
                        <i v-show="row.register==2" class="el-badge__content is-dot el-badge__content--info" title="未注册"></i>
                        {{row.status_name}}
                    </template>
                </el-table-column>
                <el-table-column label="">
                    <template slot-scope="{row}">
                        <el-button type="text" size="small" @click="globalVariateSet(row)" v-show="$auth_tag.global_variate_set && row.status==Enum.StatusDisable">编辑</el-button>
                        <el-button type="text" size="small" @click="globalVariateStatus(row, Enum.StatusActive,'激活')" v-show="$auth_tag.global_variate_status && row.status==Enum.StatusDisable">启用</el-button>
                        <el-button type="text" size="small" @click="globalVariateStatus(row, Enum.StatusDisable,'停用')" v-show="$auth_tag.global_variate_status && row.status==Enum.StatusActive">停用</el-button>
                        <el-button type="text" size="small" class="info-2" @click="globalVariateStatus(row, Enum.StatusDelete,'删除')" v-show="$auth_tag.global_variate_status && row.status==Enum.StatusDisable">删除</el-button>
                    </template>
                </el-table-column>
            </el-table>
            
            <!--设置弹窗-->
            <el-dialog title="全局变量设置" :visible.sync="global_variate.set.show" :close-on-click-modal="false" append-to-body="true" width="420px">
                <el-form :model="global_variate.set.form"  label-width="50px" size="small">
                    <el-form-item label="名称*">
                        <el-input v-model="global_variate.set.form.name" placeholder="变量名称，示例：MY_KEY"></el-input>
                        <p></p>
                    </el-form-item>
                    <el-form-item label="值">
                        <el-input type="textarea" :autosize="{minRows:1}" v-model="global_variate.set.form.value" placeholder="名称值"></el-input>
                    </el-form-item>
                    <el-form-item label="描述">
                        <el-input v-model="global_variate.set.form.remark" placeholder="名称的说明描述"></el-input>
                    </el-form-item>
                </el-form>
                <div slot="footer" class="dialog-footer" style="margin-top: -30px">
                    <el-button size="small" @click="global_variate.set.show = false">取 消</el-button>
                    <el-button size="small" type="primary" @click="globalVariateSubmit()">确 定</el-button>
                </div>
            </el-dialog>
            
        </el-tab-pane>
        
    </el-tabs>
</el-main>`,

    name: "MySetting",
    data(){
        return {
            setting:{
                style:'',
                body_style:'',
                tab_name:'preference',
            },
            preference:{
                load: false,
                set:{
                    form:{
                        pipeline:{},
                        git:{},
                        other:{},
                    },
                }
            },
            // 全局变量
            global_variate:{
                load: false,
                list:[],
                set:{
                    show: false,
                    form:{}
                }
            }
        }
    },
    // 模块初始化
    created(){
        setDocumentTitle('设置')
    },
    // 模块初始化
    mounted(){
        this.setting.style = 'height: '+(window.innerHeight-90)+'px';
        this.setting.body_style = 'height: '+(window.innerHeight-150)+'px';
        this.showTab({name:this.setting.tab_name})
    },

    // 具体方法
    methods:{
        showTab(e){
            if (e.name === 'global_variate'){
                if (!this.global_variate.load){
                    this.global_variate.load = true
                    this.globalVariateGet()
                }
            }else if (e.name === 'preference'){
                if (!this.preference.load){
                    this.preference.load = true
                    this.preferenceGet()
                }
            }
        },
        preferenceGet(){
            api.innerGet("/setting/preference_get", null, res =>{
                if (!res.status){
                    return this.$message.error(res.message);
                }
                // 数据转换
                res.data.pipeline.config_disable_action = res.data.pipeline.config_disable_action.toString()
                res.data.other.config_select_type = res.data.other.config_select_type.toString()
                res.data.git.owner_repo.forEach(function (item) {
                    if (item.repos.length === 0 || item.repos[item.repos.length-1].name !== ''){
                        item.repos.push({name:''})
                    }
                })
                if (res.data.git.owner_repo.length ===0 || res.data.git.owner_repo[res.data.git.owner_repo.length-1].owner !== ''){
                    res.data.git.owner_repo.push({owner:"",repos:[{name:''}]})
                }
                this.preference.set.form = res.data
            })
        },
        // 偏好保存
        preferenceSet(){
            // console.log("偏好设置 提交", this.preference)
            let body = copyJSON(this.preference.set.form)
            body.pipeline.interval = Number(body.pipeline.interval)
            body.pipeline.config_disable_action = Number(body.pipeline.config_disable_action)
            body.other.config_select_type = Number(body.other.config_select_type)
            body.git.owner_repo = body.git.owner_repo.filter(function (item) {
                item.repos = item.repos.filter(function (item2){
                    return item2.name !== ''
                })
                item.owner = item.owner.trim()
                return item.owner !== ""
            })

            api.innerPost("/setting/preference_set", body, res =>{
                if (!res.status){
                    return this.$message.error(res.message);
                }
                this.$message.success('保存成功')
            })
        },

        // 全局变量
        globalVariateGet(){
            api.innerGet("/global_variate/list", null, res =>{
                if (!res.status){
                    return this.$message.error(res.message);
                }
                this.global_variate.list = res.data.list
            })
        },
        // 全局变量 弹窗
        globalVariateSet(row=null){
            this.global_variate.set.show = true
            if (row == null){
                this.global_variate.set.form = {
                    id: 0,
                    name: '',
                    value: '',
                    remark: '',
                }
            }else{
                this.global_variate.set.form = copyJSON(row)
            }
        },
        // 全局变量 弹窗 保存
        globalVariateSubmit(){
            if (this.global_variate.set.form.name === ''){
                return this.$message.error('名称不能为空')
            }
            api.innerPost("/global_variate/set", this.global_variate.set.form, res =>{
                if (!res.status){
                    return this.$message.error(res.message);
                }
                this.$message.success('保存成功')
                this.globalVariateGet()
                this.global_variate.set.show = false
            })
        },
        // 全局变量 删除
        globalVariateStatus(row, status, status_name){
            this.$confirm('是否执行'+status_name, '提示', {
                confirmButtonText: '确定',
                cancelButtonText: '取消',
                type: 'warning'
            }).then(() => {
                api.innerPost("/global_variate/change_status", {id:row.id, status:status}, res =>{
                    if (!res.status){
                        return this.$message.error(res.message);
                    }
                    this.$message.success('操作成功')
                    this.globalVariateGet()
                })
            }).catch(() => {
                this.$message({
                    type: 'info',
                    message: '操作已取消'
                });
            });
        }
    }
})

Vue.component("MySetting", MySetting);