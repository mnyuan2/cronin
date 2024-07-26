var MyConfigSelect = Vue.extend({
    template: `<div class="config-select">
    <el-radio-group v-model="labelType"  @change="changeTab" size="medium">
        <el-radio-button label="1" :disabled="list.request">周期任务</el-radio-button>
        <el-radio-button label="2" :disabled="list.request">单次任务</el-radio-button>
        <el-radio-button label="5" :disabled="list.request">组件任务</el-radio-button>
    </el-radio-group>
    
    <el-row>
        <el-form :inline="true" :model="list.param" size="mini" class="search-form">
            <el-form-item label="名称">
                <el-input v-model="list.param.name" placeholder="搜索名称"></el-input>
            </el-form-item>
            <el-form-item label="协议">
                <el-select v-model="list.param.protocol" placeholder="所有" multiple>
                    <el-option v-for="item in dic.protocol" :label="item.name" :value="item.id"></el-option>
                </el-select>
            </el-form-item>
            <el-form-item label="状态">
                <el-select v-model="list.param.status" placeholder="所有" multiple>
                    <el-option v-for="item in dic.config_status" :label="item.name" :value="item.id"></el-option>
                </el-select>
            </el-form-item>
<!--            <el-form-item label="处理人">-->
<!--                <el-select v-model="list.param.handle_user_ids" placeholder="所有" multiple>-->
<!--                    <el-option v-for="item in dic.user" :label="item.name" :value="item.id"></el-option>-->
<!--                </el-select>-->
<!--            </el-form-item>-->
<!--            <el-form-item label="创建人">-->
<!--                <el-select v-model="list.param.create_user_ids" placeholder="所有" multiple>-->
<!--                    <el-option v-for="item in dic.user" :label="item.name" :value="item.id"></el-option>-->
<!--                </el-select>-->
<!--            </el-form-item>-->
            <el-form-item>
                <el-button type="primary" @click="getList">查询</el-button>
            </el-form-item>
        </el-form>
    </el-row>
    
    <el-table :data="list.items" @selection-change="selectedChange">
        <el-table-column type="selection" width="55"></el-table-column>
        <el-table-column prop="name" label="任务名称">
            <div slot-scope="{row}" class="abc" style="display: flex;">
                <span style="white-space: nowrap;overflow: hidden;text-overflow: ellipsis;">
                    <router-link target="_blank" :to="{path:'/config_detail',query:{id:row.id, type:'config'}}" class="el-link el-link--primary is-underline" :title="row.name">{{row.name}}</router-link>
                </span>
            </div>
        </el-table-column>
        <el-table-column prop="protocol_name" label="协议" width="80"></el-table-column>
        <el-table-column prop="" label="状态" width="100">
            <template slot-scope="scope">
                <el-tooltip placement="top-start">
                    <div slot="content">{{scope.row.status_dt}}  {{scope.row.status_remark}}</div>
                    <el-button :type="statusTypeName(scope.row.status)" plain size="mini" round disabled>{{scope.row.status_name}}</el-button>
                </el-tooltip>
            </template>
        </el-table-column>
        <el-table-column prop="remark" label="备注"></el-table-column>
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
            dic:{
                user: [],
                config_status: [],
                protocol: [],
            },
            list:{
                page: {},
                param:{
                    type: '1',
                    page: 1,
                    size: 10,
                    protocol: [],
                    status: [],
                    handle_user_ids:[],
                    create_user_ids:[],
                    name: '',
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
        this.getDic()
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
                    // res.data.list[i].status = res.data.list[i].status.toString()
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
        // 枚举
        getDic(){
            api.dicList([Enum.dicUser, Enum.dicProtocolType, Enum.dicConfigStatus],(res) =>{
                this.dic.user = res[Enum.dicUser]
                this.dic.config_status = res[Enum.dicConfigStatus]
                this.dic.protocol = res[Enum.dicProtocolType]
            })
        },
        close(){
            this.$emit("handleconfirm", this.doc_item.id, this.selectedData)
            this.$emit('update:visible', false) // 向外传递关闭表示
        }
    }
})

Vue.component("MyConfigSelect", MyConfigSelect);