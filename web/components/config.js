var MyConfig = Vue.extend({
    template: `<el-main>
                    <el-menu :default-active="labelType" class="el-menu-demo" mode="horizontal" @select="handleClickTypeLabel">
                        <el-menu-item index="1" :disabled="listRequest">周期任务</el-menu-item>
                        <el-menu-item index="2" :disabled="listRequest">单次任务</el-menu-item>
                        <div style="float: right">
                            <el-button type="text" @click="createShow()">添加任务</el-button>
                            <el-button type="text" @click="getRegisterList()">已注册任务</el-button>
                        </div>
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



<!--                                <el-switch-->
<!--                                        validate-event="true"-->
<!--                                        :value="scope.row.status"-->
<!--                                        active-value="2"-->
<!--                                        inactive-value="1"-->
<!--                                        @change="changeStatus(scope.$index, scope.row.id, $event)"-->
<!--                                ></el-switch>-->
                            </template>
                        </el-table-column>
                        <el-table-column label="操作">
                            <template slot-scope="{row}">
                                <el-button type="text" @click="editShow(row)">编辑</el-button>
                                <el-button type="text" @click="changeStatus(row, 2)" v-if="row.status!=2">激活</el-button>
                                <el-button type="text" @click="changeStatus(row, 1)" v-if="row.status==2">停用</el-button>
                                <el-button type="text" @click="configLogBox(row.id, row.name)">日志</el-button>
                            </template>
<!--                            // 操作，再弹窗中进行操作（小型弹窗）；-->
<!--                            // 日志，侧边栏弹窗查看完整日志；-->
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
                    
                    
                    <!-- 任务设置表单 -->
                    <el-dialog :title="setConfigTitle" :visible.sync="setConfigShow" :close-on-click-modal="false">
                        <el-form :model="form">
                            <el-form-item label="活动名称*">
                                <el-input v-model="form.name"></el-input>
                            </el-form-item>
                            <el-form-item label="类型*">
                                <el-radio v-model="form.type" label="1">周期</el-radio>
                                <el-radio v-model="form.type" label="2">单次</el-radio>
<!--                                <el-select v-model="form.type" placeholder="请选请求方式" @change="changeType">-->
<!--                                    <el-option label="周期" value="1"></el-option>-->
<!--                                    <el-option label="单次" value="2"></el-option>-->
<!--                                </el-select>-->
                            </el-form-item>
            
            
                            <el-form-item label="时间*">
                                <el-input v-show="form.type==1" v-model="form.spec" :placeholder="hintSpec"></el-input>
                                <el-date-picker 
                                    style="width: 100%"
                                    v-show="form.type==2" 
                                    v-model="form.spec" 
                                    value-format="yyyy-MM-dd HH:mm:ss"
                                    type="datetime" 
                                    placeholder="选择运行时间" 
                                    :picker-options="pickerOptions">
                                </el-date-picker>
                            </el-form-item>
                            <el-form-item>
                                <el-tabs type="border-card" v-model="form.protocol">
                                    <el-tab-pane label="http" name="1">
                                        <el-form-item label="请求地址" class="http_url_box">
                                            <el-input v-model="form.command.http.url" placeholder="请输入http:// 或 https:// 开头的完整地址">
                                                <el-select v-model="form.command.http.method" placeholder="请选请求方式" slot="prepend">
                                                    <el-option label="GET" value="GET"></el-option>
                                                    <el-option label="POST" value="POST"></el-option>
                                                </el-select>
                                            </el-input>
                                        </el-form-item>
            
                                        <el-form-item label="请求Header" class="http_header_box">
                                            <el-input v-for="(header_v,header_i) in form.command.http.header" v-model="header_v.value" placeholder="参数值">
                                                <el-input v-model="header_v.key" slot="prepend" placeholder="参数名" @input="httpHeaderInput"></el-input>
                                                <el-button slot="append" icon="el-icon-delete" @click="httpHeaderDel(header_i)"></el-button>
                                            </el-input>
            
            
            
            <!--                                输入 k,v 两个输入域，后面有一个删除按钮；当key存在值时自动新增一个空的元素组。-->
            <!--                                method与url合并到一行（样式产考postApi7）。-->
                                        </el-form-item>
                                        <el-form-item label="请求Body参数">
                                            <el-input type="textarea" v-model="form.command.http.body" rows="5" placeholder="POST请求时body参数，将通过json进行请求发起"></el-input>
                                        </el-form-item>
                                    </el-tab-pane>
                                    <el-tab-pane label="rpc" name="2">
                                        <el-form-item label="请求模式">
                                            <el-radio v-model="form.command.rpc.method" label="GRPC">GRPC</el-radio>
                                        </el-form-item>
                                        <el-form-item label="proto">
                                            <el-input type="textarea" v-model="form.command.rpc.proto" rows="5" placeholder="请输入*.proto 文件内容"></el-input>
                                        </el-form-item>
                                        <el-form-item label="地址">
                                            <el-input v-model="form.command.rpc.addr" placeholder="请输入服务地址，含端口; 示例：localhost:21014"></el-input>
                                        </el-form-item>
                                        <el-form-item label="方法">
                                            <el-input v-model="form.command.rpc.action" placeholder="请输入服务方法; 示例：user.User/Echo"></el-input>
                                        </el-form-item>
                                        <el-form-item label="请求参数">
                                            <el-input type="textarea" v-model="form.command.rpc.body" rows="5" placeholder="请输入请求参数"></el-input>
                                        </el-form-item>
                                    </el-tab-pane>
                                    <el-tab-pane label="cmd" name="3">
                                        <el-input type="textarea" v-model="form.command.cmd" rows="5" :placeholder="sys_info.cmd_name + ': 请输入命令行执行内容'"></el-input>
                                    </el-tab-pane>
                                    <el-tab-pane label="sql" name="4" label-position="left">
                                        <el-form-item label="驱动">
                                            <el-select v-model="form.command.sql.driver" placeholder="驱动">
                                                <el-option label="mysql" value="mysql"></el-option>
                                            </el-select>
                                        </el-form-item>
                                        <el-form-item label="链接">
                                            <el-select v-model="form.command.sql.source.id" placement="请选择sql链接">
                                                <el-option v-for="(dic_v,dic_k) in dic_sql_source" :label="dic_v.name" :value="dic_v.id"></el-option>
                                            </el-select>
                                            <el-button type="text" style="margin-left: 20px" @click="sqlSourceBox(true)">设置链接</el-button>
                                        </el-form-item>
                                        <el-form-item label="执行语句">
                                            <div><el-button type="text" @click="sqlSetShow(-1,'')"><i class="el-icon-plus">添加</i></el-button></div>
                                            <div v-for="(statement,sql_index) in form.command.sql.statement" style="position: relative;max-height: 200px;line-height: 133%;background: #f4f4f5;margin-bottom: 10px;padding: 6px 20px 7px 8px;border-radius: 3px;">
                                                <pre style="margin: 0;overflow: auto;"><code class="language-sql hljs" style="min-height: 50px;">{{statement}}</code></pre>
                                                <i class="el-icon-delete" style="font-size: 15px;position: absolute;top: 2px;right: 2px;cursor:pointer" @click="sqlSetDel(sql_index)"></i>
                                                <i class="el-icon-edit" style="font-size: 15px;position: absolute;top: 23px;right: 2px;cursor:pointer" @click="sqlSetShow(sql_index,statement)"></i>
                                            </div>
                                            <el-alert v-show="form.command.sql.statement.length==0" title="未添加执行sql，请添加。" type="info"></el-alert>
                                        </el-form-item>
                                        <el-form-item label="错误行为">
                                            <el-tooltip class="item" effect="dark" content="执行到错误停止，之前的成功语句保留" placement="top-start">
                                              <el-radio v-model="form.command.sql.err_action" label="1">终止任务</el-radio>
                                            </el-tooltip>
                                            <el-tooltip class="item" effect="dark" content="执行到错误跳过，继续执行后面语句" placement="top-start">
                                              <el-radio v-model="form.command.sql.err_action" label="2">跳过继续</el-radio>
                                            </el-tooltip>
                                            <el-tooltip class="item" effect="dark" content="执行到错误，之前的成功语句也整体回滚" placement="top-start">
                                              <el-radio v-model="form.command.sql.err_action" label="3">事务回滚</el-radio>
                                            </el-tooltip>
                                        </el-form-item>
                                    </el-tab-pane>
                                </el-tabs>
                            </el-form-item>
                            <el-form-item label="备注">
                                <el-input v-model="form.remark"></el-input>
                            </el-form-item>
                        </el-form>
                        <div slot="footer" class="dialog-footer">
                            <el-button @click="setConfigShow = false">取 消</el-button>
                            <el-button type="primary" @click="setCron()">确 定</el-button>
                        </div>
                    </el-dialog>
                    <!-- 任务日志弹窗 -->
                    <el-drawer :title="configLog.title" :visible.sync="configLog.show" direction="rtl" size="40%" wrapperClosable="false">
                        <my-config-log :config_id="configLog.id"></my-config-log>
                    </el-drawer>
                    <!-- 注册任务列表弹窗 -->
                    <el-drawer title="已注册任务" :visible.sync="registerListShow" direction="rtl" size="40%" wrapperClosable="false">
                        <el-table :data="registerList">
                            <el-table-column property="id" label="编号"></el-table-column>
                            <el-table-column property="spec" label="执行时间"></el-table-column>
                            <el-table-column property="update_dt" label="下一次执行"></el-table-column>
                            <el-table-column property="name" label="任务名称"></el-table-column>
                            <el-table-column property="" label="">
                                <template slot-scope="scope">
                                    <el-button type="text" @click="configLogBox(scope.row.id, scope.row.name)">日志</el-button>
                                </template>
                            </el-table-column>
                        </el-table>
                    </el-drawer>
                    <!-- sql设置弹窗 -->
                    <el-dialog :title="'sql设置-'+sqlSet.title" :visible.sync="sqlSet.show" :show-close="false" :close-on-click-modal="false">
                        <el-input type="textarea" rows="12" placeholder="请输入sql内容，批量添加时多个sql请用;分号分隔。" v-model="sqlSet.data"></el-input>
                        <span slot="footer" class="dialog-footer">
                            <el-button @click="sqlSet.show = false">取 消</el-button>
                            <el-button type="primary" @click="sqlSetConfirm()">确 定</el-button>
                        </span>
                    </el-dialog>
            
                    <!-- sql链接源管理弹窗 -->
                    <el-drawer title="sql链接管理" :visible.sync="sqlSourceBoxShow" size="40%" wrapperClosable="false" :before-close="sqlSourceBox">
                        <my-sql-source :reload_list="sqlSourceBoxShow"></my-sql-source>
                    </el-drawer>
                </el-main>`,
    name: "MyConfig",
    props: {
        dic_sql_source:[],
        sys_info:{},
    },
    data(){
        return {
            dic_sql_source:[],
            list: [],
            listPage:{
                total:0,
                page: 1,
                size: 20,
            },
            listParam:{
                type: 1,
                page: 1,
                size: 20,
            },
            listRequest: false, // 请求中标志
            configLog:{
                show: false,
                id: 0,
                title:""
            },
            registerList: [],
            registerListShow: false,
            setConfigShow: false, // 新增弹窗
            setConfigTitle: '',
            sqlSourceBoxShow: false,
            sqlSet: {
                show: false, // 是否显示
                title: '添加',
                index: -1, // 操作行号
                data: "", // 实际内容
            }, // sql设置弹窗
            form:{},
            hintSpec: "* * * * * *",
            labelType: '1',
            // 日期选择器设置
            pickerOptions: {
                disabledDate(time){
                    return time.getTime() < Date.now() - 8.64e7
                },
                selectableRange: "00:00:00 - 23:01:59",
            },
        }
    },
    // 模块初始化
    created(){
        this.form = this.initFormData()
    },
    // 模块初始化
    mounted(){

        this.getList()
    },
    watch:{
        "form.spec":{
            handler(v){
                if (this.form.type == 2){ // 年月日的变化控制
                    let cur = moment()
                    if (moment(v).format('YYYY-MM-DD') == cur.format('YYYY-MM-DD')){
                        this.pickerOptions.selectableRange = `${cur.format('HH:mm:ss')} - 23:59:59`
                    }else{
                        this.pickerOptions.selectableRange = "00:00:00 - 23:59:59"
                    }
                }
            }
        }
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
                if (res.code != "000000"){
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
        // 添加/编辑 任务
        setCron(){
            if (this.form.name == ''){
                return this.$message.error('请输入任务名称')
            }else if (this.form.spec == ''){
                return this.$message.error('请输入任务执行时间')
            }else if (!this.form.protocol){
                return this.$message.error('请选择任务协议')
            }
            // else if (this.form.command == ''){
            //     return this.$message.error('请输入命令类容')
            // }

            // 主要是强制类型
            let body = {
                id: this.form.id,
                name: this.form.name,
                type: Number(this.form.type),
                spec: this.form.spec,
                protocol: Number(this.form.protocol),
                command: this.form.command,
                remark: this.form.remark,
            }
            body.command.sql.err_action = Number(body.command.sql.err_action)
            body.command.sql.source.id = Number(body.command.sql.source.id)

            api.innerPost("/config/set", body, (res)=>{
                if (!res.status){
                    return this.$message.error(res.message)
                }
                this.setConfigShow = false;
                this.form = this.initFormData()
                this.getList(); // 刷新当前页
            })
        },
        // 添加弹窗
        createShow(){
            this.setConfigShow = true
            this.setConfigTitle = '添加任务'
            // 应该有个地方定义空结构体
            this.form = this.initFormData()
        },
        // http header 输入值变化时
        httpHeaderInput(val){
            if (val == ""){
                return
            }
            let item = this.form.command.http.header.slice(-1)[0]
            if (item == undefined || item.key != ""){
                this.form.command.http.header.push({"key":"","value":""})
            }
        },
        // http header 输入值删除
        httpHeaderDel(index){
            if ((index+1) >= this.form.command.http.header.length){
                return
            }
            this.form.command.http.header.splice(index,1)
        },
        initFormData(){
            return  {
                type: '1',
                command:{
                    http:{
                        method: 'GET',
                        header: [{"key":"","value":""}],
                        url:'',
                        body:'',
                    },
                    rpc:{
                        method: 'GRPC',
                        addr: '',
                        action: '',
                        body: ''
                    },
                    cmd:'',
                    sql:{
                        driver: "mysql",
                        source:{
                            id: "",
                            // title: "",
                            // hostname:"",
                            // port:"",
                            // username:"",
                            // password:""
                        },
                        statement:[],
                        err_action: "1",
                    },
                }
            }
        },
        // 编辑弹窗
        editShow(row){
            this.setConfigShow = true
            this.setConfigTitle = '编辑任务'
            this.form = row
            if (this.form.command.sql.driver == ""){ // 历史数据兼容
                this.form.command.sql = this.initFormData().command.sql
            }
            if (this.form.command.sql.source.id == 0){
                this.form.command.sql.source.id = ""
            }
            if (this.form.command.http.header.length == 0){
                this.form.command.http.header = this.initFormData().command.http.header
            }
            this.form.type = this.form.type.toString()
            this.form.status = this.form.status.toString() // 这里要转字符串，否则可能显示数字
            this.form.protocol = this.form.protocol.toString()
            console.log("编辑：",this.form)
            this.changeType()
        },
        // 改变状态
        changeStatus(row, newStatus){
            this.$confirm(newStatus==1? '确认关闭任务': '确认开启任务', '提示',{
                type: 'warning',
            }).then(()=>{
                // 确认操作
                api.innerPost("/config/change_status", {id:row.id,status:Number(newStatus)}, (res)=>{
                    console.log("状态改变响应",res)
                    if (res.code != '000000'){
                        return this.$message.error(res.message)
                    }
                    row.status = newStatus.toString()
                    row.status_name = newStatus == 1 ? '停用' :'激活'
                    return this.$message.success(res.message)
                })
            }).catch(()=>{
                // 取消操作
            })
        },
        // 改变类型
        // changeType: function (){
        //     if (this.form.type == 2 || this.form.type == '2'){
        //         this.hintSpec = "YYYY-mm-dd HH:MM:SS"
        //     }else{
        //         this.hintSpec = "* * * * * *"
        //     }
        // },

        getRegisterList(){
            this.registerListShow = true;
            api.innerGet("/config/register_list", {}, (res)=>{
                if (res.code != "000000"){
                    return this.$message.error(res.message);
                }
                this.registerList = res.data.list;
            })
        },
        // 成功率图标
        getTopIcon(total, ratio){
            // 没有执行
            if (total == 0){
                return "icon el-icon-remove-outline" // 0
            }
            // 成功率100
            if (ratio == 100){
                return "icon el-icon-circle-check success"
            }
            // 成功率66以上
            if (ratio < 100 && ratio > 66){
                return "icon el-icon-warning-outline warning"
            }
            // 成功率33以上
            if (ratio < 66 && ratio > 33){
                return "icon el-icon-circle-plus-outline danger"
            }
            // 成功率33一下
            return "icon el-icon-circle-close danger"
        },
        // sql source box
        sqlSourceBox(show){
            this.sqlSourceBoxShow = show == true;
            if (!this.sqlSourceBoxShow){
                this.getDicSqlSource() // 关闭弹窗要重载枚举
            }
        },
        configLogBox(id, title){
            this.configLog.id = id
            this.configLog.title = title+' 日志'
            this.configLog.show = true
        },
        // sql设置弹窗
        sqlSetShow(index, oldData){
            if (index === "" || index == null || isNaN(index)){
                console.log('sqlSetShow', index, oldData)
                return this.$message.error("索引位标志异常");
            }
            if (typeof oldData != 'string'){
                return this.$message.error("sql内容异常");
            }
            this.sqlSet.show = true
            this.sqlSet.index = Number(index)  // -1.新增、>=0.具体行的编辑
            this.sqlSet.data = oldData // 原来的内容
            this.sqlSet.title = this.sqlSet.index < 0? '添加' : '编辑';
        },
        // sql设置确定
        sqlSetConfirm(){
            if (this.sqlSet.data == ""){
                return this.$message.error("sql内容不得为空");
            }
            if (typeof this.sqlSet.index !== 'number'){
                console.log('sqlSetShow', this.sqlSet)
                return this.$message.error("索引位标志异常");
            }
            // 支持批量添加
            let temp = this.sqlSet.data.split(";")
            let datas = []
            for (let i in temp){
                if (temp[i] != ""){
                    datas.push(temp[i])
                }
            }

            if (this.sqlSet.index < 0){
                this.form.command.sql.statement.push(...datas)
            }else if (datas.length > 1){
                return this.$message.warning("不支持单sql的拆分，建议删除后批量添加");
            }else if (datas.length == 0){
                return this.$message.warning("不存在有效sql，请确认输入");
            }else{
                this.form.command.sql.statement[this.sqlSet.index] = datas[0]
            }
            this.sqlSet.show = false
            this.sqlSet.data = ""
            this.sqlSet.index = -1
        },
        // 删除sql元素
        sqlSetDel(index){
            if (index === "" || index == null || isNaN(index)){
                console.log('sqlSetDel', index)
                return this.$message.error("索引位标志异常");
            }
            this.$confirm('此操作将删除sql执行语句，是否继续？','提示',{
                type:'warning',
            }).then(()=>{
                this.form.command.sql.statement.splice(index,1)
            })
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
        // 枚举
        getDicSqlSource(){
            api.innerFoundationDic(Enum.dicSqlSource,(res) =>{
                if (!res.status){
                    return this.$message.error(res.message);
                }
                this.dic_sql_source = res.data.maps[Enum.dicSqlSource].list
            })
        },
    }
})

Vue.component("MyConfig", MyConfig);