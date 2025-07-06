var MyConfigDetail = Vue.extend({
    template: `<el-container id="config-detail">
    <!--边栏-->
    <el-aside width="330px">
        <el-descriptions class="margin-top" :column="1" size="medium" :colon="false" :label-style="{color:'#909399'}">
            <el-descriptions-item label="状态">
                <div slot="content">{{detail.status_dt}}  {{detail.status_remark}}</div>
                <el-button :type="statusTypeName(detail.status)" plain size="mini" round  @click="statusShow">{{detail.status_name}}</el-button>
            </el-descriptions-item>
            <el-descriptions-item>
                <template slot="label">状态描述</template>
                <div style="font-size: 12px;color: #909399;">{{detail.status_dt}}<el-divider direction="vertical"></el-divider>{{detail.status_remark}}</div>
            </el-descriptions-item>
            <el-descriptions-item v-if="detail.type!=5 && req.type != 'receive'">
                <template slot="label">执行时间</template>
                {{detail.spec}}
            </el-descriptions-item>
            <el-descriptions-item>
                <template slot="label">协议</template>
                {{detail.protocol_name}}
            </el-descriptions-item>
            <el-descriptions-item v-if="req.type == 'config'">
                <template slot="label">类型</template>
                {{detail.type_name}}
            </el-descriptions-item>
            <el-descriptions-item>
                <template slot="label">处理人</template>
                {{detail.handle_user_names}}
            </el-descriptions-item>
            <el-descriptions-item>
                <template slot="label">更新时间</template>
                {{detail.update_dt}}
            </el-descriptions-item>
            <el-descriptions-item>
                <template slot="label">创建时间</template>
                {{detail.create_dt}}
            </el-descriptions-item>
            <el-descriptions-item>
                <template slot="label">创建人</template>
                {{detail.create_user_name}}
            </el-descriptions-item>
            <el-descriptions-item>
                <template slot="label">审核人</template>
                {{detail.audit_user_name}}
            </el-descriptions-item>
        </el-descriptions>
    </el-aside>
    <el-main>
        <div class="title">
            <div class="left">
                <span>{{detail.name}}</span> 
                <el-dropdown placement="bottom-start" size="mini">
                    <span class="el-dropdown-link">
                        <i slot="reference" class="el-icon-connection"></i>
                    </span>
                    <el-dropdown-menu slot="dropdown">
                        <el-dropdown-item @click.native="copyToClipboard('标题&链接',detail.share_url)">复制标题&链接</el-dropdown-item>
                        <el-dropdown-item @click.native="copyToClipboard('webhook 链接',detail.webhook_url)" v-if="req.type=='receive'">复制 webhook 链接</el-dropdown-item>
                    </el-dropdown-menu>
                </el-dropdown>
                <p v-if="detail.remark">{{detail.remark}}</p>
            </div>
            <div class="right">
                <el-button type="primary" plain size="small" @click="editOpen" v-if="detail && detail.id>0">编辑</el-button>
            </div>
        </div>
        <el-tabs class="detail-wrap" v-model="tbs_name" @tab-click="handleClickTypeLabel">
            <el-tab-pane label="详细信息" name="detail">
                <el-descriptions class="margin-top" :column="2" size="medium" border labelStyle="width:80px">
                    <el-descriptions-item span="2" v-if="detail.var_fields && detail.var_fields.length">
                        <template slot="label">参数</template>
                        (<el-tooltip v-for="(var_p,var_i) in detail.var_fields" effect="dark" :content="var_p.remark" placement="top-start">
                            <code v-if="var_p.key != ''" style="padding: 0 2px;margin: 0 4px;cursor:pointer;color: #445368;background: #f9fdff;position: relative;"><span style="position: absolute;left: -6px;bottom: -2px;">{{var_i>0 ? ',': ''}}</span>{{var_p.key}}<span class="info-2">={{var_p.value}}</span></code>
                        </el-tooltip>)
                    </el-descriptions-item>
                    
                    <el-descriptions-item span="2" v-if="detail.var_params">
                        <template slot="label">参数实现</template>
                        <el-input type="textarea" v-model="detail.var_params" autosize disabled></el-input>
                    </el-descriptions-item>
                    <el-descriptions-item span="2" v-if="detail.tag_ids && detail.tag_ids.length">
                        <template slot="label">标签</template>
                        <div>{{detail.tag_names}}</div>
                    </el-descriptions-item>
                    
                    <el-descriptions-item span="2">
                        <template slot="label">细节</template>
                        <el-form size="small" disabled class="detail">
                            <div v-if="detail.protocol_name == 'http'">
                                <el-form-item label="" class="http_url_box">
                                    <el-input v-model="detail.command.http.url" placeholder="请输入http:// 或 https:// 开头的完整地址">
                                        <el-select v-model="detail.command.http.method" placeholder="请选请求方式" slot="prepend">
                                            <el-option label="GET" value="GET"></el-option>
                                            <el-option label="POST" value="POST"></el-option>
                                        </el-select>
                                    </el-input>
                                </el-form-item>
                                <el-form-item label="请求Header" class="http_header_box" v-if="detail.command.http.header.length">
                                    <el-input v-for="(header_v,header_i) in detail.command.http.header" v-model="header_v.value" v-if="header_v.key">
                                        <el-input v-model="header_v.key" slot="prepend"></el-input>
                                    </el-input>
                                </el-form-item>
                                <el-form-item label="请求Body参数" v-if="detail.command.http.body != ''">
                                    <el-input type="textarea" v-model="detail.command.http.body" autosize></el-input>
                                </el-form-item>
                            </div>
                            
                            <div v-if="detail.protocol_name == 'rpc'">
                                <el-form-item label="proto" style="margin-top: -10px" v-if="detail.command.rpc.proto">
                                    <el-input type="textarea" v-model="detail.command.rpc.proto" autosize placeholder="请输入*.proto 文件内容" style=""></el-input>
                                </el-form-item>
                                <el-form-item label="" >
                                    <el-input v-model="detail.command.rpc.addr+' / '+ detail.command.rpc.action">
                                        <el-select v-model="detail.command.rpc.method" placeholder="请选请求方式" slot="prepend" style="width: 85px">
                                            <el-option label="RPC" value="RPC"></el-option>
                                            <el-option label="GRPC" value="GRPC"></el-option>
                                        </el-select>
                                    </el-input>
                                </el-form-item>
                                <el-form-item label="请求参数">
                                    <el-input type="textarea" v-model="detail.command.rpc.body" autosize placeholder="请输入请求参数"></el-input>
                                </el-form-item>             
                            </div>
                            
                            <div v-if="detail.protocol_name == 'cmd'">
                                <p>主机<b class="b">{{detail.command.cmd.host.id==-1?'本机': dic.host_source[detail.command.cmd.host.id]}}</b>，来源 <b class="b">{{detail.command.cmd.origin}}</b>，类型 <b class="b">{{detail.command.cmd.type}}</code></p>
                                <el-input type="textarea" v-model="detail.command.cmd.statement.local" autosize></el-input>
                            </div>
                            
                            <div v-if="detail.protocol_name == 'sql'">
                                <p>驱动<b class="b">{{detail.command.sql.driver}}</b> 链接<b class="b">{{getEnumName(dic.sql_source, detail.command.sql.source.id)}}</b> 执行<b class="b">{{detail.command.sql.origin}}</b>语句：</p>
                                <div class="sql-show-warp">
                                    <!--本地sql展示-->
                                    <div class="input-box" v-if="detail.command.sql.origin == 'local'" v-for="(statement,sql_index) in detail.command.sql.statement" v-show="statement.type=='local' || statement.type==''">
                                        <pre><code>{{statement.local}}</code></pre>
                                        <span class="input-header"><i style="font-size: 16px;">{{sql_index}}</i></span>
                                    </div>
                                    <!--git sql展示-->
                                    <div class="input-box" v-if="detail.command.sql.origin == 'git'" v-for="(statement,sql_index) in detail.command.sql.statement" v-show="statement.type=='git'">
                                        <el-row v-html="statement.git.descrition"></el-row>
                                        <i style="position: absolute;right: 1px;top: 40px;font-size: 16px;">#{{sql_index}}</i>
                                    </div>
                                </div>
                                <el-row>错误时 <b class="b">{{detail.command.sql.err_action==1? '终止任务' : detail.command.sql.err_action==2? '跳过继续' : '事务回滚'}}</b> 执行间隔 <b class="b">{{detail.command.sql.interval}}/秒</b></el-row>
                            </div>
                            
                            <div v-if="detail.protocol_name == 'jenkins'">
                                <div>
                                    资源<b class="b">{{getEnumName(dic.jenkins_source, detail.command.jenkins.source.id)}}</b>
                                    项目<b class="b">{{detail.command.jenkins.name}}</b>
                                    参数：
                                </div>
                                <div>
                                    <b class="b" v-for="(header_v,header_i) in detail.command.jenkins.params" v-if="header_v.key">{{header_v.key}}<el-divider direction="vertical"></el-divider>{{header_v.value}}</b>
                                </div>
                            </div>
                            
                            <div v-if="detail.protocol_name == 'git'">
                                <p>链接 <b class="b">{{getEnumName(dic.git_source, detail.command.git.link_id)}}</b> 事件</p>
                                <div style="overflow-y: auto;max-height: 420px;">
                                    <div class="input-box" v-for="(event,git_index) in detail.command.git.events">
                                        <el-row v-html="event.desc" style="min-height: 45px;"></el-row>
                                        <span class="input-header">
                                            <i style="font-size: 16px;">{{git_index}}</i>
                                        </span>
                                        
                                    </div>
                                </div>
                            </div>
                            
                            <div v-if="detail.protocol_name == 'pipeline'">
                                <p>任务执行间隔<b class="b">{{detail.interval}}/秒</b> 任务停用时<b class="b">{{detail.config_disable_action_name}}</b> 任务执行错误时<b class="b">{{detail.config_err_action_name}}</b></p>
                                <div class="sort-drag-box">
                                    <div class="input-box" v-for="(conf,conf_index) in detail.configs" @mouseover="configDetailPanel(conf, true)"  @mouseout="configDetailPanel(conf, false)">
                                        <b class="b">{{conf.type_name}}</b>-<b>{{conf.name}}</b>-<b class="b">{{conf.protocol_name}}</b>
                                        <el-button :type="statusTypeName(conf.status)" size="mini" plain round disabled>{{conf.status_name}}</el-button>
                                        <el-divider direction="vertical"></el-divider>
                                        (<el-tooltip v-for="(var_p,var_i) in conf.var_fields" effect="dark" :content="var_p.remark" placement="top-start">
                                            <code v-if="var_p.key != ''" style="padding: 0 2px;margin: 0 4px;cursor:pointer;color: #445368;background: #f9fdff;position: relative;"><span style="position: absolute;left: -6px;bottom: -2px;">{{var_i>0 ? ',': ''}}</span>{{var_p.key}}<span class="info-2">={{var_p.value}}</span></code>
                                        </el-tooltip>)
                                        <span style="margin-left: 4px;" v-show="conf.view_panel">
                                            <i class="el-icon-view hover" @click="configDetailBox(conf)"></i>
                                        </span>
                                    </div>
                                </div>
                            </div>
                            
                            <div v-if="detail.protocol_name == 'receive'">
                                <p>webhook: <b class="b">{{detail.webhook_url}}</b> <span class="info-2">激活后有效</span></p>
                                <div class="input-box">
                                    <pre style="min-height: 200px;"><code>{{detail.receive_tmpl}}</code></pre>
                                </div>
                                <p>任务执行间隔<b class="b">{{detail.interval}}/秒</b> 任务停用时<b class="b">{{detail.config_disable_action_name}}</b> 任务执行错误时<b class="b">{{detail.config_err_action_name}}</b></p>
                                <div class="sort-drag-box">
                                    <div class="input-box" v-for="(conf,conf_index) in detail.rule_config" @mouseover="configDetailPanel(conf, true)" @mouseout="configDetailPanel(conf, false)">
                                        <div>
                                            <b class="b">{{conf.config.type_name}}</b>-<b>{{conf.config.name}}</b>-<b class="b">{{conf.config.protocol_name}}</b>
                                            <el-button :type="statusTypeName(conf.config.status)" size="mini" plain round disabled>{{conf.config.status_name}}</el-button>
                                            <el-divider direction="vertical"></el-divider>
                                            (<el-tooltip v-for="(var_p,var_i) in conf.config.var_fields" effect="dark" :content="var_p.remark" placement="top-start">
                                                <code v-if="var_p.key != ''" style="padding: 0 2px;margin: 0 4px;cursor:pointer;color: #445368;background: #f9fdff;position: relative;"><span style="position: absolute;left: -6px;bottom: -2px;">{{var_i>0 ? ',': ''}}</span>{{var_p.key}}<span class="info-2">={{var_p.value}}</span></code>
                                            </el-tooltip>)
                                            <span style="margin-left: 4px;" v-show="conf.view_panel">
                                                <i class="el-icon-view hover" @click="configDetailBox(conf.config)" title="查看任务信息"></i>
                                            </span>
                                        </div>
                                        <el-row style="margin:3px 0 0 8px">
                                            <el-col :span="12">任务：<b class="b" v-for="re in conf.rule">{{re.key}}:{{re.value}}</b></el-col>
                                            <el-col :span="12">入参：<b class="b" v-for="pm in conf.param">{{pm.key}}:{{pm.value}}</b></el-col>
                                        </el-row>
                                    </div>
                                </div>
                            </div>
                        </el-form>
                    </el-descriptions-item>
                    <el-descriptions-item span="2" v-if="detail.var_fields && detail.var_fields.length">
                        <template slot="label">消息</template>
                        <div>
                            <el-checkbox v-model="detail.empty_not_msg==1" disabled v-show="detail.empty_not_msg==1">空结果不发消息</el-checkbox>
                            <el-row class="input-box" v-for="(msg,msg_index) in detail.msg_set" v-html="msg.descrition" style="padding: 2px 4px;"></el-row>
                        </div>
                    </el-descriptions-item>
                </el-descriptions>
            </el-tab-pane>
            <el-tab-pane label="变更历史" name="change_log">
                <el-table :data="change_logs.list" stripe style="width: 100%" :span-method="changeLogsSpan">
                    <el-table-column prop="create_dt" label="变更时间" width="180"></el-table-column>
                    <el-table-column prop="create_user_name" label="变更人"></el-table-column>
                    <el-table-column prop="type_name" label="变更类型" width="80"></el-table-column>
                    <el-table-column prop="field" label="变更字段">
                        <template slot-scope="scope">
                            <table class="sub-table">
                                <tr v-for="field in scope.row.content"><td width="">{{field.field_name}}</td><td>{{field.old_val_name==''?'-':field.old_val_name}}</td><td>{{field.new_val_name==''?'-':field.new_val_name}}</td></tr>
                            </table>
                        </template>
                    </el-table-column>
                    <el-table-column label="变更前"></el-table-column>
                    <el-table-column label="变更后"></el-table-column>
                </el-table>
                <el-pagination
                    @current-change="getChangeLogList"
                    :current-page.sync="change_logs.search.page"
                    :page-size="change_logs.search.size"
                    layout="total, prev, pager, next"
                    :total="change_logs.page.total">
                </el-pagination>
            </el-tab-pane>
        </el-tabs>
        
        <el-row>
            <h4>日志</h4>
            <span><el-button size="mini" round class="el-icon-search button-icon" @click="logSearchBox"></el-button><span class="help">{{logs.form_desc}}</span></span>
            <my-config-log :search="logs.search" v-if="logs.show"></my-config-log>
        </el-row>
        
        <el-dialog title="编辑任务" :visible.sync="detail_form_box.show && req.type=='config'" :close-on-click-modal="false" class="config-form-box" :before-close="formClose">
            <my-config-form v-if="detail_form_box.show && req.type==='config'" :request="{detail:detail}" @close="formClose"></my-config-form>
        </el-dialog>
        <el-drawer title="编辑流水线" :visible.sync="detail_form_box.show && req.type=='pipeline'" size="60%" :before-close="formClose">
            <my-pipeline-form v-if="detail_form_box.show && req.type==='pipeline'" :request="{detail:detail}" @close="formClose"></my-pipeline-form>
        </el-drawer>
        <el-drawer title="编辑接收" :visible.sync="detail_form_box.show && req.type=='receive'" size="60%" :before-close="formClose">
            <my-receive-form v-if="detail_form_box.show && req.type==='receive'" :request="{detail:detail}" @close="formClose"></my-receive-form>
        </el-drawer>
        <el-dialog title="任务详情" :visible.sync="config_detail.show" :close-on-click-modal="false" class="config-form-box" :modal="false">
            <my-config-form v-if="config_detail.show" :request="{detail:config_detail.detail,disabled:true}" @close="configDetailBox()"></my-config-form>
        </el-dialog>
        <my-status-change v-if="status_box.show" :request="status_box" @close="statusShow"></my-status-change>
        
        <el-dialog title="搜索" :visible.sync="logs.form_show" width="320px" class="log_search_box">
            <el-form :model="logs.form" label-position="top" size="small">
                <el-form-item label="开始时间">
                    <el-date-picker style="width:100%"
                            format="yyyy-MM-dd HH:mm:ss"
                            value-format="yyyy-MM-dd HH:mm:ss"
                            v-model="logs.form.timestamp_start"
                            type="datetime"
                            placeholder="开始时间">
                    </el-date-picker>
                </el-form-item>
                <el-form-item label="截止时间">
                    <el-date-picker style="width:100%"
                            format="yyyy-MM-dd HH:mm:ss"
                            value-format="yyyy-MM-dd HH:mm:ss"
                            v-model="logs.form.timestamp_end"
                            type="datetime"
                            placeholder="截止时间">
                    </el-date-picker>
                </el-form-item>
                <el-form-item label="">
                    <el-row type="flex" justify="space-between">
                        <div class="left" style="margin-right: 5px">
                            <label class="el-form-item__label">最小耗时</label>
                            <el-input type="number" placeholder=">= s.秒" v-model="logs.form.duration_start"></el-input>
                        </div>
                        <div class="right"  style="margin-left: 5px">
                            <label class="el-form-item__label">最大耗时</label>
                            <el-input type="number" placeholder="<= s.秒" v-model="logs.form.duration_end"></el-input>
                        </div>
                    </el-row>
                </el-form-item>
                <el-form-item label="状态">
                    <el-select v-model="logs.form.status" placeholder="结果状态" style="width:100%" clearable>
                        <el-option label="成功" value="2"></el-option>
                        <el-option label="失败" value="1"></el-option>
                    </el-select>
                </el-form-item>
            </el-form>
            <div slot="footer" class="dialog-footer">
                <el-button @click="logs.form_show = false" size="small">取 消</el-button>
                <el-button type="primary" @click="logSearch()" size="small">确 定</el-button>
            </div>
        </el-dialog>
    </el-main>
</el-container>`,

    name: "MyConfigDetail",
    props: {
        request:{
            id:Number,
            type:String,
            entry_id:Number,
        }
    },
    data(){
        return {
            req:{id:0,type:'config',entry_id:0},
            tbs_name: 'detail',
            detail:{},
            // 流水线详情展示任务细节使用
            config_detail:{
                show:false,
                detail:{}
            },
            detail_form_box:{
                show:false,
                form:{},
            },
            // 执行日志
            logs:{
                show: false,
                form_show:false,
                form_desc: '',
                form:{
                    timestamp_start:'',
                    timestamp_end:'',
                    duration_start:'',
                    duration_end:'',
                    status:'',
                },
                search:{
                    timestamp_start:'',
                    timestamp_end:'',
                    duration_start:'',
                    duration_end:'',
                    status:'',
                    tags:null,
                },
            },
            // 变更日志
            change_logs:{
                load:false,
                search:{
                    page: 1,
                    size: 10,
                },
                list:[],
                page:{}
            },
            dic:{},
            status_box:{
                show:false,
                detail:{},
                type: '',
            },
            // 合并类型
            gitMergeTypeList: {
                "merge":"合并分支",
                "squash":"扁平化分支",
                "rebase":"变基并合并"
            },
        }
    },
    // 模块初始化
    created(){
        // console.log('config_detail', this.$route, this.$router)
        if (this.$route.query.id){
            this.req.id = Number(this.$route.query.id)
        }
        if (this.$route.query.type){
            this.req.type = this.$route.query.type.toString()
        }
        if (this.$route.query.entry_id){
            this.req.entry_id = Number(this.$route.query.entry_id)
        }
        // this.logs.search.tags =  JSON.stringify({
        //     ref_id: this.req.id,
        //     component:this.req.type
        // })
        let time = new Date()
        time.setDate(time.getDate()-7)
        this.logs.form.timestamp_start = getDatetimeString(time)
        this.logSearch()
        this.logs.search.ref_id = this.req.id
        this.logs.search.operation = 'job-task'
        if (this.req.type == 'pipeline'){
            this.logs.search.operation = 'job-pipeline'
        }else if (this.req.type == 'receive'){
            this.logs.search.operation = 'job-receive'
        }

        this.change_logs.search.ref_id= this.req.id
        this.change_logs.search.ref_type= this.req.type
    },
    // 模块初始化
    mounted(){
        this.getDetail()
        this.getDicList()
    },
    // 具体方法
    methods:{
        handleClickTypeLabel(tab, event) {
            this.tab_name = tab.name
            if (this.tab_name === 'change_log'){
                if (!this.change_logs.load){
                    this.change_logs.load = true
                    this.getChangeLogList()
                }
            }
        },
        getDetail(){
            if (!this.req.id){
                return
            }
            api.innerGet('/'+this.req.type+'/detail', {id: this.req.id}, (res)=>{
                if (!res.status){
                    return this.$message({
                        message: res.message,
                        type: 'error',
                        duration: 6000
                    })
                }
                if (this.req.type === 'config'){
                    if (res.data.command.sql.statement){
                        res.data.command.sql.statement.forEach((item)=>{
                            item = this.sqlGitBuildDesc(item)
                        })
                    }
                    if (res.data.command.git){
                        for (let i in res.data.command.git.events){
                            res.data.command.git.events[i] = this.gitBuildDesc(res.data.command.git.events[i])
                        }
                    }
                }else if (this.req.type === 'pipeline'){
                    res.data.protocol_name = 'pipeline'
                    res.data.configs.forEach((item2)=>{
                        item2.view_panel = false
                    })
                }else if (this.req.type === 'receive'){
                    res.data.protocol_name = this.req.type
                    res.data.rule_config.forEach((item2)=>{
                        item2.view_panel = false
                    })
                    res.data.webhook_url = window.location.origin+"/receive/webhook/"+(res.data.alias != '' ? res.data.alias : res.data.id)
                }
                if (res.data.msg_set){
                    res.data.msg_set.forEach((item)=>{
                        item = this.msgSetBuildDesc(item)
                    })
                }
                res.data.handle_user_names = ''
                if (res.data.handle_user_ids){
                    res.data.handle_user_ids.forEach((id)=>{
                        this.dic.user.forEach((item)=>{
                            if (item.id == id){
                                res.data.handle_user_names += item.name+','
                            }
                        })
                    })
                    res.data.handle_user_names = res.data.handle_user_names.substring(0,res.data.handle_user_names.length-1)
                }
                res.data.share_url = '【'+res.data.name+'】 '+ window.location.href

                this.detail = res.data
                setDocumentTitle(res.data.name)
                if (!this.logs.show){
                    this.logs.show = true
                    this.logs.search.env = res.data.env
                }
            })
        },
        // 获取变更日志
        getChangeLogList(page=1){
            this.change_logs.search.page = Number(page)
            api.innerGet('/change_log/list', this.change_logs.search, res=>{
                if (!res.status){
                    return this.$message.error(res.message)
                }
                this.change_logs.list = res.data.list
                this.change_logs.page = res.data.page
            })
        },
        changeLogsSpan({ row, column, rowIndex, columnIndex }){
            console.log("snap",row, column, rowIndex, columnIndex)
            if (column.property == "field"){
                return [1,3]
            }
        },
        logSearchBox(){
            this.logs.form_show = true
            let search = copyJSON(this.logs.search)
            Object.assign(this.logs.form, search)
        },
        // 日志搜索
        logSearch(){
            console.log("日志搜索")
            // this.logs.show = true
            this.logs.form_show = false
            let form = copyJSON(this.logs.form)
            let desc = ""
            if (form.timestamp_start || form.timestamp_end){
                desc += '时间：'+form.timestamp_start +' ~ '+ form.timestamp_end +'; '
            }
            if (form.duration_start || form.duration_end){
                desc += '耗时：'+form.duration_start+' ~ ' + form.duration_end +'秒; '
            }
            if (form.status){
                desc += '状态：'+(form.status==1? '失败' : '成功')
            }

            let search = copyJSON(this.logs.search)
            Object.assign(search, form)
            this.logs.search = search
            this.logs.form_desc = desc
        },
        // 枚举
        getDicList(){
            let types = [
                Enum.dicSqlSource,
                Enum.dicSqlDriver,
                Enum.dicJenkinsSource,
                Enum.dicGitSource,
                Enum.dicGitEvent,
                Enum.dicHostSource,
                Enum.dicCmdType,
                Enum.dicUser,
                Enum.dicMsg
            ]
            api.dicList(types,(res) =>{
                // this.dic_sql_source = res[Enum.dicSqlSource]
                // this.dic_jenkins_source = res[Enum.dicJenkinsSource]
                // this.dic_git_source =res[Enum.dicGitSource]
                // this.dic_git_event = res[Enum.dicGitEvent]
                // this.dic_user = res[Enum.dicUser]
                // this.dic_msg = res[Enum.dicMsg]
                // this.dic_sql_driver = res[Enum.dicSqlDriver]
                this.dic = {
                    host_source: res[Enum.dicHostSource],
                    cmd_type: res[Enum.dicCmdType],
                    sql_source: res[Enum.dicSqlSource],
                    git_source: res[Enum.dicGitSource],
                    jenkins_source: res[Enum.dicJenkinsSource],
                    msg: res[Enum.dicMsg],
                    user: res[Enum.dicUser],
                }
            })
        },
        statusShow(e=null){
            if (!this.detail.id || this.detail.id < 0){
                return
            }
            this.status_box.show = e == null || e.type == 'click'
            this.status_box.detail = this.detail
            this.status_box.type = this.req.type
            if (typeof e == 'object' && e.is_change){ // 数据变更了重新加载
                this.getDetail()
            }
        },
        // 构建消息设置描述
        msgSetBuildDesc(data){
            let statusList = [{id:1,name:"错误"}, {id:2, name:"成功"}, {id:0,name:"完成"}]
            let item1 = statusList.filter((option) => {
                return data.status.includes(option.id)
            }).map((item)=>{return item.name});
            data.status_name = item1.join(',')

            let descrition = '<i class="el-icon-bell"></i>当任务<b class="b">'+data.status_name+'</b>时'

            let item2 = this.dic.msg.find(option => option.id === data.msg_id)
            if (item2){
                data.msg_name = item2.name
                descrition += '，发送<b class="b">'+item2.name+'</b>消息'
            }
            let item3 = this.dic.user.filter((option) => {
                return data.notify_user_ids.includes(option.id);
            }).map((item)=>{return item.name})
            if (item3.length > 0){
                data.notify_users_name = item3.join(',')
                descrition += '，并且@人员<b class="b">'+data.notify_users_name+'</b>'
            }
            data.descrition = descrition
            return data
        },
        // 构建git sql 描述
        sqlGitBuildDesc(data){
            if (data.git == null){
                return data
            }
            let git = this.dic.git_source.filter(item=>{
                return item.id == data.git.link_id
            })

            data.git.descrition = '连接<span class="el-tag el-tag--small el-tag--light">'+git.length>0 ? git[0].name: '' +
                '</span> 访问<span class="el-tag el-tag--small el-tag--light">'+data.git.owner+'/'+data.git.project +'</span>'+
                '</span> 引用<span class="el-tag el-tag--small el-tag--light">'+data.git.ref +'</span> 拉取以下文件内容'+
                '<span class="el-tag el-tag--small el-tag--light">'+ (data.is_batch==1? '批量解析后' :'单文件单sql') +'</span>执行';

            if (data.git.path){
                data.git.path.forEach(function (item) {
                    if (item == ''){
                        return
                    }
                    item.split(',').forEach(function (item2) {
                        data.git.descrition += `<div style="margin: 0;padding: 0 0 0 30px;"><a href="https://gitee.com/${data.git.owner}/${data.git.project}/blob/${data.git.ref}/${item2}" target="_blank" title="点击 查看文件详情"><i class="el-icon-connection"></i></a><code>${item2}</code></div>`
                    })
                })
            }else{

            }

            console.log("sql描述",data.git.descrition)
            return data
        },
        // 构建 git 描述
        gitBuildDesc(data){
            switch (data.id){
                case 2:
                    data.desc = '完善中...'
                    break
                case 9:
                    data.desc = `<b>pr合并</b> <a href="https://gitee.com/${data.pr_merge.owner}/${data.pr_merge.repo}/pulls/${data.pr_merge.number}" target="_blank" title="点击 查看pr详情"><i class="el-icon-connection"></i></a><b class="b">${data.pr_merge.owner}/${data.pr_merge.repo}</b>/pulls/<b class="b">${data.pr_merge.number}</b> <b class="b">${this.gitMergeTypeList[data.pr_merge.merge_method]}</b>  ${data.pr_merge.prune_source_branch===true?'<b class="b">删除提交分支</b>':''}`+
                        `<br><i style="margin-left: 3em;"></i><b>${data.pr_merge.title}</b> ${data.pr_merge.description}`
                    break
                case 131:
                    data.desc = `<b>文件更新</b> <a href="https://gitee.com/${data.file_update.owner}/${data.file_update.repo}/blob/${data.file_update.branch}/${data.file_update.path}" target="_blank" title="点击 查看"><i class="el-icon-connection"></i></a> <b class="b">${data.file_update.owner}</b>/<b class="b">${data.file_update.repo}</b>/blob/<b class="b">${data.file_update.branch}</b>/<b class="b">${data.file_update.path}</b>`+
                        `<br><i style="margin-left: 2em;"></i>内容：<b class="b">${data.file_update.content}</b>`+
                        `<br><i style="margin-left: 2em;"></i>描述：<b class="b">${data.file_update.message}</b>`
                    break
                default:
                    data.desc = '未支持的事件类型'
            }
            return data
        },
        configDetailPanel(detail,show=false){
            detail.view_panel = show === true
        },
        // 流水线中的任务详情展示
        configDetailBox(detail= null){
            if (!detail){
                this.config_detail.show =false
                this.config_detail.detail = {}
                return
            }
            api.innerGet('/config/detail', {id: detail.id, var_params: this.detail.var_params}, (res)=>{
                if (!res.status){
                    return this.$message.error(res.message)
                }
                this.config_detail.show = true
                this.config_detail.detail = res.data
            },{async:false})
        },
        editOpen(){
            this.detail_form_box.show= true
        },
        formClose(e){
            console.log("close",e)
            this.detail_form_box.show = false
            if (e.is_change){
                this.getDetail()
            }
        }
    }
})

Vue.component("MyConfigDetail", MyConfigDetail);