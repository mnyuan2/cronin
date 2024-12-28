var MyConfigForm = Vue.extend({
    template: `<div class="config-form">
    <el-form :model="form" :disabled="request.disabled" size="small">
        <el-form-item label="名称*" label-width="50px" style="padding-right: 96px">
            <el-input v-model="form.name" style="margin-right: 80px"></el-input>
        </el-form-item>
        <el-form-item label="类型*" label-width="50px">
            <el-radio v-model="form.type" label="1">周期</el-radio>
            <el-radio v-model="form.type" label="2">单次</el-radio>
            <el-radio v-model="form.type" label="5">组件</el-radio>
        </el-form-item>
        <el-form-item label="时间*" label-width="50px" v-if="form.type != 5">
            <el-input v-model="form.spec" :placeholder="hintSpec" style="width: 220px" v-show="form.type==1"></el-input>
            <el-button @click="parseSpec" v-show="form.type==1">检测</el-button>
            <el-date-picker 
                style="width: 220px"
                v-show="form.type==2" 
                v-model="form.spec" 
                value-format="yyyy-MM-dd HH:mm:ss"
                type="datetime" 
                placeholder="选择运行时间" 
                :picker-options="pickerOptions">
            </el-date-picker>
        </el-form-item>
        <el-form-item label-width="50px" v-show="form.var_fields.length" style="padding-right: 96px">
            <span slot="label" style="white-space: nowrap;">
                参数
            </span>
            <div class="input-box">
                <span class="input-header"><i class="el-icon-edit" @click="biggerShow(true)"></i></span>
                (<el-tooltip v-for="(var_p,var_i) in form.var_fields" effect="dark" :content="var_p.remark" placement="top-start">
                    <code v-if="var_p.key != ''" style="padding: 0 2px;margin: 0 4px;cursor:pointer;color: #445368;background: #f9fdff;position: relative;"><span style="position: absolute;left: -6px;bottom: -2px;">{{var_i>0 ? ',': ''}}</span>{{var_p.key}}<span class="info-2">={{var_p.value}}</span></code>
                </el-tooltip>)
            </div>
            
        </el-form-item>
        <el-form-item>
            <el-button type="text" @click="biggerShow(true)" style="position: absolute;top: -36px;right: 0px">更多设置<i class="el-icon-plus"></i></el-button>
            <el-tabs type="border-card" v-model="form.protocol">
                <el-tab-pane label="http" name="1">
                    <el-form-item label="请求地址">
                        <el-input class="input-input" v-model="form.command.http.url" placeholder="请输入http:// 或 https:// 开头的完整地址">
                            <el-select v-model="form.command.http.method" placeholder="method" slot="prepend" style="width: 70px;">
                                <el-option label="GET" value="GET"></el-option>
                                <el-option label="POST" value="POST"></el-option>
                                <el-option label="PUT" value="PUT"></el-option>
                                <el-option label="DELETE" value="DELETE"></el-option>
                            </el-select>
                        </el-input>
                    </el-form-item>
                    <el-form-item label="请求Header" class="http_header_box">
                        <el-input class="input-input" v-for="(header_v,header_i) in form.command.http.header" v-model="header_v.value" placeholder="参数值">
                            <el-input v-model="header_v.key" slot="prepend" placeholder="参数名" @input="httpHeaderInput"></el-input>
                            <el-button slot="append" icon="el-icon-delete" @click="httpHeaderDel(header_i)"></el-button>
                        </el-input>
                    </el-form-item>
                    <el-form-item label="请求Body参数">
                        <el-input type="textarea" v-model="form.command.http.body" :autosize="{minRows:5}" placeholder="POST请求时body参数，将通过json进行请求发起"></el-input>
                    </el-form-item>
                </el-tab-pane>

                <el-tab-pane label="rpc" name="2">
                    <el-form-item label="请求模式">
                        <el-radio v-model="form.command.rpc.method" label="GRPC">GRPC</el-radio>
                    </el-form-item>
                    <el-form-item label="proto" style="margin-top: -10px">
                        <el-input type="textarea" v-model="form.command.rpc.proto" :autosize="{minRows:5}" placeholder="请输入*.proto 文件内容" style=""></el-input>
                    </el-form-item>
                    <el-form-item label="地址" label-width="42px">
                        <el-input v-model="form.command.rpc.addr" placeholder="请输入服务地址，含端口; 示例：localhost:21014"></el-input>
                    </el-form-item>
                    <el-form-item label="方法">
                        <el-select v-model="form.command.rpc.action" filterable clearable placeholder="填写proto信息后点击解析，可获得其方法进行选择" style="min-width: 400px">
                            <el-option v-for="item in form.command.rpc.actions" :key="item" :label="item" :value="item"></el-option>
                        </el-select>
                        <el-button @click="parseProto()">解析proto</el-button>
                    </el-form-item>
                    <el-form-item label="请求参数">
                        <el-input type="textarea" v-model="form.command.rpc.body" :autosize="{minRows:5}" placeholder="请输入请求参数"></el-input>
                    </el-form-item>
                </el-tab-pane>
                
                <el-tab-pane label="cmd" name="3">
                    <el-row style="margin: 6px 0;">
                        主机
                        <el-select v-model="form.command.cmd.host.id" style="width:180px">
                            <el-option v-for="dic_v in dic.host_source" :label="dic_v.name" :value="dic_v.id"></el-option>
                        </el-select>
                        来源
                        <el-select v-model="form.command.cmd.origin" style="width:80px">
                            <el-option label="当前" value="local"></el-option>
                            <el-option label="git" value="git"></el-option>
                        </el-select>
                        类型
                        <el-select v-model="form.command.cmd.type" clearable style="width:85px">
                            <el-option v-for="dic_v in dic.cmd_type" :label="dic_v.name" :value="dic_v.key"></el-option>
                        </el-select>
                    </el-row>
                    <el-form :model="form.command.cmd.statement.git" label-width="70px" size="small" v-if="form.command.cmd.origin=='git'">
                        <el-form-item label="连接">
                            <el-select v-model="form.command.cmd.statement.git.link_id" placement="请选择git链接">
                                <el-option v-for="(dic_v,dic_k) in dic.git_source" :label="dic_v.name" :value="dic_v.id"></el-option>
                            </el-select>
                        </el-form-item>
                        <el-form-item label="空间">
                            <el-input v-model="form.command.cmd.statement.git.owner" placeholder="仓库所属空间地址(企业、组织或个人的地址path)"></el-input>
                        </el-form-item>
                        <el-form-item label="仓库">
                            <el-input v-model="form.command.cmd.statement.git.project" placeholder="仓库路径"></el-input>
                        </el-form-item>
                        <el-form-item label="文件路径">
                            <el-input v-for="(path_v,path_i) in form.command.cmd.statement.git.path" v-model="form.command.cmd.statement.git.path[path_i]" placeholder="文件的路径">
                            </el-input>
                        </el-form-item>
                        <el-form-item label="引用">
                            <el-input v-model="form.command.cmd.statement.git.ref" placeholder="分支、tag或commit。默认: 仓库的默认分支(通常是master)"></el-input>
                        </el-form-item>
                    </el-form>
                    <el-input type="textarea" v-model="form.command.cmd.statement.local" :autosize="{minRows:5}" placeholder="请输入命令行执行内容" v-if="form.command.cmd.origin=='local'"></el-input>
                </el-tab-pane>
                
                <el-tab-pane label="sql" name="4" label-position="left">
                    <el-form-item label="驱动">
                        <el-select v-model="form.command.sql.driver" placeholder="驱动">
                            <el-option v-for="dic_v in dic.sql_driver" :label="dic_v.name" :value="dic_v.key"></el-option>
                        </el-select>
                    </el-form-item>
                    <el-form-item label="链接">
                        <el-select v-model="form.command.sql.source.id" placement="请选择sql链接">
                            <el-option v-for="(dic_v,dic_k) in dic.sql_source" v-show="dic_v.extend.driver == form.command.sql.driver" :label="dic_v.name" :value="dic_v.id"></el-option>
                        </el-select>
                        <el-button type="text" style="margin-left: 20px" @click="sqlSourceBox(true)">设置链接</el-button>
                    </el-form-item>
                    <el-form-item label="执行语句">
                        <div>
                            来源
                            <el-select v-model="form.command.sql.origin" style="width:80px">
                                <el-option label="本地" value="local"></el-option>
                                <el-option label="git" value="git"></el-option>
                            </el-select>
                            <el-button type="text" @click="sqlSetShow(-1)">添加<i class="el-icon-plus"></i></el-button>
                        </div>
                        <div class="sql-show-warp">
                            <!--本地sql展示-->
                            <div class="input-box" v-if="form.command.sql.origin == 'local'" v-for="(statement,sql_index) in form.command.sql.statement" v-show="statement.type=='local' || statement.type==''">
                                <pre><code>{{statement.local}}</code></pre>
                                <span class="input-header">
                                    <i class="el-icon-close" @click="sqlSetDel(sql_index)"></i>
                                    <i class="el-icon-edit" @click="sqlSetShow(sql_index,statement)"></i>
                                    <i style="white-space: nowrap;">{{sql_index}}</i>
                                </span>
                            </div>
                            <el-alert v-show="form.command.sql.statement.length==0" title="未添加执行sql，请添加。" type="info"></el-alert>
                            <!--git sql展示-->
                            <div class="input-box" v-if="form.command.sql.origin == 'git'" v-for="(statement,sql_index) in form.command.sql.statement" v-show="statement.type=='git'">
                                <el-row v-html="statement.git.descrition"></el-row>
                                <span class="input-header">
                                    <i class="el-icon-close" @click="sqlSetDel(sql_index)"></i>
                                    <i class="el-icon-edit" @click="sqlSetShow(sql_index,statement)"></i>
                                    <i style="white-space: nowrap;">{{sql_index}}</i>
                                </span>
                                
                            </div>
                        </div>
                    </el-form-item>
                    <el-form-item label="错误行为">
                        <el-tooltip class="item" effect="dark" content="执行到错误停止，之前的成功语句保留" placement="top-start">
                          <el-radio v-model="form.command.sql.err_action" label="1">终止任务</el-radio>
                        </el-tooltip>
                        <el-tooltip class="item" effect="dark" content="执行到错误跳过，继续执行后面语句" placement="top-start">
                          <el-radio v-model="form.command.sql.err_action" label="2">跳过继续</el-radio>
                        </el-tooltip>
                        <el-tooltip class="item" effect="dark" content="执行到错误，之前的成功语句也整体回滚" placement="top-start">
                          <el-radio v-model="form.command.sql.err_action" label="3">事务回滚</el-radio>
                        </el-tooltip>
                    </el-form-item>
                    <el-form-item label="执行间隔" v-if="form.command.sql.err_action !=3" label-width="69px">
                        <el-input type="number" v-model="form.command.sql.interval" placeholder="s秒">
                            <span slot="append">s/秒</span>
                        </el-input>
                    </el-form-item>
                </el-tab-pane>
                
                <el-tab-pane label="jenkins" name="5">
                    <el-form-item label="链接">
                        <el-select v-model="form.command.jenkins.source.id" placement="请选择链接">
                            <el-option v-for="(dic_v,dic_k) in dic.jenkins_source" :label="dic_v.name" :value="dic_v.id"></el-option>
                        </el-select>
                    </el-form-item>
                    <el-form-item label="项目">
                        <el-input v-model="form.command.jenkins.name" placeholder="jenkins job name"></el-input>
                    </el-form-item>
                    <el-form-item label="参数" class="http_header_box">
                        <el-input class="input-input" v-for="(header_v,header_i) in form.command.jenkins.params" v-model="header_v.value" placeholder="参数值">
                            <el-input v-model="header_v.key" slot="prepend" placeholder="参数名" @input="jenkinsParamInput"></el-input>
                            <el-button slot="append" icon="el-icon-delete" @click="jenkinsParamDel(header_i)"></el-button>
                        </el-input>
                    </el-form-item>
                </el-tab-pane>
                
                <el-tab-pane label="git" name="6">
                    <el-form-item label="链接">
                        <el-select v-model="form.command.git.link_id" placement="请选择链接">
                            <el-option v-for="(dic_v,dic_k) in dic.git_source" :label="dic_v.name" :value="dic_v.id"></el-option>
                        </el-select>
                    </el-form-item>
                    <el-form-item label="事件">
                        <div>
                            <el-button type="text" @click="gitSetShow(-1)">添加<i class="el-icon-plus"></i></el-button>
                        </div>
                        <div style="overflow-y: auto;max-height: 420px;">
                            <el-alert v-show="form.command.git.events.length==0" title="未添加事件，请添加。" type="info"></el-alert>
                            <!--git 事件展示-->
                            <div class="input-box" v-for="(event,git_index) in form.command.git.events">
                                <el-row class="input-body" v-html="event.desc" style="min-height: 45px;max-height: 200px;"></el-row>
                                <span class="input-header">
                                    <i class="el-icon-close" @click="gitSetDel(git_index)"></i>
                                    <i class="el-icon-edit" @click="gitSetShow(git_index,event)"></i>
                                    <i>{{git_index}}</i>
                                </span>
                            </div>
                        </div>
                    </el-form-item>
                </el-tab-pane>
            </el-tabs>
        </el-form-item>
        
        <el-form-item label-width="50px" v-show="form.after_tmpl">
            <span slot="label" style="white-space: nowrap;">
                结果<el-tooltip effect="dark" content="任务结果模板分析，点击查看更多" placement="top-start">
                    <router-link target="_blank" to="/var_params" style="color: #606266"><i class="el-icon-info"></i></router-link>
                </el-tooltip>
            </span>
            <div class="input-box">
                <span class="input-header"><i class="el-icon-edit" @click="biggerShow(true)"></i></span>
                <pre>{{form.after_tmpl}}</pre>
            </div>
        </el-form-item>
        <el-form-item label="重试" label-width="50px" v-show="form.err_retry_num>0">
            <div class="input-box">
                <span class="input-header"><i class="el-icon-edit" @click="biggerShow(true)"></i></span>
                {{form.err_retry_num}}/次，
                {{form.err_retry_mode_name}} {{form.err_retry_sleep}}/秒
            </div>
        </el-form-item>
        <el-form-item label="延迟" label-width="50px" v-show="form.after_sleep>0">
            <div class="input-box">
                <span class="input-header"><i class="el-icon-edit" @click="biggerShow(true)"></i></span>
                {{form.after_sleep}}/秒 后关闭
            </div>
        </el-form-item>
        <el-form-item label="备注" label-width="50px" v-show="form.remark">
            <div class="input-box">
                <span class="input-header"><i class="el-icon-edit" @click="biggerShow(true)"></i></span>
                {{form.remark}}
            </div>
        </el-form-item>
        <el-form-item label="标签" label-width="50px" v-show="form.tag_ids.length">
            <div class="input-box">
                <span class="input-header"><i class="el-icon-edit" @click="biggerShow(true)"></i></span>
                {{form.tag_names}}
            </div>
        </el-form-item>
        <el-form-item label-width="2px">
            <div><el-button type="text" @click="msgBoxShow(-1)">推送<i class="el-icon-plus"></i></el-button></div>
            <div class="input-box" v-for="(msg,msg_index) in form.msg_set">
                <el-row v-html="msg.descrition"></el-row>
                <span class="input-header">
                    <i class="el-icon-close" @click="msgSetDel(msg_index)"></i>
                    <i class="el-icon-edit" @click="msgBoxShow(msg_index,msg)"></i>
                </span>
            </div>
        </el-form-item>
        <div style="text-align: right;padding-bottom: 14px;">
            <el-button @click="configRun()" class="left" size="small" v-if="form.status!=Enum.StatusClosed && !request.disabled">执行一下</el-button>
            <el-button @click="close(false)" size="small" :disabled="false">取 消</el-button>
            <el-button type="success" @click="setCron()" size="small" v-if="(form.status==Enum.StatusDisable || form.status==Enum.StatusFinish || form.status==Enum.StatusError || form.status==Enum.StatusReject) && $auth_tag.config_set && !request.disabled">保存草稿</el-button>
        </div>
    </el-form>
    
    <!-- sql链接源管理弹窗 -->
    <el-drawer title="链接管理" :visible.sync="source.boxShow && !request.disabled" size="40%" wrapperClosable="false" :before-close="sqlSourceBox">
        <my-sql-source></my-sql-source>
    </el-drawer>
    <!-- sql 设置弹窗 -->
    <el-dialog :title="'sql设置-'+sqlSet.title" :visible.sync="sqlSet.show && !request.disabled" :show-close="false" :modal="false" :close-on-click-modal="false">
        <el-form :model="sqlSet.statement.git" v-if="sqlSet.source=='git'" label-width="70px" size="small">
            <el-form-item label="连接">
                <el-select v-model="sqlSet.statement.git.link_id" placement="请选择git链接">
                    <el-option v-for="(dic_v,dic_k) in dic.git_source" :label="dic_v.name" :value="dic_v.id"></el-option>
                </el-select>
            </el-form-item>
            <el-form-item label="空间*">
                <el-autocomplete v-model="sqlSet.statement.git.owner" :fetch-suggestions="preferenceGitSuggestions" placeholder="仓库所属空间地址(企业、组织或个人的地址path)"></el-autocomplete>
            </el-form-item>
            <el-form-item label="仓库*">
                <el-autocomplete v-model="sqlSet.statement.git.project" :fetch-suggestions="(text,call)=>preferenceGitSuggestions(text, call, sqlSet.statement.git.owner)" placeholder="仓库路径"></el-autocomplete>
            </el-form-item>
            <el-form-item label="文件路径">
                <el-input v-for="(path_v,path_i) in sqlSet.statement.git.path" v-model="sqlSet.statement.git.path[path_i]" placeholder="文件的路径" @input="sqlGitPathInput">
                    <template slot="prepend">{{sqlSet.statement.git.owner}}/{{sqlSet.statement.git.project}}/</template>
                    <el-button slot="append" icon="el-icon-delete" @click="sqlGitPathDel(path_i)"></el-button>
                </el-input>
            </el-form-item>
            <el-form-item label="引用">
                <el-input v-model="sqlSet.statement.git.ref" placeholder="分支、tag或commit。默认: 仓库的默认分支(通常是master)"></el-input>
            </el-form-item>
            <el-form-item label="批量解析">
                <el-tooltip class="item" effect="dark" content="文件中sql预警将根据;分号为分隔符进行多条拆分" placement="top-start">
                  <el-radio v-model="sqlSet.statement.is_batch" label="1">是</el-radio>
                </el-tooltip>
                <el-tooltip class="item" effect="dark" content="不拆分sql，单文件为一个独立sql语句" placement="top-start">
                  <el-radio v-model="sqlSet.statement.is_batch" label="2">否</el-radio>
                </el-tooltip>
            </el-form-item>
        </el-form>
        
        <el-form :model="sqlSet.statement.local"  v-if="sqlSet.source=='local'">
            <el-input type="textarea" :autosize="{minRows:12}" placeholder="请输入sql内容，批量添加时多个sql请用;分号分隔。" v-model="sqlSet.statement.local"></el-input>
        </el-form>
       
        <span slot="footer" class="dialog-footer">
            <el-button @click="sqlSet.show = false" size="small">取 消</el-button>
            <el-button type="primary" @click="sqlSetConfirm()" size="small">确 定</el-button>
        </span>
    </el-dialog>
    <!-- git 设置弹窗 -->
    <el-dialog :title="'git设置-'+sqlSet.title" :visible.sync="gitSet.show && !request.disabled" :modal="false" :show-close="false" :close-on-click-modal="false">
        <el-form :model="gitSet.event" label-width="90px" size="small">
            <el-form-item label="事件">
                <el-select v-model="gitSet.data.id" placement="请选择事件类型">
                    <el-option v-for="(dic_v,dic_k) in dic.git_event" :label="dic_v.name" :value="dic_v.id"></el-option>
                </el-select>
            </el-form-item>
            <el-divider><i class="el-icon-setting"></i></el-divider>
            
            <div v-if="gitSet.data.id==2"> pr创建 开发中 ...
                
            </div>
            
            <!-- pr详情 -->
            <div v-if="gitSet.data.id==3"> 
                <el-form-item label="空间*">
                    <el-autocomplete v-model="gitSet.data.pr_detail.owner" :fetch-suggestions="preferenceGitSuggestions" placeholder="仓库所属空间地址(企业、组织或个人的地址path)"></el-autocomplete>
                </el-form-item>
                <el-form-item label="仓库*">
                    <el-autocomplete v-model="gitSet.data.pr_detail.repo" :fetch-suggestions="(text,call)=>preferenceGitSuggestions(text, call, gitSet.data.pr_detail.owner)" placeholder="仓库路径"></el-autocomplete>
                </el-form-item>
                <el-form-item label="PR 编号*">
                    <el-input v-model="gitSet.data.pr_detail.number" placeholder="本仓库PR的序数"></el-input>
                </el-form-item>
            </div>
            <!-- pr是否合并 -->
            <div v-if="gitSet.data.id==8"> 
                <el-form-item label="空间*">
                    <el-autocomplete v-model="gitSet.data.pr_is_merge.owner" :fetch-suggestions="preferenceGitSuggestions" placeholder="仓库所属空间地址(企业、组织或个人的地址path)"></el-autocomplete>
                </el-form-item>
                <el-form-item label="仓库*">
                    <el-autocomplete v-model="gitSet.data.pr_is_merge.repo" :fetch-suggestions="(text,call)=>preferenceGitSuggestions(text, call, gitSet.data.pr_is_merge.owner)" placeholder="仓库路径"></el-autocomplete>
                </el-form-item>
                <el-form-item label="PR 编号*">
                    <el-input v-model="gitSet.data.pr_is_merge.number" placeholder="本仓库PR的序数"></el-input>
                </el-form-item>
            </div>
            
            <!-- pr合并 -->
            <div v-if="gitSet.data.id==9"> 
                <el-form-item label="空间*">
                    <el-autocomplete v-model="gitSet.data.pr_merge.owner" :fetch-suggestions="preferenceGitSuggestions" placeholder="仓库所属空间地址(企业、组织或个人的地址path)"></el-autocomplete>
                </el-form-item>
                <el-form-item label="仓库*">
                    <el-autocomplete v-model="gitSet.data.pr_merge.repo" :fetch-suggestions="(text,call)=>preferenceGitSuggestions(text, call, gitSet.data.pr_merge.owner)" placeholder="仓库路径"></el-autocomplete>
                </el-form-item>
                <el-form-item label="PR 编号*">
                    <el-input v-model="gitSet.data.pr_merge.number" placeholder="本仓库PR的序数"></el-input>
                </el-form-item>
                <el-form-item label="commit 标题">
                    <el-input v-model="gitSet.data.pr_merge.title" placeholder="合并 commit 标题，默认为 !{pr_id} {pr_title}"></el-input>
                </el-form-item>
                <el-form-item label="commit 描述">
                    <el-input type="textarea" v-model="gitSet.data.pr_merge.description" placeholder="合并 commit 描述，默认为 Merge pull request !{pr_id} from {author}/{source_branch}"></el-input>
                </el-form-item>
                <el-form-item label="合并方式">
                    <el-select v-model="gitSet.data.pr_merge.merge_method">
                        <el-option label="合并分支" value="merge">
                            <span style="float: left">合并分支</span>
                            <span style="float: right; color: #8492a6; font-size: 13px; padding-left: 26px;">源分支所有的提交都会合并到目标分支，并产生一个新的提交</span>
                        </el-option>
                        <el-option label="扁平化分支" value="squash">
                            <span style="float: left">扁平化分支</span>
                            <span style="float: right; color: #8492a6; font-size: 13px; padding-left: 26px;">源分支中的多个提交会打包成一个提交合并到目标分支</span>
                        </el-option>
                        <el-option label="变基并合并" value="rebase">
                            <span style="float: left">变基并合并</span>
                            <span style="float: right; color: #8492a6; font-size: 13px; padding-left: 26px;">来自源分支的 2 个提交将被重新定位并提交到目标分支。</span>
                        </el-option>
                    </el-select>
                </el-form-item>
                <el-form-item label="合并选项">
                    <el-checkbox v-model="gitSet.data.pr_merge.prune_source_branch">合并后删除提交分支</el-checkbox>
                </el-form-item>
            </div>
            <div v-if="gitSet.data.id==3">
                待开发...
            </div>
            <!-- 文件更新 -->
            <div v-if="gitSet.data.id==131">
                <el-form-item label="空间*">
                    <el-autocomplete v-model="gitSet.data.file_update.owner" :fetch-suggestions="preferenceGitSuggestions" placeholder="仓库所属空间地址(企业、组织或个人的地址path)"></el-autocomplete>
                </el-form-item>
                <el-form-item label="仓库*">
                    <el-autocomplete v-model="gitSet.data.file_update.repo" :fetch-suggestions="(text,call)=>preferenceGitSuggestions(text, call, gitSet.data.file_update.owner)" placeholder="仓库路径"></el-autocomplete>
                </el-form-item>
                <el-form-item label="文件*">
                    <el-input v-model="gitSet.data.file_update.path" placeholder="文件的路径">
                        <template slot="prepend">{{gitSet.data.file_update.owner}}/{{gitSet.data.file_update.repo}}/</template>
                    </el-input>
                </el-form-item>
                <el-form-item label="内容*">
                    <el-input type="textarea" :autosize="{minRows:2}" v-model="gitSet.data.file_update.content" placeholder="文件内容\n支持模板语法，文件原内容会以raw_content变量名称传入。"></el-input>
                </el-form-item>
                <el-form-item label="描述*">
                    <el-input type="textarea" v-model="gitSet.data.file_update.message" placeholder="commit 描述"></el-input>
                </el-form-item>
                <el-form-item label="分支名称">
                    <el-input v-model="gitSet.data.file_update.branch" placeholder="master"></el-input>
                </el-form-item>
            </div>
        </el-form>
       
        <span slot="footer" class="dialog-footer">
            <el-button @click="gitSet.show = false" size="small">取 消</el-button>
            <el-button type="primary" @click="gitSetConfirm()" size="small">确 定</el-button>
        </span>
    </el-dialog>
    <!-- 推送设置弹窗 -->
    <el-dialog title="推送设置" :visible.sync="msg_set_box.show && !request.disabled" :show-close="false" :close-on-click-modal="false" :modal="false">
        <el-form :model="msg_set_box.form" :inline="true" size="small">
            <el-form-item label="当执行">
                <el-select v-model="msg_set_box.form.status" multiple style="width: 143px" placeholder="状态">
                    <el-option v-for="(dic_v,dic_k) in msg_set_box.statusList" :label="dic_v.name" :value="dic_v.id"></el-option>
                </el-select>
                时
            </el-form-item>
            <el-form-item label="发送">
                <el-select v-model="msg_set_box.form.msg_id" style="width: 180px" placeholder="模板">
                    <el-option v-for="(dic_v,dic_k) in dic.msg" :label="dic_v.name" :value="dic_v.id"></el-option>
                </el-select>
                消息
            </el-form-item>
            <el-form-item label="并且@">
                <el-select v-model="msg_set_box.form.notify_user_ids" multiple style="width: 210px" placeholder="人员">
                    <el-option v-for="(dic_v,dic_k) in dic.user" :key="dic_v.id" :label="dic_v.name" :value="dic_v.id"></el-option>
                </el-select>
            </el-form-item>
        </el-form>
        <span slot="footer" class="dialog-footer">
            <el-button @click="msg_set_box.show = false" size="small">取 消</el-button>
            <el-button type="primary" @click="msgSetConfirm()" size="small">确 定</el-button>
        </span>
    </el-dialog>
    
    <!-- 更多设置 弹窗 -->
    <el-dialog title="更多设置" :visible.sync="bigger_set.show" :show-close="false" :close-on-click-modal="false" :modal="false">
        <el-form :model="bigger_set.form" size="small" label-width="69px">
            <el-form-item class="var_fields">
                <span slot="label" style="white-space: nowrap;">
                    参数
                    <el-tooltip effect="dark" content="申明过的参数可以被外部方法传入，点击查看更多" placement="top-start">
                        <router-link target="_blank" to="/var_params" style="color: #606266"><i class="el-icon-info"></i></router-link>
                    </el-tooltip>
                </span>
                <el-input v-for="(p_v,p_i) in bigger_set.form.var_fields" v-model="p_v.remark" placeholder="参数说明" class="input-input">
                    <el-input slot="prepend" v-model="p_v.key" placeholder="key" @input="e=>inputChangeArrayPush(e,p_i,bigger_set.form.var_fields)">
                        <el-input slot="suffix" v-model="p_v.value" placeholder="默认值"></el-input>
                    </el-input>
                    <el-button slot="append" icon="el-icon-delete" @click="arrayDelete(p_i, bigger_set.form.var_fields)"></el-button>
                </el-input>
            </el-form-item>
            <el-form-item label="命令内容">
                <el-skeleton :rows="3" /><!--预期占位-->
            </el-form-item>
            <el-form-item>
                <span slot="label" style="white-space: nowrap;">
                    结果<el-tooltip effect="dark" content="任务结果模板分析，点击查看更多" placement="top-start">
                        <router-link target="_blank" to="/var_params" style="color: #606266"><i class="el-icon-info"></i></router-link>
                    </el-tooltip>
                </span>
                <el-input type="textarea" v-model="bigger_set.form.after_tmpl" :autosize="{minRows:2}" placeholder="任务成功后，对结果响应文本进行二次解析；可重构响应及错误验证。\n结果变量: result·string"></el-input>
            </el-form-item>
            <el-form-item label="重试">
                <el-col :span="11">
                    <el-input type="number" v-model="bigger_set.form.err_retry_num" placeholder="失败时最大重试" class="input-input">
                        <span slot="append">次</span>
                    </el-input>
                </el-col>
                <el-col :span="2" style="text-align: center;">-</el-col>
                <el-col :span="11">
                    <el-input type="number" v-model="bigger_set.form.err_retry_sleep" placeholder="重试间隔" class="input-input">
                        <el-select v-model="bigger_set.form.err_retry_mode" slot="prepend" style="width: 80px;">
                            <el-option v-for="(dic_v,dic_k) in dic.retry_mode" :key="dic_v.id" :label="dic_v.name" :value="dic_v.id"></el-option>
                        </el-select>
                        <span slot="append">秒</span>
                    </el-input>
                </el-col>
            </el-form-item>
            <el-form-item label="延迟">
                <el-input type="number" v-model="bigger_set.form.after_sleep" placeholder="任务成功完成后，继续等待...秒">
                    <span slot="append">s/秒 后结束</span>
                </el-input>
            </el-form-item>
            <el-form-item label="备注">
                <el-input v-model="bigger_set.form.remark" placeholder="任务补充说明"></el-input>
            </el-form-item>
            <el-form-item label="标签">
                <el-select v-model="bigger_set.form.tag_ids" placeholder="无" multiple filterable multiple-limit="5" style="width:100%">
                    <el-option v-for="item in dic.tag" :label="item.name" :value="item.id">
                        <span style="float: left">{{ item.name }}</span>
                        <span style="float: right; color: #8492a6; font-size: 13px">{{ item.extend.remark }}</span>
                    </el-option>
                </el-select>
            </el-form-item>
           
        </el-form>
        <span slot="footer" class="dialog-footer">
            <el-button @click="biggerShow(false)" size="small">取 消</el-button>
            <el-button type="primary" @click="biggerConfirm()" size="small">确 定</el-button>
        </span>
    </el-dialog>
</div>`,
    name: "MyConfigForm",
    props: {
        request:{
            detail:Object, // 详情对象
            disabled: Boolean, // 禁用标识（仅查看）（而且无法触发任何二次弹窗）
        }
    },
    data(){
        return {
            sys_info:{},
            dic:{
                sql_source:[],
                jenkins_source:[],
                host_source:[],
                user:[],
                msg:[],
                tag:[],
                retry_mode: [],
            },

            form:{
                name: "",
            },
            hintSpec: "* * * * * *   秒 分 时 天 月 星期",
            source: {
                boxShow: false,
                dic_type: 0,
            },
            sqlSet: {
                show: false, // 是否显示
                title: '添加',
                index: -1, // 操作行号
                statement: {}, // 实际内容
            }, // sql设置弹窗
            gitSet:{    // git 设置弹窗
                show: false, // 是否显示
                title: '添加',
                index: -1, // 操作行号
                data:{      // 单个设置详情
                    id: '',// 事件编号
                },
                // 合并类型
                gitMergeTypeList: {
                    "merge":"合并分支",
                    "squash":"扁平化分支",
                    "rebase":"变基并合并"
                },
            },
            msg_set_box:{
                show: false, // 是否显示
                title: '添加',
                index: -1, // 操作行号
                form: {}, // 实际内容
                statusList:[{id:1,name:"错误"}, {id:2, name:"成功"}],
            },
            bigger_set:{
                show:false,
                form:{
                    // var_fields:[],
                    // after_tmpl:'',
                    // after_sleep: '',
                    // remark: ''
                },
            },
            // 日期选择器设置
            pickerOptions: {
                disabledDate(time){
                    return time.getTime() < Date.now() - 8.64e7
                },
                selectableRange: "00:00:00 - 23:01:59",
            },
            // 偏好配置
            preference:{},
        }
    },
    watch:{
        "form.spec":{
            handler(v){
                if (this.form.type == 2){ // 年月日的变化控制
                    let cur = moment()
                    if (moment(v).format('YYYY-MM-DD') == cur.format('YYYY-MM-DD')){
                        this.pickerOptions.selectableRange = `${cur.format('HH:mm:ss')} - 23:59:59`
                    }else{
                        this.pickerOptions.selectableRange = "00:00:00 - 23:59:59"
                    }
                }
            }
        },
        "gitSet.data.id":{
            handler(v) {
                console.log("gitSet.data.id 改变", v, this.gitSet.index)

                if (this.gitSet.index !== -1){
                    return
                }
                switch (v) {
                    case 2:
                        this.gitSet.data = {
                            id: v,
                            pr_create: {}
                        }
                        break
                    case 3:
                        this.gitSet.data = {
                            id: v,
                            pr_detail:{
                                owner: this.preference.git.owner ?? '', // 空间
                                repo: this.preference.git.repo ?? '', // 仓库
                                number: '',
                            }
                        }
                        break
                    case 8:
                        this.gitSet.data = {
                            id: v,
                            pr_is_merge:{
                                owner: this.preference.git.owner ?? '', // 空间
                                repo: this.preference.git.repo ?? '', // 仓库
                                number: '',
                            }
                        }
                        break
                    case 9:
                        this.gitSet.data = {
                            id: v,
                            pr_merge: {
                                owner: this.preference.git.owner ?? '', // 空间
                                repo: this.preference.git.repo ?? '', // 仓库
                                number: '',
                                title: "",
                                description: "",
                                merge_method: "merge",
                                prune_source_branch: false,
                            }
                        }
                        break
                    case 131:
                        this.gitSet.data = {
                            id: v,
                            file_update:{
                                owner: this.preference.git.owner ?? '', // 空间
                                repo: this.preference.git.repo ?? '', // 仓库
                                path: '',
                                content: '',
                                message: '',
                                branch: this.preference.git.branch ?? '',
                            }
                        }
                        break
                    default:
                        this.gitSet.data = {id:''}
                }
            }
        }
    },
    // 模块初始化
    created(){
        api.systemInfo((res)=>{
            this.sys_info = res;
        })
        this.getDicSqlSource()
        console.log("config_form request",this.request)
        if (this.request == undefined){
            this.request = {detail:{},disabled:false}
        }
        if (this.request.detail !== undefined && this.request.detail.id > 0){
            this.form = this.editShow(this.request.detail)
        }else{
            this.form = this.initFormData(this.request.detail)
        }
    },
    // 模块初始化
    mounted(){
        this.getPreference()
    },

    // 具体方法
    methods:{
        initFormData(row){
            row = row??{}
            return  {
                type: row.type ?? '1',
                protocol: '3',
                var_fields: [],
                command:{
                    http:{
                        method: 'GET',
                        header: [{}],
                        url:'',
                        body:'',
                    },
                    rpc:{
                        proto: '',
                        method: 'GRPC',
                        addr: '',
                        action: '',
                        actions: [],
                        header: [],
                        body: ''
                    },
                    cmd:{
                        host:{
                            id: -1
                        },
                        origin: 'local', // git、local
                        type: "", // cmd、bash、sh ...
                        statement:{
                            git:{
                                link_id: "",
                                owner: "",
                                project: "",
                                path: [], // 目前不支持多选，但数据结构要能为后续支持扩展
                                ref: "",
                            },
                            local: ""
                        },
                    },
                    sql:{
                        driver: "mysql",
                        source:{
                            id: "",
                            // title: "",
                            // hostname:"",
                            // port:"",
                            // username:"",
                            // password:""
                        },
                        origin: 'local',
                        statement:[],
                        err_action: "1",
                    },
                    jenkins:{
                        source:{
                            id: "",
                        },
                        name: "",
                        params: [{}]
                    },
                    git:{
                        link_id: "",
                        events: [],
                    }
                },
                after_tmpl: '',
                after_sleep: 0,
                err_retry_num: 0,
                err_retry_sleep: 0,
                err_retry_mode: 1,
                err_retry_mode_name: '',
                msg_set: [],
                tag_ids: [],
                tag_names: "",
                status: '1',
            }
        },
        // 编辑弹窗
        editShow(row){
            let form = copyJSON(row)

            if (form.var_fields == null){
                form.var_fields = []
            }else{
                form.var_fields = form.var_fields.filter(function (item){
                    return item.key !== ''
                })
            }
            if (row.command.cmd.statement.git.path == null){
                form.command.cmd.statement.git.path = this.initFormData().command.cmd.statement.git.path
            }

            if (form.command.sql.driver == ""){ // 历史数据兼容
                form.command.sql = this.initFormData().command.sql
            }
            form.command.sql.err_action = form.command.sql.err_action.toString()
            if (form.command.sql.source.id == 0){
                form.command.sql.source.id = ""
            }
            for (let i in row.command.sql.statement){
                form.command.sql.statement[i] = this.sqlGitBuildDesc(row.command.sql.statement[i])
            }

            let hl = form.command.http.header.length
            if (hl == 0){
                form.command.http.header = this.initFormData().command.http.header
            }else if (form.command.http.header[hl-1].key != ""){
                form.command.http.header.push({})
            }
            if (form.command.jenkins.source.id == 0){
                form.command.jenkins.source.id = ""
            }
            let pl = form.command.jenkins.params.length
            if (pl == 0){
                form.command.jenkins.params = this.initFormData().command.jenkins.params
            }else if (form.command.jenkins.params[pl-1].key != ""){
                form.command.jenkins.params.push({})
            }
            if (form.command.cmd.host.id == 0){
                form.command.cmd.host.id = this.initFormData().command.cmd.host.id
            }
            if (form.command.git){
                for (let i in form.command.git.events){
                    form.command.git.events[i] = this.gitBuildDesc(form.command.git.events[i])
                }
                if (form.command.git.link_id == 0){
                    form.command.git.link_id = ''
                }
            }
            for (let i in row.msg_set){
                form.msg_set[i] = this.msgSetBuildDesc(row.msg_set[i])
            }
            if (form.msg_set == null){

            }
            form.type = form.type.toString()
            form.status = form.status.toString() // 这里要转字符串，否则可能显示数字
            form.protocol = form.protocol.toString()
            console.log("编辑：",form)
            return form
        },
        // 添加/编辑 任务
        setCron(afterCall){
            if (this.request.disabled){
                return
            }else if (this.form.name == ''){
                return this.$message.error('请输入任务名称')
            }else if (this.form.spec == '' && this.form.type != 5){
                return this.$message.error('请输入任务执行时间')
            }else if (!this.form.protocol){
                return this.$message.error('请选择任务协议')
            }
            // else if (this.form.command == ''){
            //     return this.$message.error('请输入命令类容')
            // }
            let body = copyJSON(this.form)
            body.type = Number(body.type)
            body.status = Number(body.status)
            body.protocol = Number(body.protocol)
            body.command.sql.err_action = Number(body.command.sql.err_action)
            body.command.sql.interval = Number(body.command.sql.interval)
            body.command.sql.source.id = Number(body.command.sql.source.id)
            body.command.jenkins.source.id = Number(body.command.jenkins.source.id)
            body.command.cmd.statement.git.link_id = Number(body.command.cmd.statement.git.link_id)
            body.command.git.link_id = Number(body.command.git.link_id)
            body.command.http.header = body.command.http.header.filter(function (item) {
                return item['key'] !== undefined &&  item.key !== ''
            })
            body.command.jenkins.params = body.command.jenkins.params.filter(function (item) {
                return item['key'] !== undefined &&  item.key !== ''
            })

            api.innerPost("/config/set", body, (res)=>{
                if (!res.status){
                    return this.$message.error(res.message)
                }
                afterCall && afterCall(res.data.id)
                this.form = this.initFormData()
                this.close(true)
            })
        },
        // http header 输入值变化时
        httpHeaderInput(val){
            if (val == ""){
                return
            }
            let item = this.form.command.http.header.slice(-1)[0]
            if (item == undefined || item.key != ""){
                this.form.command.http.header.push({"key":"","value":""})
            }
        },
        // http header 输入值删除
        httpHeaderDel(index){
            if ((index+1) >= this.form.command.http.header.length){
                return
            }
            this.form.command.http.header.splice(index,1)
        },
        // jenkinsParam 输入值变化时
        jenkinsParamInput(val){
            if (val == ""){
                return
            }
            let item = this.form.command.jenkins.params.slice(-1)[0]
            if (item == undefined || item.key != ""){
                this.form.command.jenkins.params.push({"key":"","value":""})
            }
        },
        // jenkinsParam 输入值删除
        jenkinsParamDel(index){
            if ((index+1) >= this.form.command.jenkins.params.length){
                return
            }
            this.form.command.jenkins.params.splice(index,1)
        },
        // sqlGitPath 输入值变化时
        sqlGitPathInput(val){
            if (val == ""){
                return
            }
            let item = this.sqlSet.statement.git.path.slice(-1)[0]
            if (item == undefined || item != ""){
                this.sqlSet.statement.git.path.push("")
            }
        },
        // sqlGitPath 输入值删除
        sqlGitPathDel(index){
            if ((index+1) >= this.sqlSet.statement.git.path.length){
                return
            }
            this.sqlSet.statement.git.path.splice(index,1)
        },
        // sql source box
        sqlSourceBox(show){
            this.sqlSourceBoxShow = show == true;
            if (!this.sqlSourceBoxShow){
                this.getDicSqlSource() // 关闭弹窗要重载枚举
            }
        },
        // sql设置弹窗
        sqlSetShow(index, oldData){
            if (index === "" || index == null || isNaN(index)){
                console.log('sqlSetShow', index, oldData)
                return this.$message.error("索引位标志异常");
            }
            this.sqlSet.source = this.form.command.sql.origin
            this.sqlSet.show = true
            this.sqlSet.index = Number(index)  // -1.新增、>=0.具体行的编辑
            if (oldData == undefined){
                oldData = {
                    type: this.form.command.sql.origin,
                    git:null,
                    local: "",
                    is_batch: "1", // 批量解析：1.是（默认）、2.否
                }
            }
            if (oldData.git == null){
                oldData.git = {
                    link_id: "",
                    owner: this.preference.git.owner ?? '',
                    project: this.preference.git.repo ?? '',
                    path: [""],
                    ref: this.preference.git.branch ?? '',
                }
            }
            oldData.is_batch = oldData.is_batch.toString()
            this.sqlSet.statement = oldData

            this.sqlSet.title = this.sqlSet.index < 0? '添加' : '编辑';
        },
        // sql 设置确定
        sqlSetConfirm(){
            if (typeof this.sqlSet.index !== 'number'){
                console.log('sqlSetShow', this.sqlSet)
                return this.$message.error("索引位标志异常");
            }

            if (this.sqlSet.source == 'git'){
                if (this.sqlSet.statement.git.link_id == ""){
                    return this.$message.error("请选择连接")
                }

                if (this.sqlSet.statement.git.owner == ""){
                    return this.$message.error("空间为必填")
                }
                if (this.sqlSet.statement.git.project == ""){
                    return this.$message.error("仓库为必填")
                }
                if (this.sqlSet.statement.git.path.length == 0){
                    return this.$message.error("请输入文件的路径")
                }
                let data = copyJSON(this.sqlSet.statement);
                data.git.link_id = Number(data.git.link_id)
                data.is_batch = Number(data.is_batch)
                data.type = this.sqlSet.source
                if (data.git.ref == ""){
                    data.git.ref = 'master'
                }
                data.git.path = data.git.path.filter(function (item,index){
                    item = item.trim()
                    if (item !== ''){
                        if (item.charAt(0) === '/'){
                            item = item.substring(1)
                            console.log("xxxx",item)
                        }
                        if (item !== ''){
                            data.git.path[index] = item
                            return true
                        }
                    }
                    return false
                })

                data = this.sqlGitBuildDesc(data)
                if (this.sqlSet.index < 0){
                    this.form.command.sql.statement.push(data)
                }else{
                    this.form.command.sql.statement[this.sqlSet.index] = data
                }
                console.log("git data",data)

            }else{
                if (this.sqlSet.statement.local == ""){
                    return this.$message.error("sql内容不得为空");
                }
                // 支持批量添加
                let temp = this.sqlSet.statement.local.split(";")
                let datas = []
                for (let i in temp){
                    let val = temp[i].trim()
                    if (val != ""){
                        datas.push({"type":this.sqlSet.source, "git":{}, "local":val,"is_batch":1})
                    }
                }

                if (this.sqlSet.index < 0){
                    this.form.command.sql.statement.push(...datas)
                }else if (datas.length > 1){
                    return this.$message.warning("不支持单sql的拆分，建议删除后批量添加");
                }else if (datas.length == 0){
                    return this.$message.warning("不存在有效sql，请确认输入");
                }else{
                    this.form.command.sql.statement[this.sqlSet.index] = datas[0]
                }
            }

            this.sqlSet.show = false
            // this.sqlSet.statement = {}
            this.sqlSet.index = -1
        },
        // 删除sql元素
        sqlSetDel(index){
            if (this.request.disabled){
                return
            }
            if (index === "" || index == null || isNaN(index)){
                console.log('sqlSetDel', index)
                return this.$message.error("索引位标志异常");
            }
            this.$confirm('此操作将删除sql执行语句，是否继续？','提示',{
                type:'warning',
            }).then(()=>{
                this.form.command.sql.statement.splice(index,1)
            })
        },
        // git 设置弹窗
        gitSetShow(index, oldData){
            if (index === "" || index == null || isNaN(index)){
                return this.$message.error("索引位标志异常"+index);
            }
            let data = {id:''}
            if (oldData != undefined){
                data = copyJSON(oldData)
            }

            this.gitSet.show = true
            this.gitSet.title = this.gitSet.index < 0? '添加' : '编辑';
            this.gitSet.index = Number(index)  // -1.新增、>=0.具体行的编辑
            this.gitSet.data = data
        },
        // git 设置确定
        gitSetConfirm(){
            console.log("git 提交",this.gitSet)
            if (typeof this.gitSet.index !== 'number'){
                console.log('gitSetShow', this.sqlSet)
                return this.$message.error("索引位标志异常");
            }

            if (this.gitSet.data.id == 2){
                // 待完善...
            }else if (this.gitSet.data.id == 3){ // pr 详情
                if (this.gitSet.data.pr_detail.owner == ""){
                    return this.$message.warning("空间为必填")
                }
                if (this.gitSet.data.pr_detail.repo == ""){
                    return this.$message.warning("仓库为必填")
                }
                if (!this.gitSet.data.pr_detail.number){
                    return this.$message.warning("仓库PR编号为必填")
                }
            }else if (this.gitSet.data.id == 8){ // pr 是否 合并
                if (this.gitSet.data.pr_is_merge.owner == ""){
                    return this.$message.warning("空间为必填")
                }
                if (this.gitSet.data.pr_is_merge.repo == ""){
                    return this.$message.warning("仓库为必填")
                }
                if (!this.gitSet.data.pr_is_merge.number){
                    return this.$message.warning("仓库PR编号为必填")
                }
            }else if (this.gitSet.data.id == 9){ // pr 合并
                if (this.gitSet.data.pr_merge.owner == ""){
                    return this.$message.warning("空间为必填")
                }
                if (this.gitSet.data.pr_merge.repo == ""){
                    return this.$message.warning("仓库为必填")
                }
                if (!this.gitSet.data.pr_merge.number){
                    return this.$message.warning("仓库PR编号为必填")
                }
                if (this.gitSet.data.pr_merge.merge_method == ""){
                    return this.$message.warning("合并方式不得为空")
                }
            }else if (this.gitSet.data.id == 131){ // 文件 更新
                if (this.gitSet.data.file_update.owner == ""){
                    return this.$message.warning("空间为必填")
                }
                if (this.gitSet.data.file_update.repo == ""){
                    return this.$message.warning("仓库为必填")
                }
                if (!this.gitSet.data.file_update.path){
                    return this.$message.warning("文件路径 为必填")
                }
                if (this.gitSet.data.file_update.content == ""){
                    return this.$message.warning("文件内容 不得为空")
                }
                if (this.gitSet.data.file_update.message == ""){
                    return this.$message.warning("描述 不得为空")
                }
                if (this.gitSet.data.file_update.branch == ""){
                    this.gitSet.data.file_update.branch = 'master'
                }
            }else{
                return this.$message.error("未选择有效事件");
            }

            let data = this.gitBuildDesc(copyJSON(this.gitSet.data))
            if (this.gitSet.index < 0){
                this.form.command.git.events.push(data)
            }else{
                this.form.command.git.events[this.gitSet.index] = data
            }
            this.gitSet.id = ''
            this.gitSet.show = false
            this.gitSet.index = -1
        },
        // git 删除元素
        gitSetDel(index){
            if (this.request.disabled){
                return
            }
            if (index === "" || index == null || isNaN(index)){
                console.log('gitSetDel', index)
                return this.$message.error("索引位标志异常");
            }
            this.$confirm('此操作将删除sql执行语句，是否继续？','提示',{
                type:'warning',
            }).then(()=>{
                this.form.command.git.events.splice(index,1)
            })
        },

        // 更多设置弹窗
        biggerShow(show = true){
            if (this.request.disabled){
                return
            }
            this.bigger_set.show = show === true
            if (this.bigger_set.show){
                let form  = {
                    var_fields: copyJSON(this.form.var_fields),
                    after_tmpl: this.form.after_tmpl,
                    after_sleep: this.form.after_sleep,
                    err_retry_num: this.form.err_retry_num,
                    err_retry_sleep: this.form.err_retry_sleep,
                    err_retry_mode: this.form.err_retry_mode,
                    remark: this.form.remark,
                    tag_ids: this.form.tag_ids ?? [],
                }
                if (form.var_fields == null || form.var_fields.length == 0){
                    form.var_fields = [{key:'',value:'', remark:''}]
                }else{
                    form.var_fields.push({key:'',value:'', remark:''})
                }
                this.bigger_set.form = form
            }
        },
        // 更多设置 确认
        biggerConfirm(){
            let data = copyJSON(this.bigger_set.form)
            this.form.var_fields = data.var_fields.filter(function (item) {
                return item['key'] !== undefined &&  item.key !== ''
            })
            this.form.after_tmpl = data.after_tmpl
            this.form.after_sleep = Number(data.after_sleep)
            this.form.err_retry_num = Number(data.err_retry_num)
            this.form.err_retry_sleep = Number(data.err_retry_sleep)
            this.form.err_retry_mode = Number(data.err_retry_mode)
            this.dic.retry_mode.forEach((item)=>{
                if (item.id == this.form.err_retry_mode){
                    this.form.err_retry_mode_name = item.name
                }
            })
            this.form.remark = data.remark
            this.form.tag_ids = data.tag_ids
            this.form.tag_names = data.tag_ids.map(tag_id=>{
                const item = this.dic.tag.find(opt => opt.id==tag_id)
                return item.name
            }).join(",")

            this.biggerShow(false)
        },
        // 推送弹窗
        msgBoxShow(index, oldData){
            if (index === "" || index == null || isNaN(index)){
                console.log('msgSetShow', index, oldData)
                return this.$message.error("索引位标志异常");
            }
            if (oldData == undefined || index < 0){
                oldData = {
                    status: [],
                    msg_id: "",
                    notify_user_ids: [],
                }
            }else if (typeof oldData != 'object'){
                console.log('推送信息异常', oldData)
                return this.$message.error("推送信息异常");
            }
            this.msg_set_box.show = true
            this.msg_set_box.index = Number(index)  // -1.新增、>=0.具体行的编辑
            this.msg_set_box.title = this.msg_set_box.index < 0? '添加' : '编辑';
            this.msg_set_box.form = oldData
        },
        msgSetDel(index){
            if (this.request.disabled){
                return
            }
            if (index === "" || index == null || isNaN(index)){
                console.log('msgSetDel', index)
                return this.$message.error("索引位标志异常");
            }
            this.$confirm('确认删除推送配置','提示',{
                type:'warning',
            }).then(()=>{
                this.form.msg_set.splice(index, 1);
            })

        },
        // 推送确认
        msgSetConfirm(){
            if (this.msg_set_box.form.msg_id <= 0){
                return this.$message.warning("请选择消息模板");
            }
            let data = this.msgSetBuildDesc(this.msg_set_box.form)
            if (this.msg_set_box.index < 0){
                this.form.msg_set.push(data)
            }else{
                this.form.msg_set[this.msg_set_box.index] = data
            }
            this.msg_set_box.show = false
            this.msg_set_box.index = -1
            this.msg_set_box.form = {}
        },
        // 构建消息设置描述
        msgSetBuildDesc(data){
            let item1 = this.msg_set_box.statusList.filter((option) => {
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
            if (data.git == null || data.type !== 'git'){
                return data
            }
            let git = this.dic.git_source.filter(item=>{
                return item.id == data.git.link_id
            })

            data.git.descrition = '连接<span class="el-tag el-tag--small el-tag--light">'+git.length>0 ? git[0].name: '' +
                '</span> 访问<span class="el-tag el-tag--small el-tag--light">'+data.git.owner+'/'+data.git.project +'</span>'+
                '</span> 引用<span class="el-tag el-tag--small el-tag--light">'+data.git.ref +'</span> 拉取以下文件内容'+
                '<span class="el-tag el-tag--small el-tag--light">'+ (data.is_batch==1? '批量解析后' :'单文件单sql') +'</span>执行';

            data.git.path.forEach(function (item) {
                if (item == ''){
                    return
                }
                item.split(',').forEach(function (item2) {
                    data.git.descrition += `<div style="margin: 0;padding: 0 0 0 30px;"><a href="https://gitee.com/${data.git.owner}/${data.git.project}/blob/${data.git.ref}/${item2}" target="_blank" title="点击 查看文件详情"><i class="el-icon-connection"></i></a><code>${item2}</code></div>`
                })
            })

            console.log("sql描述",data.git.descrition)
            return data
        },
        // 构建 git 描述
        gitBuildDesc(data){
            switch (data.id){
                case 2:
                    data.desc = '完善中...'
                    break
                case 3:
                    data.desc = `<b>pr详情</b> <a href="https://gitee.com/${data.pr_detail.owner}/${data.pr_detail.repo}/pulls/${data.pr_detail.number}" target="_blank" title="点击 查看pr详情"><i class="el-icon-connection"></i></a> <b class="b">${data.pr_detail.owner}/${data.pr_detail.repo}</b>/pulls/<b class="b">${data.pr_detail.number}</b>`
                    break
                case 8:
                    data.desc = `<b>pr是否合并</b> <a href="https://gitee.com/${data.pr_is_merge.owner}/${data.pr_is_merge.repo}/pulls/${data.pr_is_merge.number}" target="_blank" title="点击 查看pr详情"><i class="el-icon-connection"></i></a> <b class="b">${data.pr_is_merge.owner}/${data.pr_is_merge.repo}</b>/pulls/<b class="b">${data.pr_is_merge.number}</b>`
                    break
                case 9:
                    data.desc = `<b>pr合并</b> <a href="https://gitee.com/${data.pr_merge.owner}/${data.pr_merge.repo}/pulls/${data.pr_merge.number}" target="_blank" title="点击 查看pr详情"><i class="el-icon-connection"></i></a> <b class="b">${data.pr_merge.owner}/${data.pr_merge.repo}</b>/pulls/<b class="b">${data.pr_merge.number}</b> <b class="b">${this.gitSet.gitMergeTypeList[data.pr_merge.merge_method]}</b>  ${data.pr_merge.prune_source_branch===true?'<b class="b">删除提交分支</b>':''}`+
                        `<br><i style="margin-left: 3em;"></i><b>${data.pr_merge.title}</b> ${data.pr_merge.description}`
                    break
                case 131:
                    data.desc = `<b>文件更新</b> <a href="https://gitee.com/${data.file_update.owner}/${data.file_update.repo}/blob/${data.file_update.branch}/${data.file_update.path}" target="_blank" title="点击 查看"><i class="el-icon-connection"></i></a> <b class="b">${data.file_update.owner}</b>/<b class="b">${data.file_update.repo}</b>/blob/<b class="b">${data.file_update.branch}</b>/<b class="b">${data.file_update.path}</b>`+
                        `<br><i style="margin-left: 2em;"></i>内容：<pre style="display: inline-flex;max-height: 140px;">${data.file_update.content}</pre>`+
                        `<br><i style="margin-left: 2em;"></i>描述：<b class="b">${data.file_update.message}</b>`
                    break
                default:
                    data.desc = '未支持的事件类型'
            }
            return data
        },
        // 枚举
        getDicSqlSource(){
            let types = [
                Enum.dicSqlSource,
                Enum.dicSqlDriver,
                Enum.dicJenkinsSource,
                Enum.dicGitSource,
                Enum.dicGitEvent,
                Enum.dicHostSource,
                Enum.dicCmdType,
                Enum.dicUser,
                Enum.dicMsg,
                Enum.dicTag,
                Enum.dicRetryMode
            ]
            api.dicList(types,(res) =>{
                this.dic.sql_source = res[Enum.dicSqlSource]
                this.dic.jenkins_source = res[Enum.dicJenkinsSource]
                this.dic.git_source =res[Enum.dicGitSource]
                this.dic.git_event = res[Enum.dicGitEvent]
                this.dic.host_source =res[Enum.dicHostSource]
                this.dic.user = res[Enum.dicUser]
                this.dic.msg = res[Enum.dicMsg]
                this.dic.tag = res[Enum.dicTag]
                this.dic.cmd_type = res[Enum.dicCmdType]
                this.dic.sql_driver = res[Enum.dicSqlDriver]
                this.dic.retry_mode = res[Enum.dicRetryMode]
            }, true)
        },
        // 解析proto内容
        parseProto(){
            let data = this.form.command.rpc.proto
            if (data == "" || data == undefined){
                return this.$message.warning("请先输入proto内容");
            }
            api.innerPost("/foundation/parse_proto", {proto:data}, (res)=>{
                if (!res.status){
                    return this.$message.error(res.message)
                }
                this.form.command.rpc.actions = res.data.actions
            })
        },
        // 解析 spec 时间
        parseSpec(){
            let data = this.form.spec
            if (!data){
                return this.$message.warning("请先输入时间");
            }
            api.innerPost("/foundation/parse_spec", {spec:data}, (res)=>{
                if (!res.status){
                    return this.$message.error(res.message)
                }
                let txt = res.data.list.join('</br>')
                // 弹窗展示结果
                this.$alert(txt, '最近5次运行时间', {
                    confirmButtonText: '确定',
                    dangerouslyUseHTMLString: true
                });
            })
        },
        // 执行一下
        configRun(){
            this.$confirm('确认现在就执行任务吗？','提示').then(()=>{
                // 主要是强制类型
                let body = copyJSON({
                    id: this.form.id,
                    name: this.form.name,
                    type: Number(this.form.type),
                    spec: this.form.spec,
                    protocol: Number(this.form.protocol),
                    command: this.form.command,
                    remark: this.form.remark,
                    after_tmpl: this.form.after_tmpl,
                    var_fields: this.form.var_fields,
                    msg_set: this.form.msg_set,
                })
                body.command.sql.err_action = Number(body.command.sql.err_action)
                body.command.sql.interval = Number(body.command.sql.interval)
                body.command.sql.source.id = Number(body.command.sql.source.id)
                body.command.jenkins.source.id = Number(body.command.jenkins.source.id)
                body.command.cmd.statement.git.link_id = Number(body.command.cmd.statement.git.link_id)
                body.command.git.link_id = Number(body.command.git.link_id)

                api.innerPost("/config/run", body, (res)=>{
                    if (!res.status){
                        return this.$message({
                            message: res.message,
                            type: 'error',
                            duration: 6000
                        })
                    }
                    return this.$message.success("ok."+res.data.result)
                })
            })
        },
        close(is_change=false){
            this.$emit('close', {is_change:is_change})
        },
        // 获取偏好
        getPreference(){
            api.innerGet("/setting/preference_get", null,res=>{
                if (!res.status){
                    return this.$message.error('偏好错误，'+res.message)
                }
                this.preference = res.data
            })
        },
        // git输入建议
        preferenceGitSuggestions(q, call, owner=null){
            let select = []
            if (owner == null){
                this.preference.git.owner_repo.forEach((item)=>{
                    select.push({value: item.owner})
                })
            }else{
                this.preference.git.owner_repo.forEach(function (item) {
                    if (item.owner === owner){
                        item.repos.forEach( (item2)=> {
                            select.push({value:item2.name})
                        })
                    }
                })
            }
            call(select)
        }
    }
})

Vue.component("MyConfigForm", MyConfigForm);