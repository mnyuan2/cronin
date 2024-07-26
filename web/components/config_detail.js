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
            <el-descriptions-item>
                <template slot="label">执行时间</template>
                {{detail.spec}}
            </el-descriptions-item>
            <el-descriptions-item>
                <template slot="label">协议</template>
                {{detail.protocol_name}}
            </el-descriptions-item>
            <el-descriptions-item>
                <template slot="label">类型</template>
                {{detail.type_name}}
            </el-descriptions-item>
            <el-descriptions-item  span="3">
                <template slot="label">消息</template>
                <div>
                   <el-row v-for="(msg,msg_index) in detail.msg_set" v-html="msg.descrition" style="background: rgb(248, 248, 249);margin-bottom: 4px;border-radius: 3px;padding: 2px 4px;"></el-row>
                </div>
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
                        <el-input type="textarea" v-model="detail.var_params" rows="3" disabled></el-input>
                    </el-descriptions-item>
                    
                    <el-descriptions-item span="2">
                        <template slot="label">细节</template>
                        <el-form size="mini" disabled class="detail">
                            <div v-if="detail.protocol_name == 'http'">
                                    <el-form-item label="" class="http_url_box" size="mini">
                                        <el-input v-model="detail.command.http.url" placeholder="请输入http:// 或 https:// 开头的完整地址">
                                            <el-select v-model="detail.command.http.method" placeholder="请选请求方式" slot="prepend">
                                                <el-option label="GET" value="GET"></el-option>
                                                <el-option label="POST" value="POST"></el-option>
                                            </el-select>
                                        </el-input>
                                    </el-form-item>
                                    <el-form-item label="请求Header" class="http_header_box" size="mini" v-if="detail.command.http.header">
                                        <el-input v-for="(header_v,header_i) in detail.command.http.header" v-model="header_v.value" v-if="header_v.key">
                                            <el-input v-model="header_v.key" slot="prepend"></el-input>
                                        </el-input>
                                    </el-form-item>
                                    <el-form-item label="请求Body参数" size="mini" v-if="detail.command.http.body != ''">
                                        <el-input type="textarea" v-model="detail.command.http.body" rows="4"></el-input>
                                    </el-form-item>
                            </div>
                            
                            <div v-if="detail.protocol_name == 'rpc'">
                                <el-form-item label="proto" style="margin-top: -10px" v-if="detail.command.rpc.proto">
                                    <el-input type="textarea" v-model="detail.command.rpc.proto" rows="3" placeholder="请输入*.proto 文件内容" style=""></el-input>
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
                                    <el-input type="textarea" v-model="detail.command.rpc.body" rows="3" placeholder="请输入请求参数"></el-input>
                                </el-form-item>             
                            </div>
                            
                            <div v-if="detail.protocol_name == 'cmd'">
                                <p>主机<b class="b">{{detail.command.cmd.host.id==-1?'本机': dic.host_source[detail.command.cmd.host.id]}}</b>，来源 <b class="b">{{detail.command.cmd.origin}}</b>，类型 <b class="b">{{detail.command.cmd.type}}</code></p>
                                <el-input type="textarea" v-model="detail.command.cmd.statement.local" ></el-input>
                            </div>
                            
                            <div v-if="detail.protocol_name == 'sql'">
                                <el-row>驱动<b class="b">{{detail.command.sql.driver}}</b> 链接<b class="b">{{getEnumName(dic.sql_source, detail.command.sql.source.id)}}</b> 执行<b class="b">{{detail.command.sql.origin}}</b>语句：</el-row>
                                <div class="sql-show-warp">
                                    <!--本地sql展示-->
                                    <div v-if="detail.command.sql.origin == 'local'" v-for="(statement,sql_index) in detail.command.sql.statement" v-show="statement.type=='local' || statement.type==''" style="position: relative;line-height: 133%;background: #f4f4f5;margin-bottom: 5px;padding: 6px 20px 7px 8px;border-radius: 3px;">
                                        <pre style="margin: 0;overflow-y: auto;max-height: 180px;min-height: 45px;"><code class="language-sql hljs">{{statement.local}}</code></pre>
                                        <i style="position: absolute;right: 1px;top: 40px;font-size: 16px;">#{{sql_index}}</i>
                                    </div>
                                    <!--git sql展示-->
                                    <div v-if="detail.command.sql.origin == 'git'" v-for="(statement,sql_index) in detail.command.sql.statement" v-show="statement.type=='git'" style="position: relative;line-height: 133%;background: #f4f4f5;margin-bottom: 5px;padding: 6px 20px 7px 8px;border-radius: 3px;">
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
                                <el-row>链接 <b class="b">{{getEnumName(dic.git_source, detail.command.git.link_id)}}</b> 事件</el-row>
                                <div style="overflow-y: auto;max-height: 420px;">
                                    <div v-for="(event,git_index) in detail.command.git.events" style="position: relative;line-height: 133%;background: #f4f4f5;margin-bottom: 5px;padding: 6px 20px 7px 8px;border-radius: 3px;">
                                        <el-row v-html="event.desc" style="min-height: 45px;"></el-row>
                                        <i style="position: absolute;right: 1px;top: 40px;font-size: 16px;">#{{git_index}}</i>
                                    </div>
                                </div>
                            </div>
                            
                            <div v-if="detail.protocol_name == 'pipeline'">
                                <div class="sort-drag-box">
                                    <div class="item-drag" 
                                        v-for="(conf,conf_index) in detail.configs" 
                                        style="position: relative;max-height: 200px;line-height: 133%;background: #f4f4f5;margin-bottom: 10px;padding: 6px 20px 7px 8px;border-radius: 3px;"
                                        @mouseover="configDetailPanel(conf, true)"
                                        @mouseout="configDetailPanel(conf, false)"
                                    >
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
                        </el-form>
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
                                <tr v-for="field in scope.row.content"><td width="">{{field.field_name}}</td><td>{{field.old_val_name}}</td><td>{{field.new_val_name}}</td></tr>
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
            <h3>执行日志</h3>
            <my-config-log :tags="logs.tags"></my-config-log>
        </el-row>
        
        <el-dialog title="编辑任务" :visible.sync="detail_form_box.show && req.type=='config'" :close-on-click-modal="false" class="config-form-box" :before-close="formClose">
            <my-config-form v-if="detail_form_box.show && req.type==='config'" :request="{detail:detail}" @close="formClose"></my-config-form>
        </el-dialog>
        <el-drawer title="编辑流水线" :visible.sync="detail_form_box.show && req.type=='pipeline'" size="60%" :before-close="formClose">
            <my-pipeline-form v-if="detail_form_box.show && req.type==='pipeline'" :request="{detail:detail}" @close="formClose"></my-pipeline-form>
        </el-drawer>
        <el-dialog title="任务详情" :visible.sync="config_detail.show" :close-on-click-modal="false" class="config-form-box" :modal="false">
            <my-config-form v-if="config_detail.show" :request="{detail:config_detail.detail,disabled:true}" @close="configDetailBox()"></my-config-form>
        </el-dialog>
        <my-status-change v-if="status_box.show" :request="status_box" @close="statusShow"></my-status-change>
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
                tags:{}
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
        this.logs.tags = {ref_id: this.req.id, component:this.req.type}
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
                res.data.share_url = '【'+res.data.name+'】'+ window.location.href

                this.detail = res.data
                setDocumentTitle(res.data.name)
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
            let item1 = statusList.find(option => option.id === data.status);
            if (item1){
                data.status_name = item1.name
            }
            let descrition = '当任务<b class="b">'+item1.name+'</b>时'

            let item2 = this.dic.msg.find(option => option.id === data.msg_id)
            if (item2){
                data.msg_name = item2.name
                descrition += '，发送<b class="b">'+item2.name+'</b>消息'
            }
            let item3 = this.dic.user.filter((option) => {
                return data.notify_user_ids.includes(option.id);
            }).map((item)=>{return item.name})
            if (item3.length > 0){
                data.notify_users_name = item3
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
                    data.git.descrition += `<div style="margin: 0;padding: 0 0 0 30px;"><a href="https://gitee.com/${data.git.owner}/${data.git.project}/blob/${data.git.ref}/${item}" target="_blank" title="点击 查看文件详情"><i class="el-icon-connection"></i></a><code>${item}</code></div>`
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
                default:
                    data.desc = '未支持的事件类型'
            }
            return data
        },
        configDetailPanel(detail,show=false){
            detail.view_panel = show === true
        },
        // 流水线中的任务详情展示
        configDetailBox(detail=null){
            if (!detail){
                this.config_detail.show =false
                this.config_detail.detail = {}
                return
            }
            this.config_detail.show = true
            this.config_detail.detail = detail
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