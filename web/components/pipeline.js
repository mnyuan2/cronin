var MyPipeline = Vue.extend({
    template: `<el-container>
        <!--边栏-->
    <el-aside width="240px" style="padding-top: 10px">
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
        <el-menu :default-active="labelType" class="el-menu-demo" mode="horizontal" @select="handleClickTypeLabel">
    <!--                        <el-menu-item index="1" :disabled="listRequest">周期任务</el-menu-item>-->
    <!--                        <el-menu-item index="2" :disabled="listRequest">单次任务</el-menu-item>-->
            <div style="float: right">
                <el-button type="text" @click="setShow()">添加流水线</el-button>
    <!--                            <el-button type="text" @click="getRegisterList()">已注册任务</el-button>-->
            </div>
        </el-menu>
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
                <template slot-scope="{row}">
                    <router-link :to="{path:'/config_detail',query:{id:row.id, type:'pipeline'}}" class="el-link el-link--primary is-underline">{{row.name}}</router-link>
                    <span v-show="row.option.name.mouse" style="margin-left: 4px"><i class="el-icon-edit hover" @click="setShow(row)"></i></span>
                </template>
            </el-table-column>
            <el-table-column prop="" label="状态" width="100">
                <template slot-scope="scope">
                    <el-tooltip placement="top-start">
                        <div slot="content">{{scope.row.status_dt}}  {{scope.row.status_remark}}</div>
                        <el-button :type="statusTypeName(scope.row.status)" plain size="mini" round @click="statusShow(scope.row, 'pipeline')">{{scope.row.status_name}}</el-botton>
                    </el-tooltip>
                </template>
            </el-table-column>
            <el-table-column prop="remark" label="备注"></el-table-column>
            <el-table-column prop="handle_user_names" label="处理人"></el-table-column>
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
        <el-drawer :title="add_box.title" :visible.sync="add_box.show" size="60%" wrapperClosable="false">
            <my-pipeline-form v-if="add_box.show" :request="{detail:add_box.detail}" @close="formClose"></my-pipeline-form>
        </el-drawer>
    
        <!--状态变更弹窗-->
        <my-status-change v-if="status_box.show" :request="status_box" @close="statusShow"></my-status-change>
    </el-main>
</el-container>
`,
    name: "MyPipeline",
    data(){
        return {
            env: {},
            dic:{
                user: [],
                msg: [],
            },
            sys_info:{},
            labelType: '2',
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
                },
                request: false, // 请求中标志
            },
            // 队列
            queue:{
                exec:[], // 执行队列
                register:[], // 注册队列
            },
            add_box: {
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
        }
    },
    // 模块初始化
    created(){
        document.title = '流水线管理';
        this.getDic()
        api.systemInfo((res)=>{
            this.sys_info = res;
        })
    },
    // 模块初始化
    mounted(){
        this.getList()
        // 添加指定事件监听
        this.env = cache.getEnv()
        this.$sse.addEventListener(this.env.env+".exec.queue", this.execQueue)
        this.$sse.addEventListener(this.env.env+'.register.queue', this.registerQueue)
    },
    beforeDestroy(){
        // 销毁指定事件监听
        this.$sse.removeEventListener(this.env.env+".exec.queue", this.execQueue)
        this.$sse.removeEventListener(this.env.env+".register.queue", this.registerQueue)
    },
    // 具体方法
    methods:{
        // 任务列表
        getList(){
            if (this.list.request){
                return this.$message.info('请求执行中,请稍等.');
            }
            this.list.request = true
            api.innerGet("/pipeline/list", this.listParam, (res)=>{
                this.list.request = false
                if (!res.status){
                    console.log("pipeline/list 错误", res)
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
            api.dicList([Enum.dicUser, Enum.dicMsg],(res) =>{
                this.dic.user = res[Enum.dicUser]
                this.dic.msg = res[Enum.dicMsg]
            })
        },
        // 添加弹窗
        setShow(row=null){
            this.add_box.show = true
            if (row == null){
                this.add_box.title= '添加流水线'
                this.add_box.detail = {}
            }else{
                this.add_box.title= '编辑流水线'
                this.add_box.detail = row
            }
        },
        formClose(e){
            console.log("close",e)
            this.add_box.show = false
            if (e.is_change){
                this.getList()
            }
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
            console.log("execQueue", e)
            let data = JSON.parse(e.data)
            if (!data){
                this.queue.exec = []
                return
            }
            this.queue.exec = data
        },
        registerQueue(e){
            console.log("registerQueue", e)
            let data = JSON.parse(e.data)
            if (!data){
                this.queue.register = []
                return
            }
            this.queue.register = data
        },
    }
})

Vue.component("MyPipeline", MyPipeline);