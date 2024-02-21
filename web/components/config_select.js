var MyConfigSelect = Vue.extend({
    template: `<div class="config-select">
    <el-menu :default-active="labelType" class="el-menu-demo" mode="horizontal" @select="handleClickTypeLabel">
        <el-menu-item index="1" :disabled="listRequest">周期任务</el-menu-item>
        <el-menu-item index="2" :disabled="listRequest">单次任务</el-menu-item>
        <el-menu-item index="5" :disabled="listRequest">组件任务</el-menu-item>
    </el-menu>
    
    <el-table :data="list">
        <el-table-column prop="id" label="编号"></el-table-column>
        <el-table-column prop="" label="执行成功率">
            <template slot-scope="scope">
                <el-tooltip placement="top-start">
                    <div slot="content">{{scope.row.topRatio}}%<br/>近期{{scope.row.top_number}}次执行，{{scope.row.top_error_number}}次失败。</div>
                    <i :class="getTopIcon(scope.row.top_number, scope.row.topRatio)"></i>
                </el-tooltip>
            </template>
        </el-table-column>
        <el-table-column prop="spec" label="执行时间"></el-table-column>
        <el-table-column prop="name" label="任务名称"></el-table-column>
        <el-table-column prop="protocol_name" label="协议"></el-table-column>
        <el-table-column prop="remark" label="备注"></el-table-column>
        <el-table-column prop="" label="状态">
            <template slot-scope="scope">
                <el-tooltip placement="top-start">
                    <div slot="content">{{scope.row.status_dt}}  {{scope.row.status_remark}}</div>
                    <span :class="statusClass(scope.row.status)">{{scope.row.status_name}}</span>
                </el-tooltip>
            </template>
        </el-table-column>
        <el-table-column label="操作">
            <template slot-scope="{row}">
                <el-button type="text" @click="editShow(row)">编辑</el-button>
                <el-button type="text" @click="changeStatus(row, 2)" v-if="row.status!=2">激活</el-button>
                <el-button type="text" @click="changeStatus(row, 1)" v-if="row.status==2">停用</el-button>
                <el-button type="text" @click="configLogBox(row.id, row.name)">日志</el-button>
            </template>
        </el-table-column>
    </el-table>
    <el-pagination
            @size-change="handleSizeChange"
            @current-change="handleCurrentChange"
            :current-page.sync="listPage.page"
            :page-size="listPage.size"
            layout="total, prev, pager, next"
            :total="listPage.total">
    </el-pagination>
</div>`,

    name: "MyConfigSelect",
    data(){
        return {
            list:{
                page: {},
                param:{
                    page: 1,
                    size: 20,
                }
            },
        }
    },
    // 模块初始化
    created(){
        this.initForm(false,"-")
    },
    // 模块初始化
    mounted(){
        console.log("sql_source mounted")
        this.getList()
    },

    // 具体方法
    methods:{
        // 任务列表
        getList(){
            if (this.listRequest){
                return this.$message.info('请求执行中,请稍等.');
            }
            this.listRequest = true
            api.innerGet("/config/list", this.listParam, (res)=>{
                this.listRequest = false
                if (!res.status){
                    console.log("config/list 错误", res)
                    return this.$message.error(res.message);
                }
                for (i in res.data.list){
                    let ratio = 0
                    if (res.data.list[i].top_number){
                        ratio = res.data.list[i].top_error_number / res.data.list[i].top_number
                    }
                    if (res.data.list[i].command.sql){
                        res.data.list[i].command.sql.err_action = res.data.list[i].command.sql.err_action.toString()
                    }
                    res.data.list[i].status = res.data.list[i].status.toString()
                    res.data.list[i].topRatio = 100 - ratio * 100
                }
                this.list = res.data.list;
                this.listPage = res.data.page;
            })
        },
        handleSizeChange(val) {
            console.log(`每页 ${val} 条`);
        },
        handleCurrentChange(val) {
            this.listParam.page = val
            this.getList()
        },
        handleClickTypeLabel(tab, event) {
            this.listParam.type = tab
            this.getList()
        },
        close(){
            this.$emit('update:visible', false) // 向外传递关闭表示
        }
    }
})

Vue.component("MyConfigSelect", MyConfigSelect);