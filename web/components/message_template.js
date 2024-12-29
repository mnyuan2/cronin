var MyMessageTemplate = Vue.extend({
    template: `<el-main>
        <el-button type="text" @click="initForm(true)" v-if="$auth_tag.message_set">新增模板</el-button>
        <el-table :data="list">
            <el-table-column property="title" label="模板名称"></el-table-column>
            <el-table-column property="sort" label="排序"></el-table-column>
            <el-table-column property="update_dt" label="更新时间"></el-table-column>
            <el-table-column label="操作">
                <template slot-scope="scope">
                    <el-button plain size="small" @click="initForm(true, scope.row)">编辑</el-button>
                    <el-button plain size="small" @click="deleteSqlSource(scope.row.id)" v-if="$auth_tag.message_set">删除</el-button>
                </template>
            </el-table-column>
        </el-table>

        <!--设置弹窗-->
        <el-dialog :title="form.box.title" :visible.sync="form.box.show" :close-on-click-modal="false" append-to-body="true">
            <el-form :model="form.data" label-position="left" label-width="100px" size="small">
                <el-form-item label="模板名称">
                    <el-input v-model="form.data.title"></el-input>
                </el-form-item>
<!--                http消息-->
                <el-form-item class="http_url_box">
                    <span slot="label" style="white-space: nowrap;">
                        地址
                        <el-tooltip effect="dark" content="支持模板语法，点击查看更多" placement="top-start">
                            <router-link target="_blank" to="/var_params" style="color: #606266"><i class="el-icon-info"></i></router-link>
                        </el-tooltip>
                    </span>
                    <el-input class="input-input" v-model="form.data.template.http.url" placeholder="请输入http:// 或 https:// 开头的完整地址">
                        <el-select v-model="form.data.template.http.method" placeholder="请选请求方式" slot="prepend" style="width: 70px;">
                            <el-option label="GET" value="GET"></el-option>
                            <el-option label="POST" value="POST"></el-option>
                        </el-select>
                    </el-input>
                </el-form-item>
                <el-form-item class="http_header_box">
                    <span slot="label" style="white-space: nowrap;">
                        Header
                        <el-tooltip effect="dark" content="支持模板语法，点击查看更多" placement="top-start">
                            <router-link target="_blank" to="/var_params" style="color: #606266"><i class="el-icon-info"></i></router-link>
                        </el-tooltip>
                    </span>
                    <el-input class="input-input" v-for="(header_v,header_i) in form.data.template.http.header" v-model="header_v.value" placeholder="参数值">
                        <el-input v-model="header_v.key" slot="prepend" placeholder="参数名" @input="httpHeaderInput"></el-input>
                        <el-button slot="append" icon="el-icon-delete" @click="httpHeaderDel(header_i)"></el-button>
                    </el-input>
                </el-form-item>
                <el-form-item label="">
                    <span slot="label" style="white-space: nowrap;">
                        Body
                        <el-tooltip effect="dark" content="支持模板语法，点击查看更多" placement="top-start">
                            <router-link target="_blank" to="/var_params" style="color: #606266"><i class="el-icon-info"></i></router-link>
                        </el-tooltip>
                    </span>
                    <el-input type="textarea" v-model="form.data.template.http.body" :autosize="{minRows:5}" placeholder="POST请求时body参数，将通过json进行请求发起"></el-input>
                </el-form-item>
                <el-form-item label="模板变量">
                    <el-table :data="varDesc">
                        <el-table-column label="字段" property="name" width="100"></el-table-column>
                        <el-table-column label="说明" property="description"></el-table-column>
                    </el-table>
                </el-form-item>
            </el-form>
            <div slot="footer" class="dialog-footer">
                <el-button @click="runForm()" style="float: left;" v-if="$auth_tag.message_set" size="small">发送测试</el-button>
                <el-button @click="initForm(false,'-')" size="small">取 消</el-button>
                <el-button type="primary" @click="submitForm()" v-if="$auth_tag.message_set" size="small">确 定</el-button>
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
                {name: 'log', description:'object 结果日志信息，包含字段：status_name·状态(成功、失败)、status_desc·状态描述、body·响应、duration·耗时/秒、create_dt·触发时间'},
                {name: 'user', description:'object 人员信息，包含字段 username·用户名、mobile·手机号 信息'},
            ],

        }
    },
    // 模块初始化
    created(){
        setDocumentTitle('消息模板')
        this.initForm(false,"-")
    },
    // 模块初始化
    mounted(){
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
            let body = copyJSON(this.form.data)
            if (body.title == ""){
                return this.$message.warning("请输入模板名称");
            }
            if (body.template.http.url == ""){
                return this.$message.warning("请输入推送请求地址");
            }
            api.innerPost("/setting/message_set", body, (res) =>{
                if (!res.status){
                    return this.$message.error(res.message)
                }
                this.initForm(false)
                this.getList()
                api.dicList([Enum.dicMsg],()=>{}, true) // 存在变化，更新缓存
            })
        },
        // 执行一下
        runForm(){
            let body = copyJSON(this.form.data)
            if (body.template.http.url == ""){
                return this.$message.warning("请输入推送地址");
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
                    template:{
                        http:{
                            method: 'POST',
                            header: [{"key":"","value":""}],
                            url:'',
                            body:'',
                        },
                    }
                }
            }
            if ( typeof data === 'object' && data["id"] != undefined &&  data.id > 0 ){
                this.form.box.title = '编辑模板'
                this.form.data = data
                if (this.form.data.template.http.header.length == 0){
                    this.form.data.template.http.header.push({"key":"","value":""})
                }
            }
        },
        // http header 输入值变化时
        httpHeaderInput(val){
            if (val == ""){
                return
            }
            let item = this.form.data.template.http.header.slice(-1)[0]
            if (item == undefined || item.key != ""){
                this.form.data.template.http.header.push({"key":"","value":""})
            }
        },
        // http header 输入值删除
        httpHeaderDel(index){
            if ((index+1) >= this.form.data.template.http.header.length){ // 不允许删除最后一个空行
                return
            }
            this.form.data.template.http.header.splice(index,1)
        },
        close(){
            this.$emit('update:visible', false) // 向外传递关闭表示
        }
    }
})

Vue.component("MyMessageTemplate", MyMessageTemplate);