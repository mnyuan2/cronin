var MyConfigSelect = Vue.extend({
    template: `<div class="config-select">
    <el-radio-group v-model="labelType"  @change="changeTab" size="medium">
        <el-radio-button label="1" :disabled="list.request">周期任务</el-radio-button>
        <el-radio-button label="2" :disabled="list.request">单次任务</el-radio-button>
        <el-radio-button label="5" :disabled="list.request">组件任务</el-radio-button>
    </el-radio-group>
    
    <el-table :data="list.items" @selection-change="selectedChange">
    <el-table-column type="selection" width="55">
    </el-table-column>
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
    </el-table>
    <el-pagination
            @size-change="handleSizeChange"
            @current-change="handleCurrentChange"
            :current-page.sync="list.page.page"
            :page-size="list.page.size"
            layout="total, prev, pager, next"
            :total="list.page.total">
    </el-pagination>
</div>`,

    name: "MyConfigSelect",
    data(){
        return {
            labelType: '1',
            list:{
                page: {},
                param:{
                    page: 1,
                    size: 10,
                },
                request: false,
            },
            selected: [],
        }
    },
    // 模块初始化
    created(){},
    // 模块初始化
    mounted(){
        console.log("sql_source mounted")
        this.getList()
    },

    // 具体方法
    methods:{
        // 任务列表
        getList(){
            if (this.list.request){
                return this.$message.info('请求执行中,请稍等.');
            }
            this.list.request = true
            api.innerGet("/config/list", this.list.param, (res)=>{
                this.list.request = false
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
                this.list.items = res.data.list;
                this.list.page = res.data.page;
            })
        },
        selectedChange(val) {
            this.selected = val;
        },
        handleSizeChange(val) {
            console.log(`每页 ${val} 条`);
        },
        handleCurrentChange(val) {
            this.list.param.page = val
            this.getList()
        },
        changeTab(tab) {
            console.log("changeTab", tab)
            this.list.param.type = tab
            this.selected = []
            this.getList()
        },
        statusClass(status){
            switch (Number(status)) {
                case 1:
                    return 'warning';
                case 2:
                    return 'primary';
                case 3:
                    return 'success';
                case 4:
                    return 'danger';
            }
        },
        close(){
            this.$emit("handleconfirm", this.doc_item.id, this.selectedData)
            this.$emit('update:visible', false) // 向外传递关闭表示
        }
    }
})

Vue.component("MyConfigSelect", MyConfigSelect);