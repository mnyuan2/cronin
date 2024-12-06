var MyConfig = Vue.extend({
    template: `<el-container>
    <!--边栏-->
    <my-sidebar></my-sidebar>
    <!--主内容-->
    <el-main>
        <el-menu :default-active="listParam.type" class="el-menu-demo" mode="horizontal" @select="handleClickTypeLabel">
            <el-menu-item index="1" :disabled="listRequest">周期任务</el-menu-item>
            <el-menu-item index="2" :disabled="listRequest">单次任务</el-menu-item>
            <el-menu-item index="5" :disabled="listRequest">组件任务</el-menu-item>
            <div style="float: right">
                <el-button type="text" @click="setShow()" v-if="$auth_tag.config_set">添加任务</el-button>
            </div>
        </el-menu>
        <el-row>
            <el-form :inline="true" :model="listParam" size="mini" class="search-form">
                <el-form-item label="名称">
                    <el-input v-model="listParam.name" placeholder="搜索名称"></el-input>
                </el-form-item>
                <el-form-item label="协议">
                    <el-select v-model="listParam.protocol" placeholder="所有" multiple style="width: 140px">
                        <el-option v-for="item in dic.protocol" :label="item.name" :value="item.id"></el-option>
                    </el-select>
                </el-form-item>
                <el-form-item label="状态">
                    <el-select v-model="listParam.status" placeholder="所有" multiple style="width: 150px">
                        <el-option v-for="item in dic.config_status" :label="item.name" :value="item.id"></el-option>
                    </el-select>
                </el-form-item>
                <el-form-item label="处理人">
                    <el-select v-model="listParam.handle_user_ids" placeholder="所有" multiple style="width: 150px">
                        <el-option v-for="item in dic.user" :label="item.name" :value="item.id"></el-option>
                    </el-select>
                </el-form-item>
                <el-form-item label="创建人">
                    <el-select v-model="listParam.create_user_ids" placeholder="所有" multiple style="width: 150px">
                        <el-option v-for="item in dic.user" :label="item.name" :value="item.id"></el-option>
                    </el-select>
                </el-form-item>
                <el-form-item label="标签">
                    <el-select v-model="listParam.tag_ids" placeholder="所有" multiple style="width: 150px">
                        <el-option v-for="item in dic.tag" :label="item.name" :value="item.id"></el-option>
                    </el-select>
                </el-form-item>
                
                <el-form-item>
                    <el-button type="primary" @click="getList">查询</el-button>
                </el-form-item>
            </el-form>
        </el-row>
        <el-table :data="list" @cell-mouse-enter="listTableCellMouse" @cell-mouse-leave="listTableCellMouse">
            <el-table-column prop="" label="成功率" width="80">
                <template slot-scope="scope">
                    <el-tooltip placement="top-start">
                        <div slot="content">{{scope.row.topRatio}}%<br/>近期{{scope.row.top_number}}次执行，{{scope.row.top_error_number}}次失败。</div>
                        <i :class="getTopIcon(scope.row.top_number, scope.row.topRatio)"></i>
                    </el-tooltip>
                </template>
            </el-table-column>
            <el-table-column prop="spec" label="执行时间" v-show="listParam.type!=5" width="160"></el-table-column>
            <el-table-column prop="name" label="任务名称">
                <div slot-scope="{row}" class="abc" style="display: flex;">
                    <span style="white-space: nowrap;overflow: hidden;text-overflow: ellipsis;">
                        <router-link :to="{path:'/config_detail',query:{id:row.id, type:'config'}}" class="el-link el-link--primary is-underline" :title="row.name">{{row.name}}</router-link>
                    </span>
                    <span v-show="row.option.name.mouse" style="margin-left: 4px;white-space: nowrap;">
                        <i  class="el-icon-edit hover" @click="setShow(row)" title="编辑"></i>
                        <i class="el-icon-notebook-2 hover" @click="configLogBox(row)" title="日志"></i>
                    </span>
                </div>
            </el-table-column>
            <el-table-column prop="protocol_name" label="协议" width="80"></el-table-column>
            <el-table-column prop="" label="状态" width="100">
                <template slot-scope="scope">
                    <el-tooltip placement="top-start">
                        <div slot="content">{{scope.row.status_dt}}  {{scope.row.status_remark}}</div>
                        <el-button :type="statusTypeName(scope.row.status)" plain size="mini" round @click="statusShow(scope.row, 'config')">{{scope.row.status_name}}</el-botton>
                    </el-tooltip>
                </template>
            </el-table-column>
            <el-table-column prop="remark" label="备注"></el-table-column>
            <el-table-column prop="handle_user_names" label="处理人" width="120"></el-table-column>
            <el-table-column prop="create_user_name" label="创建人" width="80"></el-table-column>
            <el-table-column prop="tag_names" label="标签"></el-table-column>
        </el-table>
        <el-pagination
                @size-change="handleSizeChange"
                @current-change="handleCurrentChange"
                :current-page.sync="listPage.page"
                :page-sizes="[10, 50, 100]"
                :page-size="listPage.size"
                layout="total, sizes, prev, pager, next"
                :total="listPage.total">
        </el-pagination>
    </el-main>
    
    <!-- 任务设置表单 -->
    <el-dialog :title="add_box.title" :visible.sync="add_box.show" :close-on-click-modal="false" class="config-form-box":before-close="formClose">
        <my-config-form v-if="add_box.show" :request="{detail:add_box.detail}" @close="formClose"></my-config-form>
    </el-dialog>
    <!-- 任务日志弹窗 -->
    <el-drawer :title="config_log_box.title" :visible.sync="config_log_box.show" direction="rtl" size="40%" wrapperClosable="false" :before-close="configLogBoxClose">
        <my-config-log :tags="config_log_box.tags"></my-config-log>
    </el-drawer>
    <my-status-change v-if="status_box.show" :request="status_box" @close="statusShow"></my-status-change>

</el-container>`,
    name: "MyConfig",
    data(){
        return {
            dic:{
                config_status:[],
                user:[],
                protocol: [],
                tag: [],
            },
            // dic_sql_source:[],
            // dic_sql_driver:[],
            // dic_msg: [],
            // dic_jenkins_source: [],
            // dic_git_source: [],
            // dic_git_event:[],
            // dic_host_source: [],
            // dic_cmd_type: [],
            sys_info:{},
            auth_tags:{},
            list: [],
            listPage:{
                total:0,
                page: 1,
                size: 10,
            },
            listParam:{
                type: '1',
                page: 1,
                size: 10,
                protocol: [],
                status: [],
                handle_user_ids:[],
                create_user_ids:[],
                name: '',
            },
            listRequest: false, // 请求中标志
            config_log_box:{
                show: false,
                title:'',
                tags: {},
            },
            registerList: [],
            registerListShow: false,
            add_box: {
                show: false,
                title: '',
                detail:{},
            },
            sqlSourceBoxShow: false,
            sqlSet: {  // sql 设置弹窗
                show: false, // 是否显示
                title: '添加',
                index: -1, // 操作行号
                source: '', // sql 来源
                statement:{}
            },
            gitSet:{    // git 设置弹窗
                show: false, // 是否显示
                title: '添加',
                index: -1, // 操作行号
                data:{      // 单个设置详情
                    id: '',// 事件编号
                },
                // 合并类型
                gitMergeTypeList: {
                    "merge":"合并分支",
                    "squash":"扁平化分支",
                    "rebase":"变基并合并"
                },
            },
            msgSet:{
                show: false, // 是否显示
                title: '添加',
                index: -1, // 操作行号
                data: {}, // 实际内容
                statusList:[{id:1,name:"错误"}, {id:2, name:"成功"}], // {id:0,name:"完成"} 有歧义，以后待定
            },
            // 状态弹窗
            status_box:{
                show:false,
                detail:{},
                type: '',
            },
            form:{},
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
        setDocumentTitle('任务管理')
        api.systemInfo((res)=>{
            this.sys_info = res;
        })
        this.getDicSqlSource()
        this.loadParams(getHashParams(window.location.hash))
    },
    // 模块初始化
    mounted(){
        this.getList()
        this.auth_tags = cache.getAuthTags()
    },
    beforeDestroy(){},

    // 具体方法
    methods:{
        loadParams(param){
            if (typeof param !== 'object'){return}
            if (param.type){this.listParam.type = param.type.toString()}
            if (param.page){this.listParam.page = Number(param.page)}
            if (param.size){this.listParam.size = Number(param.size)}
            if (param.name){this.listParam.name = param.name.toString()}
            if (param.protocol){this.listParam.protocol = param.protocol.map(Number)}
            if (param.status){this.listParam.status = param.status.map(Number)}
            if (param.handle_user_ids){this.listParam.handle_user_ids = param.handle_user_ids.map(Number)}
            if (param.create_user_ids){this.listParam.create_user_ids = param.create_user_ids.map(Number)}
        },
        // 任务列表
        getList(){
            if (this.listRequest){
                return this.$message.info('请求执行中,请稍等.');
            }
            replaceHash('/config', this.listParam)
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
                    // if (res.data.list[i].command.sql){
                    //     res.data.list[i].command.sql.err_action = res.data.list[i].command.sql.err_action.toString()
                    //     res.data.list[i].command.sql.interval = res.data.list[i].command.sql.interval.toString()
                    // }
                    // res.data.list[i].status = res.data.list[i].status.toString()
                    res.data.list[i].topRatio = 100 - ratio * 100
                    // 处理人
                    let handle_user_names = ''
                    if (res.data.list[i].handle_user_ids){
                        res.data.list[i].handle_user_ids.forEach((id)=>{
                            this.dic.user.forEach((item)=>{
                                if (item.id == id){
                                    handle_user_names += item.name+','
                                }
                            })
                        })
                        res.data.list[i].handle_user_names = handle_user_names.substring(0,handle_user_names.length-1)
                    }
                    // 创建人
                    this.dic.user.forEach((item)=>{
                        if (item.id == res.data.list[i].create_user_id){
                            res.data.list[i].create_user_name = item.name
                        }
                    })
                    // 前端用设置
                    res.data.list[i].option = {
                        name:{
                            mouse:false
                        },
                    }
                }
                this.list = res.data.list;
                this.listPage = res.data.page;
            })
        },
        handleSizeChange(val) {
            console.log(`每页 ${val} 条`);
            this.listParam.size = val
        },
        handleCurrentChange(val) {
            this.listParam.page = val
            this.getList()
        },
        handleClickTypeLabel(tab, event) {
            this.listParam.type = tab
            this.listParam.page = 1
            this.listParam.total = 0
            this.getList()
        },

        // 列表表格列鼠标悬浮事件
        listTableCellMouse(row, column, cell, event){
            // console.log("列表", row, column, cell, event)
            if (row.option[column.property]){
                if (event.type === 'mouseenter'){
                    row.option[column.property].mouse = true
                }else if (event.type === 'mouseleave'){
                    row.option[column.property].mouse = false
                }
            }
        },

        // 添加弹窗
        setShow(row=null){
            this.add_box.show = true
            if (row == null){
                this.add_box.title= '添加任务'
                this.add_box.detail = {
                    type: this.listParam.type.toString()
                }
            }else{
                this.add_box.title= '编辑任务'
                api.innerGet('/config/detail', {id: row.id}, (res)=>{
                    if (!res.status){
                        return this.$message.error(res.message)
                    }
                    this.add_box.detail = res.data
                },{async:false})
            }
        },
        formClose(e){
            console.log("close",e)
            this.add_box.show = false
            if (e.is_change){
                this.getList()
            }
        },
        // http header 输入值变化时
        // httpHeaderInput(val){
        //     if (val == ""){
        //         return
        //     }
        //     let item = this.form.command.http.header.slice(-1)[0]
        //     if (item == undefined || item.key != ""){
        //         this.form.command.http.header.push({"key":"","value":""})
        //     }
        // },
        // http header 输入值删除
        // httpHeaderDel(index){
        //     if ((index+1) >= this.form.command.http.header.length){
        //         return
        //     }
        //     this.form.command.http.header.splice(index,1)
        // },
        // jenkinsParam 输入值变化时
        // jenkinsParamInput(val){
        //     if (val == ""){
        //         return
        //     }
        //     let item = this.form.command.jenkins.params.slice(-1)[0]
        //     if (item == undefined || item.key != ""){
        //         this.form.command.jenkins.params.push({"key":"","value":""})
        //     }
        // },
        // jenkinsParam 输入值删除
        // jenkinsParamDel(index){
        //     if ((index+1) >= this.form.command.jenkins.params.length){
        //         return
        //     }
        //     this.form.command.jenkins.params.splice(index,1)
        // },
        // sqlGitPath 输入值变化时
        // sqlGitPathInput(val){
        //     if (val == ""){
        //         return
        //     }
        //     let item = this.sqlSet.statement.git.path.slice(-1)[0]
        //     if (item == undefined || item != ""){
        //         this.sqlSet.statement.git.path.push("")
        //     }
        // },
        // sqlGitPath 输入值删除
        // sqlGitPathDel(index){
        //     if ((index+1) >= this.sqlSet.statement.git.path.length){
        //         return
        //     }
        //     this.sqlSet.statement.git.path.splice(index,1)
        // },
        // 提交审核
        // auditedBox(type){
        //     if (type == 'open'){
        //         this.auditor_box.show = true
        //         this.auditor_box.user_id = ''
        //     }else if (type == 'close'){
        //         this.auditor_box.show = false
        //         this.auditor_box.user_id = ''
        //     }else if (type == 'submit'){
        //
        //         this.setCron( (id)=> {
        //             api.innerPost("/config/change_status", {id:id,status:5, auditor_user_id: Number(this.auditor_box.user_id)}, (res)=>{
        //                 if (!res.status){
        //                     return this.$message.error(res.message)
        //                 }
        //                 // row.status = '5'
        //                 // row.status_name = '待审核'
        //                 this.auditor_box.show = false
        //                 return this.$message.success(res.message)
        //             })
        //         })
        //     }else{
        //         return this.$message.warning('错误操作')
        //     }
        // },
        // 驳回
        // rejectBox(row){
        //     this.$prompt('请输入驳回理由', '驳回说明', {
        //         confirmButtonText: '确定',
        //         cancelButtonText: '取消',
        //         inputPattern: /^.+$/,
        //         inputErrorMessage: '驳回理由必填'
        //     }).then(({ value }) => {
        //         api.innerPost("/config/change_status?auth_type=audit", {id:row.id, status: 6, status_remark:value}, (res)=>{
        //             if (!res.status){
        //                 return this.$message.error(res.message)
        //             }
        //             row.status = '6'
        //             row.status_name = '驳回'
        //             return this.$message.success(res.message)
        //         })
        //     }).catch(() => {});
        // },

        // 改变状态
        // changeStatus(row, newStatus, newStatusName){
        //     this.$confirm('确认'+newStatusName+'任务', '提示', {
        //         type: 'warning',
        //     }).then(()=>{
        //         let path = "/config/change_status"
        //         if (newStatus == 2){
        //             path += "?auth_type=audit"
        //         }
        //         api.innerPost(path, {id:row.id,status:Number(newStatus)}, (res)=>{
        //             if (!res.status){
        //                 return this.$message.error(res.message)
        //             }
        //             row.status = newStatus.toString()
        //             row.status_name = newStatus == 1 ? '停用' :'激活'
        //             return this.$message.success(res.message)
        //         })
        //     }).catch(()=>{
        //         // 取消操作
        //     })
        // },
        // 改变类型
        // changeType: function (){
        //     if (this.form.type == 2 || this.form.type == '2'){
        //         this.hintSpec = "YYYY-mm-dd HH:MM:SS"
        //     }else{
        //         this.hintSpec = "* * * * * *"
        //     }
        // },

        // getRegisterList(){
        //     this.registerListShow = true;
        //     api.innerGet("/config/register_list", {}, (res)=>{
        //         if (!res.status){
        //             console.log("config/register_list 错误", res)
        //             return this.$message.error(res.message);
        //         }
        //         this.registerList = res.data.list;
        //     })
        // },
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
        // sqlSourceBox(show){
        //     this.sqlSourceBoxShow = show == true;
        //     if (!this.sqlSourceBoxShow){
        //         this.getDicSqlSource() // 关闭弹窗要重载枚举
        //     }
        // },
        configLogBox(item){
            let tags = {ref_id:item.id, component:"config"}
            this.config_log_box.tags = tags
            this.config_log_box.title = item.name+' 日志'
            this.config_log_box.show = true
        },
        configLogBoxClose(done){
            this.config_log_box.show = false;
            this.config_log_box.id = 0;
            this.config_log_box.title = ' 日志'
            this.config_log_box.tags = {}
        },
        // sql设置弹窗
        // sqlSetShow(index, oldData){
        //     if (index === "" || index == null || isNaN(index)){
        //         console.log('sqlSetShow', index, oldData)
        //         return this.$message.error("索引位标志异常");
        //     }
        //     this.sqlSet.source = this.form.command.sql.origin
        //     this.sqlSet.show = true
        //     this.sqlSet.index = Number(index)  // -1.新增、>=0.具体行的编辑
        //     if (oldData == undefined){
        //         oldData = {
        //             type: this.form.command.sql.origin,
        //             git:null,
        //             local: "",
        //             is_batch: "1", // 批量解析：1.是（默认）、2.否
        //         }
        //     }
        //     if (oldData.git == null){
        //         oldData.git = {
        //             link_id: "",
        //             owner: "",
        //             project: "",
        //             path: [""],
        //             ref: "",
        //         }
        //     }
        //     oldData.is_batch = oldData.is_batch.toString()
        //     this.sqlSet.statement = oldData
        //
        //     this.sqlSet.title = this.sqlSet.index < 0? '添加' : '编辑';
        // },
        // sql 设置确定
        // sqlSetConfirm(){
        //     if (typeof this.sqlSet.index !== 'number'){
        //         console.log('sqlSetShow', this.sqlSet)
        //         return this.$message.error("索引位标志异常");
        //     }
        //
        //     if (this.sqlSet.source == 'git'){
        //         if (this.sqlSet.statement.git.link_id == ""){
        //             return this.$message.error("请选择连接")
        //         }
        //
        //         if (this.sqlSet.statement.git.owner == ""){
        //             return this.$message.error("仓库空间为必填")
        //         }
        //         if (this.sqlSet.statement.git.project == ""){
        //             return this.$message.error("项目名称为必填")
        //         }
        //         if (this.sqlSet.statement.git.path.length == 0){
        //             return this.$message.error("请输入文件的路径")
        //         }
        //         let data = copyJSON(this.sqlSet.statement);
        //         data.git.link_id = Number(data.git.link_id)
        //         data.is_batch = Number(data.is_batch)
        //         data.type = this.sqlSet.source
        //         if (data.git.ref == ""){
        //             data.git.ref = 'master'
        //         }
        //         data = this.sqlGitBuildDesc(data)
        //         if (this.sqlSet.index < 0){
        //             this.form.command.sql.statement.push(data)
        //         }else{
        //             this.form.command.sql.statement[this.sqlSet.index] = data
        //         }
        //         console.log("git data",data)
        //
        //     }else{
        //         if (this.sqlSet.statement.local == ""){
        //             return this.$message.error("sql内容不得为空");
        //         }
        //         // 支持批量添加
        //         let temp = this.sqlSet.statement.local.split(";")
        //         let datas = []
        //         for (let i in temp){
        //             let val = temp[i].trim()
        //             if (val != ""){
        //                 datas.push({"type":this.sqlSet.source, "git":{}, "local":val,"is_batch":1})
        //             }
        //         }
        //
        //         if (this.sqlSet.index < 0){
        //             this.form.command.sql.statement.push(...datas)
        //         }else if (datas.length > 1){
        //             return this.$message.warning("不支持单sql的拆分，建议删除后批量添加");
        //         }else if (datas.length == 0){
        //             return this.$message.warning("不存在有效sql，请确认输入");
        //         }else{
        //             this.form.command.sql.statement[this.sqlSet.index] = datas[0]
        //         }
        //     }
        //
        //     this.sqlSet.show = false
        //     // this.sqlSet.statement = {}
        //     this.sqlSet.index = -1
        // },
        // 删除sql元素
        // sqlSetDel(index){
        //     if (index === "" || index == null || isNaN(index)){
        //         console.log('sqlSetDel', index)
        //         return this.$message.error("索引位标志异常");
        //     }
        //     this.$confirm('此操作将删除sql执行语句，是否继续？','提示',{
        //         type:'warning',
        //     }).then(()=>{
        //         this.form.command.sql.statement.splice(index,1)
        //     })
        // },
        // git 设置弹窗
        // gitSetShow(index, oldData){
        //     if (index === "" || index == null || isNaN(index)){
        //         return this.$message.error("索引位标志异常"+index);
        //     }
        //     let data = {id:''}
        //     if (oldData != undefined){
        //         data = copyJSON(oldData)
        //     }
        //
        //     this.gitSet.show = true
        //     this.gitSet.title = this.gitSet.index < 0? '添加' : '编辑';
        //     this.gitSet.index = Number(index)  // -1.新增、>=0.具体行的编辑
        //     this.gitSet.data = data
        // },
        // git 设置确定
        // gitSetConfirm(){
        //     console.log("git 提交",this.gitSet)
        //     if (typeof this.gitSet.index !== 'number'){
        //         console.log('gitSetShow', this.sqlSet)
        //         return this.$message.error("索引位标志异常");
        //     }
        //
        //     if (this.gitSet.data.id == 2){
        //         // 待完善...
        //     }else if (this.gitSet.data.id == 9){
        //         if (this.gitSet.data.pr_merge.owner == ""){
        //             return this.$message.error("仓库空间为必填")
        //         }
        //         if (this.gitSet.data.pr_merge.repo == ""){
        //             return this.$message.error("项目名称为必填")
        //         }
        //         if (!this.gitSet.data.pr_merge.number){
        //             return this.$message.error("仓库PR编号为必填")
        //         }
        //         if (this.gitSet.data.pr_merge.merge_method == ""){
        //             return this.$message.error("合并方式不得为空")
        //         }
        //         this.gitSet.data.pr_merge.number = Number(this.gitSet.data.pr_merge.number)
        //     }else{
        //         return this.$message.error("未选择有效事件");
        //     }
        //
        //     let data = this.gitBuildDesc(copyJSON(this.gitSet.data))
        //     if (this.gitSet.index < 0){
        //         this.form.command.git.events.push(data)
        //     }else{
        //         this.form.command.git.events[this.gitSet.index] = data
        //     }
        //     this.gitSet.id = ''
        //     this.gitSet.show = false
        //     this.gitSet.index = -1
        // },
        // git 删除元素
        // gitSetDel(index){
        //     if (index === "" || index == null || isNaN(index)){
        //         console.log('gitSetDel', index)
        //         return this.$message.error("索引位标志异常");
        //     }
        //     this.$confirm('此操作将删除sql执行语句，是否继续？','提示',{
        //         type:'warning',
        //     }).then(()=>{
        //         this.form.command.git.events.splice(index,1)
        //     })
        // },

        // 推送弹窗
        // msgBoxShow(index, oldData){
        //     if (index === "" || index == null || isNaN(index)){
        //         console.log('msgSetShow', index, oldData)
        //         return this.$message.error("索引位标志异常");
        //     }
        //     if (oldData == undefined || index < 0){
        //         oldData = {
        //             status: 1,
        //             msg_id: "",
        //             notify_user_ids: [],
        //         }
        //     }else if (typeof oldData != 'object'){
        //         console.log('推送信息异常', oldData)
        //         return this.$message.error("推送信息异常");
        //     }
        //     this.msgSet.show = true
        //     this.msgSet.index = Number(index)  // -1.新增、>=0.具体行的编辑
        //     this.msgSet.title = this.msgSet.index < 0? '添加' : '编辑';
        //     this.msgSet.data = oldData
        // },
        // msgSetDel(index){
        //     if (index === "" || index == null || isNaN(index)){
        //         console.log('msgSetDel', index)
        //         return this.$message.error("索引位标志异常");
        //     }
        //     this.$confirm('确认删除推送配置','提示',{
        //         type:'warning',
        //     }).then(()=>{
        //         this.form.msg_set.splice(index, 1);
        //     })
        //
        // },
        // 推送确认
        // msgSetConfirm(){
        //     if (this.msgSet.data.msg_id <= 0){
        //         return this.$message.warning("请选择消息模板");
        //     }
        //     let data = this.msgSetBuildDesc(this.msgSet.data)
        //
        //     if (this.msgSet.index < 0){
        //         this.form.msg_set.push(data)
        //     }else{
        //         this.form.msg_set[this.msgSet.index] = data
        //     }
        //     this.msgSet.show = false
        //     this.msgSet.index = -1
        //     this.msgSet.data = {}
        // },
        // 构建消息设置描述
        // msgSetBuildDesc(data){
        //     let item1 = this.msgSet.statusList.find(option => option.id === data.status);
        //     if (item1){
        //         data.status_name = item1.name
        //     }
        //     let descrition = '当任务<span class="el-tag el-tag--small el-tag--light">'+item1.name+'</span>时'
        //
        //     let item2 = this.dic_msg.find(option => option.id === data.msg_id)
        //     if (item2){
        //         data.msg_name = item2.name
        //         descrition += '，发送<span class="el-tag el-tag--small el-tag--light">'+item2.name+'</span>消息'
        //     }
        //     let item3 = this.dic_user.filter((option) => {
        //         return data.notify_user_ids.includes(option.id);
        //     }).map((item)=>{return item.name})
        //     if (item3.length > 0){
        //         data.notify_users_name = item3
        //         descrition += '，并且@人员<span class="el-tag el-tag--small el-tag--light">'+data.notify_users_name+'</span>'
        //     }
        //     data.descrition = descrition
        //     return data
        // },
        // 构建git sql 描述
        // sqlGitBuildDesc(data){
        //     if (data.git == null){
        //         return data
        //     }
        //     let git = this.dic_git_source.filter(item=>{
        //         return item.id == data.git.link_id
        //     })
        //
        //     data.git.descrition = '连接<span class="el-tag el-tag--small el-tag--light">'+git.length>0 ? git[0].name: '' +
        //         '</span> 访问<span class="el-tag el-tag--small el-tag--light">'+data.git.owner+'/'+data.git.project +'</span>'+
        //         '</span> 引用<span class="el-tag el-tag--small el-tag--light">'+data.git.ref +'</span> 拉取以下文件内容'+
        //         '<span class="el-tag el-tag--small el-tag--light">'+ (data.is_batch==1? '批量解析后' :'单文件单sql') +'</span>执行'
        //     return data
        // },
        // 构建 git 描述
        // gitBuildDesc(data){
        //     let left = '<span class="el-tag el-tag--small el-tag--light">'
        //     let right = '</span>'
        //     switch (data.id){
        //         case 2:
        //             data.summary = '完善中...'
        //             break
        //         case 9:
        //             data.summary = `<b>pr合并</b> <a href="https://gitee.com/${data.pr_merge.owner}/${data.pr_merge.repo}/pulls/${data.pr_merge.number}" target="_blank" title="点击 查看pr详情"><i class="el-icon-connection"></i></a> ${left}${data.pr_merge.owner}/${data.pr_merge.repo}${right}/pulls/${left}${data.pr_merge.number}${right} ${left}${this.gitSet.gitMergeTypeList[data.pr_merge.merge_method]}${right}  ${data.pr_merge.prune_source_branch===true?left+'删除提交分支'+right:''}`+
        //                 `<br><i style="margin-left: 3em;"></i><b>${data.pr_merge.title}</b> ${data.pr_merge.description}`
        //             break
        //         default:
        //             data.summary = '未支持的事件类型'
        //     }
        //     return data
        // },
        statusShow(e, type=''){
            this.status_box.show = e.id > 0
            if (this.status_box.show){
                this.status_box.detail = e
                this.status_box.type = type
            }else{
                this.status_box.detail = {}
                this.status_box.type = ''
            }
            if (typeof e == 'object' && e.is_change===true){ // 数据变更了重新加载
                this.getList()
            }
        },
        // 枚举
        getDicSqlSource(){
            let types = [
                // Enum.dicSqlSource,
                // Enum.dicSqlDriver,
                Enum.dicConfigStatus,
                Enum.dicProtocolType,
                // Enum.dicJenkinsSource,
                // Enum.dicGitSource,
                // Enum.dicGitEvent,
                // Enum.dicHostSource,
                // Enum.dicCmdType,
                Enum.dicUser,
                Enum.dicTag,
                // Enum.dicMsg
            ]
            api.dicList(types,(res) =>{
                // this.dic_sql_source = res[Enum.dicSqlSource]
                // this.dic_jenkins_source = res[Enum.dicJenkinsSource]
                // this.dic_git_source =res[Enum.dicGitSource]
                // this.dic_git_event = res[Enum.dicGitEvent]
                // this.dic_host_source =res[Enum.dicHostSource]
                this.dic.user = res[Enum.dicUser]
                this.dic.tag = res[Enum.dicTag]
                // this.dic_msg = res[Enum.dicMsg]
                // this.dic_cmd_type = res[Enum.dicCmdType]
                // this.dic_sql_driver = res[Enum.dicSqlDriver]
                this.dic.config_status = res[Enum.dicConfigStatus]
                this.dic.protocol = res[Enum.dicProtocolType]

            })
        },

    }
})

Vue.component("MyConfig", MyConfig);