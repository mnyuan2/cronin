var MyLogs = Vue.extend({
    template: `<el-container>
    <!--边栏-->
    <el-aside style="padding-top: 20px" class="my-aside-2">
        <el-form :model="form" label-position="top" size="small">
            <el-form-item label="操作">
                <el-select v-model="form.operation" placeholder="结果状态" style="width:100%">
                    <el-option label="job-task" value="job-task"><span>job-task</span><span style="float: right; color: #8492a6;">任务</span></el-option>
                    <el-option label="job-pipeline" value="job-pipeline"><span>job-pipeline</span><span style="float: right; color: #8492a6;">流水线</span></el-option>
                    <el-option label="job-receive" value="job-receive"><span>job-receive</span><span style="float: right; color: #8492a6;">接收</span></el-option>
                </el-select>
            </el-form-item>
            <el-form-item label="任务">
                <el-select v-model="form.ref_id" placeholder="结果状态" style="width:100%" clearable filterable>
                    <el-option v-for="item in dic.log_name" :label="item.name" :value="item.id" v-if="item.extend.operation == form.operation"></el-option>
                </el-select>
            </el-form-item>
            
            <el-form-item label="状态">
                <el-select v-model="form.status" placeholder="结果状态" style="width:100%" clearable>
                    <el-option label="成功" value="2"></el-option>
                    <el-option label="失败" value="1"></el-option>
                </el-select>
            </el-form-item>
            <el-form-item label="开始时间">
                <el-date-picker style="width:100%" clearable="false"
                        format="yyyy-MM-dd HH:mm:ss"
                        value-format="yyyy-MM-dd HH:mm:ss"
                        v-model="form.timestamp_start"
                        type="datetime"
                        placeholder="开始时间">
                </el-date-picker>
            </el-form-item>
            <el-form-item label="截止时间">
                <el-date-picker style="width:100%" clearable="false"
                        format="yyyy-MM-dd HH:mm:ss"
                        value-format="yyyy-MM-dd HH:mm:ss"
                        v-model="form.timestamp_end"
                        type="datetime"
                        placeholder="截止时间">
                </el-date-picker>
            </el-form-item>
            <el-form-item label="">
                <el-row type="flex" justify="space-between">
                    <div class="left" style="margin-right: 5px">
                        <label class="el-form-item__label">最小耗时</label>
                        <el-input type="number" placeholder=">= s.秒" v-model="form.duration_start"></el-input>
                    </div>
                    <div class="right"  style="margin-left: 5px">
                        <label class="el-form-item__label">最大耗时</label>
                        <el-input type="number" placeholder="<= s.秒" v-model="form.duration_end"></el-input>
                    </div>
                </el-row>
            </el-form-item>
            <el-form-item style="margin-top: 30px;">
                <el-row type="flex" justify="space-between">
                    <el-button type="success" @click="logSearch()" size="small">搜索</el-button>
                </el-row>
            </el-form-item>
        </el-form>
        
    </el-aside>
    
    <!--主内容-->
    <el-main>
        <my-config-log :search="search" show-link></my-config-log>
    </el-main>
</el-container>
`,
    name: "MyLogs",
    data(){
        return {
            dic:{},
            form:{
                operation:'',
                timestamp_start:'',
                timestamp_end:'',
                duration_start:'',
                duration_end:'',
                status:'',
            },
            search:{},
        }
    },
    // 模块初始化
    created(){
        setDocumentTitle('日志')
        this.getDic()
        // api.systemInfo((res)=>{
        //     this.sys_info = res;
        // })
        this.loadParams(getHashParams(window.location.hash))
    },
    // 模块初始化
    mounted(){},
    beforeDestroy(){},
    // 具体方法
    methods:{
        loadParams(param){
            // if (typeof param !== 'object'){return}
            // if (param.type){this.listParam.type = param.type.toString()}
            // if (param.page){this.listParam.page = Number(param.page)}
            // if (param.size){this.listParam.size = Number(param.size)}
            // if (param.name){this.listParam.name = param.name.toString()}
            // if (param.protocol){this.listParam.protocol = param.protocol.map(Number)}
            // if (param.status){this.listParam.status = param.status.map(Number)}
            // if (param.handle_user_ids){this.listParam.handle_user_ids = param.handle_user_ids.map(Number)}
            // if (param.create_user_ids){this.listParam.create_user_ids = param.create_user_ids.map(Number)}

            this.form.operation = 'job-task'
            const end = new Date();
            const start = new Date();
            start.setDate(start.getDate() - 7);
            this.form.timestamp_start = getDateString(start)+' 00:00:00'
            this.form.timestamp_end = getDateString(end) + ' 23:59:59'
        },
        // 日志搜索
        logSearch(){
            console.log("日志搜索")
            let form = copyJSON(this.form)
            let search = copyJSON(this.search)
            Object.assign(search, form)
            this.search = search
        },
        // 枚举
        getDic(){
            api.dicList([Enum.dicLogName],(res) =>{
                this.dic.log_name = res[Enum.dicLogName]
            })
        },
    }
})

Vue.component("MyLogs", MyLogs);