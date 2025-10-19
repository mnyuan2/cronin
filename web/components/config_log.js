var MyConfigLog = Vue.extend({
    template: `<div class="config-log">
<!-- 第一版方案
    <el-table :data="list" :empty-text="list_init? '暂无数据' : '请点击搜索查询数据'">
        <el-table-column property="timestamp" label="开始时间" width="160"></el-table-column>
        <el-table-column property="operation" label="操作" width="200"></el-table-column>
        <el-table-column property="ref_name" label="任务" width="300"></el-table-column>
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
    -->
    <!-- 第二版方案 -->
    <ul class="list">
        <li v-for="item in list">
            <el-card class="box-card" shadow="hover">
                <div slot="header" @click="traceBox(item.trace_id)">
                    <span :class="item.status == 1 ? 'danger el-icon-warning': 'el-icon-empty'"></span>
                    <span class="h4" style="max-width: 73%;overflow: clip;white-space: nowrap;">{{item.ref_name}}</span>
                    <span style="margin: 0 5px;font-weight: 600;" class="info-2">{{item.operation}}</span>
                    <span v-if="showLink && item.ref_id" class="panel">
                        <i class="el-icon-link hover" @click.stop="jumpDetail(item)" title="任务详情"></i>
                    </span>
                    <span style="float: right">
                        {{durationTransform(item.duration)}}
                    </span>
                </div>
                <el-row>
                    <el-col :span="3">
                        <el-tag size="small" effect="plain" type="info">total {{item.span_total}}</el-tag>
                    </el-col>
                    <el-col :span="16">
                        <el-tag v-for="group in item.span_group" size="small" effect="plain" type="info">{{group.key}}({{group.value}})</el-tag>
                    </el-col>
                    <el-col :span="5" style="text-align: right;white-space: nowrap;">
                        {{item.timestamp}}
                    </el-col>
                </el-row>
            </el-card>
        </li>
    </ul>
    <el-pagination background v-show="page.total"
        @size-change="handleSizeChange"
        @current-change="handleCurrentChange"
        :current-page="page.page"
        :page-sizes="[15, 100, 500, 1000, 2000]"
        :page-size="page.size"
        layout="total, sizes, prev, pager, next"
        :total="page.total">
    </el-pagination>
    
    <!-- 踪迹弹窗 -->
    <el-drawer title="日志踪迹" :visible.sync="trace.show" direction="rtl" size="70%" wrapperClosable="false" :before-close="traceBox" append-to-body>
        <my-trace :trace_id="trace.id"></my-trace>
    </el-drawer>
</div>`,
    name: "MyConfigLog",
    props: {
        // tags:Object,
        search: Object,
        showLink: Boolean,
    },
    data(){
        return {
            // tags:{},
            list_init: false,
            list_empty_txt: "",
            search:{},
            page:{
                page: 1,
                size: 15,
                total: 0,
            },
            list:[],
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
        // tags:{
        //     immediate: true, // 解决首次负值不触发的情况
        //     handler: function (newVal,oldVal){
        //         if (Object.keys(newVal).length){
        //             this.logByConfig({tags: JSON.stringify(newVal)})
        //         }
        //     },
        // },
        search:{
            immediate: true, // 解决首次负值不触发的情况
            handler: function (newVal,oldVal){
                if (newVal && Object.keys(newVal).length){
                    this.page.page = 1
                    this.logByConfig(newVal)
                }
            },
        }
    },

    // 具体方法
    methods:{
        // 配置日志
        logByConfig(body){
            body["limit"] = this.page.size
            body["page"] = this.page.page
            api.innerGet("/log/list", body, (res)=>{
                if (!res.status){
                    console.log("log/list 错误", res)
                    return this.$message.error(res.message);
                }
                this.list = res.data.list;
                this.page = res.data.page
                this.list_init = true;
            })
        },
        handleSizeChange(val) {
            this.page.size = val
        },
        handleCurrentChange(val) {
            this.page.page = val
            this.logByConfig(this.search)
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
        jumpDetail(row){
            if (row.ref_id){
                let t = 'config'
                if (row.operation == 'job-pipeline'){
                    t = 'pipeline'
                }else if (row.operation == 'job-receive'){
                    t = 'receive'
                }
                window.open(`/index?env=prod#/config_detail?id=${row.ref_id}&type=${t}`)
            }
        },
    }
})

Vue.component("MyConfigLog", MyConfigLog);