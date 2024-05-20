var MyRole = Vue.extend({
    template: `<el-row class="my-work">
<!--
        先规划一下页面布局（整体仿造tapd，还有我们自己的手机端），
        # 各个环境名称（如果没有对应任务就不在哪上）（默认只有第一个标签页打开）
            流水线和任务混合在一起（sql写法要研究一下）
                点击跳转到详情页（新的标签页）
-->
        <el-col :span="3">
            <el-row>
                <h3>角色列表</h3>
                <button>添加角色</button>
            </el-row>
            <el-row>
                ...
            </el-row>
        </el-col>
        <el-col>
            <el-row>Header</el-row>
            <el-row>Main</el-row>
        </el-col>
        
        <el-dialog :title="'sql设置-'+sqlSet.title" :visible.sync="sqlSet.show" :show-close="false" :close-on-click-modal="false">
            <el-form>
                <el-form-item label="名称">
                    <el-input></el-input>
                </el-form-item>
                <el-form-item label="备注">
                    <el-input></el-input>
                </el-form-item>
            </el-form>
            <span slot="footer" class="dialog-footer">
                <el-button @click="dialogVisible = false">取 消</el-button>
                <el-button type="primary" @click="dialogVisible = false">确 定</el-button>
            </span>
        </el-dialog>
    </el-row>`,

    name: "MyRole",
    props: {
        data_id:Number
    },
    data(){
        return {
            group_list: [],
            form:{
                
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

Vue.component("MyRole", MyRole);