var MyMessageTemplate = Vue.extend({
    template: `<el-main>
        <el-button type="text" @click="initForm(true)">新增模板</el-button>
        <el-table :data="list">
            <el-table-column property="title" label="模板名称"></el-table-column>
            <el-table-column property="sort" label="排序"></el-table-column>
            <el-table-column property="update_dt" label="更新时间"></el-table-column>
            <el-table-column label="操作">
                <template slot-scope="scope">
                    <el-button plain @click="initForm(true, scope.row)">编辑</el-button>
                    <el-button plain @click="deleteSqlSource(scope.row.id)">删除</el-button>
                </template>
            </el-table-column>
        </el-table>

        <!--设置弹窗-->
        <el-dialog :title="form.box.title" :visible.sync="form.box.show" :close-on-click-modal="false" append-to-body="true">
            <el-form :model="form.data" label-position="left" label-width="100px">
                <el-form-item label="模板名称">
                    <el-input v-model="form.data.title"></el-input>
                </el-form-item>
                <el-form-item label="请求地址" class="http_url_box">
                    <el-input v-model="form.data.http.url" placeholder="请输入http:// 或 https:// 开头的完整地址">
                        <el-select v-model="form.data.http.method" placeholder="请选请求方式" slot="prepend">
                            <el-option label="GET" value="GET"></el-option>
                            <el-option label="POST" value="POST"></el-option>
                        </el-select>
                    </el-input>
                </el-form-item>
                <el-form-item label="请求Header" class="http_header_box">
                    <el-input v-for="(header_v,header_i) in form.data.http.header" v-model="header_v.value" placeholder="参数值">
                        <el-input v-model="header_v.key" slot="prepend" placeholder="参数名" @input="httpHeaderInput"></el-input>
                        <el-button slot="append" icon="el-icon-delete" @click="httpHeaderDel(header_i)"></el-button>
                    </el-input>
                </el-form-item>
                <el-form-item label="请求Body">
                    <el-input type="textarea" v-model="form.data.http.body" rows="5" placeholder="POST请求时body参数，将通过json进行请求发起"></el-input>
                </el-form-item>
                <el-form-item label="模板变量">
                  <el-alert title="请求地址、请求Header、请求Body 中使用[[var_name]] 双中括号包含变量名称，消息推送时会被实际值替换。" type="info" :closable="false"></el-alert>
                    <el-table :data="varDesc">
                        <el-table-column label="name" property="name" width="100"></el-table-column>
                        <el-table-column label="说明" property="description"></el-table-column>
                    </el-table>
                </el-form-item>
            </el-form>
            <div slot="footer" class="dialog-footer">
                <el-button @click="runForm()" style="float: left;">发送测试</el-button>
                <el-button @click="initForm(false,'-')">取 消</el-button>
                <el-button type="primary" @click="submitForm()">确 定</el-button>
            </div>
        </el-dialog>
    </el-main>`,

    name: "MyMessageTemplate",
    data(){
        return {
            list:[],
            page:{},
            listParam:{
                page: 1,
                size: 20,
            },
            form:{}, // 表单
            varDesc: [
                {name: 'env', description:'string 环境名称'},
                {name: 'config', description:'object 任务配置信息，包含字段：name·任务名称、protocol_name·协议类型'},
                {name: 'log', description:'object 结果日志信息，包含字段：status_name·状态(成功、失败)、status_desc·状态描述、body·结果日志、duration·耗时/秒、create_dt·触发时间(yyyy-mm-dd HH:ii:ss)'},
                {name: 'user', description:'object 人员信息，包含字段 username·用户名、mobie·手机号 信息'},
            ],

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
        // 列表
        getList(){
            api.innerGet("/setting/message_list", this.listParam, (res)=>{
                if (!res.status){
                    return this.$message.error(res.message);
                }
                this.list = res.data.list;
                this.page = res.data.page;
            })
        },
        handleSizeChange(val) {
            console.log(`每页 ${val} 条`);
        },
        handleCurrentChange(val) {
            console.log(`当前页: ${val}`);
            this.listParam.page = val
            this.getList()
        },
        submitForm(){
            let body = this.form.data
            if (body.title == ""){
                return this.$message.warning("请输入模板名称");
            }
            if (body.http.url == ""){
                return this.$message.warning("请输入推送请求地址");
            }
            api.innerPost("/setting/message_set", body, (res) =>{
                console.log("sql源设置响应",res)
                if (!res.status){
                    return this.$message.error(res.message)
                }
                this.initForm(false)
                this.getList()
            })
        },
        // 执行一下
        runForm(){
            let body = this.form.data
            if (body.title == ""){
                return this.$message.warning("请输入模板名称");
            }
            if (body.http.url == ""){
                return this.$message.warning("请输入推送请求地址");
            }
            api.innerPost("/setting/message_run", body, (res) =>{
                if (!res.status){
                    return this.$message({
                        message: res.message,
                        type: 'error',
                        duration: 6000
                    })
                }
                return this.$message.success("ok."+res.data.result)
            })
        },
        // 初始化表单数据
        initForm(show, data){
            this.form = {
                box:{
                    show: show == true,
                    title: "设置消息模板",
                },
                data: {
                    id: 0,
                    title:"",
                    sort: 1,
                    http:{
                        method: 'POST',
                        header: [{"key":"","value":""}],
                        url:'',
                        body:'',
                    },
                }
            }
            if ( typeof data === 'object' && data["id"] != undefined &&  data.id > 0 ){
                this.form.box.title = '编辑sql连接'
                this.form.data = data
                console.log("编辑源",data)
            }
        },
        // http header 输入值变化时
        httpHeaderInput(val){
            if (val == ""){
                return
            }
            let item = this.form.data.http.header.slice(-1)[0]
            if (item == undefined || item.key != ""){
                this.form.data.http.header.push({"key":"","value":""})
            }
        },
        // http header 输入值删除
        httpHeaderDel(index){
            if ((index+1) >= this.form.data.http.header.length){
                return
            }
            this.form.data.http.header.splice(index,1)
        },
        close(){
            this.$emit('update:visible', false) // 向外传递关闭表示
        }
    }
})

Vue.component("MyMessageTemplate", MyMessageTemplate);