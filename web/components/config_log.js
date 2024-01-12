var MyConfigLog = Vue.extend({
    template: `<div>
                        <el-table :data="list">
                            <el-table-column property="create_dt" label="记录时间" width="200"></el-table-column>
                            <el-table-column property="status_name" label="状态" width="80">
                                <template slot-scope="scope">
                                    <el-tag :type="scope.row.status == 1 ? 'danger' : 'success'">{{scope.row.status_name}}</el-tag>
                                </template>
                            </el-table-column>
                            <el-table-column property="duration" label="耗时/秒" width="80"></el-table-column>
                            <el-table-column property="" label="详情">
                                <template slot-scope="scope">
<!--                                    <el-popover trigger="hover" placement="left">-->
<!--                                        <div>-->
                                        {{scope.row.body.substring(0, 105)+'...'}}
<!--                                        </div>-->
                                        <el-button type="text" slot="reference" @click="logBodyShow(scope.row)">更多</el-button>
<!--                                    </el-popover>-->
                                </template>
                            </el-table-column>
                        </el-table>
                        
                      
                        <!-- 日志详情弹窗 -->
                        <el-dialog title="日志详情" :visible.sync="body.show" size="30%" style="white-space: pre;" append-to-body="true">
                            <pre><code>{{body.data.replace(/\\\\n/g, '\\n').replace(/\\\\t/g, '\\t')}}</code></pre>
                        </el-dialog>
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
                data:"",
                show:false
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
                if (newVal > 0){
                    this.logByConfig(newVal)
                }
            },
        }
    },

    // 具体方法
    methods:{
        // 配置日志
        logByConfig(id){
            api.innerGet("/log/by_config", {conf_id:id, limit:15}, (res)=>{
                if (!res.status){
                    console.log("log/by_config 错误", res)
                    return this.$message.error(res.message);
                }
                this.list = res.data.list;
            })
        },
        // 日志描述展示
        logBodyShow(logItem){
            this.body.data = logItem.body
            this.body.show = true
            // this.$alert('<div style="white-space: pre;">'+logItem.body+'<div>', '日志详情',{
            //     dangerouslyUseHTMLString: true
            // })
        },
    }
})

Vue.component("MyConfigLog", MyConfigLog);