var MySetting = Vue.extend({
    template: `<el-main class="my-setting">
    <el-tabs tab-position="left" :style="setting.style" stretch>
        <el-tab-pane class="setting-wrap" :style="setting.body_style">
            <span slot="label" style="padding-right: 30px;"><i class="el-icon-collection-tag"></i> &nbsp; 偏好设置</span>
            <el-form :model="preference_form" label-width="100px" size="small">
                <el-card class="box-card" shadow="always">
                    <div slot="header" class="clearfix">
                        <span class="h3">流水线默认值</span>
                    </div>
                    <div class="text item">
                        <el-form-item label="执行间隔">
                            <el-input v-model="preference_form.pipeline.interval"><span slot="append">s/秒</span></el-input>
                        </el-form-item>
                        <el-form-item label="任务停用">
                            <el-tooltip class="item" effect="dark" content="存在停止、错误状态任务时流水线整体停止" placement="top-start">
                                <el-radio v-model="preference_form.pipeline.config_disable_action" label="1">停止</el-radio>
                            </el-tooltip>
                            <el-tooltip class="item" effect="dark" content="跳过停用、错误状态任务" placement="top-start">
                                <el-radio v-model="preference_form.pipeline.config_disable_action" label="2">跳过</el-radio>
                            </el-tooltip>
                            <el-tooltip class="item" effect="dark" content="执行停用、错误状态任务" placement="top-start">
                                <el-radio v-model="preference_form.pipeline.config_disable_action" label="3">执行</el-radio>
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
                            <el-table :data="preference_form.git.owner_repo">
                                <el-table-column prop="owner" label="空间" width="180">
                                    <template slot-scope="scope">
                                        <el-input v-model="scope.row.owner" size="small" @input="e=>inputChangeArrayPush(e,scope.$index,preference_form.git.owner_repo, {owner:'',repos:[{name:''}]})" placeholder="空间路径"></el-input>
                                    </template>
                                </el-table-column>
                                <el-table-column prop="repos" label="仓库" width="180">
                                    <template slot-scope="scope">
                                        <el-input v-for="(v, i) in scope.row.repos" v-model="v.name" size="small" @input="e=>inputChangeArrayPush(e,i,preference_form.git.owner_repo[scope.$index].repos, {name:''})" placeholder="仓库路径"></el-input>
                                    </template>
                                </el-table-column>
                            </el-table>
                        </el-form-item>
                        <el-form-item label="默认空间">
                            <el-input v-model="preference_form.git.owner" placeholder="空间路径默认值"></el-input>
                        </el-form-item>
                        <el-form-item label="默认仓库">
                            <el-input v-model="preference_form.git.repo" placeholder="仓库路径默认值"></el-input>
                        </el-form-item>
                        <el-form-item label="默认分支">
                            <el-input v-model="preference_form.git.branch" placeholder="仓库分支默认值，通常为 master"></el-input>
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
                            <el-select v-model="preference_form.other.config_select_type" placement="默认类型">
                                <el-option label="周期任务" value="1"></el-option>
                                <el-option label="单次任务" value="2"></el-option>
                                <el-option label="组件任务" value="5"></el-option>
                            </el-select>
                            <span class="info-2"> &nbsp; 任务选择弹窗默认标签页</span>
                        </el-form-item>
                    </div>
                </el-card>
                <div class="footer" v-if="auth_tags.preference_set">
                    <el-button type="primary" @click="preferenceSet" size="medium">保 存</el-button>
                </div>
            </el-form>
        </el-tab-pane>
        
        <el-tab-pane label="" class="setting-wrap" :style="setting.body_style">-</el-tab-pane>
        
    </el-tabs>
</el-main>`,

    name: "MySetting",
    data(){
        return {
            auth_tags:{},
            setting:{
                style:'',
                body_style:''
            },
            preference_form:{
                pipeline:{},
                git:{},
                other:{},
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
        this.preferenceGet()
        this.auth_tags = cache.getAuthTags()
    },

    // 具体方法
    methods:{
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
                this.preference_form = res.data
            })
        },
        // 偏好保存
        preferenceSet(){
            // console.log("偏好设置 提交", this.preference_form)
            let body = copyJSON(this.preference_form)
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
    }
})

Vue.component("MySetting", MySetting);