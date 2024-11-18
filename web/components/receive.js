var MyReceive = Vue.extend({
    template: `<el-container>
        <!--边栏-->
    <el-aside style="padding-top: 10px">
        <el-card class="aside-card">
            <div slot="header">
                <span class="h3">执行任务</span>
            </div>
            <ol>
                <li v-for="item in queue.exec">
                    <router-link :to="{path:'/config_detail', query:{id:item.ref_id, type:item.ref_type, entry_id:item.entry_id}}" class="el-link el-link--default is-underline">{{item.name}}</router-link>
                    <p style="margin: 0;color: #909399;line-height: 100%;font-size: 12px;">
                        ({{durationTransform(item.duration, 's')}}) 
                        <el-popconfirm :title="'确定停用 '+item.name+' 任务吗？'" @confirm="jobStop(item)"><i slot="reference" class="el-icon-circle-close stop"></i></el-popconfirm>
                    </p>
                </li>
            </ul>
            <div v-show="!queue.exec.length">-</div>
        </el-card>
        <el-card class="aside-card">
            <div slot="header">
                <span class="h3">注册任务</span>
            </div>
            <ol>
                <li v-for="item in queue.register">
                    <router-link :to="{path:'/config_detail', query:{id:item.ref_id, type:item.ref_type, entry_id:item.entry_id}}" class="el-link el-link--default is-underline">{{item.name}}</router-link>
                </li>
            </ol>
            <div v-show="!queue.register.length">-</div>
        </el-card>
    </el-aside>
    <!--主内容-->
    <el-main>
        <el-menu class="el-menu-demo" mode="horizontal" @select="handleClickTypeLabel">
            <div style="float: right">
                <el-button type="text" @click="setShow()" v-if="$auth_tag.receive_set">添加接收规则</el-button>
            </div>
        </el-menu>
        <el-row>
            <el-form :inline="true" :model="list.param" size="small" class="search-form">
                <el-form-item label="名称">
                    <el-input v-model="list.param.name" placeholder="搜索名称"></el-input>
                </el-form-item>
                <el-form-item label="状态">
                    <el-select v-model="list.param.status" placeholder="所有" multiple>
                        <el-option v-for="item in dic.config_status" :label="item.name" :value="item.id"></el-option>
                    </el-select>
                </el-form-item>
                <el-form-item label="处理人">
                    <el-select v-model="list.param.handle_user_ids" placeholder="所有" multiple>
                        <el-option v-for="item in dic.user" :label="item.name" :value="item.id"></el-option>
                    </el-select>
                </el-form-item>
                <el-form-item label="创建人">
                    <el-select v-model="list.param.create_user_ids" placeholder="所有" multiple>
                        <el-option v-for="item in dic.user" :label="item.name" :value="item.id"></el-option>
                    </el-select>
                </el-form-item>
                <el-form-item>
                    <el-button type="primary" @click="getList">查询</el-button>
                </el-form-item>
            </el-form>
        </el-row>
        <el-table :data="list.items" @cell-mouse-enter="listTableCellMouse" @cell-mouse-leave="listTableCellMouse">
            <el-table-column prop="" label="成功率" width="80">
                <template slot-scope="scope">
                    <el-tooltip placement="top-start">
                        <div slot="content">{{scope.row.topRatio}}%<br/>近期{{scope.row.top_number}}次执行，{{scope.row.top_error_number}}次失败。</div>
                        <i :class="getTopIcon(scope.row.top_number, scope.row.topRatio)"></i>
                    </el-tooltip>
                </template>
            </el-table-column>
            <el-table-column prop="spec" label="执行时间" width="160"></el-table-column>
            <el-table-column prop="name" label="任务名称">
                <div slot-scope="{row}" style="display: flex;">
                    <span style="white-space: nowrap;overflow: hidden;text-overflow: ellipsis;">
                        <router-link :to="{path:'/config_detail',query:{id:row.id, type:'receive'}}" class="el-link el-link--primary is-underline" :title="row.name">{{row.name}}</router-link>
                    </span>
                    <span v-show="row.option.name.mouse" style="margin-left: 4px;white-space: nowrap;">
                        <i  class="el-icon-edit hover" @click="setShow(row)" title="编辑"></i>
                        <i class="el-icon-notebook-2 hover" @click="configLogBox(row)" title="日志"></i>
                    </span>
                </div>
            </el-table-column>
            <el-table-column prop="" label="状态" width="100">
                <template slot-scope="scope">
                    <el-tooltip placement="top-start">
                        <div slot="content">{{scope.row.status_dt}}  {{scope.row.status_remark}}</div>
                        <el-button :type="statusTypeName(scope.row.status)" plain size="mini" round @click="statusShow(scope.row, 'receive')">{{scope.row.status_name}}</el-botton>
                    </el-tooltip>
                </template>
            </el-table-column>
            <el-table-column prop="remark" label="备注"></el-table-column>
            <el-table-column prop="handle_user_names" label="处理人" width="120"></el-table-column>
            <el-table-column prop="create_user_name" label="创建人" width="80"></el-table-column>
        </el-table>
        <el-pagination
                @size-change="handleSizeChange"
                @current-change="handleCurrentChange"
                :current-page.sync="list.page.page"
                :page-size="list.page.size"
                layout="total, prev, pager, next"
                :total="list.page.total">
        </el-pagination>
        
        
        <!-- 流水线设置表单 -->
        <el-drawer :title="set_box.title" :visible.sync="set_box.show" size="60%" wrapperClosable="false">
            <my-receive-form v-if="set_box.show" :request="{detail:set_box.detail}" @close="formClose"></my-receive-form>
        </el-drawer>
        <!-- 任务日志弹窗 -->
        <el-drawer :title="config_log_box.title" :visible.sync="config_log_box.show" direction="rtl" size="40%" wrapperClosable="false" :before-close="configLogBoxClose">
            <my-config-log :tags="config_log_box.tags"></my-config-log>
        </el-drawer>
        <!--状态变更弹窗-->
        <my-status-change v-if="status_box.show" :request="status_box" @close="statusShow"></my-status-change>
    </el-main>
</el-container>
`,
    name: "MyReceive",
    data(){
        return {
            env: {},
            dic:{
                user: [],
                msg: [],
                config_status: [],
            },
            sys_info:{},
            // 列表数据
            list: {
                items: [],
                page: {
                    total:0,
                    page: 1,
                    size: 20,
                },
                param:{
                    type: 2,
                    page: 1,
                    size: 20,
                    status: [],
                    handle_user_ids:[],
                    create_user_ids:[],
                    name: '',
                },
                request: false, // 请求中标志
            },
            // 队列
            queue:{
                exec:[], // 执行队列
                register:[], // 注册队列
            },
            set_box: {
                show:false,
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
    created(){
        setDocumentTitle('接收规则管理')
        this.getDic()
        api.systemInfo((res)=>{
            this.sys_info = res;
        })
        this.loadParams(getHashParams(window.location.hash))
    },
    // 模块初始化
    mounted(){
        this.getList()
        // 添加指定事件监听
        this.env = cache.getEnv()
        // this.$sse.addEventListener(this.env.env+".exec.queue", this.execQueue)
        // this.$sse.addEventListener(this.env.env+'.register.queue', this.registerQueue)
    },
    beforeDestroy(){
        // 销毁指定事件监听
        // this.$sse.removeEventListener(this.env.env+".exec.queue", this.execQueue)
        // this.$sse.removeEventListener(this.env.env+".register.queue", this.registerQueue)
    },
    // 具体方法
    methods:{
        loadParams(param){
            if (typeof param !== 'object'){return}
            if (param.type){this.list.param.type = Number(param.type)}
            if (param.page){this.list.param.page = Number(param.page)}
            if (param.size){this.list.param.size = Number(param.size)}
            if (param.name){this.list.param.name = param.name.toString()}
            if (param.status){this.list.param.status = param.status.map(Number)}
            if (param.handle_user_ids){this.list.param.handle_user_ids = param.handle_user_ids.map(Number)}
            if (param.create_user_ids){this.list.param.create_user_ids = param.create_user_ids.map(Number)}
        },
        // 任务列表
        getList(){
            if (this.list.request){
                return this.$message.info('请求执行中,请稍等.');
            }
            replaceHash('/receive', this.list.param)
            this.list.request = true
            api.innerGet("/receive/list", this.list.param, (res)=>{
                this.list.request = false
                if (!res.status){
                    return this.$message.error(res.message);
                }
                for (i in res.data.list){
                    let ratio = 0
                    if (res.data.list[i].top_number){
                        ratio = res.data.list[i].top_error_number / res.data.list[i].top_number
                    }
                    res.data.list[i].status = res.data.list[i].status.toString()
                    res.data.list[i].topRatio = 100 - ratio * 100
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
                    // 前端用设置
                    res.data.list[i].option = {
                        name:{
                            mouse:false
                        },
                    }
                }
                this.list.items = res.data.list;
                this.list.page = res.data.page;
            })
        },
        handleSizeChange(val) {
            console.log(`每页 ${val} 条`);
        },
        handleCurrentChange(val) {
            this.list.param.page = val
            this.getList()
        },
        handleClickTypeLabel(tab, event) {
            this.list.param.type = tab
            this.getList()
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
        // 状态弹窗
        statusShow(e, type=''){
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
        // 枚举
        getDic(){
            api.dicList([Enum.dicUser, Enum.dicMsg, Enum.dicConfigStatus],(res) =>{
                this.dic.user = res[Enum.dicUser]
                this.dic.msg = res[Enum.dicMsg]
                this.dic.config_status = res[Enum.dicConfigStatus]
            })
        },
        // 添加弹窗
        setShow(row=null){
            this.set_box.show = true
            if (row == null){
                this.set_box.title= '添加接收规则'
                this.set_box.detail = {}
            }else{
                this.set_box.title= '编辑接收规则'
                api.innerGet('/receive/detail', {id: row.id}, (res)=>{
                    if (!res.status){
                        return this.$message.error(res.message)
                    }
                    this.set_box.detail = res.data
                },{async:false})
            }
        },
        formClose(e){
            console.log("close",e)
            this.set_box.show = false
            if (e.is_change){
                this.getList()
            }
        },
        configLogBox(item){
            let tags = {ref_id:item.id, component:""}
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
        // 停止执行任务
        jobStop(row){
            api.innerPost("/job/stop",row, res=>{
                if (!res.status){
                    return this.$message.error(res.message)
                }
                this.$message.success('操作成功')
            })
        },
        // 消息监听处理
        execQueue(e){
            // console.log("execQueue", e)
            let data = JSON.parse(e.data)
            if (!data){
                this.queue.exec = []
                return
            }
            this.queue.exec = data
        },
        registerQueue(e){
            // console.log("registerQueue", e)
            let data = JSON.parse(e.data)
            if (!data){
                this.queue.register = []
                return
            }
            this.queue.register = data
        },
    }
})

Vue.component("MyReceive", MyReceive);