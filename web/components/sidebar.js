var MySidebar = Vue.extend({
    template: `<el-aside style="padding-top: 10px">
        <el-card class="aside-card">
            <div slot="header">
                <span class="h3">执行任务</span>
            </div>
            <ol style="max-height: 320px;">
                <li v-for="item in queue.exec">
                    <span v-html="taskItemIcon(item)"></span>
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
            <ol style="max-height: 700px;">
                <li v-for="item in queue.register">
                    <span v-html="taskItemIcon(item)"></span>
                    <router-link :to="{path:'/config_detail', query:{id:item.ref_id, type:item.ref_type, entry_id:item.entry_id}}" class="el-link el-link--default is-underline">{{item.name}}</router-link>
                </li>
            </ol>
            <div v-show="!queue.register.length">-</div>
        </el-card>
    </el-aside>`,

    name: "MySidebar",
    data(){
        return {
            env: {},
            queue:{
                exec:[], // 执行队列
                register:[], // 注册队列
            },
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
        taskItemIcon(row){
            if (row.ref_type == 'config'){
                return '<i class="task-item-icon" style="background: #28ab80;">c</i>'
            }else if (row.ref_type == 'pipeline'){
                return '<i class="task-item-icon" style="background: #5c88c5;">p</i>'
            }else if(row.ref_type == 'receive'){
                return '<i class="task-item-icon" style="background: #182b50;">r</i>'
            }
        },
    }
})

Vue.component("MySidebar", MySidebar);