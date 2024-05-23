var MyWork = Vue.extend({
    template: `<el-main class="my-work">
<!--
        先规划一下页面布局（整体仿造tapd，还有我们自己的手机端），
        # 各个环境名称（如果没有对应任务就不在哪上）（默认只有第一个标签页打开）
            流水线和任务混合在一起（sql写法要研究一下）
                点击跳转到详情页（新的标签页）
-->
        <el-card shadow="never" body-style="padding:12px">
            <el-button size="mini" round plain @click="statusOption(5)" :type="options.status==5?'primary': ''">待办</el-button>
            <el-button size="mini" round plain @click="statusOption(0)" :type="options.status==0?'primary': ''">所有</el-button>
        </el-card>
        <el-row v-for="(group_v,group_k) in group_list">
            <el-row class="header">
                <span class="my-table-icon" @click="showTable(group_v)">
                    <i :class="group_v.show ? 'el-icon-caret-bottom' : 'el-icon-caret-right'"></i>
                </span>
                
                {{group_v.env_title}} <el-divider direction="vertical"></el-divider> 
                {{group_v.join_type=='config'? '任务' : '流水线'}} <span style="color: #72767b">（{{group_v.total}}）</span>
            </el-row>
            <el-row class="body" v-if="group_v.show">
                <el-table :data="group_v.list" v-loading="group_v.loading">
                    <el-table-column prop="spec" label="执行时间"></el-table-column>
                    <el-table-column prop="name" label="任务名称"></el-table-column>
                    <el-table-column prop="protocol_name" label="协议"></el-table-column>
                    <el-table-column prop="remark" label="备注"></el-table-column>
                    <el-table-column prop="status_user_name" label="操作人"></el-table-column>
                    <el-table-column prop="" label="状态">
                        <template slot-scope="scope">
                            <el-tooltip placement="top-start">
                                <div slot="content">{{scope.row.status_dt}}  {{scope.row.status_remark}}</div>
                                <span >{{scope.row.status_name}}</span>
                            </el-tooltip>
                        </template>
                    </el-table-column>
                </el-table>
            </el-row>
        </el-row>
    </el-main>`,

    name: "MyWork",
    props: {
        data_id:Number
    },
    data(){
        return {
            group_list: [],
            options:{
                status: 5, // 默认待办
            }
        }
    },
    // 模块初始化
    created(){},
    // 模块初始化
    mounted(){
        this.getTables()
    },
    // 具体方法
    methods:{
        getTables(){
            api.innerGet('/work/table',{status:this.options.status}, res =>{
                if (!res.status){
                    return this.$message.error(res.message)
                }
                let list = []
                res.data.list.forEach((item,index)=>{
                    item.show = false
                    item.list = null
                    item.listpage = {}
                    item.loading = false
                    list.push(item)
                })
                this.group_list = list
            })
        },
        // 显示工作表
        showTable(row){
            row.show = !row.show
            let path = ""
            if (row.join_type == 'config'){
                path = "/config/list"
            }else if (row.join_type == 'pipeline'){
                path = "/pipeline/list"
            }else{
                return this.$message.warning('未支持的表格类型')
            }
            let param =  {page:1,size:10,env:row.env, status: this.options.status}

            if (row.show && row.list == null){
                row.loading = true
                row.list = []
                api.innerGet(path, param, (res)=>{
                    row.loading = false
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
                    }
                    row.list = res.data.list;
                    row.listPage = res.data.page;
                })
            }
        },
        // 状态切换
        statusOption(status){
            if (this.options.status == status){
                return
            }
            this.options.status = status
            this.getTables()
        }
    }
})

Vue.component("MyWork", MyWork);