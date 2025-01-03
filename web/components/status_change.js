var MyStatusChange = Vue.extend({
    template: `<el-dialog :title="info.name" :visible.sync="request.show" :show-close="false" class="status-change-warp" width="540px">
        <el-tabs tab-position="left">
            <el-tab-pane v-if="status_show.audited">
                <span slot="label" ><el-button plain size="mini" round>待审核</el-botton></span>
                <el-form ref="form" :model="form" label-width="100px" size="small">
                    <el-form-item label="处理人">
                        <el-select v-model="form.handle_user_ids" multiple style="width:100%">
                            <el-option v-for="(dic_v,dic_k) in dic.user" :key="dic_v.id" :label="dic_v.name" :value="dic_v.id"></el-option>
                        </el-select>
                    </el-form-item>
                    <el-form-item label="备注">
                        <el-input type="textarea" v-model="form.status_remark" rows="5"></el-input>
                    </el-form-item>
                    <el-form-item class="status-change-footer">
                        <el-button @click="close()">取消</el-button>
                        <el-button type="primary" @click="statusSubmit(Enum.StatusAudited)">确定</el-button>
                    </el-form-item>
                </el-form>
            </el-tab-pane>
            <el-tab-pane v-if="status_show.reject">
                <span slot="label"><el-button plain size="mini" round>驳回</el-botton></span>
                <el-form ref="form" :model="form" label-width="100px" size="small">
                    <el-form-item label="处理人">
                        <el-select v-model="form.handle_user_ids" multiple style="width:100%">
                            <el-option v-for="(dic_v,dic_k) in dic.user" :key="dic_v.id" :label="dic_v.name" :value="dic_v.id"></el-option>
                        </el-select>
                    </el-form-item>
                    <el-form-item label="审核备注">
                        <el-input type="textarea" v-model="form.status_remark" rows="5"></el-input>
                    </el-form-item>
                    <el-form-item class="status-change-footer">
                        <el-button @click="close()">取消</el-button>
                        <el-button type="primary" @click="statusSubmit(Enum.StatusReject)">确定</el-button>
                    </el-form-item>
                </el-form>
            </el-tab-pane>
            <el-tab-pane v-if="status_show.active">
                <span slot="label"><el-button plain size="mini" round>{{info.status==Enum.StatusAudited ?'通过':'激活'}}</el-botton></span>
                <el-form ref="form" :model="form" label-width="100px" size="small">
                    <el-form-item label="处理人">
                        <el-select v-model="form.handle_user_ids" multiple style="width:100%">
                            <el-option v-for="(dic_v,dic_k) in dic.user" :key="dic_v.id" :label="dic_v.name" :value="dic_v.id"></el-option>
                        </el-select>
                    </el-form-item>
                    <el-form-item label="审核备注">
                        <el-input type="textarea" v-model="form.status_remark" rows="5"></el-input>
                    </el-form-item>
                    <el-form-item class="status-change-footer">
                        <el-button @click="close()">取消</el-button>
                        <el-button type="primary" @click="statusSubmit(Enum.StatusActive)">确定</el-button>
                    </el-form-item>
                </el-form>
            </el-tab-pane>
            <el-tab-pane v-if="status_show.disable">
                <span slot="label"><el-button plain size="mini" round>{{info.status==Enum.StatusClosed ?'草稿':'停用'}}</el-botton></span>
                <el-form ref="form" :model="form" label-width="100px" size="small">
                    <el-form-item label="处理人">
                        <el-select v-model="form.handle_user_ids" multiple style="width:100%">
                            <el-option v-for="(dic_v,dic_k) in dic.user" :key="dic_v.id" :label="dic_v.name" :value="dic_v.id"></el-option>
                        </el-select>
                    </el-form-item>
                    <el-form-item label="备注">
                        <el-input type="textarea" v-model="form.status_remark" rows="5"></el-input>
                    </el-form-item>
                    <el-form-item class="status-change-footer">
                        <el-button @click="close()">取消</el-button>
                        <el-button type="primary" @click="statusSubmit(Enum.StatusDisable)">确定</el-button>
                    </el-form-item>
                </el-form>
            </el-tab-pane>
            <el-tab-pane v-if="status_show.closed">
                <span slot="label"><el-button plain size="mini" round>已关闭&删除</el-botton></span>
                <el-form ref="form" :model="form" label-width="100px" size="small">
                    <el-form-item label="备注">
                        <el-input type="textarea" v-model="form.status_remark" rows="5"></el-input>
                    </el-form-item>
                    <el-form-item class="status-change-footer">
                        <el-button @click="close()">取消</el-button>
                        <el-button type="primary" @click="statusSubmit(Enum.StatusClosed)">确定</el-button>
                    </el-form-item>
                </el-form>
            </el-tab-pane>
        </el-tabs>
    </el-dialog>`,

    // StatusDisable: 1, // 草稿
    // StatusAudited: 5, // 待审核
    // StatusReject: 6, // 驳回
    // StatusActive: 2, // 激活
    // StatusFinish: 3, // 完成
    // StatusError: 4, // 错误
    // StatusDelete: 9, // 删除

    name: "MyStatusChange",
    props: {
        request:{
            show:Boolean,
            detail:Object,
            type: String,
        }
    },
    data(){
        return {
            type: '',
            tab_name: '',
            form:{
                id:0,
                status: '',
                handle_user_ids:[],
                status_remark: '',
            },
            info:{status:0},
            dic:{
                user:[]
            },
            status_show:{
                audited: false,
                reject: false,
                active: false,
                disable: false,
                closed: false,
            }
        }
    },
    // 模块初始化
    created(){
    },
    // 模块初始化
    mounted(){
        if (this.request.detail){
            this.info = {
                id: this.request.detail.id,
                status: this.request.detail.status,
                name: this.request.detail.name,
            }
        }
        if (this.request.type){
            this.type = this.request.type
        }
        this.getDicSqlSource()
        this.initStatusShow()
        console.log("状态变更 start",this.request, this.$auth_tag)
        window.copyToClipboard = this.copyToClipboard; // 事件绑定在window对象中后才可以被触发
    },

    // 具体方法
    methods:{
        close(is_change=false){
            this.$emit('close', {is_change:is_change})
        },
        // 枚举
        getDicSqlSource(){
            let types = [
                // Enum.dicSqlSource,
                // Enum.dicSqlDriver,
                // Enum.dicJenkinsSource,
                // Enum.dicGitSource,
                // Enum.dicGitEvent,
                // Enum.dicHostSource,
                // Enum.dicCmdType,
                Enum.dicUser,
                // Enum.dicMsg
            ]
            api.dicList(types,(res) =>{
                // this.dic.sql_source = res[Enum.dicSqlSource]
                // this.dic.jenkins_source = res[Enum.dicJenkinsSource]
                // this.dic.git_source =res[Enum.dicGitSource]
                // this.dic.git_event = res[Enum.dicGitEvent]
                // this.dic.host_source =res[Enum.dicHostSource]
                this.dic.user = res[Enum.dicUser]
                // this.dic.msg = res[Enum.dicMsg]
                // this.dic.cmd_type = res[Enum.dicCmdType]
                // this.dic.sql_driver = res[Enum.dicSqlDriver]
            })
        },
        statusSubmit(status){
            let sucMsg = {message:'操作成功',type: 'success'}
            let body = copyJSON(this.form)
            body.status = status
            body.id = this.info.id
            let path = ''
            if (this.type == 'config'){
                path = '/config/change_status'
            }else if (this.type == 'pipeline'){
                path = '/pipeline/change_status'
            }else if (this.type == 'receive'){
                path = '/receive/change_status'
            }else{
                return this.$message.warning('业务类型错误！')
            }
            if (body.status === Enum.StatusActive || body.status === Enum.StatusReject){
                path += '?auth_type=audit'
            }else if (body.status === Enum.StatusAudited){
                let url = '【'+this.info.name+'】'+ window.location.protocol + "//"+window.location.host+'/index?env='+cache.getEnv().env+'#/config_detail?id='+body.id+'&type='+this.type
                sucMsg.dangerouslyUseHTMLString = true
                sucMsg.message += `, <a href="javascript:;" onclick="copyToClipboard('标题&链接','${url}')">复制标题&链接</a>`
            }

            api.innerPost(path, body, res=>{
                if (!res.status){
                    return this.$message.error(res.message)
                }
                this.$message(sucMsg)
                this.close(true)
            })
        },
        initStatusShow(){
            if (this.info.status==Enum.StatusDisable || this.info.status==Enum.StatusReject || this.info.status==Enum.StatusFinish){
                if ((this.$auth_tag.config_status && this.type=='config') ||
                    (this.$auth_tag.pipeline_status && this.type=='pipeline') ||
                    (this.$auth_tag.receive_status && this.type=='receive')){
                    this.status_show.audited = true // 待审核 需要具有状态操作权限
                }
            }
            if (this.info.status==Enum.StatusAudited){
                if ((this.$auth_tag.config_audit && this.type=='config') ||
                    (this.$auth_tag.pipeline_audit && this.type=='pipeline') ||
                    (this.$auth_tag.receive_audit && this.type=='receive')){
                    this.status_show.reject = true // 驳回  需要具有审核权限
                }
            }
            if (this.info.status==Enum.StatusDisable || this.info.status==Enum.StatusReject || this.info.status==Enum.StatusFinish || this.info.status==Enum.StatusError || this.info.status==Enum.StatusAudited){
                if ((this.$auth_tag.config_audit && this.type=='config') ||
                    (this.$auth_tag.pipeline_audit && this.type=='pipeline') ||
                    (this.$auth_tag.receive_audit && this.type=='receive')){
                    this.status_show.active = true // 激活/通过  需要具有审核权限
                }
            }
            if (this.info.status==Enum.StatusActive || this.info.status==Enum.StatusClosed){
                if ((this.$auth_tag.config_status && this.type=='config') ||
                    (this.$auth_tag.pipeline_status && this.type=='pipeline') ||
                    (this.$auth_tag.receive_status && this.type=='receive')){
                    this.status_show.disable = true // 草稿/停用 需要具有状态操作权限
                }
            }
            if (this.info.status==Enum.StatusDisable || this.info.status==Enum.StatusReject || this.info.status==Enum.StatusFinish || this.info.status==Enum.StatusError){
                if ((this.$auth_tag.config_status && this.type=='config') ||
                    (this.$auth_tag.pipeline_status && this.type=='pipeline') ||
                    (this.$auth_tag.receive_status && this.type=='receive')){
                    this.status_show.closed = true // 预删除锁定，可以转为草稿后进行编辑
                }
            }
        }
    }
})

Vue.component("MyStatusChange", MyStatusChange);