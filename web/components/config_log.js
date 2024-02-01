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
                <el-button type="text" @click="getTrace(scope.row.trace_id)">查看</el-button>
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
            traces:[], // 踪迹
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
        // 踪迹展示
        getTrace(traceId){
            api.innerGet("/log/traces", {trace_id: traceId}, (res)=>{
                if (!res.status){
                    console.log("log/traces 错误", res)
                    return this.$message.error(res.message);
                }
                let trace = [];
                res.data.list[0].Spans.forEach(function (item, index, raw) {
                    if (item == 0){
                        if (item){
                            
                        }
                    }
                })




                // 剩下就是我要怎么来组织这个树了！

                // this.traces = res.data.list;
            })
        }

    }
})

Vue.component("MyConfigLog", MyConfigLog);