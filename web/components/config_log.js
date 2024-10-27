var MyConfigLog = Vue.extend({
    template: `<div class="config-log">
    <el-table :data="list">
        <el-table-column label="开始时间" width="160">
            <template slot-scope="scope">
                {{getDatetimeString(new Date(scope.row.timestamp/1000))}}
            </template>
        </el-table-column>
        <el-table-column property="operation" label="操作" width="200"></el-table-column>
        <el-table-column property="status_name" label="状态" width="70">
            <template slot-scope="scope">
                <el-tooltip placement="top-start">
                    <div slot="content">{{scope.row.status_desc}}</div>
                    <span :class="scope.row.status == 1 ? 'danger' : 'success'">{{scope.row.status_name}}</span>
                </el-tooltip>
            </template>
        </el-table-column>
        <el-table-column label="耗时" width="100">
            <template slot-scope="scope">
                {{durationTransform(scope.row.duration)}}
            </template>
        </el-table-column>
        <el-table-column property="" label="详情">
            <template slot-scope="scope">
                <el-button type="text" @click="traceBox(scope.row.trace_id)">查看</el-button>
            </template>
        </el-table-column>
    </el-table>
    
    <!-- 踪迹弹窗 -->
    <el-drawer title="日志踪迹" :visible.sync="trace.show" direction="rtl" size="70%" wrapperClosable="false" :before-close="traceBox" append-to-body>
        <my-trace :trace_id="trace.id"></my-trace>
    </el-drawer>
</div>`,
    name: "MyConfigLog",
    props: {
        tags:Object
    },
    data(){
        return {
            tags:{},
            list:[],// 日志列表，没有分页；
            trace:{
                id: "",
                show: false
            }
        }
    },
    // 模块初始化
    created(){},
    // 模块初始化
    mounted(){},
    watch:{
        tags:{
            immediate: true, // 解决首次负值不触发的情况
            handler: function (newVal,oldVal){
                if (Object.keys(newVal).length){
                    this.logByConfig(newVal)
                }
            },
        }
    },

    // 具体方法
    methods:{
        // 配置日志
        logByConfig(tags){
            let body = {tags: JSON.stringify(tags), limit:15}
            api.innerGet("/log/list", body, (res)=>{
                if (!res.status){
                    console.log("log/list 错误", res)
                    return this.$message.error(res.message);
                }
                this.list = res.data.list;
            })
        },
        // 踪迹盒子
        traceBox(id){
            if (typeof id === "string" && id != ""){
                this.trace.show = true;
                this.trace.id = id;
            }else{
                this.trace.show = false;
                this.trace.id = "";
            }
        },

    }
})

Vue.component("MyConfigLog", MyConfigLog);