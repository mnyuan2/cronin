var MySidebar = Vue.extend({
    template: `<el-aside style="padding-top: 20px">
        <el-card class="aside-card">
            <div slot="header">
                <span class="h3">执行任务</span>
            </div>
            <ol style="max-height: 320px;">
                <li v-for="item in queue.exec">
                    <span v-html="getTaskIcon(item.ref_type)"></span>
                    <router-link :to="{path:'/config_detail', query:{id:item.ref_id, type:item.ref_type, entry_id:item.entry_id}}" class="el-link el-link--default is-underline">{{item.name}}</router-link>
                    <p style="margin: 0;color: #909399;line-height: 100%;font-size: 12px;">
                        <i class="el-icon-loading" title="查看日志" style="cursor: pointer;margin-right: 2px" @click="showLog(item)"></i>
                        ({{durationTransform(item.duration, 's')}}) 
                        <el-popconfirm :title="'确定停止执行 '+item.name+' 任务吗？'" @confirm="jobStop(item)"><i slot="reference" class="el-icon-circle-close stop"></i></el-popconfirm>
                    </p>
                </li>
            </ul>
            <div v-show="!queue.exec.length">-</div>
        </el-card>
        <el-card class="aside-card">
            <div slot="header">
                <span class="h3">注册任务</span>
            </div>
            <ol style="max-height: 700px;">
                <li v-for="item in queue.register">
                    <span v-html="getTaskIcon(item.ref_type)"></span>
                    <router-link :to="{path:'/config_detail', query:{id:item.ref_id, type:item.ref_type, entry_id:item.entry_id}}" class="el-link el-link--default is-underline">{{item.name}}</router-link>
                </li>
            </ol>
            <div v-show="!queue.register.length">-</div>
        </el-card>
        <!-- 踪迹弹窗 -->
        <el-drawer title="日志踪迹" :visible.sync="trace.show" direction="rtl" size="70%" wrapperClosable="false" :before-close="closeLog" append-to-body>
            <my-trace :job="trace.job" v-if="trace.show"></my-trace>
        </el-drawer>
    </el-aside>`,

    name: "MySidebar",
    data(){
        return {
            env: {},
            queue:{
                exec:[], // 执行队列
                register:[], // 注册队列
            },
            trace:{
                show: false,
                job: {},
            }
        }
    },
    // 模块初始化
    created(){
    },
    // 模块初始化
    mounted(){
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
        // 停止执行任务
        jobStop(row){
            api.innerPost("/job/stop",row, res=>{
                if (!res.status){
                    return this.$message.error(res.message)
                }
                this.$message.success('操作成功')
            })
        },
        // 查看日志
        showLog(row){
            console.log("showLog", row)
            this.trace.show = true
            this.trace.job = {
                ref_id: row.ref_id,
                entry_id: row.entry_id,
                trace_id: row.trace_id,
            }
        },
        closeLog(e){
            console.log("closeLog", e)
            this.trace.show = false
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

Vue.component("MySidebar", MySidebar);