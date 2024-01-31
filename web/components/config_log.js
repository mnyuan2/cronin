var MyConfigLog = Vue.extend({
    template: `<div>
    <el-table :data="list">
        <el-table-column property="timestamp" label="时间" width="200"></el-table-column>
        <el-table-column property="operation" label="操作" width="200"></el-table-column>
        <el-table-column property="status_name" label="状态" width="80">
            <template slot-scope="scope">
                <el-tooltip placement="top-start">
                    <div slot="content">{{scope.row.status_desc}}</div>
                    <span :class="scope.row.status == 1 ? 'danger' : 'success'">{{scope.row.status_name}}</span>
                </el-tooltip>
            </template>
        </el-table-column>
        <el-table-column property="duration" label="耗时/秒" width="80"></el-table-column>
        <el-table-column property="" label="详情">
            <template slot-scope="scope">
                <el-button type="text" @click="logBodyShow(scope.row)">查看</el-button>
            </template>
        </el-table-column>
    </el-table>
    
               
</div>`,
    name: "MyConfigLog",
    props: {
        config_id:Number,
    },
    data(){
        return {
            config_id:0,
            list:[],// 日志列表，没有分页；
            logBody: "",
            body:{
                status: 0,
                status_desc:"",
                data:"",
                show:false,
                msg_data: []
            }
        }
    },
    // 模块初始化
    created(){},
    // 模块初始化
    mounted(){},
    watch:{
        config_id:{
            immediate: true, // 解决首次负值不触发的情况
            handler: function (newVal,oldVal){
                console.log("config_log config_id",newVal, oldVal)
                if (newVal != 0){
                    this.logByConfig(newVal)
                }
            },
        }
    },

    // 具体方法
    methods:{
        // 配置日志
        logByConfig(id){
            api.innerGet("/log/list", {tags: JSON.stringify({config_id: id}), limit:15}, (res)=>{
                if (!res.status){
                    console.log("log/list 错误", res)
                    return this.$message.error(res.message);
                }
                this.list = res.data.list;
            })
        },
        // 日志描述展示
        logBodyShow(logItem){
            this.body.data = logItem.body
            this.body.status = logItem.status
            this.body.status_desc = logItem.status_desc
            if (logItem.msg_body != null && logItem.msg_body.length > 0){
                this.body.msg_data = logItem.msg_body
            }else{
                this.body.msg_data = []
            }
            this.body.show = true
            // this.$alert('<div style="white-space: pre;">'+logItem.body+'<div>', '日志详情',{
            //     dangerouslyUseHTMLString: true
            // })
        },
    }
})

Vue.component("MyConfigLog", MyConfigLog);