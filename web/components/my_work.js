var MyWork = Vue.extend({
    template: `<el-main class="my-work">
<!--
        先规划一下页面布局（整体仿造tapd，还有我们自己的手机端），
        # 各个环境名称（如果没有对应任务就不在哪上）（默认只有第一个标签页打开）
            流水线和任务混合在一起（sql写法要研究一下）
                点击跳转到详情页（新的标签页）
-->
        <el-card shadow="never" body-style="padding:12px">
            <el-button size="mini" round plain @click="tabChangeName('todo')" :class="options.tab=='todo'?'active': ''">待办</el-button>
            <el-button size="mini" round plain @click="tabChangeName('created')" :class="options.tab=='created'?'active': ''">创建</el-button>
            <el-button size="mini" round plain @click="tabChangeName('draft')" :class="options.tab=='draft'?'active': ''">草稿</el-button>
            
        </el-card>
        
        <el-row v-for="(group_v,group_k) in group_list">
            <el-row class="header">
                <span class="my-table-icon" @click="showTable(group_v)">
                    <i :class="group_v.show ? 'el-icon-caret-bottom' : 'el-icon-caret-right'"></i>
                </span>
                
                {{group_v.env_title}} <el-divider direction="vertical"></el-divider> 
                {{group_v.type=='config'? '任务' : '流水线'}} <span style="color: #72767b">（{{group_v.total}}）</span>
            </el-row>
            <el-row class="body" v-if="group_v.show">
                <el-table :data="group_v.list" v-loading="group_v.loading" @cell-mouse-enter="listTableCellMouse" @cell-mouse-leave="listTableCellMouse">
                    <el-table-column prop="spec" label="执行时间"></el-table-column>
                    <el-table-column prop="name" label="任务名称">
                        <div slot-scope="{row}" class="abc" style="display: flex;">
                            <span style="white-space: nowrap;overflow: hidden;text-overflow: ellipsis;">
                                <router-link target="_blank" :to="{path:'/config_detail',query:{id:row.id, type:group_v.type}}" class="el-link el-link--primary is-underline" :title="row.name">{{row.name}}</router-link>
                            </span>
                            <span v-show="row.option.name.mouse" style="margin-left: 4px;white-space: nowrap;">
                                <i  class="el-icon-edit hover" @click="setShow(row, group_v.type)" title="编辑"></i>
                                <i class="el-icon-notebook-2 hover" @click="configLogBox(row, group_v.type)" title="日志"></i>
                            </span>
                        </div>
                    </el-table-column>
                    <el-table-column prop="protocol_name" label="协议"></el-table-column>
                    <el-table-column prop="" label="状态">
                        <template slot-scope="scope">
                            <el-tooltip placement="top-start">
                                <div slot="content">{{scope.row.status_dt}}  {{scope.row.status_remark}}</div>
                                <el-button :type="statusTypeName(scope.row.status)" plain size="mini" round @click="statusShow(scope.row, group_v.type)">{{scope.row.status_name}}</el-botton>
                            </el-tooltip>
                        </template>
                    </el-table-column>
                    <el-table-column prop="remark" label="备注"></el-table-column>
                    <el-table-column prop="handle_user_names" label="处理人"></el-table-column>
                    <el-table-column prop="create_user_name" label="创建人"></el-table-column>
                </el-table>
            </el-row>
        </el-row>
        
        <el-dialog title="编辑任务" :visible.sync="detail_form_box.show && detail_form_box.type=='config'" :close-on-click-modal="false" class="config-form-box" :before-close="formClose">
            <my-config-form v-if="detail_form_box.show && detail_form_box.type==='config'" :request="{detail:detail_form_box.detail}" @close="formClose"></my-config-form>
        </el-dialog>
        <el-drawer title="编辑流水线" :visible.sync="detail_form_box.show && detail_form_box.type=='pipeline'" size="60%" :before-close="formClose">
            <my-pipeline-form v-if="detail_form_box.show && detail_form_box.type==='pipeline'" :request="{detail:detail_form_box.detail}" @close="formClose"></my-pipeline-form>
        </el-drawer>
        <!-- 任务日志弹窗 -->
        <el-drawer :title="config_log_box.title" :visible.sync="config_log_box.show" direction="rtl" size="40%" wrapperClosable="false" :before-close="configLogBoxClose">
            <my-config-log :tags="config_log_box.tags"></my-config-log>
        </el-drawer>
        <my-status-change v-if="status_box.show" :request="status_box" @close="statusShow"></my-status-change>
    </el-main>`,

    name: "MyWork",
    props: {
        data_id:Number
    },
    data(){
        return {
            dic:{user:[]},
            user: {},
            group_list: [],
            options:{
                tab: 'todo', // 默认待办
            },
            detail_form_box: {
                show: false,
                type: '',
                title: '',
                detail:{},
            },
            // 状态弹窗
            status_box:{
                show:false,
                detail:{},
                type: '',
            },
            config_log_box:{
                show: false,
                title:'',
                tags: {},
            },
        }
    },
    // 模块初始化
    created(){},
    // 模块初始化
    mounted(){
        setDocumentTitle('我的工作')
        this.user = cache.getUser()
        this.getDicList()
        this.getTables()
    },
    // 具体方法
    methods:{
        getTables(){
            api.innerGet('/work/table',{tab:this.options.tab}, res =>{
                if (!res.status){
                    return this.$message.error(res.message)
                }
                let list = []
                res.data.list.forEach((item,index)=>{
                    item.show = false
                    item.list = null
                    item.listpage = {}
                    item.loading = false
                    list.push(item)
                })
                list.some((item)=>{
                    if (item.total>0){
                        this.showTable(item)
                        return true
                    }
                })
                this.group_list = list
            },{async:false})
        },
        // 显示工作表
        showTable(row){
            row.show = !row.show
            let path = ""
            if (row.type == 'config'){
                path = "/config/list"
            }else if (row.type == 'pipeline'){
                path = "/pipeline/list"
            }else{
                return this.$message.warning('未支持的表格类型')
            }
            let param =  {page:1,size:10,env:row.env}
            switch (this.options.tab){
                case "todo":
                    param['status'] = [Enum.StatusAudited, Enum.StatusReject, Enum.StatusError]
                    param['handle_user_ids'] = [this.user.id]
                    break
                case "created":
                    param['create_user_ids'] = [this.user.id]
                    break
                case "draft": // 创建人或处理人
                    param['status'] = [Enum.StatusDisable]
                    param['create_or_handle_user_id'] = this.user.id
                    break
                default:
                    return this.$message.warning('未选择工作类型')
            }

            if (row.show && row.list == null){
                row.loading = true
                row.list = []
                api.innerGet(path, param, (res)=>{
                    row.loading = false
                    if (!res.status){
                        return this.$message.error(res.message);
                    }
                    for (i in res.data.list){
                        // 处理人
                        let handle_user_names = ''
                        if (res.data.list[i].handle_user_ids){
                            res.data.list[i].handle_user_ids.forEach((id)=>{
                                this.dic.user.forEach((item)=>{
                                    if (item.id == id){
                                        handle_user_names += item.name+','
                                    }
                                })
                            })
                            res.data.list[i].handle_user_names = handle_user_names.substring(0,handle_user_names.length-1)
                        }
                        // 创建人
                        this.dic.user.forEach((item)=>{
                            if (item.id == res.data.list[i].create_user_id){
                                res.data.list[i].create_user_name = item.name
                            }
                        })
                        // 前端用设置
                        res.data.list[i].option = {
                            name:{
                                mouse:false
                            },
                        }
                    }
                    row.list = res.data.list;
                    row.listPage = res.data.page;
                })
            }
        },
        // 状态切换
        tabChangeName(tab){
            if (this.options.tab == tab){
                return
            }
            this.options.tab = tab
            this.getTables()
        },
        // 枚举
        getDicList(){
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
                this.dic = {
                    // host_source: res[Enum.dicHostSource],
                    // cmd_type: res[Enum.dicCmdType],
                    // sql_source: res[Enum.dicSqlSource],
                    // git_source: res[Enum.dicGitSource],
                    // jenkins_source: res[Enum.dicJenkinsSource],
                    // msg: res[Enum.dicMsg],
                    user: res[Enum.dicUser],
                }
            })
        },
        // 添加弹窗
        setShow(row=null,type=''){
            this.detail_form_box.show = true
            this.detail_form_box.type = type
            this.detail_form_box.title= '编辑任务'
            api.innerGet('/'+type+'/detail', {id: row.id}, (res)=>{
                if (!res.status){
                    return this.$message.error(res.message)
                }
                this.detail_form_box.detail = res.data
            },{async:false})

        },
        formClose(e){
            console.log("close",e)
            this.detail_form_box.show = false
            if (e.is_change){
                this.getList()
            }
        },
        // 列表表格列鼠标悬浮事件
        listTableCellMouse(row, column, cell, event){
            // console.log("列表", row, column, cell, event)
            if (row.option[column.property]){
                if (event.type === 'mouseenter'){
                    row.option[column.property].mouse = true
                }else if (event.type === 'mouseleave'){
                    row.option[column.property].mouse = false
                }
            }
        },
        statusShow(e, type){
            this.status_box.show = e.id > 0
            if (this.status_box.show){
                this.status_box.detail = e
                this.status_box.type = type
            }else{
                this.status_box.detail = {}
                this.status_box.type = ''
            }
            if (typeof e == 'object' && e.is_change===true){ // 数据变更了重新加载
                this.getList()
            }
        },
        configLogBox(item, type){
            let tags = {ref_id:item.id, component:type}
            this.config_log_box.tags = tags
            this.config_log_box.title = item.name+' 日志'
            this.config_log_box.show = true
        },
        configLogBoxClose(done){
            this.config_log_box.show = false;
            this.config_log_box.id = 0;
            this.config_log_box.title = ' 日志'
            this.config_log_box.tags = {}
        },
    }
})

Vue.component("MyWork", MyWork);