var MyConfig = Vue.extend({
    template: `<el-main>
                    <el-menu :default-active="labelType" class="el-menu-demo" mode="horizontal" @select="handleClickTypeLabel">
                        <el-menu-item index="1" :disabled="listRequest">周期任务</el-menu-item>
                        <el-menu-item index="2" :disabled="listRequest">单次任务</el-menu-item>
                        <el-menu-item index="5" :disabled="listRequest">组件任务</el-menu-item>
                        <div style="float: right">
                            <el-button type="text" @click="createShow()">添加任务</el-button>
                            <el-button type="text" @click="getRegisterList()">已注册任务</el-button>
                        </div>
                    </el-menu>
                    <el-table :data="list">
                        <el-table-column prop="id" label="编号"></el-table-column>
                        <el-table-column prop="" label="执行成功率">
                            <template slot-scope="scope">
                                <el-tooltip placement="top-start">
                                    <div slot="content">{{scope.row.topRatio}}%<br/>近期{{scope.row.top_number}}次执行，{{scope.row.top_error_number}}次失败。</div>
                                    <i :class="getTopIcon(scope.row.top_number, scope.row.topRatio)"></i>
                                </el-tooltip>
                            </template>
                        </el-table-column>
                        <el-table-column prop="spec" label="执行时间"></el-table-column>
                        <el-table-column prop="name" label="任务名称"></el-table-column>
                        <el-table-column prop="protocol_name" label="协议"></el-table-column>
                        <el-table-column prop="remark" label="备注"></el-table-column>
                        <el-table-column prop="" label="状态">
                            <template slot-scope="scope">
                                <el-tooltip placement="top-start">
                                    <div slot="content">{{scope.row.status_dt}}  {{scope.row.status_remark}}</div>
                                    <span :class="statusClass(scope.row.status)">{{scope.row.status_name}}</span>
                                </el-tooltip>
                            </template>
                        </el-table-column>
                        <el-table-column label="操作">
                            <template slot-scope="{row}">
                                <el-button type="text" @click="editShow(row)">编辑</el-button>
                                <el-button type="text" @click="changeStatus(row, 2)" v-if="row.status!=2">激活</el-button>
                                <el-button type="text" @click="changeStatus(row, 1)" v-if="row.status==2">停用</el-button>
                                <el-button type="text" @click="configLogBox(row)">日志</el-button>
                            </template>
                        </el-table-column>
                    </el-table>
                    <el-pagination
                            @size-change="handleSizeChange"
                            @current-change="handleCurrentChange"
                            :current-page.sync="listPage.page"
                            :page-size="listPage.size"
                            layout="total, prev, pager, next"
                            :total="listPage.total">
                    </el-pagination>
                    
                    
                    <!-- 任务设置表单 -->
                    <el-dialog :title="setConfigTitle" :visible.sync="setConfigShow" :close-on-click-modal="false" class="config-form-box">
                        <el-form :model="form">
                            <el-form-item label="名称*" label-width="50px">
                                <el-input v-model="form.name"></el-input>
                            </el-form-item>
                            <el-form-item label="类型*" label-width="50px">
                                <el-radio v-model="form.type" label="1">周期</el-radio>
                                <el-radio v-model="form.type" label="2">单次</el-radio>
                                <el-radio v-model="form.type" label="5">组件</el-radio>
                            </el-form-item>
                            <el-form-item label="时间*" label-width="50px" v-if="form.type != 5">
                                <el-input v-show="form.type==1" v-model="form.spec" :placeholder="hintSpec"></el-input>
                                <el-date-picker 
                                    style="width: 100%"
                                    v-show="form.type==2" 
                                    v-model="form.spec" 
                                    value-format="yyyy-MM-dd HH:mm:ss"
                                    type="datetime" 
                                    placeholder="选择运行时间" 
                                    :picker-options="pickerOptions">
                                </el-date-picker>
                            </el-form-item>
                            <el-form-item prop="plate" label-width="50px" size="mini" class="http_header_box">
                                <span slot="label" style="white-space: nowrap;">
                                    参数
                                    <el-tooltip effect="dark" content="申明过的参数可以被外部方法传入，点击查看更多" placement="top-start">
                                    <router-link target="_blank" to="/var_params" style="color: #606266"><i class="el-icon-info"></i></router-link>
                                    </el-tooltip>
                                </span>
                                <el-input v-for="(p_v,p_i) in form.var_fields" v-model="p_v.value" placeholder="参数说明">
                                    <el-input v-model="p_v.key" slot="prepend" placeholder="参数 key" @input="e=>inputChangeArrayPush(e,p_i,form.var_fields)"></el-input>
                                    <el-button slot="append" icon="el-icon-delete" @click="arrayDelete(p_i, form.var_fields)"></el-button>
                                </el-input>
                            </el-form-item>
                            
                            <el-form-item>
                                <el-tabs type="border-card" v-model="form.protocol">
                                    <el-tab-pane label="http" name="1">
                                        <el-form-item label="请求地址" class="http_url_box" size="mini">
                                            <el-input v-model="form.command.http.url" placeholder="请输入http:// 或 https:// 开头的完整地址">
                                                <el-select v-model="form.command.http.method" placeholder="请选请求方式" slot="prepend">
                                                    <el-option label="GET" value="GET"></el-option>
                                                    <el-option label="POST" value="POST"></el-option>
                                                </el-select>
                                            </el-input>
                                        </el-form-item>
            
                                        <el-form-item label="请求Header" class="http_header_box" size="mini">
                                            <el-input v-for="(header_v,header_i) in form.command.http.header" v-model="header_v.value" placeholder="参数值">
                                                <el-input v-model="header_v.key" slot="prepend" placeholder="参数名" @input="httpHeaderInput"></el-input>
                                                <el-button slot="append" icon="el-icon-delete" @click="httpHeaderDel(header_i)"></el-button>
                                            </el-input>
                                        </el-form-item>
                                        <el-form-item label="请求Body参数" size="mini">
                                            <el-input type="textarea" v-model="form.command.http.body" rows="5" placeholder="POST请求时body参数，将通过json进行请求发起"></el-input>
                                        </el-form-item>
                                    </el-tab-pane>
                                    <el-tab-pane label="rpc" name="2">
                                        <el-form-item label="请求模式" size="mini">
                                            <el-radio v-model="form.command.rpc.method" label="GRPC">GRPC</el-radio>
                                        </el-form-item>
                                        <el-form-item label="proto" size="mini" style="margin-top: -10px">
                                            <el-input type="textarea" v-model="form.command.rpc.proto" rows="5" placeholder="请输入*.proto 文件内容" style=""></el-input>
                                        </el-form-item>
                                        <el-form-item label="地址" label-width="42px" size="mini">
                                            <el-input v-model="form.command.rpc.addr" placeholder="请输入服务地址，含端口; 示例：localhost:21014"></el-input>
                                        </el-form-item>
                                        <el-form-item label="方法" size="mini">
                                            <el-select v-model="form.command.rpc.action" filterable clearable placeholder="填写proto信息后点击解析，可获得其方法进行选择" style="min-width: 400px">
                                                <el-option v-for="item in form.command.rpc.actions" :key="item" :label="item" :value="item"></el-option>
                                            </el-select>
                                            <el-button @click="parseProto()">解析proto</el-button>
                                        </el-form-item>
                                        <el-form-item label="请求参数" size="mini">
                                            <el-input type="textarea" v-model="form.command.rpc.body" rows="5" placeholder="请输入请求参数"></el-input>
                                        </el-form-item>
                                    </el-tab-pane>
                                    <el-tab-pane label="cmd" name="3">
                                        主机
                                        <el-select v-model="form.command.cmd.host.id" size="mini" style="width:180px">
                                            <el-option v-for="dic_v in dic_host_source" :label="dic_v.name" :value="dic_v.id"></el-option>
                                        </el-select>
                                        来源
                                        <el-select v-model="form.command.cmd.origin" size="mini" style="width:80px">
                                            <el-option label="当前" value="local"></el-option>
                                            <el-option label="git" value="git"></el-option>
                                        </el-select>
                                        类型
                                        <el-select v-model="form.command.cmd.type" size="mini" clearable style="width:80px">
                                            <el-option v-for="dic_v in dic_cmd_type" :label="dic_v.name" :value="dic_v.key"></el-option>
                                        </el-select>
                                        <el-card class="box-card" shadow="hover" v-if="form.command.cmd.origin=='git'">
                                            <el-form :model="form.command.cmd.statement.git" label-width="70px" size="small">
                                                <el-form-item label="连接">
                                                    <el-select v-model="form.command.cmd.statement.git.link_id" placement="请选择git链接">
                                                        <el-option v-for="(dic_v,dic_k) in dic_git_source" :label="dic_v.name" :value="dic_v.id"></el-option>
                                                    </el-select>
                                                </el-form-item>
                                                <el-form-item label="仓库空间">
                                                    <el-input v-model="form.command.cmd.statement.git.owner" placeholder="仓库所属空间地址(企业、组织或个人的地址path)"></el-input>
                                                </el-form-item>
                                                <el-form-item label="项目名称">
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
                                        </el-card>
                                        <el-card class="box-card"  shadow="hover" v-if="form.command.cmd.origin=='local'">
                                            <el-input type="textarea" v-model="form.command.cmd.statement.local" rows="5" placeholder="请输入命令行执行内容"></el-input>
                                        </el-card>
                                    </el-tab-pane>
                                    <el-tab-pane label="sql" name="4" label-position="left">
                                        <el-form-item label="驱动" size="mini">
                                            <el-select v-model="form.command.sql.driver" placeholder="驱动">
                                                <el-option v-for="dic_v in dic_sql_driver" :label="dic_v.name" :value="dic_v.key"></el-option>
                                            </el-select>
                                        </el-form-item>
                                        <el-form-item label="链接" size="mini">
                                            <el-select v-model="form.command.sql.source.id" placement="请选择sql链接">
                                                <el-option v-for="(dic_v,dic_k) in dic_sql_source" v-show="dic_v.extend.driver == form.command.sql.driver" :label="dic_v.name" :value="dic_v.id"></el-option>
                                            </el-select>
                                            <el-button type="text" style="margin-left: 20px" @click="sqlSourceBox(true)">设置链接</el-button>
                                        </el-form-item>
                                        <el-form-item label="执行语句" size="mini">
                                            <div>
                                                来源
                                                <el-select v-model="form.command.sql.origin" size="mini" style="width:80px">
                                                    <el-option label="本地" value="local"></el-option>
                                                    <el-option label="git" value="git"></el-option>
                                                </el-select>
                                                <el-button type="text" @click="sqlSetShow(-1)">添加<i class="el-icon-plus"></i></el-button>
                                            </div>
                                            <div style="overflow-y: auto;max-height: 420px;">
                                                <!--本地sql展示-->
                                                <div v-if="form.command.sql.origin == 'local'" v-for="(statement,sql_index) in form.command.sql.statement" v-show="statement.type=='local' || statement.type==''" style="position: relative;line-height: 133%;background: #f4f4f5;margin-bottom: 5px;padding: 6px 20px 7px 8px;border-radius: 3px;">
                                                    <pre style="margin: 0;overflow-y: auto;max-height: 180px;min-height: 45px;"><code class="language-sql hljs">{{statement.local}}</code></pre>
                                                    <i class="el-icon-close" style="font-size: 15px;position: absolute;top: 2px;right: 2px;cursor:pointer" @click="sqlSetDel(sql_index)"></i>
                                                    <i class="el-icon-edit" style="font-size: 15px;position: absolute;top: 23px;right: 2px;cursor:pointer" @click="sqlSetShow(sql_index,statement)"></i>
                                                    <i style="position: absolute;right: 1px;top: 40px;font-size: 16px;">#{{sql_index}}</i>
                                                </div>
                                                <el-alert v-show="form.command.sql.statement.length==0" title="未添加执行sql，请添加。" type="info"></el-alert>
                                                <!--git sql展示-->
                                                <div v-if="form.command.sql.origin == 'git'" v-for="(statement,sql_index) in form.command.sql.statement" v-show="statement.type=='git'" style="position: relative;line-height: 133%;background: #f4f4f5;margin-bottom: 5px;padding: 6px 20px 7px 8px;border-radius: 3px;">
                                                    <el-row v-html="statement.git.descrition"></el-row>
                                                    <pre v-for="path_item in statement.git.path" style="margin: 0;padding-left: 30px;"><code class="">{{path_item}}</code></pre>
                                                    <i class="el-icon-close" style="font-size: 15px;position: absolute;top: 2px;right: 2px;cursor:pointer" @click="sqlSetDel(sql_index)"></i>
                                                    <i class="el-icon-edit" style="font-size: 15px;position: absolute;top: 23px;right: 2px;cursor:pointer" @click="sqlSetShow(sql_index,statement)"></i>
                                                    <i style="position: absolute;right: 1px;top: 40px;font-size: 16px;">#{{sql_index}}</i>
                                                </div>
                                            </div>
                                        </el-form-item>
                                        <el-form-item label="错误行为" size="mini">
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
                                        <el-form-item label="执行间隔" size="mini" v-if="form.command.sql.err_action !=3" label-width="69px">
                                            <el-input type="number" v-model="form.command.sql.interval" placeholder="s秒">
                                                <span slot="append">s/秒</span>
                                            </el-input>
                                        </el-form-item>
                                    </el-tab-pane>
                                    
                                    <el-tab-pane label="jenkins" name="5">
                                        <el-form-item label="链接">
                                            <el-select v-model="form.command.jenkins.source.id" placement="请选择链接">
                                                <el-option v-for="(dic_v,dic_k) in dic_jenkins_source" :label="dic_v.name" :value="dic_v.id"></el-option>
                                            </el-select>
                                        </el-form-item>
                                        <el-form-item label="项目">
                                            <el-input v-model="form.command.jenkins.name" placeholder="jenkins job name"></el-input>
                                        </el-form-item>
                                        <el-form-item label="参数" class="http_header_box">
                                            <el-input v-for="(header_v,header_i) in form.command.jenkins.params" v-model="header_v.value" placeholder="参数值">
                                                <el-input v-model="header_v.key" slot="prepend" placeholder="参数名" @input="jenkinsParamInput"></el-input>
                                                <el-button slot="append" icon="el-icon-delete" @click="jenkinsParamDel(header_i)"></el-button>
                                            </el-input>
                                        </el-form-item>
                                    </el-tab-pane>
                                    <el-tab-pane label="git" name="6">
                                        <el-form-item label="链接">
                                            <el-select v-model="form.command.git.link_id" placement="请选择链接">
                                                <el-option v-for="(dic_v,dic_k) in dic_git_source" :label="dic_v.name" :value="dic_v.id"></el-option>
                                            </el-select>
                                        </el-form-item>
                                        <el-form-item label="事件">
                                            <div>
                                                <el-button type="text" @click="gitSetShow(-1)">添加<i class="el-icon-plus"></i></el-button>
                                            </div>
                                            <div style="overflow-y: auto;max-height: 420px;">
                                                <el-alert v-show="form.command.git.events.length==0" title="未添加事件，请添加。" type="info"></el-alert>
                                                <!--git 事件展示-->
                                                <div v-for="(event,git_index) in form.command.git.events" style="position: relative;line-height: 133%;background: #f4f4f5;margin-bottom: 5px;padding: 6px 20px 7px 8px;border-radius: 3px;">
                                                    <el-row v-html="event.summary" style="min-height: 45px;"></el-row>
                                                    <i class="el-icon-close" style="font-size: 15px;position: absolute;top: 2px;right: 2px;cursor:pointer" @click="gitSetDel(git_index)"></i>
                                                    <i class="el-icon-edit" style="font-size: 15px;position: absolute;top: 23px;right: 2px;cursor:pointer" @click="gitSetShow(git_index,event)"></i>
                                                    <i style="position: absolute;right: 1px;top: 40px;font-size: 16px;">#{{git_index}}</i>
                                                </div>
                                            </div>
                                        </el-form-item>
                                    </el-tab-pane>
                                </el-tabs>
                            </el-form-item>
                            <el-form-item label="备注" label-width="46px">
                                <el-input v-model="form.remark"></el-input>
                            </el-form-item>
                            <el-form-item label-width="2px">
                                <div><el-button type="text" @click="msgBoxShow(-1)">推送<i class="el-icon-plus"></i></el-button></div>
                                <div v-for="(msg,msg_index) in form.msg_set" style="position: relative;max-height: 200px;line-height: 133%;background: #f4f4f5;margin-bottom: 10px;padding: 6px 20px 7px 8px;border-radius: 3px;">
                                    <el-row v-html="msg.descrition"></el-row>
                                    <i class="el-tag__close el-icon-close" style="font-size: 15px;position: absolute;top: 2px;right: 2px;cursor:pointer" @click="msgSetDel(msg_index)"></i>
                                    <i class="el-icon-edit" style="font-size: 15px;position: absolute;top: 23px;right: 2px;cursor:pointer" @click="msgBoxShow(msg_index,msg)"></i>
                                </div>
                            </el-form-item>
                        </el-form>
                        <div slot="footer" class="dialog-footer">
                            <el-button @click="configRun()" class="left">执行一下</el-button>
                            <el-button @click="setConfigShow = false">取 消</el-button>
                            <el-button type="primary" @click="setCron()">确 定</el-button>
                        </div>
                    </el-dialog>
                    <!-- 任务日志弹窗 -->
                    <el-drawer :title="configLog.title" :visible.sync="configLog.show" direction="rtl" size="40%" wrapperClosable="false" :before-close="configLogBoxClose">
                        <my-config-log :tags="configLog.tags"></my-config-log>
                    </el-drawer>
                    <!-- 注册任务列表弹窗 -->
                    <el-drawer title="已注册任务" :visible.sync="registerListShow" direction="rtl" size="40%" wrapperClosable="false">
                        <el-table :data="registerList">
                            <el-table-column property="id" label="编号">
                                <template slot-scope="scope">
                                    {{scope.row.id}}/{{scope.row.entry_id}}
                                </template>
                            </el-table-column>
                            <el-table-column property="spec" label="执行时间"></el-table-column>
                            <el-table-column property="update_dt" label="下一次执行"></el-table-column>
                            <el-table-column property="name" label="任务名称"></el-table-column>
                            <el-table-column property="" label="">
                                <template slot-scope="scope">
                                    <el-button type="text" @click="configLogBox(scope.row)">日志</el-button>
                                </template>
                            </el-table-column>
                        </el-table>
                    </el-drawer>
                    <!-- sql 设置弹窗 -->
                    <el-dialog :title="'sql设置-'+sqlSet.title" :visible.sync="sqlSet.show" :show-close="false" :close-on-click-modal="false">
                        <el-form :model="sqlSet.statement.git" v-if="sqlSet.source=='git'" label-width="70px" size="small">
                            <el-form-item label="连接">
                                <el-select v-model="sqlSet.statement.git.link_id" placement="请选择git链接">
                                    <el-option v-for="(dic_v,dic_k) in dic_git_source" :label="dic_v.name" :value="dic_v.id"></el-option>
                                </el-select>
                            </el-form-item>
                            <el-form-item label="仓库空间">
                                <el-input v-model="sqlSet.statement.git.owner" placeholder="仓库所属空间地址(企业、组织或个人的地址path)"></el-input>
                            </el-form-item>
                            <el-form-item label="项目名称">
                                <el-input v-model="sqlSet.statement.git.project" placeholder="仓库路径"></el-input>
                            </el-form-item>
                            <el-form-item label="文件路径">
                                <el-input v-for="(path_v,path_i) in sqlSet.statement.git.path" v-model="sqlSet.statement.git.path[path_i]" placeholder="文件的路径" @input="sqlGitPathInput">
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
                            <el-input type="textarea" rows="12" placeholder="请输入sql内容，批量添加时多个sql请用;分号分隔。" v-model="sqlSet.statement.local"></el-input>
                        </el-form>
                       
                        <span slot="footer" class="dialog-footer">
                            <el-button @click="sqlSet.show = false">取 消</el-button>
                            <el-button type="primary" @click="sqlSetConfirm()">确 定</el-button>
                        </span>
                    </el-dialog>
                    <!-- git 设置弹窗 -->
                    <el-dialog :title="'git设置-'+sqlSet.title" :visible.sync="gitSet.show" :show-close="false" :close-on-click-modal="false">
                        <el-form :model="gitSet.event" label-width="90px" size="small">
                            <el-form-item label="事件">
                                <el-select v-model="gitSet.data.id" placement="请选择事件类型">
                                    <el-option v-for="(dic_v,dic_k) in dic_git_event" :label="dic_v.name" :value="dic_v.id"></el-option>
                                </el-select>
                            </el-form-item>
                            <el-divider><i class="el-icon-setting"></i></el-divider>
                            
                            <div v-if="gitSet.data.id==2"> pr创建 开发中 ...
                                
                            </div>
                            
                            <div v-if="gitSet.data.id==9"> <!-- pr合并 -->
                                <el-form-item>
                                    <el-alert title="合并pr前请确认已完成了审查与测试确认！" type="warning" show-icon :closable="false"></el-alert>
                                </el-form-item>
                                <el-form-item label="仓库空间*">
                                    <el-input v-model="gitSet.data.pr_merge.owner" placeholder="仓库所属空间地址(企业、组织或个人的地址path)"></el-input>
                                </el-form-item>
                                <el-form-item label="项目名称*">
                                    <el-input v-model="gitSet.data.pr_merge.repo" placeholder="仓库路径"></el-input>
                                </el-form-item>
                                <el-form-item label="PR 编号*">
                                    <el-input type="number" v-model="gitSet.data.pr_merge.number" placeholder="本仓库PR的序数"></el-input>
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
                            <div v-if="gitSet.id==3">
                                待开发...
                            </div>
                        </el-form>
                       
                        <span slot="footer" class="dialog-footer">
                            <el-button @click="gitSet.show = false">取 消</el-button>
                            <el-button type="primary" @click="gitSetConfirm()">确 定</el-button>
                        </span>
                    </el-dialog>
            
                    <!-- sql链接源管理弹窗 -->
                    <el-drawer title="sql链接管理" :visible.sync="sqlSourceBoxShow" size="40%" wrapperClosable="false" :before-close="sqlSourceBox">
                        <my-source></my-source>
                    </el-drawer>
                    <!-- 推送设置弹窗 -->
                    <el-dialog title="推送设置" :visible.sync="msgSet.show" :show-close="false" :close-on-click-modal="false">
                        <el-form :model="msgSet" :inline="true" size="mini">
                            <el-form-item label="当结果">
                                <el-select v-model="msgSet.data.status" style="width: 90px">
                                    <el-option v-for="(dic_v,dic_k) in msgSet.statusList" :label="dic_v.name" :value="dic_v.id"></el-option>
                                </el-select>
                                时
                            </el-form-item>
                            <el-form-item label="发送">
                                <el-select v-model="msgSet.data.msg_id">
                                    <el-option v-for="(dic_v,dic_k) in dic_msg" :label="dic_v.name" :value="dic_v.id"></el-option>
                                </el-select>
                                消息
                            </el-form-item>
                            <el-form-item label="并且@用户">
                                <el-select v-model="msgSet.data.notify_user_ids" multiple="true">
                                    <el-option v-for="(dic_v,dic_k) in dic_user" :key="dic_v.id" :label="dic_v.name" :value="dic_v.id"></el-option>
                                </el-select>
                            </el-form-item>
                        </el-form>
                        <span slot="footer" class="dialog-footer">
                            <el-button @click="msgSet.show = false">取 消</el-button>
                            <el-button type="primary" @click="msgSetConfirm()">确 定</el-button>
                        </span>
                    </el-dialog>
                </el-main>`,
    name: "MyConfig",
    data(){
        return {
            dic_sql_source:[],
            dic_sql_driver:[],
            dic_user: [],
            dic_msg: [],
            dic_jenkins_source: [],
            dic_git_source: [],
            dic_git_event:[],
            dic_host_source: [],
            dic_cmd_type: [],
            sys_info:{},
            list: [],
            listPage:{
                total:0,
                page: 1,
                size: 10,
            },
            listParam:{
                type: 1,
                page: 1,
                size: 15,
            },
            listRequest: false, // 请求中标志
            configLog:{
                show: false,
                id: 0,
                title:""
            },
            registerList: [],
            registerListShow: false,
            setConfigShow: false, // 新增弹窗
            setConfigTitle: '',
            sqlSourceBoxShow: false,
            sqlSet: {  // sql 设置弹窗
                show: false, // 是否显示
                title: '添加',
                index: -1, // 操作行号
                source: '', // sql 来源
                statement:{}
            },
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
            msgSet:{
                show: false, // 是否显示
                title: '添加',
                index: -1, // 操作行号
                data: {}, // 实际内容
                statusList:[{id:1,name:"错误"}, {id:2, name:"成功"}], // {id:0,name:"完成"} 有歧义，以后待定
            },
            form:{},
            hintSpec: "* * * * * *",
            labelType: '1',
            // 日期选择器设置
            pickerOptions: {
                disabledDate(time){
                    return time.getTime() < Date.now() - 8.64e7
                },
                selectableRange: "00:00:00 - 23:01:59",
            },

        }
    },
    // 模块初始化
    created(){
        this.form = this.initFormData()
        api.systemInfo((res)=>{
            this.sys_info = res;
        })
        this.getDicSqlSource()
        console.log("config created")
    },
    // 模块初始化
    mounted(){
        console.log("config mounted")
        this.getList()
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
                    case 9:
                        this.gitSet.data = {
                            id: v,
                            pr_merge: {
                                owner: "", // 空间
                                repo: "", // 仓库
                                number: null,
                                title: "",
                                description: "",
                                merge_method: "merge",
                                prune_source_branch: false,
                            }
                        }
                        break
                    default:
                        this.gitSet.data = {id:''}
                }
            }
        }
    },

    // 具体方法
    methods:{
        // 任务列表
        getList(){
            if (this.listRequest){
                return this.$message.info('请求执行中,请稍等.');
            }
            this.listRequest = true
            api.innerGet("/config/list", this.listParam, (res)=>{
                this.listRequest = false
                if (!res.status){
                    console.log("config/list 错误", res)
                    return this.$message.error(res.message);
                }
                for (i in res.data.list){
                    let ratio = 0
                    if (res.data.list[i].top_number){
                        ratio = res.data.list[i].top_error_number / res.data.list[i].top_number
                    }
                    if (res.data.list[i].command.sql){
                        res.data.list[i].command.sql.err_action = res.data.list[i].command.sql.err_action.toString()
                        res.data.list[i].command.sql.interval = res.data.list[i].command.sql.interval.toString()
                    }
                    res.data.list[i].status = res.data.list[i].status.toString()
                    res.data.list[i].topRatio = 100 - ratio * 100
                }
                this.list = res.data.list;
                this.listPage = res.data.page;
            })
        },
        handleSizeChange(val) {
            console.log(`每页 ${val} 条`);
        },
        handleCurrentChange(val) {
            this.listParam.page = val
            this.getList()
        },
        handleClickTypeLabel(tab, event) {
            this.listParam.type = tab
            this.getList()
        },
        // 添加/编辑 任务
        setCron(){
            if (this.form.name == ''){
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

            api.innerPost("/config/set", body, (res)=>{
                if (!res.status){
                    console.log("config/set 错误", res)
                    return this.$message.error(res.message)
                }
                this.setConfigShow = false;
                this.form = this.initFormData()
                this.getList(); // 刷新当前页
            })
        },
        // 添加弹窗
        createShow(){
            this.setConfigShow = true
            this.setConfigTitle = '添加任务'
            // 应该有个地方定义空结构体
            this.form = this.initFormData()
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
        initFormData(){
            return  {
                type: this.listParam.type.toString(),
                protocol: '3',
                var_fields: [{}],
                command:{
                    http:{
                        method: 'GET',
                        header: [{"key":"","value":""}],
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
                                path: [""], // 目前不支持多选，但数据结构要能为后续支持扩展
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
                        params: [{"key":"","value":""}]
                    },
                    git:{
                        link_id: "",
                        events: [],
                    }
                },
                msg_set: []
            }
        },
        // 编辑弹窗
        editShow(row){
            this.setConfigShow = true
            this.setConfigTitle = '编辑任务'
            this.form = copyJSON(row)

            if (this.form.var_fields == null || this.form.var_fields.length == 0){
                this.form.var_fields = [{}]
            }
            if (row.command.cmd.statement.git.path == null){
                this.form.command.cmd.statement.git.path = this.initFormData().command.cmd.statement.git.path
            }

            if (this.form.command.sql.driver == ""){ // 历史数据兼容
                this.form.command.sql = this.initFormData().command.sql
            }
            if (this.form.command.sql.source.id == 0){
                this.form.command.sql.source.id = ""
            }
            for (let i in row.command.sql.statement){
                this.form.command.sql.statement[i] = this.sqlGitBuildDesc(row.command.sql.statement[i])
            }

            if (this.form.command.http.header.length == 0){
                this.form.command.http.header = this.initFormData().command.http.header
            }
            if (this.form.command.jenkins.source.id == 0){
                this.form.command.jenkins.source.id = ""
            }
            if (this.form.command.jenkins.params.length == 0){
                this.form.command.jenkins.params = this.initFormData().command.jenkins.params
            }
            if (this.form.command.cmd.host.id == 0){
                this.form.command.cmd.host.id = this.initFormData().command.cmd.host.id
            }
            if (this.form.command.git){
                for (let i in this.form.command.git.events){
                    this.form.command.git.events[i] = this.gitBuildDesc(this.form.command.git.events[i])
                }
                if (this.form.command.git.link_id == 0){
                    this.form.command.git.link_id = ''
                }
            }
            for (let i in row.msg_set){
                this.form.msg_set[i] = this.msgSetBuildDesc(row.msg_set[i])
            }
            if (this.form.msg_set == null){

            }
            this.form.type = this.form.type.toString()
            this.form.status = this.form.status.toString() // 这里要转字符串，否则可能显示数字
            this.form.protocol = this.form.protocol.toString()
            console.log("编辑：",this.form)
        },
        // 改变状态
        changeStatus(row, newStatus){
            this.$confirm(newStatus==1? '确认关闭任务': '确认开启任务', '提示',{
                type: 'warning',
            }).then(()=>{
                // 确认操作
                api.innerPost("/config/change_status", {id:row.id,status:Number(newStatus)}, (res)=>{
                    if (!res.status){
                        return this.$message.error(res.message)
                    }
                    row.status = newStatus.toString()
                    row.status_name = newStatus == 1 ? '停用' :'激活'
                    return this.$message.success(res.message)
                })
            }).catch(()=>{
                // 取消操作
            })
        },
        // 改变类型
        // changeType: function (){
        //     if (this.form.type == 2 || this.form.type == '2'){
        //         this.hintSpec = "YYYY-mm-dd HH:MM:SS"
        //     }else{
        //         this.hintSpec = "* * * * * *"
        //     }
        // },

        getRegisterList(){
            this.registerListShow = true;
            api.innerGet("/config/register_list", {}, (res)=>{
                if (!res.status){
                    console.log("config/register_list 错误", res)
                    return this.$message.error(res.message);
                }
                this.registerList = res.data.list;
            })
        },
        // 成功率图标
        getTopIcon(total, ratio){
            // 没有执行
            if (total == 0){
                return "icon el-icon-remove-outline" // 0
            }
            // 成功率100
            if (ratio == 100){
                return "icon el-icon-circle-check success"
            }
            // 成功率66以上
            if (ratio < 100 && ratio > 66){
                return "icon el-icon-warning-outline warning"
            }
            // 成功率33以上
            if (ratio < 66 && ratio > 33){
                return "icon el-icon-circle-plus-outline danger"
            }
            // 成功率33一下
            return "icon el-icon-circle-close danger"
        },
        // sql source box
        sqlSourceBox(show){
            this.sqlSourceBoxShow = show == true;
            if (!this.sqlSourceBoxShow){
                this.getDicSqlSource() // 关闭弹窗要重载枚举
            }
        },
        configLogBox(item){
            let tags = {ref_id:item.id, component:"job"}
            if (item.protocol == 99){
                tags.component = 'pipeline'
            }
            this.configLog.id = item.id
            this.configLog.tags = tags
            this.configLog.title = item.name+' 日志'
            this.configLog.show = true
        },
        configLogBoxClose(done){
            this.configLog.show = false;
            this.configLog.id = 0;
            this.configLog.title = ' 日志'
            this.configLog.tags = {}
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
                    owner: "",
                    project: "",
                    path: [""],
                    ref: "",
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
                    return this.$message.error("仓库空间为必填")
                }
                if (this.sqlSet.statement.git.project == ""){
                    return this.$message.error("项目名称为必填")
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
            }else if (this.gitSet.data.id == 9){
                if (this.gitSet.data.pr_merge.owner == ""){
                    return this.$message.error("仓库空间为必填")
                }
                if (this.gitSet.data.pr_merge.repo == ""){
                    return this.$message.error("项目名称为必填")
                }
                if (!this.gitSet.data.pr_merge.number){
                    return this.$message.error("仓库PR编号为必填")
                }
                if (this.gitSet.data.pr_merge.merge_method == ""){
                    return this.$message.error("合并方式不得为空")
                }
                this.gitSet.data.pr_merge.number = Number(this.gitSet.data.pr_merge.number)
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

        // 推送弹窗
        msgBoxShow(index, oldData){
            if (index === "" || index == null || isNaN(index)){
                console.log('msgSetShow', index, oldData)
                return this.$message.error("索引位标志异常");
            }
            if (oldData == undefined || index < 0){
                oldData = {
                    status: 1,
                    msg_id: "",
                    notify_user_ids: [],
                }
            }else if (typeof oldData != 'object'){
                console.log('推送信息异常', oldData)
                return this.$message.error("推送信息异常");
            }
            this.msgSet.show = true
            this.msgSet.index = Number(index)  // -1.新增、>=0.具体行的编辑
            this.msgSet.title = this.msgSet.index < 0? '添加' : '编辑';
            this.msgSet.data = oldData
        },
        msgSetDel(index){
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
            if (this.msgSet.data.msg_id <= 0){
                return this.$message.warning("请选择消息模板");
            }
            let data = this.msgSetBuildDesc(this.msgSet.data)

            if (this.msgSet.index < 0){
                this.form.msg_set.push(data)
            }else{
                this.form.msg_set[this.msgSet.index] = data
            }
            this.msgSet.show = false
            this.msgSet.index = -1
            this.msgSet.data = {}
        },
        // 构建消息设置描述
        msgSetBuildDesc(data){
            let item1 = this.msgSet.statusList.find(option => option.id === data.status);
            if (item1){
                data.status_name = item1.name
            }
            let descrition = '当任务<span class="el-tag el-tag--small el-tag--light">'+item1.name+'</span>时'

            let item2 = this.dic_msg.find(option => option.id === data.msg_id)
            if (item2){
                data.msg_name = item2.name
                descrition += '，发送<span class="el-tag el-tag--small el-tag--light">'+item2.name+'</span>消息'
            }
            let item3 = this.dic_user.filter((option) => {
                return data.notify_user_ids.includes(option.id);
            }).map((item)=>{return item.name})
            if (item3.length > 0){
                data.notify_users_name = item3
                descrition += '，并且@人员<span class="el-tag el-tag--small el-tag--light">'+data.notify_users_name+'</span>'
            }
            data.descrition = descrition
            return data
        },
        // 构建git sql 描述
        sqlGitBuildDesc(data){
            if (data.git == null){
                return data
            }
            let git = this.dic_git_source.filter(item=>{
                return item.id == data.git.link_id
            })

            data.git.descrition = '连接<span class="el-tag el-tag--small el-tag--light">'+git.length>0 ? git[0].name: '' +
                '</span> 访问<span class="el-tag el-tag--small el-tag--light">'+data.git.owner+'/'+data.git.project +'</span>'+
                '</span> 引用<span class="el-tag el-tag--small el-tag--light">'+data.git.ref +'</span> 拉取以下文件内容'+
                '<span class="el-tag el-tag--small el-tag--light">'+ (data.is_batch==1? '批量解析后' :'单文件单sql') +'</span>执行'
            return data
        },
        // 构建 git 描述
        gitBuildDesc(data){
            let left = '<span class="el-tag el-tag--small el-tag--light">'
            let right = '</span>'
            switch (data.id){
                case 2:
                    data.summary = '完善中...'
                    break
                case 9:
                    data.summary = `<b>pr合并</b> ${left}${data.pr_merge.owner}/${data.pr_merge.repo}${right} pr ${left}${data.pr_merge.number}${right} ${left}${this.gitSet.gitMergeTypeList[data.pr_merge.merge_method]}${right}  ${data.pr_merge.prune_source_branch===true?left+'删除提交分支'+right:''}`+
                        `<br><i style="margin-left: 3em;"></i><b>${data.pr_merge.title}</b> ${data.pr_merge.description}`
                    break
                default:
                    data.summary = '未支持的事件类型'
            }
            return data
        },
        statusClass(status){
            switch (Number(status)) {
                case 1:
                    return 'warning';
                case 2:
                    return 'primary';
                case 3:
                    return 'success';
                case 4:
                    return 'danger';
            }
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
                Enum.dicMsg
            ]
            api.dicList(types,(res) =>{
                this.dic_sql_source = res[Enum.dicSqlSource]
                this.dic_jenkins_source = res[Enum.dicJenkinsSource]
                this.dic_git_source =res[Enum.dicGitSource]
                this.dic_git_event = res[Enum.dicGitEvent]
                this.dic_host_source =res[Enum.dicHostSource]
                this.dic_user = res[Enum.dicUser]
                this.dic_msg = res[Enum.dicMsg]
                this.dic_cmd_type = res[Enum.dicCmdType]
                this.dic_sql_driver = res[Enum.dicSqlDriver]
            })
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
        // 执行一下
        configRun(){

            // 主要是强制类型
            let body = copyJSON({
                id: this.form.id,
                name: this.form.name,
                type: Number(this.form.type),
                spec: this.form.spec,
                protocol: Number(this.form.protocol),
                command: this.form.command,
                remark: this.form.remark,
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
        }
    }
})

Vue.component("MyConfig", MyConfig);